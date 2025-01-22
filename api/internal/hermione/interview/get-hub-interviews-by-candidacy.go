package interview

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func GetHubInterviewsByCandidacy(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetHubInterviewsByCandidacy")
	}
}
