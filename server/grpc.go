package server

import "context"

type GrpcServer struct {
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	return nil
}

func (s *GrpcServer) Stop(ctx context.Context) error {
	return nil
}
