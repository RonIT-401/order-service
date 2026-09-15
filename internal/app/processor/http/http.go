package rprocessor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RonIT-401/order-service/internal/app/config/section"
	rhandler "github.com/RonIT-401/order-service/internal/app/handler/http"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHTTP(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.NoRoute(handleNotFound)

	vGenericRegHealthCheck(router, hHealth)

	for _, route := range router.Routes() {
		log.Printf("Route registered: %s %s", route.Method, route.Path)
	}

	addr := fmt.Sprintf(":%d", cfg.ListenPort)

	return &httpProc{
		server: http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
		addr: addr,
	}
}

func (p *httpProc) Serve() error {
	log.Printf("Starting HTTP server on %s", p.addr)
	return p.server.ListenAndServe()
}
