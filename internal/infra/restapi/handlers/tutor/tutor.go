package tutor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"google.golang.org/api/option"
)

type service interface {
}

type Tutor struct {
}

func NewTutor(svc service) *Tutor {
	return &Tutor{}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams) middleware.Responder {
	_, err := verifyToken(params.HTTPRequest)
	if err != nil {
		fmt.Println(err)
	}

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply: "Hello, world!",
	})
}

func verifyToken(r *http.Request) (*auth.Token, error) {
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

	// Extract the ID Token from the Authorization header
	idToken := r.Header.Get("Authorization")
	// Usually, the token is prefixed with "Bearer ", you might need to remove this part.

	idToken = strings.TrimPrefix(idToken, "Bearer ")

	// Verify the ID Token
	decodedToken, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}

	return decodedToken, nil
}
