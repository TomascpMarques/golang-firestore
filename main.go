package main

import (
	"context"

	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	services "github.com/Tomascpmarques/golang-firestore/services"
	gin "github.com/gin-gonic/gin"
)

const PROJECT_ID = "fir-go-4a841"

func GetClient(ctx context.Context) (client *firestore.Client) {
	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func main() {
	ctx := context.Background()

	setFirestoreEmulatorHost()

	client := GetClient(ctx)
	defer client.Close()

	r := gin.Default()

	services.SetupApiRoutes(r, client)

	r.Run(":8191")
}

func setFirestoreEmulatorHost() {
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

}
