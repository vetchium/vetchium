package employerauth

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func GetOnboardStatus(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getOnboardStatusReq employer.GetOnboardStatusRequest
		err := json.NewDecoder(r.Body).Decode(&getOnboardStatusReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getOnboardStatusReq) {
			return
		}

		var status employer.OnboardStatus

		employerObj, err := h.DB().GetEmployer(
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

		switch employerObj.EmployerState {
		case db.OnboardPendingEmployerState:
			status = employer.DomainVerifiedOnboardPending
		case db.OnboardedEmployerState:
			status = employer.DomainOnboarded
		default:
			h.Err(
				"unknown employer state",
				"client_id",
				getOnboardStatusReq.ClientID,
				"state",
				employerObj.EmployerState,
			)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		resp := employer.GetOnboardStatusResponse{Status: status}
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
		resp := employer.GetOnboardStatusResponse{
			Status: employer.DomainNotVerified,
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
		resp := employer.GetOnboardStatusResponse{
			Status: employer.DomainNotVerified,
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

	resp := employer.GetOnboardStatusResponse{
		Status: employer.DomainVerifiedOnboardPending,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.Err("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
