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
			status, err = h.newDomainProcess(r.Context(), req.ClientID)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		} else {
			h.logger.Error("failed to get employer", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	} else {
		status = libvetchi.OnboardStatus(employer.OnboardStatus)
	}

	resp := libvetchi.GetOnboardStatusResponse{Status: status}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *Hermione) newDomainProcess(
	ctx context.Context,
	domain string,
) (libvetchi.OnboardStatus, error) {
	url := "vetchiadmin." + domain
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		h.logger.Debug("lookup TXT records", "domain", domain, "error", err)
		return libvetchi.DomainNotVerified, nil
	}

	admin := ""

	if len(txtRecords) > 0 {
		admin = txtRecords[0]
	}

	if admin == "" {
		return libvetchi.DomainNotVerified, nil
	}

	err = h.db.CreateEmployer(ctx, db.Employer{
		ClientID:      domain,
		OnboardStatus: string(libvetchi.DomainVerifiedOnboardPending),
		OnboardAdmin:  admin,
	})
	if err != nil {
		h.logger.Error("create employer failed", "domain", domain, "error", err)
		return libvetchi.DomainNotVerified, err
	}

	return libvetchi.DomainVerifiedOnboardPending, nil
}
