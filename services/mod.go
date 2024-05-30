package services

import (
	"context"
	"log"

	firestore "cloud.google.com/go/firestore"
	storage "cloud.google.com/go/storage"
	chat "github.com/Tomascpmarques/golang-firestore/services/chat"
	storage_api "github.com/Tomascpmarques/golang-firestore/services/storage"
	gin "github.com/gin-gonic/gin"
)

func SetupApiRoutes(router *gin.Engine, client *firestore.Client, project_id string) {
	ctx := context.Background()
	bucket_client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("falha ao iniciar client para bucket storage: %s", err.Error())
	}

	group := router.Group("/api")

	group.POST("/rooms/create", func(ctx *gin.Context) { chat.CreateRoom(ctx, client) })
	group.POST("/rooms/post", func(ctx *gin.Context) { chat.PublishMessageToRoom(ctx, client) })
	group.GET("/rooms/messages", func(ctx *gin.Context) { chat.RetrieveMessagesFromRoom(ctx, client) })

	group.POST("/storage/create", func(ctx *gin.Context) { storage_api.NewBucket(ctx, project_id, bucket_client) })
	group.POST("/storage/upload/avatar", func(ctx *gin.Context) { storage_api.AddAvatar(ctx, project_id, bucket_client) })
}
