package hermione

import "net/http"

func (h *Hermione) getCostCenters(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}
