package router

import (
	"net/http"
	"ws_server/services"
)

func Router() *http.ServeMux {
	server := services.NewServer()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", server.HandleClient)
	return mux
}
