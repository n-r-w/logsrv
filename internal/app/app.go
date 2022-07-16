// Package app ...
package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n-r-w/httpserver"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/di"
	grpcsrv "github.com/n-r-w/logsrv/internal/presenter/grpc"
	"github.com/n-r-w/postgres"
)

const version = "1.0.6"

func Start(cfg *config.Config, logger lg.Logger) {
	logger.Info("logsrv %s", version)

	// инициализация DI контейнера
	con, _, err := di.NewContainer(logger, cfg, postgres.Url(cfg.DatabaseURL),
		[]postgres.Option{
			postgres.MaxConns(cfg.MaxDbSessions),
			postgres.MaxMaxConnIdleTime(time.Duration(cfg.MaxDbSessionIdleTime) * time.Second),
		},
	)
	if err != nil {
		logger.Err(err)
		return
	}

	// запуск прослоек буферизации записи в БД
	// con.WDispatch.Start()
	con.WBuf.Start()

	// запускаем http сервер
	httpServer := httpserver.New(con.Router.Handler(), logger,
		httpserver.Address(con.Config.RestHost, con.Config.RestPort),
		httpserver.ReadTimeout(time.Millisecond*time.Duration(con.Config.HttpReadTimeout)), // меняет также ReadHeaderTimeout, IdleTimeout
		httpserver.WriteTimeout(time.Millisecond*time.Duration(con.Config.HttpWriteTimeout)),
		httpserver.ShutdownTimeout(time.Second*time.Duration(con.Config.HttpShutdownTimeout)),
	)

	// запускаем grpc сервер
	grpcServer := grpcsrv.NewGrpcServer(logger, con.Config.GrpcHost, con.Config.GrpcPort)

	// ждем сигнал от сервера или нажатия ctrl+c
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		logger.Info("shutdown...")
	case err := <-httpServer.Notify():
		if err != nil {
			logger.Error("http server notification: %v", err)
		}
	case err := <-grpcServer.Notify():
		if err != nil {
			logger.Error("grpc server notification: %v", err)
		}
	}

	// ждем завершения
	grpcServer.Shutdown()
	logger.Info("grpc shutdown ok")

	errHttp := httpServer.Shutdown()
	con.WBuf.Stop()

	if errHttp != nil {
		logger.Error("http shutdown error: %v", errHttp)
	} else {
		logger.Info("http shutdown ok")
	}

}
