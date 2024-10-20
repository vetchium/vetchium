package hermione

import (
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
)

func (h *Hermione) EmployerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		orgUser, err := h.db.AuthOrgUser(r.Context(), authHeader)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.log.Info("Authenticated org user", "orgUser", orgUser)

		next.ServeHTTP(w, r)
	})
}
