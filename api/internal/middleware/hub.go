package middleware

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

// Guard provides Authentication and Authorization on the /hub/* routes.
// It checks the tier of the hub user and whether the user is allowed to access
// the route. If the user is not allowed to access the route, it returns a 403
// Forbidden error.
func (m *Middleware) Guard(
	route string,
	handlerFunc http.HandlerFunc,
	allowedTiers []hub.HubUserTier,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Dbg("Entered hubAuth middleware")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Dbg("No auth header")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")

		hubUser, err := m.db.AuthHubUser(r.Context(), authHeader)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				m.log.Dbg("No hub user")
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			m.log.Err("Failed to auth hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), HubUserCtxKey, hubUser)

		if len(allowedTiers) > 0 &&
			!slices.Contains(allowedTiers, hubUser.Tier) {
			m.log.Dbg("User is not allowed to access this route")
			http.Error(w, "", http.StatusForbidden)
			return
		}

		handlerFunc.ServeHTTP(w, r.WithContext(ctx))
	})
}
