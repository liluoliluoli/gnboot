package server

import (
	"net/http"

	"gnboot/internal/service"
)

func HealthHandler(svc *service.GnbootService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pub/healthcheck", svc.HealthCheck)
	return mux
}
