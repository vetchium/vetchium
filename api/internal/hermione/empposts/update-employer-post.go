package empposts

import (
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
)

func UpdateEmployerPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateEmployerPost")
		http.Error(
			w,
			"History of changed posts support not implemented yet",
			http.StatusNotImplemented,
		)
	}
}
