package hubusers

import (
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
)

func CheckHandleAvailability(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}
