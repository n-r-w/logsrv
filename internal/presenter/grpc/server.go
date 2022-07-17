// Package grpcsrv ...
package grpcsrv

import (
	"fmt"
	"net"

	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/presenter"
	grpc_gen "github.com/n-r-w/logsrv/internal/presenter/grpc/generated/proto"
	"github.com/n-r-w/nerr"
	"google.golang.org/grpc"
)

type Service struct {
	grpc_gen.UnimplementedLogsrvServer

	cfg       *config.Config
	logger    lg.Logger
	logRepo   presenter.LogInterface
	presenter *presenter.Service
	srv       *grpc.Server
	notify    chan error
}

func New(logger lg.Logger, logRepo presenter.LogInterface, cfg *config.Config, presenter *presenter.Service, host, port string) *Service {
	s := &Service{
		cfg:       cfg,
		logger:    logger,
		logRepo:   logRepo,
		presenter: presenter,
		srv:       grpc.NewServer(),
		notify:    make(chan error, 1),
	}

	s.srv = grpc.NewServer()
	grpc_gen.RegisterLogsrvServer(s.srv, s)

	go s.start(host, port)

	return s
}

func (s *Service) Notify() <-chan error {
	return s.notify
}

func (s *Service) start(host, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.notify <- nerr.New("net.Listen error", nerr.New(err))
	}

	s.logger.Info("grpc server started on %s", addr)
	if err := s.srv.Serve(lis); err != nil {
		s.notify <- nerr.New("GrpcServer error", err)
	}
	close(s.notify)
}

func (s *Service) Shutdown() {
	s.srv.GracefulStop()
}
