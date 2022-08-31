package ginserver

import (
	"fmt"
	"myself_framwork/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tlsserverconfig = true

var (
	serverfileca         string
	serverfileprivatekey string
	serverfilepubkey     string
	servertls12client    bool
)

func GinServerUp(listenAddr string, router *gin.Engine) error {
	cfg := *utils.GetServerConfig()
	fmt.Println("[TLS.1.2]:", cfg.Servertls12client)
	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           router,
		TLSConfig:         utils.GetServerTlsConfig(),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
	}

	if cfg.Servertls12client == "ON" {
		return srv.ListenAndServeTLS(cfg.Serverfilepubkey, cfg.Serverfileprivatekey)
	}
	return srv.ListenAndServe()
}
