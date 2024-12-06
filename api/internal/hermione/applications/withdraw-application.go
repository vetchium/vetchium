package applications

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func WithdrawApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered WithdrawApplication")

	}
}
