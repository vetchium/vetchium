package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

const (
	OrgUserIDHeader    = "X-Vetchi-OrgUserID"
	OrgUserRolesHeader = "X-Vetchi-OrgUserRoles"
	EmployerIDHeader   = "X-Vetchi-EmployerID"
)

type Middleware struct {
	db  db.DB
	log *slog.Logger
}

func NewMiddleware(db db.DB, log *slog.Logger) *Middleware {
	return &Middleware{db: db, log: log}
}

func (m *Middleware) employerAuth(next http.Handler) http.Handler {
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
		r.Header.Set(EmployerIDHeader, orgUser.EmployerID.String())

		roles := strings.Builder{}
		for _, role := range orgUser.OrgUserRoles {
			roles.WriteString(string(role))
			roles.WriteString(";")
		}
		r.Header.Set(OrgUserRolesHeader, roles.String())

		next.ServeHTTP(w, r)
	})
}

// Protect middleware only checks with the roles of the OrgUser. Whether the
// OrgUser belongs to the Org or not, should be verified via EmployerAuth. We
// should not take the OrgID/EmployerID on any of the request bodies but should
// be derived from the OrgUserToken.
func (m *Middleware) Protect(
	route string,
	handler http.Handler,
	allowedRoles []string,
) {
	http.Handle(route, m.employerAuth(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve roles from the header set by EmployerAuth
			userRolesHeader := r.Header.Get(OrgUserRolesHeader)
			if userRolesHeader == "" {
				m.log.Error("No roles found in header")
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			if strings.Contains(userRolesHeader, string(vetchi.Admin)) {
				// Admin can do anything within Org
				handler.ServeHTTP(w, r)
				return
			}

			// Split the roles into a slice
			userRoles := strings.Split(userRolesHeader, ";")

			// Check if the user has any of the allowed roles
			if !IsAllowed(allowedRoles, userRoles) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the actual handler if roles are sufficient
			handler.ServeHTTP(w, r)
		}),
	))
}

func IsAllowed(allowedRoles, userRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		for _, userRole := range userRoles {
			if allowedRole == userRole {
				return true
			}
		}
	}
	return false
}
