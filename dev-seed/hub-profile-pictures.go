package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"net/textproto"

	"github.com/fatih/color"
)

func uploadProfilePicture(avatarPath string, authToken string) error {
	// Read the avatar file
	imageData, err := os.ReadFile(avatarPath)
	if err != nil {
		return fmt.Errorf("failed to read avatar file: %v", err)
	}

	// Create a buffer to store the multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Determine content type based on file extension
	contentType := "image/jpeg" // default
	if strings.HasSuffix(avatarPath, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(avatarPath, ".webp") {
		contentType = "image/webp"
	}

	// Create a form file field with name "image" and content type as expected by the API
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", contentType)
	h.Set(
		"Content-Disposition",
		fmt.Sprintf(
			`form-data; name="image"; filename="%s"`,
			filepath.Base(avatarPath),
		),
	)
	part, err := writer.CreatePart(h)
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	// Write the image data to the form file field
	if _, err := part.Write(imageData); err != nil {
		return fmt.Errorf("failed to write image data: %v", err)
	}

	// Close the multipart writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(
		"POST",
		serverURL+"/hub/upload-profile-picture",
		&body,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type header and auth token
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	return nil
}

func uploadHubUserProfilePictures() {
	for _, user := range hubUsers {
		// 90% chance of having a profile picture
		if rand.Float32() < 0.9 {
			// Get the auth token from the session map
			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			avatarPath := fmt.Sprintf("avatar%d.jpg", rand.Intn(18)+1)
			if err := uploadProfilePicture(avatarPath, authToken); err != nil {
				log.Fatalf("upload profile picture fail %s: %v", user.Name, err)
				return
			}
			color.Magenta("added profile picture for %s", user.Name)
		}
	}
}
