package tutor

import (
	"context"
	"errors"
	"io"
	"time"

	service2 "github.com/babadro/tutor/internal/core/service"
	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog/hlog"
)

type service interface {
	SendMessage(ctx context.Context, message, userID string, timestamp int64, chatID string) (string, swagger.Chat, error)
	SendVoiceMessage(
		ctx context.Context, voiceMsgBytes []byte, userID, chatID string, timestamp int64,
	) (models.SendVoiceMessageResult, error)
	GetChatMessages(
		ctx context.Context, chatID string, userID string, limit int32, timestamp int64,
	) ([]*swagger.ChatMessage, error)
	GetChats(ctx context.Context, userID string, limit int32, timestamp int64) ([]*swagger.Chat, error)
}

type Tutor struct {
	svc service
}

func NewTutor(svc service) *Tutor {
	return &Tutor{svc: svc}
}

func (t *Tutor) SendChatMessage(
	params operations.SendChatMessageParams, principal *models.Principal,
) middleware.Responder {
	if *params.Body.Text == "" {
		hlog.FromRequest(params.HTTPRequest).Error().Msg("Empty message")
		return operations.NewSendChatMessageBadRequest()
	}

	reply, createdChat, err := t.svc.SendMessage(
		params.HTTPRequest.Context(), *params.Body.Text, principal.UserID, *params.Body.Timestamp, params.Body.ChatID,
	)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send message")
	}

	var chat *swagger.Chat
	if createdChat.ChatID != "" {
		chat = &createdChat
	}

	return operations.NewSendChatMessageOK().WithPayload(&operations.SendChatMessageOKBody{
		Reply:     reply,
		Chat:      chat,
		Timestamp: time.Now().UnixMilli(),
	})
}

func (t *Tutor) SendVoiceMessage(
	params operations.SendVoiceMessageParams, principal *models.Principal,
) middleware.Responder {
	// todo check if the userID matches with the chatID, otherwise return unauthorized

	// read file to []byte
	voiceMsgBytes, err := io.ReadAll(params.File)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to read voice message")
		return operations.NewSendVoiceMessageBadRequest()
	}

	// log file length of readcloser
	hlog.FromRequest(params.HTTPRequest).Info().Msgf("File length: %d", len(voiceMsgBytes))

	chatID := ""
	if params.ChatID != nil {
		chatID = *params.ChatID
	}

	result, err := t.svc.SendVoiceMessage(params.HTTPRequest.Context(), voiceMsgBytes, principal.UserID, chatID, params.Timestamp)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send voice message")
		return operations.NewSendVoiceMessageBadRequest()
	}

	var chat *swagger.Chat
	if result.CreatedChat.ChatID != "" {
		chat = &result.CreatedChat
	}

	return operations.NewSendVoiceMessageOK().WithPayload(&operations.SendVoiceMessageOKBody{
		UsrAudio:   result.UserAudioURL,
		UsrTxt:     result.UserText,
		UsrTime:    params.Timestamp,
		ReplyAudio: result.LLMAudioURL,
		ReplyTxt:   result.LLMText,
		ReplyTime:  result.LLMTimestamp,
		Chat:       chat,
	})
}

func (t *Tutor) GetChatMessages(
	params operations.GetChatMessagesParams, principal *models.Principal,
) middleware.Responder {
	messages, err := t.svc.GetChatMessages(
		params.HTTPRequest.Context(), params.ChatID, principal.UserID, *params.Limit, *params.Timestamp,
	)
	if err != nil {
		if errors.Is(err, service2.ErrUserNotAuthorizedToViewThisChat) {
			return operations.NewGetChatMessagesUnauthorized()
		}

		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to get chat messages")

		return operations.NewGetChatMessagesBadRequest()
	}

	return operations.NewGetChatMessagesOK().WithPayload(&operations.GetChatMessagesOKBody{Messages: messages})
}

func (t *Tutor) GetChats(params operations.GetChatsParams, principal *models.Principal) middleware.Responder {
	chats, err := t.svc.GetChats(params.HTTPRequest.Context(), principal.UserID, *params.Limit, *params.Timestamp)
	if err != nil {
		hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to get chats")
		return operations.NewGetChatsBadRequest()
	}

	return operations.NewGetChatsOK().WithPayload(&operations.GetChatsOKBody{Chats: chats})
}
