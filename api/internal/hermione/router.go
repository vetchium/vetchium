package hermione

import (
	"fmt"
	"net/http"
)

func (h *Hermione) Run() error {
	RegisterEmployerRoutes(h)
	RegisterHubRoutes(h)

	port := fmt.Sprintf(":%d", h.Config().Port)
	return http.ListenAndServe(port, nil)
}
