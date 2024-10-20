package hermione

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
)

const (
	OrgUserIDHeader   = "X-Vetchi-OrgUserID"
	OrgUserRoleHeader = "X-Vetchi-OrgUserRole"
)

type Middleware struct {
	db  db.DB
	log slog.Logger
}

func NewMiddleware(db db.DB, log slog.Logger) *Middleware {
	return &Middleware{db: db, log: log}
}

func (m *Middleware) EmployerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		orgUser, err := m.db.AuthOrgUser(r.Context(), authHeader)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.log.Debug("Authenticated org user", "orgUser", orgUser)

		r.Header.Set(OrgUserIDHeader, orgUser.ID.String())
		r.Header.Set(OrgUserRoleHeader, string(orgUser.OrgUserRole))

		next.ServeHTTP(w, r)
	})
}
