package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/babadro/tutor/internal/models"
	"github.com/tmc/langchaingo/llms"
)

type Service struct {
	llm          llm
	openaiClient *azopenai.Client
}

func NewService(llm llm, openaiClient *azopenai.Client) *Service {
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

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to file from body: %s", err.Error())
	}

	// download voice message
	req := azopenai.AudioTranscriptionOptions{
		File:           audioBytes,
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),
		DeploymentName: to.Ptr("whisper-1"),
	}

	voiceResponseTranscript, err := s.openaiClient.GetAudioTranscription(ctx, req, nil)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to transcribe voice message: %s", err.Error())
	}

	if voiceResponseTranscript.Text == nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("voiceResponseTranscript.Text is nil")
	}

	responseTranscript := *voiceResponseTranscript.Text

	llmTextResponse, err := s.llm.Call(ctx, responseTranscript)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get llmTextResponse from llm: %s", err.Error())
	}

	return models.SendVoiceMessageResult{
		VoiceMessageURL:         "https://firebasestorage.googleapis.com/v0/b/tutor-fq8fmu.appspot.com/o/How-are-you.mp3?alt=media&token=b1339bfe-6cf8-44ae-bda5-6bb6dcb43c5b",
		VoiceResponseURL:        "https://firebasestorage.googleapis.com/v0/b/tutor-fq8fmu.appspot.com/o/I-m-fine-thank-you.mp3?alt=media&token=ce525989-7c66-4ee3-8e21-68e4619fccc0",
		VoiceMessageTranscript:  responseTranscript,
		VoiceResponseTranscript: llmTextResponse,
	}, nil
}
