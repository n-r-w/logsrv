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
	"github.com/n-r-w/logsrv/internal/repo/psql"
	"github.com/n-r-w/logsrv/internal/repo/wbuf"
	"github.com/n-r-w/logsrv/internal/repo/wdispatch"
	"github.com/n-r-w/postgres"
)

// Injectors from wire.go:

// NewContainer - создание DI контейнера с помощью google wire
func NewContainer(logger lg.Logger, config2 *config.Config, dbUrl postgres.Url, dbOptions []postgres.Option) (*Container, func(), error) {
	postgresPostgres, err := postgres.New(dbUrl, logger, dbOptions...)
	if err != nil {
		return nil, nil, err
	}
	logRepo := psql.NewLog(postgresPostgres, config2)
	dispatcher := wdispatch.New(logRepo, logger, config2)
	wBuf := wbuf.New(logRepo, logger, config2)
	routerData := httprouter.New(logger)
	presenterPresenter, err := presenter.New(routerData, logRepo, config2)
	if err != nil {
		return nil, nil, err
	}
	container := &Container{
		Logger:    logger,
		Config:    config2,
		DB:        postgresPostgres,
		LogRepo:   logRepo,
		WDispatch: dispatcher,
		WBuf:      wBuf,
		Router:    routerData,
		Presenter: presenterPresenter,
	}
	return container, func() {
	}, nil
}

// wire.go:

type Container struct {
	Logger    lg.Logger
	Config    *config.Config
	DB        *postgres.Postgres
	LogRepo   *psql.LogRepo
	WDispatch *wdispatch.Dispatcher
	WBuf      *wbuf.WBuf
	Router    *httprouter.RouterData
	Presenter *presenter.Presenter
}
