package service

import (
	"context"

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
