package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"firebase.google.com/go/storage"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/babadro/tutor/internal/models"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
)

type Service struct {
	llm           llm
	azureOpenai   *azopenai.Client
	baranovOpenai *openai.Client
	storageClient *storage.Client
}

func NewService(
	llm llm,
	azureOpenai *azopenai.Client,
	baranovOpenai *openai.Client,
	storageClient *storage.Client,
) *Service {
	return &Service{
		llm:           llm,
		azureOpenai:   azureOpenai,
		baranovOpenai: baranovOpenai,
		storageClient: storageClient,
	}
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

	voiceResponseTranscript, err := s.azureOpenai.GetAudioTranscription(ctx, req, nil)
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

	textToSpeechReq := openai.CreateSpeechRequest{
		Model:          "tts-1",
		Input:          "I like to play football.",
		Voice:          openai.VoiceFable,
		ResponseFormat: "mp3",
	}

	voiceResponse, err := s.baranovOpenai.CreateSpeech(ctx, textToSpeechReq)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get speech from text: %s", err.Error())
	}

	// convert voiceResponse to []byte
	voiceResponseBytes, err := io.ReadAll(voiceResponse)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to read voice response: %s", err.Error())
	}

	bucket, err := s.storageClient.DefaultBucket()
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get default bucket: %s", err.Error())
	}

	fileName := fmt.Sprintf("users/%s/backend_uploads/%d", userEmail, time.Now().UnixNano())

	object := bucket.Object(fileName)
	wc := object.NewWriter(ctx)

	if _, err = wc.Write(voiceResponseBytes); err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to write voice response to storage: %s", err.Error())
	}

	if err = wc.Close(); err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to close writer: %s", err.Error())
	}

	voiceResponseURL := generateFirebaseStorageURL(object.BucketName(), object.ObjectName())

	return models.SendVoiceMessageResult{
		// todo delete voiceMessageURL as we get it from the request
		VoiceMessageURL:         "https://firebasestorage.googleapis.com/v0/b/tutor-fq8fmu.appspot.com/o/How-are-you.mp3?alt=media&token=b1339bfe-6cf8-44ae-bda5-6bb6dcb43c5b",
		VoiceResponseURL:        voiceResponseURL,
		VoiceMessageTranscript:  responseTranscript,
		VoiceResponseTranscript: llmTextResponse,
	}, nil
}

func generateFirebaseStorageURL(bucketName, filePath string) string {
	baseURL := "https://firebasestorage.googleapis.com/v0/b/"
	storagePath := "o/"

	// URL encode the file path to handle special characters
	// Note: This is a simplified approach. You might need a more robust way to URL-encode paths,
	// especially if they contain slashes (/) or other special characters.
	encodedFilePath := url.PathEscape(filePath)

	// Construct the full URL
	fullURL := fmt.Sprintf("%s%s/%s%s?alt=media", baseURL, bucketName, storagePath, encodedFilePath)

	return fullURL
}
