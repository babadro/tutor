package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"google.golang.org/api/iterator"
)

var (
	ErrUserNotAuthorizedToViewThisChat = errors.New("user is not authorized to view this chat")
)

type Service struct {
	llm             llm
	azureOpenai     *azopenai.Client
	baranovOpenai   *openai.Client
	storageClient   *storage.Client
	firestoreClient *firestore.Client
}

func NewService(
	llm llm,
	azureOpenai *azopenai.Client,
	baranovOpenai *openai.Client,
	storageClient *storage.Client,
	firestoreClient *firestore.Client,
) *Service {
	return &Service{
		llm:             llm,
		azureOpenai:     azureOpenai,
		baranovOpenai:   baranovOpenai,
		storageClient:   storageClient,
		firestoreClient: firestoreClient,
	}
}

type llm interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error)
	CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
}

func (s *Service) SendMessage(
	ctx context.Context, message, userID string, timestamp int64, chatID string,
) (string, swagger.Chat, error) {
	var newlyCreatedChat swagger.Chat

	if chatID == "" {
		newChat, _, err := s.firestoreClient.
			Collection("chats").
			Add(ctx, map[string]interface{}{
				"user_id": userID,
				"time":    timestamp,
				"title":   message,
			})

		if err != nil {
			return "", swagger.Chat{}, fmt.Errorf("unable to create chat: %s", err.Error())
		}

		newlyCreatedChat = swagger.Chat{
			ChatID: newChat.ID,
			Time:   timestamp,
			Title:  message,
		}

		chatID = newChat.ID
	}

	_, _, err := s.firestoreClient.Collection("messages").
		Add(ctx, map[string]interface{}{
			"text":    message,
			"user_id": userID,
			"chat_id": chatID,
			"time":    timestamp,
		})

	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to add user's message to firestore: %s", err.Error())
	}

	aiResponse, err := s.llm.Call(ctx, message)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to get AI response from llm: %s", err.Error())
	}

	_, _, err = s.firestoreClient.Collection("messages").
		Add(ctx, map[string]interface{}{
			"text":    aiResponse,
			"chat_id": chatID,
			"time":    time.Now().UnixMilli(),
		})

	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to add AI response message to firestore: %s", err.Error())
	}

	return aiResponse, newlyCreatedChat, nil
}

func (s *Service) SendVoiceMessage(
	ctx context.Context, userAudio []byte, userID string,
) (models.SendVoiceMessageResult, error) {
	/* todo uncomment when front end is ready
	req := azopenai.AudioTranscriptionOptions{
		File:           userAudio,
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),
		DeploymentName: to.Ptr("whisper-1"),
	}

	userTranscript, err := s.azureOpenai.GetAudioTranscription(ctx, req, nil)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to transcribe voice message: %s", err.Error())
	}

	if userTranscript.Text == nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("userTranscript.Text is nil")
	}

	userText := *userTranscript.Text

	llmReplyText, err := s.llm.Call(ctx, userText)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get llmReplyText from llm: %s", err.Error())
	}

	textToSpeechReq := openai.CreateSpeechRequest{
		Model:          "tts-1",
		Input:          llmReplyText,
		Voice:          openai.VoiceFable,
		ResponseFormat: "mp3",
	}

	llmAudio, err := s.baranovOpenai.CreateSpeech(ctx, textToSpeechReq)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get speech from text: %s", err.Error())
	}

	// convert llmAudio to []byte
	llmReplyAudio, err := io.ReadAll(llmAudio)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to read voice response: %s", err.Error())
	}

	// for llm audio we use mp3
	voiceResponseName := fmt.Sprintf("users/%s/backend_uploads/%d.mp3", userID, time.Now().UnixNano())

	llmReplyAudioURL, err := s.uploadFileToStorage(ctx, llmReplyAudio, voiceResponseName)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to upload voice response to storage: %s", err.Error())
	}

	// for user audio we use webm
	voiceMessageName := fmt.Sprintf("users/%s/backend_uploads/%d.webm", userID, time.Now().UnixNano())
	userAudioURL, err := s.uploadFileToStorage(ctx, userAudio, voiceMessageName)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to upload voice message to storage: %s", err.Error())
	}

	*/

	return models.SendVoiceMessageResult{
		UserAudioURL:     userAudioURL,
		LLMReplyAudioURL: llmReplyAudioURL,
		UserText:         userText,
		LLMText:          llmReplyText,
	}, nil
}

func (s *Service) uploadFileToStorage(
	ctx context.Context, fileBytes []byte, fileName string,
) (string, error) {
	bucket, err := s.storageClient.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("unable to get default bucket: %s", err.Error())
	}

	object := bucket.Object(fileName)
	wc := object.NewWriter(ctx)

	if _, err = wc.Write(fileBytes); err != nil {
		return "", fmt.Errorf("unable to write file to storage: %s", err.Error())
	}

	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("unable to close writer: %s", err.Error())
	}

	return generateFirebaseStorageURL(object.BucketName(), object.ObjectName()), nil
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

type ChatMessage struct {
	ChatID string `firestore:"chat_id"`
	Text   string `firestore:"text"`
	Time   int64  `firestore:"time"`
	UserID string `firestore:"user_id"`
}

func (s *Service) GetChatMessages(
	ctx context.Context, chatID string, userID string, limit int32, timestamp int64,
) ([]*swagger.ChatMessage, error) {
	docRef := s.firestoreClient.Collection("chats").Doc(chatID)

	// Attempt to retrieve the document
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chat by id: %s", err.Error())
	}

	// Read the user_id property from the document
	userIDinDB, err := docSnapshot.DataAt("user_id")
	if err != nil {
		return nil, fmt.Errorf("failed to read user_id from chat: %s", err.Error())
	}

	// Assert the type of userID if necessary, assuming it's a string
	userIDinDBStr, ok := userIDinDB.(string)
	if !ok {
		return nil, fmt.Errorf("user_id in chat is not a string")
	}

	if userIDinDBStr != userID {
		return nil, ErrUserNotAuthorizedToViewThisChat
	}

	var messages []ChatMessage

	query := s.firestoreClient.Collection("messages").
		Where("chat_id", "==", chatID).
		Where("time", ">=", timestamp).
		OrderBy("time", firestore.Asc).
		Limit(int(limit))

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		var doc *firestore.DocumentSnapshot
		doc, err = iter.Next()

		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return nil, fmt.Errorf("unable to get messages from firestore: %s", err.Error())
		}

		var message ChatMessage
		if err = doc.DataTo(&message); err != nil {
			return nil, fmt.Errorf("unable to get message data: %s", err.Error())
		}

		messages = append(messages, message)
	}

	// convert messages to swagger.ChatMessage
	var swaggerMessages []*swagger.ChatMessage
	for _, message := range messages {
		swaggerMessages = append(swaggerMessages, &swagger.ChatMessage{
			Text:              message.Text,
			Timestamp:         message.Time,
			UserID:            message.UserID,
			IsFromCurrentUser: message.UserID != "",
		})
	}

	return swaggerMessages, nil
}

type chat struct {
	Timestamp int64  `firestore:"time"`
	Title     string `firestore:"title"`
}

func (s *Service) GetChats(ctx context.Context, userID string, limit int32, timestamp int64) ([]*swagger.Chat, error) {
	query := s.firestoreClient.Collection("chats").
		Where("user_id", "==", userID).
		Where("time", ">=", timestamp).
		OrderBy("time", firestore.Desc).
		Limit(int(limit))

	iter := query.Documents(ctx)
	defer iter.Stop()

	var chats []*swagger.Chat

	for {
		doc, err := iter.Next()

		var chatModel chat

		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return nil, fmt.Errorf("unable to get chats from firestore: %s", err.Error())
		}

		if err = doc.DataTo(&chatModel); err != nil {
			return nil, fmt.Errorf("unable to get chat data: %s", err.Error())
		}

		chats = append(chats, &swagger.Chat{
			ChatID: doc.Ref.ID,
			Title:  chatModel.Title,
			Time:   chatModel.Timestamp,
		})
	}

	return chats, nil
}
