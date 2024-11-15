package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ChangeOpeningState(h wand.Wand) http.HandlerFunc {
	type Key struct {
		From vetchi.OpeningState
		To   vetchi.OpeningState
	}
	validTransitions := map[Key]bool{
		// Basically one cannot revert to draft state after publishing and
		// close state is final

		{From: vetchi.DraftOpening, To: vetchi.ActiveOpening}: true,
		{From: vetchi.DraftOpening, To: vetchi.ClosedOpening}: true,

		{From: vetchi.ActiveOpening, To: vetchi.SuspendedOpening}: true,
		{From: vetchi.ActiveOpening, To: vetchi.ClosedOpening}:    true,

		{From: vetchi.SuspendedOpening, To: vetchi.ActiveOpening}: true,
		{From: vetchi.SuspendedOpening, To: vetchi.ClosedOpening}: true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangeOpeningState")
		var changeOpeningStateReq vetchi.ChangeOpeningStateRequest
		err := json.NewDecoder(r.Body).Decode(&changeOpeningStateReq)
		if err != nil {
			h.Dbg("failed to decode change opening state request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &changeOpeningStateReq) {
			h.Dbg("invalid", "changeOpeningStateReq", changeOpeningStateReq)
			return
		}

		if !validTransitions[Key{
			From: changeOpeningStateReq.FromState,
			To:   changeOpeningStateReq.ToState,
		}] {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		err = h.DB().ChangeOpeningState(r.Context(), changeOpeningStateReq)
		if err != nil {
			h.Dbg("failed to change opening state", "error", err)

			if errors.Is(err, db.ErrInternal) {
				http.Error(w, "", http.StatusInternalServerError)
			} else if errors.Is(err, db.ErrStateMismatch) {
				http.Error(w, "", http.StatusConflict)
			} else if errors.Is(err, db.ErrNoOpening) {
				http.Error(w, "", http.StatusNotFound)
			} else {
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("changed state", "openingID", changeOpeningStateReq.OpeningID)
		w.WriteHeader(http.StatusOK)
	}
}
