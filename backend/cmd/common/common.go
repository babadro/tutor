package common

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

const (
	Mp3AudioPathTemplate = "prepared_messages/%d.mp3"
)

type envVars struct {
	OpenaiAPIKey  string `env:"OPENAI_API_KEY,required"`
	StorageBucket string `env:"STORAGE_BUCKET,required"`
}

type Clients struct {
	BaranovOpenai   *openai.Client
	StorageClient   *storage.Client
	FirestoreClient *firestore.Client
}

func InitClients() Clients {
	if err := godotenv.Load(".env.secrets", ".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// set env
	var envs envVars
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	ctx := context.Background()

	// Initialize Firebase SDK
	firebaseConfig := &firebase.Config{
		StorageBucket: envs.StorageBucket,
	}
	opt := option.WithCredentialsFile("secrets/tutor.json")
	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig, opt)
	if err != nil {
		log.Fatalf("unable to init firebase app: %s", err.Error())
	}

	baranovClient := openai.NewClient(envs.OpenaiAPIKey)
	storageClient, err := firebaseApp.Storage(ctx)
	if err != nil {
		log.Fatalf("unable to init storage client: %s", err.Error())
	}

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to init firestore client: %s", err.Error())
	}

	return Clients{
		BaranovOpenai:   baranovClient,
		StorageClient:   storageClient,
		FirestoreClient: firestoreClient,
	}
}
