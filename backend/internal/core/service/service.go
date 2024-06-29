package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"google.golang.org/api/iterator"
)

var (
	ErrUserNotAuthorizedToViewThisChat = errors.New("user is not authorized to view this chat")
)

const (
	mp3AudioPathTemplate  = "users/%s/backend_uploads/%d.mp3"
	webmAudioPathTemplate = "users/%s/backend_uploads/%d.webm"
)

type Service struct {
	llm             llm
	azureOpenai     *azopenai.Client
	baranovOpenai   *openai.Client
	storageClient   *storage.Client
	firestoreClient *firestore.Client
	prompts         []string
}

func NewService(
	llm llm,
	azureOpenai *azopenai.Client,
	baranovOpenai *openai.Client,
	storageClient *storage.Client,
	firestoreClient *firestore.Client,
	prompts []string,
) *Service {
	return &Service{
		llm:             llm,
		azureOpenai:     azureOpenai,
		baranovOpenai:   baranovOpenai,
		storageClient:   storageClient,
		firestoreClient: firestoreClient,
		prompts:         prompts,
	}
}

type llm interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
	GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)
}

func (s *Service) SendMessage(
	ctx context.Context, message, userID string, timestamp int64, chatID string,
) (string, swagger.Chat, error) {
	chatType := GeneralChatType // todo
	newlyCreatedChat, err := s.saveMessageToDB(ctx, message, userID, chatID, "", chatType, timestamp)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to save message to db: %s", err.Error())
	}

	if chatID == "" {
		chatID = newlyCreatedChat.ChatID
	}

	aiResponse, err := s.llm.Call(ctx, message)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to get AI response from llm: %s", err.Error())
	}

	_, err = s.saveMessageToDB(ctx, aiResponse, "", chatID, "", chatType, time.Now().UnixMilli())
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to add AI response message to firestore: %s", err.Error())
	}

	return aiResponse, newlyCreatedChat, nil
}

// saveMessageToDB saves the message to the database and returns the newly created chat if it was created.
func (s *Service) saveMessageToDB(
	ctx context.Context, message, userID, chatID, audioURL string, chatTyp ChatType, timestamp int64,
) (swagger.Chat, error) {
	var newlyCreatedChat swagger.Chat

	if chatID == "" {
		chatTitle := cutChatTitle(message)

		newChat, _, err := s.firestoreClient.
			Collection("chats").
			Add(ctx, map[string]interface{}{
				"user_id": userID,
				"time":    timestamp,
				"title":   chatTitle,
				"type":    chatTyp,
			})

		if err != nil {
			return swagger.Chat{}, fmt.Errorf("unable to create chat: %s", err.Error())
		}

		newlyCreatedChat = swagger.Chat{
			ChatID: newChat.ID,
			Time:   timestamp,
			Title:  chatTitle,
			Typ:    swagger.ChatType(chatTyp),
		}

		chatID = newChat.ID
	}

	dbMessage := map[string]interface{}{
		"text":    message,
		"chat_id": chatID,
		"time":    timestamp,
	}

	if userID != "" {
		dbMessage["user_id"] = userID
	}

	if audioURL != "" {
		dbMessage["audio_url"] = audioURL
	}

	_, _, err := s.firestoreClient.Collection("messages").Add(ctx, dbMessage)

	if err != nil {
		return swagger.Chat{}, fmt.Errorf("unable to add user's message to firestore: %s", err.Error())
	}

	return newlyCreatedChat, nil
}

func cutChatTitle(str string) string {
	if len(str) > 20 {
		return str[:20]
	}

	return str
}

func (s *Service) SendVoiceMessage(
	ctx context.Context, userAudio []byte, userID string, chatID string, userMsgTimestamp int64,
) (models.SendVoiceMessageResult, error) {
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

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, userText),
	}

	llmReplyText, llmReplyAudioURL, err := s.generateTextAndAudioContent(ctx, content, userID, "gpt-3")
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get AI text and audio reply: %s", err.Error())
	}

	// for user audio we use webm
	userAudioName := fmt.Sprintf(webmAudioPathTemplate, userID, time.Now().UnixNano())

	userAudioURL, err := s.uploadFileToStorage(ctx, userAudio, userAudioName)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to upload voice message to storage: %s", err.Error())
	}

	chatType := GeneralChatType // todo

	createdChat, err := s.saveMessageToDB(ctx, userText, userID, chatID, userAudioURL, chatType, userMsgTimestamp)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to save user message to db: %s", err.Error())
	}

	if chatID == "" {
		chatID = createdChat.ChatID
	}

	llmReplyTimestamp := time.Now().UnixMilli()

	_, err = s.saveMessageToDB(ctx, llmReplyText, "", chatID, llmReplyAudioURL, chatType, llmReplyTimestamp)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to save llm reply to db: %s", err.Error())
	}

	return models.SendVoiceMessageResult{
		UserAudioURL: userAudioURL,
		LLMAudioURL:  llmReplyAudioURL,
		UserText:     userText,
		LLMText:      llmReplyText,
		CreatedChat:  createdChat,
		LLMTimestamp: llmReplyTimestamp,
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
	ChatID   string `firestore:"chat_id"`
	Text     string `firestore:"text"`
	Time     int64  `firestore:"time"`
	UserID   string `firestore:"user_id"`
	AudioURL string `firestore:"audio_url"`
}

func (c *ChatMessage) toSwagger() swagger.ChatMessage {
	return swagger.ChatMessage{
		Text:              c.Text,
		Timestamp:         c.Time,
		UserID:            c.UserID,
		IsFromCurrentUser: c.UserID != "",
		AudioURL:          c.AudioURL,
	}
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
			AudioURL:          message.AudioURL,
		})
	}

	return swaggerMessages, nil
}

type chat struct {
	ID               string
	UserID           string   `firestore:"user_id"`
	Timestamp        int64    `firestore:"time"`
	Title            string   `firestore:"title"`
	PreparedMessages []string `firestore:"prep_msgs"`
	Type             ChatType `firestore:"type"`
}

func (c *chat) toSwagger() swagger.Chat {
	return swagger.Chat{
		ChatID:           c.ID,
		Time:             c.Timestamp,
		Title:            c.Title,
		Typ:              swagger.ChatType(c.Type),
		PreparedMessages: c.PreparedMessages,
	}
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
			Title:  cutChatTitle(chatModel.Title),
			Time:   chatModel.Timestamp,
		})
	}

	return chats, nil
}

type ChatType int16

const (
	UnknownChatType                       = ChatType(0)
	GeneralChatType                       = ChatType(1)
	JobInterviewSeparateQuestionsChatType = ChatType(2)
)

func (s *Service) CreateChat(
	ctx context.Context, userID string, chatType ChatType, timestamp int64,
) (swagger.Chat, error) {
	if chatType != JobInterviewSeparateQuestionsChatType {
		return swagger.Chat{}, fmt.Errorf("can't create chat for type: %d", chatType)
	}

	createdChat, firstQuestion, err := s.createSeparateJobQuestionsChat(ctx, userID, timestamp)
	if err != nil {
		return swagger.Chat{}, fmt.Errorf("unable to create chat: %s", err.Error())
	}

	_, err = s.saveMessageToDB(
		ctx, firstQuestion.GermanText, userID, createdChat.ChatID, firstQuestion.GermanAudio, chatType, timestamp)

	if err != nil {
		return swagger.Chat{}, fmt.Errorf("unable to save first message to db: %s", err.Error())
	}

	return createdChat, nil
}

func (s *Service) createSeparateJobQuestionsChat(
	ctx context.Context, userID string, timestamp int64,
) (swagger.Chat, models.PreparedMessage, error) {
	query := s.firestoreClient.
		Collection("prepared_messages").
		Where("typ", "==", models.JobInterviewQuestion)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messageIDs []string

	for {
		doc, err := iter.Next()

		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return swagger.Chat{}, models.PreparedMessage{},
				fmt.Errorf("unable to get prepared_messages from firestore: %s", err.Error())
		}

		messageIDs = append(messageIDs, doc.Ref.ID)
	}

	rand.Shuffle(len(messageIDs), func(i, j int) {
		messageIDs[i], messageIDs[j] = messageIDs[j], messageIDs[i]
	})

	firstMsgID := messageIDs[0]

	// get first message from firestore
	doc, err := s.firestoreClient.Collection("prepared_messages").Doc(firstMsgID).Get(ctx)
	if err != nil {
		return swagger.Chat{}, models.PreparedMessage{},
			fmt.Errorf("unable to get first message from firestore: %s", err.Error())
	}

	var firstMsg models.PreparedMessage
	if err = doc.DataTo(&firstMsg); err != nil {
		return swagger.Chat{}, models.PreparedMessage{},
			fmt.Errorf("unable to get first message data: %s", err.Error())
	}

	createdChat := swagger.Chat{
		PreparedMessages: messageIDs,
		Time:             timestamp,
		Title:            cutChatTitle(firstMsg.GermanText),
		Typ:              swagger.ChatType(JobInterviewSeparateQuestionsChatType),
	}

	newChat, _, err := s.firestoreClient.
		Collection("chats").
		Add(ctx, chat{
			UserID:           userID,
			Type:             JobInterviewSeparateQuestionsChatType,
			Timestamp:        time.Now().UnixMilli(),
			Title:            cutChatTitle(firstMsg.GermanText),
			PreparedMessages: messageIDs,
		})

	if err != nil {
		return swagger.Chat{}, models.PreparedMessage{},
			fmt.Errorf("unable to create chat: %s", err.Error())
	}

	createdChat.ChatID = newChat.ID

	return createdChat, firstMsg, nil
}

func (s *Service) generateTextAndAudioContent(
	ctx context.Context, content []llms.MessageContent, userID, model string,
) (string, string, error) {
	resp, err := s.llm.GenerateContent(ctx, content,
		llms.WithMaxTokens(200),
		llms.WithModel(model),
	)

	if err != nil {
		return "", "", fmt.Errorf("unable to get content from llm: %s", err.Error())
	}

	if len(resp.Choices) == 0 {
		return "", "", errors.New("no choices in response")
	}

	text := resp.Choices[0].Content

	/* todo unkomment when front is ready
	textToSpeechReq := openai.CreateSpeechRequest{
		Model:          "tts-1",
		Input:          text,
		Voice:          openai.VoiceFable,
		ResponseFormat: "mp3",
	}

	llmAudio, err := s.baranovOpenai.CreateSpeech(ctx, textToSpeechReq)
	if err != nil {
		return "", "", fmt.Errorf("unable to get speech from text: %s", err.Error())
	}

	// convert llmAudio to []byte
	audio, err := io.ReadAll(llmAudio)
	if err != nil {
		return "", "", fmt.Errorf("unable to read voice message: %s", err.Error())
	}

	*/

	audio, err := os.ReadFile("cmd/server/sound_example.mp3")
	if err != nil {
		return "", "", fmt.Errorf("unable to read voice response: %s", err.Error())
	}

	// for llm audio we use mp3
	audioName := fmt.Sprintf(mp3AudioPathTemplate, userID, time.Now().UnixNano())

	audioURL, err := s.uploadFileToStorage(ctx, audio, audioName)
	if err != nil {
		return "", "", fmt.Errorf("unable to upload voice message to storage: %s", err.Error())
	}

	return text, audioURL, nil
}
