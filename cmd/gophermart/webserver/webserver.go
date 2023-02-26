package webserver

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/webserver"
)

type Webserver struct {
	Server *http.Server
	Engine *gin.Engine

	log *zap.Logger
}

func New(
	cfg webserver.Config,
	log *zap.Logger,
) *Webserver {
	e := gin.New()

	e.RedirectTrailingSlash = true
	e.UnescapePathValues = true
	e.HandleMethodNotAllowed = true

	e.Use(ginzap.Ginzap(log, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(log, true))

	return &Webserver{
		Server: &http.Server{ //nolint:exhaustruct
			Addr:    cfg.Address,
			Handler: e,

			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		Engine: e,
		log:    log,
	}
}
