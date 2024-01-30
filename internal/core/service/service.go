package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/babadro/tutor/internal/models"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
)

type Service struct {
	llm          llm
	openaiClient *openai.Client
}

func NewService(llm llm, openaiClient *openai.Client) *Service {
	return &Service{llm: llm, openaiClient: openaiClient}
}

type llm interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error)
	CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
}

func (s *Service) SendMessage(ctx context.Context, message string) (string, error) {
	return s.llm.Call(ctx, message)
}

func (s *Service) SendVoiceMessage(ctx context.Context, voiceMsgFileUrl string, userEmail string) (models.SendVoiceMessageResult, error) {
	resp, err := http.Get(voiceMsgFileUrl)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to download voice message: %s", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to download voice message: %s", resp.Status)
	}

	// download voice message
	req := openai.AudioRequest{
		Model:  "whisper-1",
		Reader: resp.Body,
	}

	transcriptionResp, err := s.openaiClient.CreateTranscription(ctx, req)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to transcribe voice message: %s", err.Error())
	}

	voiceMessageTranscript, err := s.llm.Call(ctx, transcriptionResp.Text)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to voiceMessageTranscript from llm: %s", err.Error())
	}

	return models.SendVoiceMessageResult{
		VoiceMessageURL:         "https://firebasestorage.googleapis.com/v0/b/tutor-fq8fmu.appspot.com/o/How-are-you.mp3?alt=media&token=b1339bfe-6cf8-44ae-bda5-6bb6dcb43c5b",
		VoiceResponseURL:        "https://firebasestorage.googleapis.com/v0/b/tutor-fq8fmu.appspot.com/o/I-m-fine-thank-you.mp3?alt=media&token=ce525989-7c66-4ee3-8e21-68e4619fccc0",
		VoiceMessageTranscript:  voiceMessageTranscript,
		VoiceResponseTranscript: "I'm fine, thank you.",
	}, nil
}
