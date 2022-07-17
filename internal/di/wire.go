//go:build wireinject
// +build wireinject

// Package di. Автоматическое внедрение зависимостей
package di

import (
	"github.com/google/wire"
	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/logsrv/internal/presenter/rest"
	"github.com/n-r-w/logsrv/internal/repo/psql"
	"github.com/n-r-w/logsrv/internal/repo/wbuf"
	"github.com/n-r-w/postgres"
)

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

// NewContainer - создание DI контейнера с помощью google wire
func NewContainer(logger lg.Logger, config *config.Config, dbUrl postgres.Url, dbOptions []postgres.Option) (*Container, func(), error) {
	panic(wire.Build(
		postgres.New,

		wire.Bind(new(presenter.LogInterface), new(*wbuf.Service)),
		psql.New,

		wire.Bind(new(wbuf.WBufInterface), new(*psql.Service)),
		wbuf.New,

		wire.Bind(new(httprouter.Router), new(*httprouter.Service)),
		httprouter.New,

		presenter.New,
		rest.New,

		wire.Struct(new(Container), "*"),
	))
}
