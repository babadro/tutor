package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
	"unicode"

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

const llmModel = "gpt-4o-mini"

func (s *Service) SendMessage(
	ctx context.Context, userText, userID string, timestamp int64, chatID string,
) (string, swagger.Chat, error) {
	chatType := models.GeneralChatType
	if chatID != "" {
		userChat, err := s.getChatIfUserAutorized(ctx, chatID, userID)
		if err != nil {
			return "", swagger.Chat{}, fmt.Errorf("unable to get chat: %w", err)
		}

		if chatType, err = models.GetChatTypeFromNumber(userChat.Type); err != nil {
			return "", swagger.Chat{}, fmt.Errorf("unable to get chat type: %s", err.Error())
		}
	}

	newlyCreatedChat, err := s.saveMessageToDB(ctx, userText, userID, chatID, "", chatType, timestamp)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to save userText to db: %s", err.Error())
	}

	if chatID == "" {
		chatID = newlyCreatedChat.ChatID
	}

	llmIn, err := s.getLlmInput(userText, chatType)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to get llmInput: %s", err.Error())
	}

	aiResponse, err := s.generateTextContent(ctx, llmIn, llmModel)
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to get AI response from llm: %s", err.Error())
	}

	_, err = s.saveMessageToDB(ctx, aiResponse, "", chatID, "", chatType, time.Now().UnixMilli())
	if err != nil {
		return "", swagger.Chat{}, fmt.Errorf("unable to add AI response userText to firestore: %s", err.Error())
	}

	return aiResponse, newlyCreatedChat, nil
}

// saveMessageToDB saves the message to the database and returns the newly created chat if it was created.
func (s *Service) saveMessageToDB(
	ctx context.Context, message, userID, chatID, audioURL string, chatTyp models.ChatType, timestamp int64,
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
	ctx context.Context,
	userAudio []byte,
	userID string,
	chatID string,
	userMsgTimestamp int64,
	typ models.VoiceMsgType,
) (models.SendVoiceMessageResult, error) {
	chatType := models.GeneralChatType
	if chatID != "" {
		userChat, err := s.getChatIfUserAutorized(ctx, chatID, userID)
		if err != nil {
			return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get chat: %w", err)
		}

		if chatType, err = models.GetChatTypeFromNumber(userChat.Type); err != nil {
			return models.SendVoiceMessageResult{},
				fmt.Errorf("unable to get chat type: %s", err.Error())
		}
	}

	req := azopenai.AudioTranscriptionOptions{
		File:           userAudio,
		DeploymentName: to.Ptr("whisper-1"),
		Language:       getTranscriptLang(chatType),
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),
	}

	userTranscript, err := s.azureOpenai.GetAudioTranscription(ctx, req, nil)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to transcribe voice message: %s", err.Error())
	}

	if userTranscript.Text == nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("userTranscript.Text is nil")
	}

	userText := *userTranscript.Text

	// for user audio we use webm
	userAudioName := fmt.Sprintf(webmAudioPathTemplate, userID, time.Now().UnixNano())

	userAudioURL, err := s.uploadFileToStorage(ctx, userAudio, userAudioName)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to upload voice message to storage: %s", err.Error())
	}

	createdChat, err := s.saveMessageToDB(ctx, userText, userID, chatID, userAudioURL, chatType, userMsgTimestamp)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to save user message to db: %s", err.Error())
	}

	if chatID == "" {
		chatID = createdChat.ChatID
	}

	if typ == models.AwaitingCompletionVoiceMsgType {
		// then we should not answer to the message now
		return models.SendVoiceMessageResult{
			UserAudioURL: userAudioURL,
			UserText:     userText,
			CreatedChat:  createdChat,
		}, nil
	}

	llmIn, err := s.getLlmInput(userText, chatType)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get llmInput: %s", err.Error())
	}

	llmReplyText, llmReplyAudioURL, err := s.generateTextAndAudioContent(ctx, llmIn, userID, llmModel)
	if err != nil {
		return models.SendVoiceMessageResult{}, fmt.Errorf("unable to get AI text and audio reply: %s", err.Error())
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

// The primary spoken language of the audio data to be transcribed, supplied as a two-letter ISO-639-1 language code such
// as 'en' or 'fr'.
func getTranscriptLang(chatType models.ChatType) *string {
	if chatType == models.JobInterviewSeparateQuestionsChatType {
		return to.Ptr("de")
	}

	return nil
}

type llmInput struct {
	content   []llms.MessageContent
	maxTokens int
}

func (s *Service) getLlmInput(userTxt string, chatType models.ChatType,
) (llmInput, error) {
	if chatType == models.JobInterviewSeparateQuestionsChatType {
		content, err := s.getInterviewSeparateQuestionsContent(userTxt)
		if err != nil {
			return llmInput{}, fmt.Errorf("unable to get interview separate questions content: %s", err.Error())
		}

		return llmInput{content: content, maxTokens: 200}, nil
	}

	return llmInput{
		content: []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, userTxt)},
	}, nil
}

func (s *Service) getInterviewSeparateQuestionsContent(userTxt string) ([]llms.MessageContent, error) {
	return []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem,
			"You are a German language teacher. You are preparing for a job interview in German. "+
				"Correct the following text in case of mistakes. "+
				"If the text contains only minor mistakes, just say that it's okay, "+
				"and we can move on to the next interview question."),
		llms.TextParts(llms.ChatMessageTypeHuman, userTxt),
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
	_, err := s.getChatIfUserAutorized(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get chat: %w", err)
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
	UserID           string          `firestore:"user_id"`
	Timestamp        int64           `firestore:"time"`
	Title            string          `firestore:"title"`
	PreparedMessages []string        `firestore:"prep_msgs"`
	Type             models.ChatType `firestore:"type"`
	CurrQuestionIDx  int32           `firestore:"curr_m"`
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
			ChatID:            doc.Ref.ID,
			Title:             cutChatTitle(chatModel.Title),
			Time:              chatModel.Timestamp,
			CurrentMessageIDx: chatModel.CurrQuestionIDx,
			Typ:               swagger.ChatType(chatModel.Type),
		})
	}

	return chats, nil
}

func (s *Service) CreateChat(
	ctx context.Context, userID string, chatType models.ChatType, timestamp int64,
) (swagger.Chat, error) {
	if chatType != models.JobInterviewSeparateQuestionsChatType {
		return swagger.Chat{}, fmt.Errorf("can't create chat for type: %d", chatType)
	}

	createdChat, firstQuestion, err := s.createSeparateJobQuestionsChat(ctx, userID, timestamp)
	if err != nil {
		return swagger.Chat{}, fmt.Errorf("unable to create chat: %s", err.Error())
	}

	_, err = s.saveMessageToDB(
		ctx, firstQuestion.Text, "", createdChat.ChatID, firstQuestion.Audio, chatType, timestamp)

	if err != nil {
		return swagger.Chat{}, fmt.Errorf("unable to save first message to db: %s", err.Error())
	}

	return createdChat, nil
}

func (s *Service) createSeparateJobQuestionsChat(
	ctx context.Context, userID string, timestamp int64,
) (swagger.Chat, models.Variation, error) {
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

			return swagger.Chat{}, models.Variation{},
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
		return swagger.Chat{}, models.Variation{},
			fmt.Errorf("unable to get first message from firestore: %s", err.Error())
	}

	var firstMsg models.PreparedMessage
	if err = doc.DataTo(&firstMsg); err != nil {
		return swagger.Chat{}, models.Variation{},
			fmt.Errorf("unable to get first message data: %s", err.Error())
	}

	// pick random variation of the question
	variation := firstMsg.Variations[rand.Intn(len(firstMsg.Variations))]

	title := cutChatTitle(variation.Text)

	createdChat := swagger.Chat{
		PreparedMessages: messageIDs,
		Time:             timestamp,
		Title:            title,
		Typ:              swagger.ChatType(models.JobInterviewSeparateQuestionsChatType),
	}

	newChat, _, err := s.firestoreClient.
		Collection("chats").
		Add(ctx, chat{
			UserID:           userID,
			Type:             models.JobInterviewSeparateQuestionsChatType,
			Timestamp:        time.Now().UnixMilli(),
			Title:            title,
			PreparedMessages: messageIDs,
		})

	if err != nil {
		return swagger.Chat{}, models.Variation{},
			fmt.Errorf("unable to create chat: %s", err.Error())
	}

	createdChat.ChatID = newChat.ID

	return createdChat, variation, nil
}

func (s *Service) generateTextContent(
	ctx context.Context, in llmInput, model string) (string, error) {
	options := []llms.CallOption{llms.WithModel(model)}

	if in.maxTokens != 0 {
		options = append(options, llms.WithMaxTokens(in.maxTokens))
	}

	resp, err := s.llm.GenerateContent(ctx, in.content, options...)

	if err != nil {
		return "", fmt.Errorf("unable to get content from llm: %s", err.Error())
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in response")
	}

	return resp.Choices[0].Content, nil
}

func (s *Service) generateTextAndAudioContent(
	ctx context.Context, llmIn llmInput, userID, model string,
) (string, string, error) {
	text, err := s.generateTextContent(ctx, llmIn, model)
	if err != nil {
		return "", "", fmt.Errorf("unable to get text from llm: %s", err.Error())
	}

	audioURL, err := s.generateAudioContent(ctx, userID, text)
	if err != nil {
		return "", "", fmt.Errorf("unable to get audio from text: %s", err.Error())
	}

	return text, audioURL, nil
}

func (s *Service) generateAudioContent(ctx context.Context, userID, text string) (string, error) {
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
		return "", fmt.Errorf("unable to read voice response: %s", err.Error())
	}

	// for llm audio we use mp3
	audioName := fmt.Sprintf(mp3AudioPathTemplate, userID, time.Now().UnixNano())

	audioURL, err := s.uploadFileToStorage(ctx, audio, audioName)
	if err != nil {
		return "", fmt.Errorf("unable to upload voice message to storage: %s", err.Error())
	}

	return audioURL, nil
}

func (s *Service) GoToMessage(
	ctx context.Context, userID, chatID string, messageIDx int32,
) (swagger.ChatMessage, error) {
	userChat, err := s.getChatIfUserAutorized(ctx, chatID, userID)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get chat: %w", err)
	}

	if int(messageIDx) >= len(userChat.PreparedMessages) {
		text := "Herzlichen Glückwunsch, Sie haben alle Interviewfragen beantwortet!"
		timestamp := time.Now().UnixMilli()
		audioUrl, err := s.generateAudioContent(ctx, userID, text)
		if err != nil {
			return swagger.ChatMessage{}, fmt.Errorf("unable to get audio from text: %s", err.Error())
		}

		_, err = s.saveMessageToDB(ctx, text, "", chatID, audioUrl, userChat.Type, timestamp)
		if err != nil {
			return swagger.ChatMessage{}, fmt.Errorf("unable to save message to db: %s", err.Error())
		}

		return swagger.ChatMessage{
			AudioURL:          audioUrl,
			IsFromCurrentUser: false,
			Text:              text,
			Timestamp:         timestamp,
			UserID:            userID,
		}, nil
	}

	messageID := userChat.PreparedMessages[messageIDx]

	doc, err := s.firestoreClient.Collection("prepared_messages").Doc(messageID).Get(ctx)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get prepared message by id: %s", err.Error())
	}

	var message models.PreparedMessage
	if err = doc.DataTo(&message); err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get prepared message data: %s", err.Error())
	}

	timestamp := time.Now().UnixMilli()
	_, err = s.saveMessageToDB(
		ctx, message.GermanText, "", chatID, message.GermanAudio, userChat.Type, timestamp,
	)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to save message to db: %s", err.Error())
	}

	// update messageIDx in chat
	_, err = s.firestoreClient.Collection("chats").
		Doc(chatID).
		Update(ctx, []firestore.Update{
			{Path: "curr_m", Value: messageIDx},
		})

	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to update chat: %s", err.Error())
	}

	return swagger.ChatMessage{
		AudioURL:          message.GermanAudio,
		IsFromCurrentUser: false,
		Text:              message.GermanText,
		Timestamp:         timestamp,
		UserID:            userID,
	}, nil
}

// getUsersChatIfAutorized checks if the user is authorized to view the chat.
func (s *Service) getChatIfUserAutorized(
	ctx context.Context, chatID, userID string,
) (chat, error) {
	docRef := s.firestoreClient.Collection("chats").Doc(chatID)

	// Attempt to retrieve the document
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return chat{}, fmt.Errorf("failed to retrieve chat by id: %s", err.Error())
	}

	var chatModel chat

	if err = docSnapshot.DataTo(&chatModel); err != nil {
		return chat{}, fmt.Errorf("unable to get chat data: %s", err.Error())
	}

	if chatModel.UserID != userID {
		return chat{}, ErrUserNotAuthorizedToViewThisChat
	}

	return chatModel, nil
}

func (s *Service) DeleteChat(ctx context.Context, chatID, userID string) error {
	_, err := s.getChatIfUserAutorized(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("unable to get chat: %w", err)
	}

	// delete all messages in the chat
	iter := s.firestoreClient.Collection("messages").Where("chat_id", "==", chatID).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return fmt.Errorf("unable to get messages from firestore: %s", err.Error())
		}

		if _, err = doc.Ref.Delete(ctx); err != nil {
			return fmt.Errorf("unable to delete message: %s", err.Error())
		}
	}

	_, err = s.firestoreClient.Collection("chats").Doc(chatID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete chat: %s", err.Error())
	}

	return nil
}

func (s *Service) AnswerToMessages(ctx context.Context, chatID, userID string) (swagger.ChatMessage, error) {
	userChat, err := s.getChatIfUserAutorized(ctx, chatID, userID)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get chat: %w", err)
	}

	const limitLastMessages = 20

	query := s.firestoreClient.Collection("messages").
		Where("chat_id", "==", chatID).
		OrderBy("time", firestore.Desc).
		Limit(limitLastMessages)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var texts []string

	for {
		var doc *firestore.DocumentSnapshot
		doc, err = iter.Next()

		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return swagger.ChatMessage{}, fmt.Errorf("unable to get messages from firestore: %s", err.Error())
		}

		var message ChatMessage
		if err = doc.DataTo(&message); err != nil {
			return swagger.ChatMessage{}, fmt.Errorf("unable to get message data: %s", err.Error())
		}

		if message.UserID == "" {
			iter.Stop()
			break
		}

		texts = append(texts, message.Text)
	}

	if len(texts) == 0 {
		return swagger.ChatMessage{}, errors.New("no messages found")
	}

	slices.Reverse(texts)

	in, err := s.getLlmInput(concatTexts(texts), userChat.Type)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get llmInput: %s", err.Error())
	}

	llmReplyText, llmReplyAudioURL, err := s.generateTextAndAudioContent(ctx, in, userID, llmModel)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to get AI text and audio reply: %s", err.Error())
	}

	timestamp := time.Now().UnixMilli()

	_, err = s.saveMessageToDB(ctx, llmReplyText, "", chatID, llmReplyAudioURL, userChat.Type, timestamp)
	if err != nil {
		return swagger.ChatMessage{}, fmt.Errorf("unable to save llm reply to db: %s", err.Error())
	}

	return swagger.ChatMessage{
		AudioURL:          llmReplyAudioURL,
		IsFromCurrentUser: false,
		Text:              llmReplyText,
		Timestamp:         timestamp,
	}, nil
}

func concatTexts(arr []string) string {
	var b []byte
	for _, s := range arr {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}

		if len(b) > 0 {
			b = append(b, ' ')
		}

		last := s[len(s)-1]
		if unicode.IsLetter(rune(last)) {
			s += "."
		}

		b = append(b, s...)
	}

	return string(b)
}
