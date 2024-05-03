package http

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ImageHandler struct {
	MongoDb   *mongo.Client
	MongoOpt  *options.BucketOptions
	ImageRepo *repo.ImageRepoImp
}

func (h ImageHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	bucket, err := gridfs.NewBucket(
		h.MongoDb.Database("image-server"), h.MongoOpt,
	)
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	filename := time.Now().Format(time.RFC3339) + "_" + header.Filename
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(buf.Bytes())
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fileId, err := json.Marshal(uploadStream.FileID)
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"fileId": strings.Trim(string(fileId), `"`), "fileSize": fileSize})

}

func (h ImageHandler) DownloadImage(c *gin.Context) {
	imageId := c.Query("imageId")

	objID, err := primitive.ObjectIDFromHex(imageId)
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	bucket, _ := gridfs.NewBucket(
		h.MongoDb.Database("image-server"), h.MongoOpt,
	)

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(objID, &buf)
	if err != nil {
		log.GetLog().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	contentType := http.DetectContentType(buf.Bytes())

	c.Writer.Header().Add("Content-Type", contentType)
	c.Writer.Header().Add("Content-Length", strconv.Itoa(len(buf.Bytes())))

	c.Writer.Write(buf.Bytes())
}

type RequestSubmitImage struct {
	CafeID  int32  `json:"cafe_id"`
	ImageID string `json:"image"`
}

func (h ImageHandler) SubmitImage(context *gin.Context) {
	var request RequestSubmitImage
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image := &models.Image{
		ID:        request.ImageID,
		Reference: request.CafeID,
	}

	err := h.ImageRepo.Create(context, image)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Image submitted successfully"})
}
