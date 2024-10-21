package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/storage"
	"github.com/babadro/tutor/cmd/common"
	"github.com/babadro/tutor/internal/models"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type interviewQuestion struct {
	BaseText   string   `json:"base_text"`
	Variations []string `json:"variations"`
}

func main() {
	cl := common.InitClients()
	ctx := context.Background()

	fileContent, err := os.ReadFile("questions_with_variations.json")
	if err != nil {
		log.Fatalf("unable to read file: %s", err.Error())
	}

	var interviewQuestions []interviewQuestion
	err = json.Unmarshal(fileContent, &interviewQuestions)
	if err != nil {
		log.Fatalf("unable to unmarshal file content: %s", err.Error())
	}

	for _, question := range interviewQuestions {
		message := models.PreparedMessage{
			Type:     models.JobInterviewQuestion,
			BaseText: question.BaseText,
		}

		for i, variation := range question.Variations {
			if i == 5 {
				break // so far 5 variations are enough
			}

			audio, err := getAudio(ctx, variation, cl.BaranovOpenai)
			if err != nil {
				log.Fatalf("unable to get audio: %s", err.Error())
			}

			audioName := fmt.Sprintf(common.Mp3AudioPathTemplate, time.Now().UnixNano())
			audioURL, err := uploadFileToStorage(ctx, audio, audioName, cl.StorageClient)
			if err != nil {
				log.Fatalf("unable to upload audio to storage: %s", err.Error())
			}

			message.Variations = append(message.Variations, models.Variation{
				Language: "de",
				Text:     variation,
				Audio:    audioURL,
				ID:       uuid.New().String(),
			})
		}

		err = saveDocToFirestore(ctx, cl.FirestoreClient, message)
		if err != nil {
			log.Fatalf("unable to save doc to firestore: %s", err.Error())
		}
	}
}

func saveDocToFirestore(ctx context.Context, cl *firestore.Client, doc any) error {
	_, _, err := cl.Collection("prepared_messages").Add(ctx, doc)
	if err != nil {
		return fmt.Errorf("unable to save doc to firestore: %s", err.Error())
	}

	return nil
}

func getAudio(ctx context.Context, text string, cl *openai.Client) ([]byte, error) {
	textToSpeechReq := openai.CreateSpeechRequest{
		Model:          openai.TTSModel1,
		Input:          text,
		Voice:          openai.VoiceNova,
		ResponseFormat: "mp3",
		Speed:          1,
	}

	llmAudio, err := cl.CreateSpeech(ctx, textToSpeechReq)
	if err != nil {
		return nil, fmt.Errorf("unable to get speech from text: %s", err.Error())
	}

	// convert llmAudio to []byte
	audio, err := io.ReadAll(llmAudio)
	if err != nil {
		return nil, fmt.Errorf("unable to read voice message: %s", err.Error())
	}

	return audio, nil
}

func uploadFileToStorage(
	ctx context.Context, fileBytes []byte, fileName string, stCl *storage.Client,
) (string, error) {
	bucket, err := stCl.DefaultBucket()
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
