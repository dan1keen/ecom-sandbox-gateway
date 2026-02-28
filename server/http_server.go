package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NewHTTPServer создает и возвращает HTTP сервер с конфигом
func NewHTTPServer(router *gin.Engine, port string, readTimeout, writeTimeout, idleTimeout time.Duration) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
