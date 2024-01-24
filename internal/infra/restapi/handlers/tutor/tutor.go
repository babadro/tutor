package tutor

import (
	"fmt"

	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/go-openapi/runtime/middleware"
)

type service interface {
}

type Tutor struct {
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
