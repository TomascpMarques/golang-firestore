package chat

import (
	"context"
	"log"
	"net/http"
	"time"

	firestore "cloud.google.com/go/firestore"
	messages "github.com/Tomascpmarques/golang-firestore/models/messages"
	rooms "github.com/Tomascpmarques/golang-firestore/models/rooms"
	gin "github.com/gin-gonic/gin"
	iterator "google.golang.org/api/iterator"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type targetRoom struct {
	Name     string `form:"name" json:"name" query:"name"`
	Quantity int    `form:"num" json:"quantity" query:"num"`
}

func RetrieveMessagesFromRoom(c *gin.Context, client *firestore.Client) {
	var targetRoom targetRoom

	if c.ShouldBind(&targetRoom) != nil {
		log.Println("Falha de parsing do pedido")
		c.JSON(http.StatusBadRequest, gin.H{
			"estado": "",
			"razao":  "pedido invalido/mal definido",
		})
		return
	}

	if targetRoom.Quantity > 35 {
		targetRoom.Quantity = 35
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*35)
	defer cancel()

	docs := client.Collection("messages").Doc(targetRoom.Name).Collection("messages").Limit(targetRoom.Quantity).Documents(ctx)

	messages := make([]messages.StoredMessage, targetRoom.Quantity)

	i := 0
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Printf("Erro ao obter mensagem: %s\n", err.Error())
			continue
		}
		if err := doc.DataTo(&messages[i]); err != nil {
			log.Printf("Erro lidar com parsing da mensagem: %s\n", err.Error())
			continue
		}
		i += 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sucesso",
		"valores": messages[0:i], // Truncate the messages, so no empty messages are included
		"query":   targetRoom,
	})
}

func PublishMessageToRoom(c *gin.Context, client *firestore.Client) {
	var newMessage messages.NewMessage

	if c.ShouldBind(&newMessage) != nil {
		log.Println("Falha ao criar uma nova sala")
		c.JSON(http.StatusBadRequest, gin.H{
			"razao": "pedido invalido",
		})
		return
	}

	// Verificar se a sala que vai receber a mensagem já existe
	_, statusForTopic := client.
		Collection("rooms").
		Doc("topic").
		Collection("definitions").
		Doc(newMessage.Room).
		Get(context.Background())

	_, statusForNoTopic := client.
		Collection("rooms").
		Doc("no_topic").
		Collection("definitions").
		Doc(newMessage.Room).
		Get(context.Background())

	topicRoomCheckExists := status.Code(statusForTopic) == codes.NotFound
	noTopicRoomCheckExists := status.Code(statusForNoTopic) == codes.NotFound

	if topicRoomCheckExists && noTopicRoomCheckExists {
		log.Println("Tentativa de criar mensagem para sala que não existe")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "Tentativa de enviar mensagem para uma sala inexistente",
		})
		return
	}
	// --------------------------------------------------------

	message := messages.StoredMessage{
		Message: messages.Message{
			Sender:     newMessage.Sender,
			Content:    newMessage.Content,
			SenderType: newMessage.SenderType,
		},
		Sent: time.Now(),
	}

	if newMessage.SenderType == "" {
		message.SenderType = "user"
	}

	// Add room to respective topic collection
	ctx := context.Background()

	_, _, err := client.Collection("messages").Doc(newMessage.Room).Collection("messages").Add(ctx, message)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"razao": "erro interno",
		})
		return
	}
}

func CreateRoom(c *gin.Context, client *firestore.Client) {
	var newRoom rooms.Room

	if c.ShouldBind(&newRoom) != nil {
		log.Println("Falha ao criar uma nova sala")
		c.JSON(http.StatusBadRequest, gin.H{
			"razao": "pedido invalido",
		})
		return
	}

	room := rooms.RoomDefined{
		Room: rooms.Room{
			Name:        newRoom.Name,
			Owner:       newRoom.Owner,
			Description: newRoom.Description,
			Category:    newRoom.Category,
		},
		MsgCount: 0, // If undefined it will be 0 by default
	}

	topicVariant := "no_topic"
	if newRoom.Category != "" {
		topicVariant = "topic"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	_, err := client.Collection("rooms").Doc(topicVariant).Collection("definitions").Doc(room.Name).Create(ctx, room)

	if err != nil {
		log.Printf("Erro ao criar sala: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "not created",
			"reason": "room already exists",
		})
		return
	}
	// -------------------------------------

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
