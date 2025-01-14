package granger

import (
	"io"
	"net/http"
	"os"

	"github.com/psankar/vetchi/api/internal/util"
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

	// Generate a unique ID based on content hash
	resumeSHA := util.GenerateResumeID(body)
	storageDir := util.GetResumeStorageDir("/resumes", resumeSHA)
	filepath := util.GetResumeStoragePath("/resumes", resumeSHA)

	// Create the directory structure if it doesn't exist
	err = os.MkdirAll(storageDir, 0755)
	if err != nil {
		g.log.Err("failed to create directory structure", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Try to create the file - if it already exists, it's likely the same content
	file, err := os.OpenFile(
		filepath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		if os.IsExist(err) {
			// File already exists - since we use content hash, this means
			// the same resume was uploaded before
			g.log.Dbg("resume already exists, reusing", "filepath", filepath)
			w.Write([]byte(resumeSHA))
			return
		}
		g.log.Err("failed to create resume file", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Write the content and close the file
	_, err = file.Write(body)
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
	w.Write([]byte(resumeSHA))
}
