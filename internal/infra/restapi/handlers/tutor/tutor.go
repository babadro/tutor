package tutor

import (
	"context"
	"fmt"

	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/go-openapi/runtime/middleware"

	hlog "github.com/rs/zerolog/hlog"
)

type service interface {
	SendMessage(ctx context.Context, message string) (string, error)
}

type Tutor struct {
	svc service
}

func NewTutor(svc service) *Tutor {
	return &Tutor{svc: svc}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams, principal *models.Principal) middleware.Responder {
	fmt.Println("Got message from flutter: ", params.Body.Message)

	reply, err := t.svc.SendMessage(params.HTTPRequest.Context(), params.Body.Message)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send message")
	} else {
		fmt.Println("Reply from LLM: ", reply)
	}

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply: reply,
	})
}
