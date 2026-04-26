package grpc

import (
	"context"
	"fmt"

	pb "github.com/logservice/github.com/logservice/pkg/proto"
	
	"github.com/logservice/internal/model"
	"github.com/logservice/internal/repo"
)

type Server struct {
	pb.UnimplementedLogServiceServer
	repo repo.LogRepository
}

func NewServer(r repo.LogRepository) *Server {
	return &Server{repo: r}
}

func (s *Server) IngestLog(ctx context.Context, req *pb.IngestRequest) (*pb.LogEntry, error) {

	// validation
	if req.Level == "" || req.Service == "" || req.Message == "" {
		return nil, fmt.Errorf("missing fields")
	}

	log := &model.Log{
		Level:   req.Level,
		Service: req.Service,
		Message: req.Message,
	}

	err := s.repo.Insert(ctx, log)
	if err != nil {
		return nil, err
	}

	// map to proto
	return &pb.LogEntry{
		Id:        log.ID.String(),
		Level:     log.Level,
		Service:   log.Service,
		Message:   log.Message,
		CreatedAt: log.CreatedAt.String(),
	}, nil
}


func (s *Server) QueryLogs(req *pb.QueryRequest, stream pb.LogService_QueryLogsServer) error {

	logs, err := s.repo.List(stream.Context(), repo.Filter{
		Level:   req.Level,
		Service: req.Service,
		Limit:   int(req.Limit),
	})
	if err != nil {
		return err
	}

	for _, l := range logs {
		err := stream.Send(&pb.LogEntry{
			Id:        l.ID.String(),
			Level:     l.Level,
			Service:   l.Service,
			Message:   l.Message,
			CreatedAt: l.CreatedAt.String(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}