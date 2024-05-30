package main

import (
	"context"
	"net"
	"time"

	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	services "github.com/Tomascpmarques/golang-firestore/services"
	cors "github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
)

const PROJECT_ID = "cloudcomputing-424022"

func GetClient(ctx context.Context) (client *firestore.Client) {
	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func main() {
	// godotenv.Load()
	ctx := context.Background()

	// setFirestoreEmulatorHost()

	client := GetClient(ctx)
	defer client.Close()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"https://localhost:8191/api/storage/"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		/* AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		}, */
		MaxAge: 12 * time.Hour,
	}))
	r.MaxMultipartMemory = 10 << 20 // 8 MiB

	services.SetupApiRoutes(r, client, PROJECT_ID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8191"
	}
	listener, err := net.Listen("tcp", "[::]:"+port)
	if err != nil {
		log.Fatalf("Erro ao obter endereço")
	}

	r.RunListener(listener)

}

func setFirestoreEmulatorHost() {
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

}
