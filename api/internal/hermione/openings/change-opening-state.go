package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func ChangeOpeningState(h wand.Wand) http.HandlerFunc {
	type Key struct {
		From common.OpeningState
		To   common.OpeningState
	}
	validTransitions := map[Key]bool{
		// Basically one cannot revert to draft state after publishing and
		// close state is final

		{From: common.DraftOpening, To: common.ActiveOpening}: true,
		{From: common.DraftOpening, To: common.ClosedOpening}: true,

		{From: common.ActiveOpening, To: common.SuspendedOpening}: true,
		{From: common.ActiveOpening, To: common.ClosedOpening}:    true,

		{From: common.SuspendedOpening, To: common.ActiveOpening}: true,
		{From: common.SuspendedOpening, To: common.ClosedOpening}: true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangeOpeningState")
		var changeOpeningStateReq employer.ChangeOpeningStateRequest
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
