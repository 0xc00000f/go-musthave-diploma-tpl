package service

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/webserver"
)

type APIService struct {
	cfg *config.Config

	webserver *webserver.Webserver
	storage   *storage.Storage

	logger *zap.Logger
}

func New(cfg *config.Config) *APIService {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	api := &APIService{ //nolint:exhauststruct
		cfg:    cfg,
		logger: logger,
	}

	api.setupDB()
	api.setupWebserver()

	return api
}

func (api *APIService) setupDB() {
	db, err := storage.New(api.cfg.Pgsql)
	if err != nil {
		panic(err)
	}

	db.Logger = api.logger

	api.storage = db
}

func (api *APIService) setupWebserver() {
	api.webserver = webserver.New(api.cfg.Webserver, api.logger)
}

func (api *APIService) CreateHTTPEndpoints() {
	api.webserver.Engine.GET("ping", handlers.Ping())
}

func (api *APIService) Run() {
	api.logger.Info("server is about to listen", zap.String("addr", api.webserver.Server.Addr))

	if err := api.webserver.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		api.logger.Fatal("server listen failed")
	}
}
