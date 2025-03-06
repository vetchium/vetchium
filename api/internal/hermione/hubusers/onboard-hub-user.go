package hubusers

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/hub"
	"golang.org/x/crypto/bcrypt"
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

		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(onboardHubUserRequest.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			h.Err("failed to generate password hash", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		generatedHandle, err := h.DB().
			OnboardHubUser(r.Context(), db.OnboardHubUserReq{
				InviteToken:         onboardHubUserRequest.Token,
				FullName:            onboardHubUserRequest.FullName,
				PasswordHash:        string(passwordHash),
				Tier:                onboardHubUserRequest.SelectedTier,
				ResidentCountryCode: onboardHubUserRequest.ResidentCountryCode,

				// TODO: Remove hard-coded language
				PreferredLanguage: vetchi.PreferredLanguage,

				ShortBio: onboardHubUserRequest.ShortBio,
				LongBio:  onboardHubUserRequest.LongBio,
			})
		if err != nil {
			h.Err("failed to onboard hub user", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		onBoardHubUserReponse := hub.OnboardHubUserResponse{
			GeneratedHandle: generatedHandle,
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
