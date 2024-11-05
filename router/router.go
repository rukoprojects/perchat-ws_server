package router

import (
	"net/http"
	"ws_server/services"
)

func Router(addr string) *http.ServeMux {
	server := services.NewServer(addr)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", server.HandleClient)
	return mux
}
