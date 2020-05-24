package common

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/joho/godotenv"
)

/*
DatastoreClient exposed value
*/
var DatastoreClient *datastore.Client

/*
Ctx current contect
*/
var Ctx = context.Background()

func init() {
	var err error
	godotenvErr := godotenv.Load()
	if godotenvErr != nil {
		log.Fatal("Error loading .env file")
	}

	projectID := os.Getenv("GCLOUD_DATASET_ID")
	DatastoreClient, err = datastore.NewClient(Ctx, projectID)
	if err != nil {
		fmt.Printf("datastore connect error %v", err)
	}
}
