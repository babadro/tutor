package main

import (
	"context"
	"errors"
	"log"

	"github.com/babadro/tutor/cmd/common"
	"google.golang.org/api/iterator"
)

func main() {
	cl := common.InitClients()

	// remove all documents from the collection
	ctx := context.Background()
	iter := cl.FirestoreClient.Collection("prepared_messages").Documents(ctx)

	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			log.Fatalf("unable to get messages from firestore: %s", err.Error())
		}

		if _, err = doc.Ref.Delete(ctx); err != nil {
			log.Fatalf("unable to delete message: %s", err.Error())
		}
	}
}
