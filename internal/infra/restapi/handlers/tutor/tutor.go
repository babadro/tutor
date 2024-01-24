package tutor

import (
	"context"
	"fmt"

	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/go-openapi/runtime/middleware"
	"github.com/tmc/langchaingo/llms"
)

type service interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

type Tutor struct {
	svc service
}

func NewTutor(svc service) *Tutor {
	return &Tutor{}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams, principal *models.Principal) middleware.Responder {
	fmt.Println(principal.Email)

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply: "Hello, world!",
	})
}
