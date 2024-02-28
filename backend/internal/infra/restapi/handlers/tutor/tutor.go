package tutor

import (
	"context"

	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog/hlog"
)

type service interface {
	SendMessage(ctx context.Context, message string) (string, error)
	SendVoiceMessage(ctx context.Context, voiceMsgFileUrl string, userID string) (models.SendVoiceMessageResult, error)
}

type Tutor struct {
	svc service
}

func NewTutor(svc service) *Tutor {
	return &Tutor{svc: svc}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams, principal *models.Principal) middleware.Responder {
	if params.Body.Message == "" {
		hlog.FromRequest(params.HTTPRequest).Error().Msg("Empty message")
		return operations.NewSendChatMessageBadRequest()
	}

	reply, err := t.svc.SendMessage(params.HTTPRequest.Context(), params.Body.Message)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send message")
	}

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply: reply,
	})
}

func (t *Tutor) SendVoiceMessage(params operations.SendVoiceMessageParams, principal *models.Principal) middleware.Responder {
	voiceMessage := params.Body.VoiceMessageURL

	// log voice message
	hlog.FromRequest(params.HTTPRequest).Info().Msgf("Voice message: %s", voiceMessage)

	result, err := t.svc.SendVoiceMessage(params.HTTPRequest.Context(), voiceMessage, principal.UserID)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send voice message")
		return operations.NewSendVoiceMessageBadRequest()
	}

	return operations.NewSendVoiceMessageOK().WithPayload(&operations.SendVoiceMessageOKBody{
		VoiceMessageTranscript:  result.VoiceMessageTranscript,
		VoiceMessageURL:         result.VoiceMessageURL,
		VoiceResponseTranscript: result.VoiceResponseTranscript,
		VoiceResponseURL:        result.VoiceResponseURL,
	})
}

func (t *Tutor) GetChatMessages(params operations.GetChatMessagesParams, principal *models.Principal) middleware.Responder {
	// return mocked messages
	messages := []*swagger.ChatMessage{
		{
			IsFromAi:  true,
			Text:      "Hello, I'm a tutor bot. How can I help you?",
			Timestamp: 1631535500,
		},
		{
			IsFromAi:  false,
			Text:      "I need help with my homework",
			Timestamp: 1631535510,
			UserID:    "user1",
		},
		{
			IsFromAi:  true,
			Text:      "Sure, I can help you with that. What's the problem?",
			Timestamp: 1631535520,
		},
		{
			IsFromAi:  false,
			Text:      "I don't understand the question",
			Timestamp: 1631535530,
			UserID:    "user1",
		},
		{
			IsFromAi:  true,
			Text:      "Let me see the question",
			Timestamp: 1631535540,
		},
		{
			IsFromAi:  false,
			Text:      "Here it is",
			Timestamp: 1631535550,
			UserID:    "user1",
		},
		{
			IsFromAi:  true,
			Text:      "I see. The question is asking for the derivative of the function. Let me calculate that for you",
			Timestamp: 1631535560,
		},
	}

	return operations.NewGetChatMessagesOK().WithPayload(messages)
}
