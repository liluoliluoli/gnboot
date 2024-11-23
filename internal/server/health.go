package server

import (
	"net/http"

	"gnboot/internal/adaptor"
)

func HealthHandler(svc *adaptor.GnbootService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pub/healthcheck", svc.HealthCheck)
	return mux
}
