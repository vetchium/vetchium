package hubusers

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func SetHandle(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}
