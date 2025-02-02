package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func CreateOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered CreateOpening")
		var createOpeningReq employer.CreateOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&createOpeningReq)
		if err != nil {
			h.Dbg("failed to decode create opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &createOpeningReq) {
			h.Dbg("validation failed", "createOpeningReq", createOpeningReq)
			return
		}
		h.Dbg("validated", "createOpeningReq", createOpeningReq)

		if createOpeningReq.YoeMax < createOpeningReq.YoeMin {
			h.Dbg("yoe_max < yoe min", "createOpeningReq", createOpeningReq)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(common.ValidationErrors{
				Errors: []string{"yoe_min", "yoe_max"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		// Validate tags
		totalTags := len(createOpeningReq.Tags) + len(createOpeningReq.NewTags)
		if totalTags == 0 {
			h.Dbg("no tags specified", "createOpeningReq", createOpeningReq)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(common.ValidationErrors{
				Errors: []string{"tags", "new_tags"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		if totalTags > 3 {
			h.Dbg(
				"too many tags specified",
				"createOpeningReq",
				createOpeningReq,
			)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(common.ValidationErrors{
				Errors: []string{"tags", "new_tags"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		if createOpeningReq.Salary != nil {
			if createOpeningReq.Salary.MinAmount > createOpeningReq.Salary.MaxAmount {
				h.Dbg("salary min > max", "createOpeningReq", createOpeningReq)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"salary"},
				})
				if err != nil {
					h.Err("failed to encode validation errors", "error", err)
				}
				return
			}

			if createOpeningReq.Salary.Currency == "" {
				h.Dbg("currency is empty", "createOpeningReq", createOpeningReq)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"currency"},
				})
				if err != nil {
					h.Err("failed to encode validation errors", "error", err)
				}
				return
			}
		}

		if len(createOpeningReq.RemoteCountryCodes) == 0 &&
			len(createOpeningReq.LocationTitles) == 0 {
			h.Dbg(
				"neither remote countries nor locations specified",
				"createOpeningReq",
				createOpeningReq,
			)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(common.ValidationErrors{
				Errors: []string{"remote_country_codes", "location_titles"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		openingID, err := h.DB().CreateOpening(r.Context(), createOpeningReq)
		if err != nil {
			if errors.Is(err, db.ErrNoRecruiter) ||
				errors.Is(err, db.ErrNoLocation) ||
				errors.Is(err, db.ErrNoHiringManager) ||
				errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("location or team or recruiter not found", "error", err)
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}

			h.Err("failed to create opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created opening", "openingID", openingID)
		err = json.NewEncoder(w).Encode(employer.CreateOpeningResponse{
			OpeningID: openingID,
		})
		if err != nil {
			h.Err("failed to encode create opening response", "error", err)
			return
		}
	}
}
