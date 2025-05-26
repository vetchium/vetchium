package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/typespec/common"
)

type Middleware struct {
	db  db.DB
	log util.Logger
}

func NewMiddleware(db db.DB, log util.Logger) *Middleware {
	return &Middleware{db: db, log: log}
}

// Protect provides Authentication and Authorization on the /employer/* routes.
// For Hub related endpoints use the Guard function in middleware/hub.go
// Protect middleware only checks with the roles of the OrgUser. Whether the
// OrgUser belongs to the Org or not, should be verified via EmployerAuth. We
// should not take the OrgID/EmployerID on any of the request bodies but should
// be derived from the OrgUserToken.
func (m *Middleware) Protect(
	route string,
	handlerFunc http.HandlerFunc,
	allowedRoles []common.OrgUserRole,
) {
	http.Handle(
		route,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.log.Dbg("Entered Protect middleware")

			// Authentication part (inlined from employerAuth)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				m.log.Dbg("No auth header")
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			authHeader = strings.TrimPrefix(authHeader, "Bearer ")

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
			r = r.WithContext(ctx)

			// Authorization part
			if len(allowedRoles) == 0 {
				m.log.Err("No allowed roles for endpoint", "endpoint", route)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			if hasRoles(allowedRoles, []common.OrgUserRole{common.AnyOrgUser}) {
				handlerFunc(w, r)
				return
			}

			if !hasRoles(orgUser.OrgUserRoles, allowedRoles) {
				m.log.Inf(
					"User does not have required roles",
					"userRoles", orgUser.OrgUserRoles,
					"allowedRoles", allowedRoles,
				)
				http.Error(w, "", common.ErrEmployerRBAC)
				return
			}

			// Call the actual handler if roles are sufficient
			handlerFunc(w, r)
		}),
	)
}

func hasRoles(
	orgUserRoles []common.OrgUserRole,
	allowedRoles []common.OrgUserRole,
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
