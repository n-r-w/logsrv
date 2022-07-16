//go:build wireinject
// +build wireinject

// Package di. Автоматическое внедрение зависимостей
package di

import (
	"github.com/google/wire"
	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/presenter/rest"
	"github.com/n-r-w/logsrv/internal/repo/psql"
	"github.com/n-r-w/logsrv/internal/repo/wbuf"
	"github.com/n-r-w/postgres"
)

type Container struct {
	Logger  lg.Logger
	Config  *config.Config
	DB      *postgres.Postgres
	LogRepo *psql.LogRepo
	WBuf    *wbuf.WBuf
	Router  *httprouter.RouterData
	Rest    *rest.Service
}

// NewContainer - создание DI контейнера с помощью google wire
func NewContainer(logger lg.Logger, config *config.Config, dbUrl postgres.Url, dbOptions []postgres.Option) (*Container, func(), error) {
	panic(wire.Build(
		postgres.New,

		wire.Bind(new(rest.LogInterface), new(*wbuf.WBuf)),
		psql.NewLog,

		wire.Bind(new(wbuf.WBufInterface), new(*psql.LogRepo)),
		wbuf.New,

		wire.Bind(new(httprouter.Router), new(*httprouter.RouterData)),
		httprouter.New,

		rest.New,

		wire.Struct(new(Container), "*"),
	))
}
