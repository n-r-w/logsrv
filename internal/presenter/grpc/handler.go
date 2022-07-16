package grpcsrv

import (
	"context"

	grpc_gen "github.com/n-r-w/logsrv/internal/presenter/grpc/generated/proto"
)

func (s *GrpcServer) SendLog(context.Context, *grpc_gen.SendOptions) (*grpc_gen.Error, error) {
	return &grpc_gen.Error{
		Code: 0,
		Text: "",
	}, nil
}

func (s *GrpcServer) SearchLog(context.Context, *grpc_gen.SearchOptions) (*grpc_gen.SearchLogReply, error) {
	return &grpc_gen.SearchLogReply{
		Error: &grpc_gen.Error{
			Code: 0,
			Text: "",
		},
		Records: []*grpc_gen.LogRecord{},
	}, nil
}
