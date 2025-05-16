package hubemp

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/config"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
	"github.com/vetchium/vetchium/typespec/libgranger"
)

func GetEmployerDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request hub.GetEmployerDetailsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: Get the employer name

		grangerReq := libgranger.GetEmployerCountsRequest{
			Domain: request.Domain,
		}

		grangerReqBytes, err := json.Marshal(grangerReq)
		if err != nil {
			h.Dbg("failed to marshal request", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var employerCounts libgranger.EmployerCounts
		grangerResp, err := http.Post(
			config.GrangerBaseURL+"/internal/get-employer-counts",
			"application/json",
			bytes.NewBuffer(grangerReqBytes),
		)
		if err != nil {
			h.Dbg("failed to post request", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer grangerResp.Body.Close()

		if grangerResp.StatusCode != http.StatusOK {
			h.Dbg(
				"failed to get employer counts",
				"status",
				grangerResp.StatusCode,
			)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(grangerResp.Body).Decode(&employerCounts); err != nil {
			h.Dbg("failed to decode employer counts", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		hubEmployerDetails := hub.HubEmployerDetails{
			Name:                   request.Domain,
			VerifiedEmployeesCount: employerCounts.VerifiedEmployeesCount,
			ActiveOpeningsCount:    employerCounts.ActiveOpeningsCount,
		}

		if err := json.NewEncoder(w).Encode(hubEmployerDetails); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
