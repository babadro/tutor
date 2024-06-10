package models

import "github.com/babadro/tutor/internal/models/swagger"

type SendVoiceMessageResult struct {
	UserAudioURL string
	UserText     string
	LLMAudioURL  string
	LLMText      string
	LLMTimestamp int64
	CreatedChat  swagger.Chat
}
