package http

import (
	"gateway/internal/config"
	"gateway/internal/http/middleware"
	"net/http"

	"gateway/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	registerHealthRoutes(router)
	registerProxyRoutes(router, cfg)

	return router
}

func registerHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func registerProxyRoutes(r *gin.Engine, cfg *config.Config) {
	authProxy := handlers.NewProxyHandler(cfg.AuthServiceURL)
	r.POST("/auth/login", middleware.Unauthenticated(cfg.JWTSecret), authProxy.Login())
	r.POST("/auth/register", middleware.Unauthenticated(cfg.JWTSecret), authProxy.Register())

	userProxy := handlers.NewProxyHandler(cfg.UserServiceURL)
	r.Any("users/*proxyPath", middleware.JWT(cfg.JWTSecret), userProxy.Handle())
}
