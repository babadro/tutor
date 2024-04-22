package tutor

import (
	"context"
	"time"

	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog/hlog"
)

type service interface {
	SendMessage(ctx context.Context, message string, userID string, timestamp int64, chatID string) (string, error)
	SendVoiceMessage(ctx context.Context, voiceMsgFileUrl string, userID string) (models.SendVoiceMessageResult, error)
	GetChatMessages(ctx context.Context, chatID string, limit int32, timestamp int64) ([]*swagger.ChatMessage, error)
	GetChats(ctx context.Context, userID string, limit int32) ([]*swagger.Chat, error)
}

type Tutor struct {
	svc service
}

func NewTutor(svc service) *Tutor {
	return &Tutor{svc: svc}
}

func (t *Tutor) SendChatMessage(params operations.SendChatMessageParams, principal *models.Principal) middleware.Responder {
	if *params.Body.Text == "" {
		hlog.FromRequest(params.HTTPRequest).Error().Msg("Empty message")
		return operations.NewSendChatMessageBadRequest()
	}

	reply, err := t.svc.SendMessage(
		params.HTTPRequest.Context(), *params.Body.Text, principal.UserID, *params.Body.Timestamp, params.Body.ChatID,
	)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send message")
	}

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply:     reply,
		Timestamp: time.Now().UnixMilli(),
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
	/*
		messages := []*swagger.ChatMessage{
			{
				IsFromCurrentUser: false,
				Text:              "Hello, I'm a tutor bot. How can I help you?",
				Timestamp:         1631535500,
			},
			{
				IsFromCurrentUser: true,
				Text:              "I need help with my homework",
				Timestamp:         1631535510,
				UserID:            "user1",
			},
			{
				IsFromCurrentUser: false,
				Text:              "Sure, I can help you with that. What's the problem?",
				Timestamp:         1631535520,
			},
			{
				IsFromCurrentUser: true,
				Text:              "I don't understand the question",
				Timestamp:         1631535530,
				UserID:            "user1",
			},
			{
				IsFromCurrentUser: false,
				Text:              "Let me see the question",
				Timestamp:         1631535540,
			},
			{
				IsFromCurrentUser: true,
				Text:              "Here it is",
				Timestamp:         1631535550,
				UserID:            "user1",
			},
			{
				IsFromCurrentUser: false,
				Text:              "I see. The question is asking for the derivative of the function. Let me calculate that for you",
				Timestamp:         1631535560,
			},
		}
	*/
	messages, err := t.svc.GetChatMessages(params.HTTPRequest.Context(), params.ChatID, *params.Limit, *params.Timestamp)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to get chat messages")
		return operations.NewGetChatMessagesBadRequest()
	}

	//for _, message := range messages {
	//	hlog.FromRequest(params.HTTPRequest).Info().Msgf("Message: %v", *message)
	//}

	return operations.NewGetChatMessagesOK().WithPayload(&operations.GetChatMessagesOKBody{Messages: messages})
}

func (t *Tutor) GetChats(params operations.GetChatsParams, principal *models.Principal) middleware.Responder {
	chats, err := t.svc.GetChats(params.HTTPRequest.Context(), principal.UserID, *params.Limit)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to get chats")
		return operations.NewGetChatsBadRequest()
	}

	for _, chat := range chats {
		hlog.FromRequest(params.HTTPRequest).Info().Msgf("Chat: %s", chat.ChatID)
	}

	return operations.NewGetChatsOK().WithPayload(&operations.GetChatsOKBody{Chats: chats})
}
