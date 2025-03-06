package hubusers

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func OnboardHubUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered OnboardHubUser")

		var onboardHubUserRequest hub.OnboardHubUserRequest
		if err := json.NewDecoder(r.Body).Decode(&onboardHubUserRequest); err != nil {
			h.Err("Failed to decode onboardHubUserRequest", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &onboardHubUserRequest) {
			h.Dbg("validation failed", "request", onboardHubUserRequest)
			return
		}
		h.Dbg("validated", "onboardHubUserRequest", onboardHubUserRequest)

		onBoardHubUserReponse, err := h.DB().
			OnboardHubUser(r.Context(), onboardHubUserRequest)
		if err != nil {
			h.Err("failed to onboard hub user", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("onboarded", "onBoardHubUserReponse", onBoardHubUserReponse)

		err = json.NewEncoder(w).Encode(onBoardHubUserReponse)
		if err != nil {
			h.Err("failed to encode onboardHubUserReponse", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
