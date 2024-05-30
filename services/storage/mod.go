package storage

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func createStorageBucket(ctx context.Context, client *storage.Client, projectID, bucketName string) (err error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	storage_metadata := &storage.BucketAttrs{
		StorageClass: "STANDARD",
		Location:     "europe-west2",
	}

	bucket := client.Bucket(bucketName)
	if err = bucket.Create(timeout, projectID, storage_metadata); err != nil {
		log.Fatalf("Erro: %s", err.Error())
		return
	}

	return
}

type NewBucketQueryString struct {
	Name string `form:"name" query:"name"`
}

func NewBucket(ctx *gin.Context, project_id string, storage_client *storage.Client) {
	var newBucketData NewBucketQueryString

	if err := ctx.BindQuery(&newBucketData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "bad request data",
		})
		return
	}

	name := newBucketData.Name + project_id

	err := createStorageBucket(ctx, storage_client, project_id, name)
	if err != nil {
		log.Printf("Erro ao criar um novo bucket: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "erro",
			"reason": "internal error",
		})
		return
	}
}

type NewAvatarDataQueryString struct {
	Public bool   `form:"public"`
	Name   string `form:"name"`
}

func AddAvatar(ctx *gin.Context, project_id string, bucket_client *storage.Client) {
	var newAvatarMetadata NewAvatarDataQueryString

	if err := ctx.ShouldBindQuery(&newAvatarMetadata); err != nil {
		log.Printf("Erro: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "Request não contêm os dados necessários",
		})
		return
	}
	log.Printf("VALORES: %v", newAvatarMetadata)

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "Ficheiros recebidos não válidos",
		})
		return
	}

	log.Printf("Values: %v", file.Header)

	if file.Header.Get("Content-Type") != "image/jpeg" {
		log.Printf("Content-Type is: %s", file.Header.Get("Content-Type"))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "Conteúdo não é do tipo JPEG",
		})
		return
	}

	file_reader, err := file.Open()
	if err != nil {
		log.Fatalf("ERRO AO ABRIR O FICHEIRO")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "erro",
			"reason": "Erro interno",
		})
		return
	}

	c, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	object := bucket_client.Bucket("avatar-images-" + project_id).Object("users/avatar/" + newAvatarMetadata.Name)

	object = object.If(storage.Conditions{DoesNotExist: true})

	wc := object.NewWriter(c)
	wc.ChunkSize = 1

	if _, err = io.Copy(wc, file_reader); err != nil {
		log.Printf("Error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "falha ao copiar o conteúdo",
		})
		return
	}

	if err := wc.Close(); err != nil {
		log.Printf("Writer.Close: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"reason": "falha ao copiar os dados",
		})
		return
	}

	reader_role_crowd := storage.AllUsers
	if newAvatarMetadata.Public {
		if err := object.ACL().Set(c, reader_role_crowd, storage.RoleReader); err != nil {
			log.Fatalf("Erro: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "erro",
				"reason": "falha ao tornar público",
			})
			return
		}
	}
}
