// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/babadro/tutor/internal/infra/restapi/handlers/tutor"
	"github.com/babadro/tutor/internal/infra/restapi/middlewares"
	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/caarlos0/env"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"google.golang.org/api/option"
)

//go:generate swagger generate server --target ../../../../tutor --name Tutor --spec ../../../swagger.yaml --model-package internal/models/swagger --server-package internal/infra/restapi --principal interface{} --exclude-main

type envVars struct {
	NgrokAgentAddr string `env:"NGROK_AGENT_ADDR,required"`
}

func configureFlags(api *operations.TutorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TutorAPI) http.Handler {
	l := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	var envs envVars
	if err := env.Parse(&envs); err != nil {
		l.Fatal().Msgf("Unable to parse env vars: %v\n", err)
	}

	_ = context.Background()

	tutorAPI := tutor.NewTutor(nil)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	firebaseClient, err := initFirebaseClient()
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init firebase client")
	}

	api.KeyAuth = func(token string) (*models.Principal, error) {
		token = strings.TrimPrefix(token, "Bearer ")

		// Verify the ID Token
		decodedToken, err := firebaseClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			return nil, fmt.Errorf("error verifying ID token: %v", err)
		}

		email, ok := decodedToken.Claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("error getting email from token")
		}

		return &models.Principal{
			Email: email,
		}, nil
	}

	api.SendChatMessageHandler = operations.SendChatMessageHandlerFunc(tutorAPI.SendChatMessage)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares(l)))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(_ *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
	_, _, _ = s, scheme, addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(l zerolog.Logger) middleware.Builder {
	return alice.New(
		middlewares.Logging(l),
	).Then
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func initFirebaseClient() (*auth.Client, error) {
	// Initialize Firebase SDK
	opt := option.WithCredentialsFile("/app/secrets/tutor.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	// Get a Firebase Auth client from the Firebase App
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	return client, nil
}
