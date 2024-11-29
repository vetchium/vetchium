package granger

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type UploadResumeRequest struct {
	Resume string `json:"resume"`
}

func (g *Granger) uploadResumeHandler(w http.ResponseWriter, r *http.Request) {
	g.log.Dbg("Entered uploadResumeHandler")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.log.Err("failed to read request body", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	decodedResume, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		g.log.Err("failed to decode base64 resume", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var file *os.File
	var filename string

	// Try up to 3 times to create a file with a unique name
	for i := 0; i < 3; i++ {
		uniqueID := uuid.New().String()
		filename = fmt.Sprintf("/resumes/%s", uniqueID)

		// Atomically create and open the file
		var err error
		file, err = os.OpenFile(
			filename,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL,
			0644,
		)
		if err != nil {
			if os.IsExist(err) {
				continue // File exists, try again with a new UUID
			}
			g.log.Err("failed to create resume file", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		break // Successfully created the file
	}

	if file == nil {
		g.log.Err("failed to create resume file after 3 attempts")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Write the content and close the file
	_, err = file.Write(decodedResume)
	closeErr := file.Close()
	if err != nil {
		g.log.Err("failed to write resume content", "error", err)
		os.Remove(filename)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if closeErr != nil {
		g.log.Err("failed to close resume file", "error", closeErr)
		os.Remove(filename)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// File saved successfully
	g.log.Dbg("saved resume file", "filename", filename)
	w.Write([]byte(filename))
}
