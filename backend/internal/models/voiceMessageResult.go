package models

type SendVoiceMessageResult struct {
	UserAudioURL     string
	UserText         string
	LLMReplyAudioURL string
	LLMText          string
}
