package granger

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
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
	var filepath, filename string

	// Try up to 3 times to create a file with a unique name
	for i := 0; i < 3; i++ {
		filename = util.RandomUniqueID(vetchi.ResumeIDLenBytes)
		filepath = fmt.Sprintf("/resumes/%s", filename)

		// Atomically create and open the file
		var err error
		file, err = os.OpenFile(
			filepath,
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
		os.Remove(filepath)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if closeErr != nil {
		g.log.Err("failed to close resume file", "error", closeErr)
		os.Remove(filepath)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	g.log.Dbg("saved resume file", "filepath", filepath)

	// filename should be returned to hermione.
	// filepath has the volume mount as per the granger deployment spec,
	// while hermione may mount in a different path
	w.Write([]byte(filename))
}
