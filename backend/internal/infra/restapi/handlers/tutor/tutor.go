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
		ctx context.Context, voiceMsgFileURL string, userID string,
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

	// log file length of readcloser
	fileLength, _ := io.Copy(io.Discard, params.File)
	hlog.FromRequest(params.HTTPRequest).Info().Msgf("File length: %d", fileLength)

	// log chatID
	hlog.FromRequest(params.HTTPRequest).Info().Msgf("ChatID: %s", params.ChatID)

	//voiceMessage := params.File

	// log voice message
	//hlog.FromRequest(params.HTTPRequest).Info().Msgf("Voice message: %s", voiceMessage)
	//
	//result, err := t.svc.SendVoiceMessage(params.HTTPRequest.Context(), voiceMessage, principal.UserID)
	//if err != nil {
	//	hlog.FromRequest(params.HTTPRequest).Error().Err(err).Msg("Unable to send voice message")
	//	return operations.NewSendVoiceMessageBadRequest()
	//}

	return operations.NewSendVoiceMessageOK().WithPayload(&operations.SendVoiceMessageOKBody{
		UsrAudio:   "audio.mp3",            // todo
		UsrTxt:     "text",                 // todo
		UsrTime:    time.Now().UnixMilli(), // todo
		ReplyAudio: "audio.mp3",            // todo
		ReplyTxt:   "text",                 // todo
		ReplyTime:  time.Now().UnixMilli(), // todo
		Chat:       nil,                    // todo
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
