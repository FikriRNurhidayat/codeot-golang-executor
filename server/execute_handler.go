package server

import (
	"context"

	"github.com/fikrirnurhidayat/codeot-golang-executor/domain/dto"
	"github.com/fikrirnurhidayat/codeotapis/proto"
)

func (s *Server) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
	result, err := s.executionService.Execute(ctx, &dto.ExecutionParams{
		EID:   req.GetEid(),
		Code:  req.GetCode(),
		Stdin: req.GetStdin(),
	})

	if err != nil {
		return nil, err
	}

	res := &proto.ExecuteResponse{
		Eid:      req.GetEid(),
		Stdout:   result.Stdout,
		ExitCode: result.ExitCode,
	}

	return res, nil
}
