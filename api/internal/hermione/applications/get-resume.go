package applications

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func GetResume(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}
}
