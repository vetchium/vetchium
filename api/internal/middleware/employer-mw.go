package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type userCtx int

const (
	OrgUserCtxKey userCtx = iota
	HubUserCtxKey
)

type Middleware struct {
	db  db.DB
	log util.Logger
}

func NewMiddleware(db db.DB, log util.Logger) *Middleware {
	return &Middleware{db: db, log: log}
}

func (m *Middleware) employerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Dbg("Entered employerAuth middleware")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Dbg("No auth header")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		orgUser, err := m.db.AuthOrgUser(r.Context(), authHeader)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				m.log.Dbg("No org user")
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			m.log.Err("Failed to auth org user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.log.Dbg("Authenticated org user", "orgUser", orgUser)

		ctx := context.WithValue(r.Context(), OrgUserCtxKey, orgUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Protect middleware only checks with the roles of the OrgUser. Whether the
// OrgUser belongs to the Org or not, should be verified via EmployerAuth. We
// should not take the OrgID/EmployerID on any of the request bodies but should
// be derived from the OrgUserToken.
func (m *Middleware) Protect(
	route string,
	handlerFunc http.HandlerFunc,
	allowedRoles []vetchi.OrgUserRole,
) {
	http.Handle(route, m.employerAuth(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.log.Dbg("Entered Protect middleware")
			ctx := r.Context()
			orgUser, ok := ctx.Value(OrgUserCtxKey).(db.OrgUserTO)
			if !ok {
				m.log.Err("Failed to get orgUser from context")
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			if !hasRoles(orgUser.OrgUserRoles, allowedRoles) {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			// Call the actual handler if roles are sufficient
			handlerFunc(w, r)
		}),
	))
}

func hasRoles(
	orgUserRoles []vetchi.OrgUserRole,
	allowedRoles []vetchi.OrgUserRole,
) bool {
	// TODO: Can potentially cache this ifF there is a performance issue
	for _, orgUserRole := range orgUserRoles {
		for _, role := range allowedRoles {
			if role == orgUserRole {
				return true
			}
		}
	}

	return false
}

func (m *Middleware) HubWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Dbg("Entered hubAuth middleware")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Dbg("No auth header")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		hubUser, err := m.db.AuthHubUser(r.Context(), authHeader)
		if err != nil {
			m.log.Err("Failed to auth hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), HubUserCtxKey, hubUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
