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
	"github.com/n-r-w/logsrv/internal/repo/psql"
	"github.com/n-r-w/logsrv/internal/repo/wbuf"
	"github.com/n-r-w/logsrv/internal/repo/wdispatch"
	"github.com/n-r-w/postgres"
)

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

// NewContainer - создание DI контейнера с помощью google wire
func NewContainer(logger lg.Logger, config *config.Config, dbUrl postgres.Url, dbOptions []postgres.Option) (*Container, func(), error) {
	panic(wire.Build(
		postgres.New,

		wire.Bind(new(presenter.LogInterface), new(*psql.LogRepo)),
		psql.NewLog,

		// Репозиторий для работы с логами, буфер для асинхронной записи в БД и диспетчер для параллельной записи в БД, реализуют интерфейс
		// usecase.LogInterface. Поэтому можно писать как сразу в logRepo, так и через буфер или диспетчер
		// В данном случае запись идет по цепочке ->WBuf->Dispatcher->LogRepo
		// WBuf накапливает входящие записи в буфере и когда он достигает определенного уровня, сбрасывает их на запись в Dispatcher
		// Dispatcher создает очередь задач на запись в БД через пул коннектов
		wdispatch.New,
		wire.Bind(new(wbuf.WBufInterface), new(*psql.LogRepo)),
		wbuf.New,

		wire.Bind(new(httprouter.Router), new(*httprouter.RouterData)),
		httprouter.New,

		presenter.New,

		wire.Struct(new(Container), "*"),
	))
}
