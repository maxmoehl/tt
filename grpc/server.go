package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc.proto

type Server struct {
	UnimplementedTtServer
}

func (s *Server) StartTimer(ctx context.Context, in *StartParameters) (*Timer, error) {
	if in == nil {
		return nil, fmt.Errorf("start: %w: paramteres are nil", tt.ErrInvalidData)
	}
	timestamp, err := time.Parse(time.RFC3339, in.Timestamp)
	if err != nil {
		return nil, err
	}
	timer, err := tt.Start(in.Project, in.Task, in.Tags, timestamp, 0)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	return fromTt(timer), nil
}

func (s *Server) StopTimer(ctx context.Context, in *StopParameters) (*Timer, error) {
	return nil, tt.ErrNotImplemented
}

func (s *Server) ResumeTimer(ctx context.Context, in *ResumeParameters) (*Timer, error) {
	return nil, tt.ErrNotImplemented
}
