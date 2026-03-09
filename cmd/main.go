package main

import (
	"gateway/config"
	internalHttp "gateway/internal/http"
	"gateway/server"
)

func main() {
	cfg := config.LoadConfig()

	router := internalHttp.SetupRouter(cfg)

	srv := server.NewHTTPServer(router, cfg.Port, cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout)

	server.StartServer(srv, cfg)

	// --- Graceful shutdown ---
	server.WaitForShutdown(srv, cfg)
}
