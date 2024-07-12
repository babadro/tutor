package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/babadro/tutor/internal/models"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

const (
	mp3AudioPathTemplate = "prepared_messages/%d.mp3"
)

type envVars struct {
	OpenaiAPIKey  string `env:"OPENAI_API_KEY,required"`
	StorageBucket string `env:"STORAGE_BUCKET,required"`
}

type clients struct {
	baranovOpenai   *openai.Client
	storageClient   *storage.Client
	firestoreClient *firestore.Client
}

func main() {
	if err := godotenv.Load(".env.secrets", ".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// set env
	var envs envVars
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	ctx := context.Background()

	cl, err := initClients(ctx, envs.StorageBucket, envs.OpenaiAPIKey)
	if err != nil {
		log.Fatalf("unable to init clients: %s", err.Error())
	}

	for _, txt := range interviewQuestions {
		audio, err := getAudio(ctx, txt, cl.baranovOpenai)
		if err != nil {
			log.Fatalf("unable to get audio: %s", err.Error())
		}

		audioName := fmt.Sprintf(mp3AudioPathTemplate, time.Now().UnixNano())
		audioURL, err := uploadFileToStorage(ctx, audio, audioName, cl.storageClient)
		if err != nil {
			log.Fatalf("unable to upload audio to storage: %s", err.Error())
		}

		message := models.PreparedMessage{
			Type:        models.JobInterviewQuestion,
			GermanText:  txt,
			GermanAudio: audioURL,
		}

		err = saveDocToFirestore(ctx, cl.firestoreClient, message)
		if err != nil {
			log.Fatalf("unable to save doc to firestore: %s", err.Error())
		}
	}
}

func initClients(ctx context.Context, storageBucket, openaiAPIKey string) (clients, error) {
	// Initialize Firebase SDK
	firebaseConfig := &firebase.Config{
		StorageBucket: storageBucket,
	}
	opt := option.WithCredentialsFile("secrets/tutor.json")
	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		return clients{}, fmt.Errorf("unable to init firebase app: %s", err.Error())
	}

	baranovClient := openai.NewClient(openaiAPIKey)
	storageClient, err := firebaseApp.Storage(ctx)
	if err != nil {
		return clients{}, fmt.Errorf("unable to init storage client: %s", err.Error())
	}

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return clients{}, fmt.Errorf("unable to init firestore client: %s", err.Error())
	}

	return clients{
		baranovOpenai:   baranovClient,
		storageClient:   storageClient,
		firestoreClient: firestoreClient,
	}, nil
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
