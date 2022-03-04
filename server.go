package tt

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative gRPC.proto

type GrpcServer struct {
	UnimplementedTtServer
}

func (s *GrpcServer) StartTimer(ctx context.Context, in *StartParameters) (*GRpcTimer, error) {
	if in == nil {
		return nil, fmt.Errorf("start: %w: paramteres are nil", ErrInvalidData)
	}
	t := &GRpcTimer{
		Id:      uuid.Must(uuid.NewRandom()).String(),
		Project: in.Project,
		Task:    in.Task,
		Tags:    in.Tags,
		Start:   in.Timestamp,
		Stop:    "",
	}
	err := t.Validate()
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	return t, nil
}

func (s *GrpcServer) StopTimer(ctx context.Context, in *StopParameters) (*GRpcTimer, error) {
	return nil, ErrNotImplemented
}

func (s *GrpcServer) ResumeTimer(ctx context.Context, in *ResumeParameters) (*GRpcTimer, error) {
	return nil, ErrNotImplemented
}
