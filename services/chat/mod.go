package chat

import (
	"context"
	"log"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"
)

type newRoom struct {
	Name        string `form:"name" json:"name"`
	Owner       string `form:"owner" json:"owner"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category"`
}

type room struct {
	Name        string `form:"name" json:"name"`
	Owner       string `form:"owner" json:"owner"`
	Description string `form:"description" json:"description"`
	Category    string `form:"category" json:"category"`
	MsgCount    int    `form:"msg_count" json:"msg_count"`
}

func CreateRoom(c *gin.Context, client *firestore.Client) {
	var newRoom newRoom

	if c.ShouldBind(&newRoom) != nil {
		log.Println("Falha ao criar uma nova sala")
		c.JSON(400, gin.H{
			"razao": "pedido invalido",
		})
		return
	}

	room := room{
		Name:        newRoom.Name,
		Owner:       newRoom.Owner,
		Description: newRoom.Description,
		Category:    newRoom.Category,
	}

	var roomDocForTopic *firestore.DocumentRef

	if newRoom.Category == "" {
		roomDocForTopic = client.Collection("rooms").Doc("no_topic")
	} else if newRoom.Category != "" {
		roomDocForTopic = client.Collection("rooms").Doc("topic")
	}

	// client.Collection("users").Doc("us-na").Collection("admins").Add(ctx, exampleUser)
	// Add room to respective topic collection
	ctx := context.Background()
	_, _, err := roomDocForTopic.Collection(newRoom.Category).Add(ctx, room)
	if err != nil {
		log.Println(err.Error())
	}

	c.JSON(200, gin.H{
		"hello": "world",
	})
}
