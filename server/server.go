package server

import (
	"github.com/fikrirnurhidayat/codeot-golang-executor/domain/service"
	"github.com/fikrirnurhidayat/codeotapis/proto"
)

type Server struct {
	proto.UnimplementedGoexeServer

	executionService service.ExecutionService
}

type ServerOpts func(*Server)

func New(sets ...ServerOpts) *Server {
	s := new(Server)

	for _, set := range sets {
		set(s)
	}

	return s
}

func WithExecutionService(executionService service.ExecutionService) ServerOpts {
	return func(s *Server) {
		s.executionService = executionService
	}
}
