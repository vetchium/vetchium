package employersettings

import (
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func GetCoolOffPeriod(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCoolOffPeriod")
		period, err := h.DB().GetCoolOffPeriod(r.Context())
		if err != nil {
			h.Dbg("failed to get cool off period", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", period)))
	}
}
