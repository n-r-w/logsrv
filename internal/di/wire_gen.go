// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/logsrv/internal/presenter/rest"
	"github.com/n-r-w/logsrv/internal/repo/psql"
	"github.com/n-r-w/logsrv/internal/repo/wbuf"
	"github.com/n-r-w/postgres"
)

// Injectors from wire.go:

// NewContainer - создание DI контейнера с помощью google wire
func NewContainer(logger lg.Logger, config2 *config.Config, dbUrl postgres.Url, dbOptions []postgres.Option) (*Container, func(), error) {
	service, err := postgres.New(dbUrl, logger, dbOptions...)
	if err != nil {
		return nil, nil, err
	}
	psqlService := psql.New(service, logger, config2)
	wbufService := wbuf.New(psqlService, logger, config2)
	httprouterService := httprouter.New(logger)
	presenterService := presenter.New(config2)
	restService := rest.New(httprouterService, wbufService, config2, presenterService)
	container := &Container{
		Logger:    logger,
		Config:    config2,
		DB:        service,
		LogRepo:   psqlService,
		WBuf:      wbufService,
		Router:    httprouterService,
		Rest:      restService,
		Presenter: presenterService,
	}
	return container, func() {
	}, nil
}

// wire.go:

type Container struct {
	Logger    lg.Logger
	Config    *config.Config
	DB        *postgres.Service
	LogRepo   *psql.Service
	WBuf      *wbuf.Service
	Router    *httprouter.Service
	Rest      *rest.Service
	Presenter *presenter.Service
}
