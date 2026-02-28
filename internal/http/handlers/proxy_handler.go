package handlers

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	loginRoutePrefix    string = "/login"
	registerRoutePrefix string = "/register"
)

type ProxyHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxyHandler(targetUrl string) *ProxyHandler {
	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatalf("Failed to parse target url: %v", err)
	}

	return &ProxyHandler{
		target: u,
		proxy:  httputil.NewSingleHostReverseProxy(u),
	}
}

func (p *ProxyHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = loginRoutePrefix
		c.Request.Host = p.target.Host
		p.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *ProxyHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = registerRoutePrefix
		c.Request.Host = p.target.Host
		p.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *ProxyHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		c.Request.Host = p.target.Host
		p.proxy.ServeHTTP(c.Writer, c.Request)
	}
}
