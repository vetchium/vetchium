package hermione

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

func (h *Hermione) getOnboardStatus(w http.ResponseWriter, r *http.Request) {
	var req libvetchi.GetOnboardStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var status libvetchi.OnboardStatus

	employer, err := h.db.GetEmployer(r.Context(), req.ClientID)
	if err != nil {
		if errors.Is(err, db.ErrNoEmployer) {
			// Unregistered domain. Check if vetchiadmin TXT record is present
			h.newDomainProcess(r.Context(), w, req.ClientID)
			return
		} else {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	switch employer.EmployerState {
	case db.OnboardPendingEmployerState:
		status = libvetchi.DomainVerifiedOnboardPending
	case db.OnboardedEmployerState:
		status = libvetchi.DomainOnboarded
	case db.DeboardedEmployerState:
		status = libvetchi.DomainNotVerified
	default:
		h.log.Error(
			"unknown employer state",
			"client_id",
			req.ClientID,
			"state",
			employer.EmployerState,
		)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp := libvetchi.GetOnboardStatusResponse{Status: status}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Hermione) newDomainProcess(
	ctx context.Context,
	w http.ResponseWriter,
	domain string,
) {
	url := "vetchiadmin." + domain
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		h.log.Debug("lookup TXT records", "domain", domain, "error", err)
		resp := libvetchi.GetOnboardStatusResponse{
			Status: libvetchi.DomainNotVerified,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
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
		resp := libvetchi.GetOnboardStatusResponse{
			Status: libvetchi.DomainNotVerified,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			h.log.Error("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		return
	}

	err = h.db.InitEmployerAndDomain(ctx, db.Employer{
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

	resp := libvetchi.GetOnboardStatusResponse{
		Status: libvetchi.DomainVerifiedOnboardPending,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
