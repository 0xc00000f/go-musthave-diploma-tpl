package service

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/auth"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers/balance"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers/balance/withdraw"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers/orders"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers/user"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/webserver"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/must"
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

	api := &APIService{ //nolint:exhaustruct
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
	userStorage := must.OK(api.storage.Users())
	orderStorage := must.OK(api.storage.Orders())

	authService := auth.New(api.cfg.Auth)

	api.webserver.Engine.GET("/ping", handlers.Ping())

	api.webserver.Engine.POST("/api/user/register", user.RegisterUser(userStorage, authService))
	api.webserver.Engine.POST("/api/user/login", user.AuthUser(userStorage, authService))

	authOnly := api.webserver.Engine.Group("")
	authOnly.Use(authService.AuthMiddleware())

	authOnly.POST("/api/user/orders", orders.CreateOrder(orderStorage))
	authOnly.GET("/api/user/orders", orders.FetchOrder(orderStorage))

	authOnly.GET("/api/user/balance", balance.FetchUserInfo(orderStorage))

	authOnly.POST("/api/user/balance/withdraw", withdraw.Request(orderStorage))
	authOnly.GET("/api/user/balance/withdrawals", withdraw.FetchUserInfo(orderStorage))
}

func (api *APIService) Run() {
	api.logger.Info("server is about to listen", zap.String("addr", api.webserver.Server.Addr))

	if err := api.webserver.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		api.logger.Fatal("server listen failed")
	}
}
