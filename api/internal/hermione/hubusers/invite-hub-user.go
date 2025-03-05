package hubusers

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func InviteHubUser(w wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
