package services

import (
	firestore "cloud.google.com/go/firestore"
	chat "github.com/Tomascpmarques/golang-firestore/services/chat"
	gin "github.com/gin-gonic/gin"
)

func SetupApiRoutes(router *gin.Engine, client *firestore.Client) {
	group := router.Group("/api")
	group.POST("/rooms/create", func(ctx *gin.Context) { chat.CreateRoom(ctx, client) })
}
