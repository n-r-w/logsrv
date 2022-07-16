// Package grpcsrv ...
package grpcsrv

import (
	"fmt"
	"net"

	"github.com/n-r-w/lg"
	grpc_gen "github.com/n-r-w/logsrv/internal/presenter/grpc/generated/proto"
	"github.com/n-r-w/nerr"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpc_gen.UnimplementedLogsrvServer

	logger lg.Logger
	srv    *grpc.Server
	notify chan error
}

func NewGrpcServer(logger lg.Logger, host, port string) *GrpcServer {

	s := &GrpcServer{
		logger: logger,
		srv:    grpc.NewServer(),
		notify: make(chan error, 1),
	}

	s.srv = grpc.NewServer()
	grpc_gen.RegisterLogsrvServer(s.srv, s)

	go s.start(host, port)

	return s
}

func (s *GrpcServer) Notify() <-chan error {
	return s.notify
}

func (s *GrpcServer) start(host, port string) {
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

func (s *GrpcServer) Shutdown() {
	s.srv.GracefulStop()
}
