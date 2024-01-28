package service

import (
	"context"

	"github.com/babadro/tutor/internal/models"
	"github.com/tmc/langchaingo/llms"
)

type Service struct {
	llm llm
}

func NewService(llm llm) *Service {
	return &Service{llm: llm}
}

type llm interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error)
	CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
}

func (s *Service) SendMessage(ctx context.Context, message string) (string, error) {
	return s.llm.Call(ctx, message)
}

func (s *Service) SendVoiceMessage(ctx context.Context, voiceMsgFile []byte, userEmail string) (models.SendVoiceMessageResult, error) {
	return models.SendVoiceMessageResult{
		VoiceMessageURL:         "https://example.com/voice_message.wav",
		VoiceResponseURL:        "https://example.com/voice_response.wav",
		VoiceMessageTranscript:  "Hello, how are you?",
		VoiceResponseTranscript: "I'm fine, thank you.",
	}, nil
}
