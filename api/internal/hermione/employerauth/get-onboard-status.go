package employerauth

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func GetOnboardStatus(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getOnboardStatusReq vetchi.GetOnboardStatusRequest
		err := json.NewDecoder(r.Body).Decode(&getOnboardStatusReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getOnboardStatusReq) {
			return
		}

		var status vetchi.OnboardStatus

		employer, err := h.DB().GetEmployer(
			r.Context(),
			getOnboardStatusReq.ClientID,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoEmployer) {
				// Unregistered domain. Check if vetchiadmin TXT record is present
				newDomainProcess(
					r.Context(),
					w,
					getOnboardStatusReq.ClientID,
					h,
				)
				return
			} else {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		switch employer.EmployerState {
		case db.OnboardPendingEmployerState:
			status = vetchi.DomainVerifiedOnboardPending
		case db.OnboardedEmployerState:
			status = vetchi.DomainOnboarded
		case db.DeboardedEmployerState:
			status = vetchi.DomainNotVerified
		default:
			h.Err(
				"unknown employer state",
				"client_id",
				getOnboardStatusReq.ClientID,
				"state",
				employer.EmployerState,
			)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		resp := vetchi.GetOnboardStatusResponse{Status: status}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

// TODO: This is badly written and can be refactored to use r.Context() better
func newDomainProcess(
	ctx context.Context,
	w http.ResponseWriter,
	domain string,
	h wand.Wand,
) {
	url := "vetchiadmin." + domain
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		h.Dbg("lookup TXT records", "domain", domain, "error", err)
		resp := vetchi.GetOnboardStatusResponse{
			Status: vetchi.DomainNotVerified,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	admin := ""
	if len(txtRecords) > 0 {
		admin = txtRecords[0]
	}

	if admin == "" {
		resp := vetchi.GetOnboardStatusResponse{
			Status: vetchi.DomainNotVerified,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	err = h.DB().InitEmployerAndDomain(ctx, db.Employer{
		ClientIDType:      db.DomainClientIDType,
		OnboardAdminEmail: admin,
		EmployerState:     db.OnboardPendingEmployerState,
	}, db.Domain{
		DomainName:  domain,
		DomainState: db.VerifiedDomainState,
	})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp := vetchi.GetOnboardStatusResponse{
		Status: vetchi.DomainVerifiedOnboardPending,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Err("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
