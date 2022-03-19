package server

import (
  "context"

  "google.golang.org/grpc/status"
  "google.golang.org/grpc/codes"
  "github.com/fikrirnurhidayat/codeot-golang-executor/app/service"
  "github.com/fikrirnurhidayat/codeotapis/proto"
)

type server struct {
  proto.UnimplementedGoexeServer

  executionService service.ExecutionService
}

func New() *server {
  return &server{
    executionService: service.NewExecutionService(),
  }
}

func (s *server) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
  stdout, exitCode, err := s.executionService.Execute(ctx, req.GetEid(), req.GetCode(), req.GetStdin())

  if err != nil {
    return nil, status.Error(codes.Internal, err.Error())
  }

  res := proto.ExecuteResponse{
    Stdout: stdout,
    ExitCode: exitCode,
  }

  return &res, nil
}
