package granger

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/typespec/libgranger"
)

func (g *Granger) getEmployerCounts(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request libgranger.GetEmployerCountsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	activeOpeningsCounts, ok := g.getEmployerActiveJobCount(request.Domain)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	verifiedEmployeesCounts, ok := g.getEmployerEmployeeCount(request.Domain)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	counts := libgranger.EmployerCounts{
		ActiveOpeningsCount:    activeOpeningsCounts,
		VerifiedEmployeesCount: verifiedEmployeesCounts,
	}

	err := json.NewEncoder(w).Encode(counts)
	if err != nil {
		g.log.Err("JSON encoding failed", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (g *Granger) getEmployerActiveJobCount(domain string) (uint32, bool) {
	count, ok := g.employerActiveJobCountCache.Get(domain)
	if ok {
		g.log.Dbg("cache hit", "domain", domain, "count", count)
		return count, true
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	count, err := g.db.GetEmployerActiveJobCount(ctx, domain)
	if err != nil {
		g.log.Dbg("failed to get employer active job count", "error", err)
		return 0, false
	}

	g.log.Dbg("active job count", "domain", domain, "count", count)
	g.employerActiveJobCountCache.Set(domain, count, 32)
	return count, true
}

func (g *Granger) getEmployerEmployeeCount(domain string) (uint32, bool) {
	count, ok := g.employerEmployeeCountCache.Get(domain)
	if ok {
		return count, true
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	count, err := g.db.GetEmployerEmployeeCount(ctx, domain)
	if err != nil {
		g.log.Dbg("failed to get employer employee count", "error", err)
		return 0, false
	}

	g.log.Dbg("employee count", "domain", domain, "count", count)
	g.employerEmployeeCountCache.Set(domain, count, 32)
	return count, true
}
