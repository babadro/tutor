package models

import "github.com/babadro/tutor/internal/models/swagger"

type SendVoiceMessageResult struct {
	UserAudioURL     string
	UserText         string
	LLMReplyAudioURL string
	LLMText          string
	CreatedChat      swagger.Chat
}
