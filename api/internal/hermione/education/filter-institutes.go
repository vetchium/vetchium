package education

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func FilterInstitutes(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}
