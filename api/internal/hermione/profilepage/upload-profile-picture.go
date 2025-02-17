package profilepage

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func UploadProfilePicture(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("UploadProfilePicture")
	}
}
