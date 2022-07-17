package grpcsrv

import (
	"context"
	"encoding/json"

	"github.com/gogo/status"
	"github.com/n-r-w/logsrv/internal/entity"
	grpc_gen "github.com/n-r-w/logsrv/internal/presenter/grpc/generated/proto"
	"github.com/n-r-w/nerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) SendLog(ctx context.Context, data *grpc_gen.SendOptions) (*emptypb.Empty, error) {
	if err := s.presenter.CheckRights(data.Token, true, false); err != nil {
		err = nerr.New(err)
		s.logger.Err(err)
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	var records []entity.LogRecord

	for _, r := range data.Records {
		req := entity.LogRecord{
			ID:          0,
			LogTime:     r.LogTime.AsTime(),
			Service:     r.Service,
			Source:      r.Source,
			Category:    r.Category,
			Level:       r.Level,
			Session:     r.Session,
			Info:        r.Info,
			Properties:  r.Properties,
			Url:         r.Url,
			HttpType:    r.HttpType,
			HttpCode:    int(r.HttpCode),
			ErrorCode:   int(r.ErrorCode),
			HttpHeaders: r.HttpHeaders,
			Body:        r.Body,
		}

		records = append(records, req)
	}

	if err := s.logRepo.Insert(records); err != nil {
		err = nerr.New(err)
		s.logger.Err(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) SearchLog(ctx context.Context, data *grpc_gen.SearchOptions) (*grpc_gen.SearchLogReply, error) {
	if err := s.presenter.CheckRights(data.Token, false, true); err != nil {
		err = nerr.New(err)
		s.logger.Err(err)
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	req := entity.SearchRequest{
		AndOp: data.And,
	}

	for _, c := range data.Criteria {
		criteria := entity.SearchCriteria{
			AndOp:       c.And,
			From:        c.From.AsTime(),
			To:          c.To.AsTime(),
			Service:     c.Service,
			Source:      c.Source,
			Category:    c.Category,
			Level:       c.Level,
			Session:     c.Session,
			Info:        c.Info,
			Properties:  c.Properties,
			Url:         c.Url,
			HttpType:    c.HttpType,
			HttpCode:    int(c.HttpCode),
			ErrorCode:   int(c.ErrorCode),
			HttpHeaders: c.HttpHeaders,
			BodyValues:  map[string]json.RawMessage{},
			Body:        c.Body,
		}

		for k, v := range c.BodyValues {
			criteria.BodyValues[k] = v
		}

		req.Criteria = append(req.Criteria, criteria)
	}

	records, _, err := s.logRepo.Find(req)
	if err != nil {
		err = nerr.New(err)
		s.logger.Err(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	respond := []*grpc_gen.LogRecord{}
	for _, r := range records {
		respond = append(respond,
			&grpc_gen.LogRecord{
				RecTime:     timestamppb.New(r.RecordTime),
				LogTime:     timestamppb.New(r.LogTime),
				Service:     r.Service,
				Source:      r.Source,
				Category:    r.Category,
				Level:       r.Level,
				Session:     r.Session,
				Info:        r.Info,
				Url:         r.Url,
				HttpType:    r.HttpType,
				HttpCode:    int32(r.HttpCode),
				ErrorCode:   int32(r.ErrorCode),
				HttpHeaders: r.HttpHeaders,
				Properties:  r.Properties,
				Body:        r.Body,
			})
	}

	return &grpc_gen.SearchLogReply{
		Records: respond,
	}, nil
}
