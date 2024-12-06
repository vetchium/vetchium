package hubauth

import (
	"net/http"
	"strings"

	"github.com/psankar/vetchi/api/internal/wand"
)

func Logout(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered hub logout")

		token := r.Header.Get("Authorization")
		if token == "" {
			h.Dbg("no token found")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		err := h.DB().Logout(r.Context(), token)
		if err != nil {
			h.Dbg("failed to logout", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
