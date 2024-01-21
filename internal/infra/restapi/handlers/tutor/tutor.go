package tutor

import (
	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type service interface {
}

type Tutor struct {
}

func NewTutor(svc service) *Tutor {
	return &Tutor{}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams) middleware.Responder {
	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply: "Hello, world!",
	})
}
