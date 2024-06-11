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
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
	baranovOpenai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"

	"github.com/babadro/tutor/internal/core/service"
	"github.com/tmc/langchaingo/llms/openai"
)

//go:generate swagger generate server --target ../../../../tutor --name Tutor --spec ../../../swagger.yaml --model-package internal/models/swagger --server-package internal/infra/restapi --principal interface{} --exclude-main

type envVars struct {
	NgrokAgentAddr string   `env:"NGROK_AGENT_ADDR,required"`
	AllowedUsers   []string `env:"ALLOWED_USERS,required"`
	OpenaiAPIKey   string   `env:"OPENAI_API_KEY,required"`
	StorageBucket  string   `env:"STORAGE_BUCKET,required"`
}

func configureFlags(_ *operations.TutorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TutorAPI) http.Handler {
	ctx := context.Background()

	l := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	var envs envVars
	if err := env.Parse(&envs); err != nil {
		l.Fatal().Msgf("Unable to parse env vars: %v\n", err)
	}

	// configure the api here
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Initialize Firebase SDK
	firebaseConfig := &firebase.Config{
		StorageBucket: envs.StorageBucket,
	}
	opt := option.WithCredentialsFile("/app/secrets/tutor.json")
	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig, opt)

	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init firebase app")
	}

	// Get a Firebase Auth client from the Firebase App
	firebaseAuthClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init firebase client")
	}

	api.KeyAuth = getKeyAuthFunc(l, firebaseAuthClient, envs.AllowedUsers)

	llm, err := openai.New()
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init llm client")
	}

	openAICredential := azcore.NewKeyCredential(envs.OpenaiAPIKey)

	openaiClient, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", openAICredential, nil)
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init openai client")
	}

	baranovClient := baranovOpenai.NewClient(envs.OpenaiAPIKey)

	storageClient, err := firebaseApp.Storage(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init storage client")
	}

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to init firestore client")
	}

	prompts, err := readPrompts()
	if err != nil {
		l.Fatal().Err(err).Msg("Unable to read prompts")
	}

	tutorService := service.NewService(
		llm, openaiClient, baranovClient, storageClient, firestoreClient, prompts,
	)
	tutorAPI := tutor.NewTutor(tutorService)

	api.SendChatMessageHandler = operations.SendChatMessageHandlerFunc(tutorAPI.SendChatMessage)
	api.SendVoiceMessageHandler = operations.SendVoiceMessageHandlerFunc(tutorAPI.SendVoiceMessage)
	api.GetChatMessagesHandler = operations.GetChatMessagesHandlerFunc(tutorAPI.GetChatMessages)
	api.GetChatsHandler = operations.GetChatsHandlerFunc(tutorAPI.GetChats)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		err = firestoreClient.Close()
		if err != nil {
			l.Error().Err(err).Msg("Unable to close firestore client")
		}
	}

	return setupGlobalMiddleware(l, api.Serve(setupMiddlewares()))
}

func readPrompts() ([]string, error) {
	b, err := os.ReadFile("/app/secrets/prompts/job_interview/1.txt")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	return []string{string(b)}, nil
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
func setupMiddlewares() middleware.Builder {
	return nil
}

// The middleware configuration happens before anything,
// this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(
	l zerolog.Logger, handler http.Handler,
) http.Handler {
	return alice.New(
		middlewares.Logging(l),
		middlewares.Cors,
	).Then(handler)
}

func getKeyAuthFunc(
	l zerolog.Logger, firebaseAuthClient *auth.Client, allowedUsers []string,
) func(string) (*models.Principal, error) {
	return func(token string) (*models.Principal, error) {
		token = strings.TrimPrefix(token, "Bearer ")

		// Verify the ID Token
		decodedToken, err := firebaseAuthClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			l.Error().Msgf("error verifying ID token: %s", err.Error())
			return nil, fmt.Errorf("error verifying ID token: %s", err.Error())
		}

		email, ok := decodedToken.Claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("error getting email from token")
		}

		for _, user := range allowedUsers {
			if user == email {
				return &models.Principal{
					Email:  email,
					UserID: decodedToken.UID,
				}, nil
			}
		}

		l.Error().Msgf("Unauthorized user: %s", email)

		return nil, errors.New(http.StatusUnauthorized, "Unauthorized")
	}
}
