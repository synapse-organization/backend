package utils

import (
	"barista/pkg/log"
	"barista/pkg/models"
	bb "bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func TestUploadImage(cafes []models.Cafe) []string {
	filePath := "./pkg/utils/coffee.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.GetLog().Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.GetLog().Errorf("Unable to read file: %v", err)
	}

	var fileIDs []string
	cafesLen := len(cafes)

	myIP, _ := getBaseURL()

	for i := 0; i < cafesLen+40; i++ {
		body := &bb.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", filepath.Base(filePath))
		if err != nil {
			log.GetLog().Errorf("Unable to create form file: %v", err)
			continue
		}

		_, err = io.Copy(part, bb.NewReader(fileContent))
		if err != nil {
			log.GetLog().Errorf("Unable to copy file content: %v", err)
			continue
		}

		err = writer.Close()
		if err != nil {
			log.GetLog().Errorf("Unable to close writer: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", myIP+"/api/image/upload", body)
		if err != nil {
			log.GetLog().Errorf("Unable to create request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.GetLog().Errorf("Unable to send request: %v", err)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.GetLog().Errorf("Unable to read response body: %v", err)
			continue
		}

		var response struct {
			FileID   string `json:"fileId"`
			FileSize int    `json:"fileSize"`
		}

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			log.GetLog().Errorf("Unable to unmarshal response: %v", err)
			continue
		}

		fileIDs = append(fileIDs, strings.Trim(response.FileID, `"`))
	}

	return fileIDs
}
