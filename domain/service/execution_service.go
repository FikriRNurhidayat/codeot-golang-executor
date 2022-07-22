package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fikrirnurhidayat/codeot-golang-executor/domain/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type ExecutionService interface {
	Execute(ctx context.Context, params *dto.ExecutionParams) (*dto.ExecutionResult, error)
}

type ExecutionServiceImpl struct {
	logger grpclog.LoggerV2
}

func NewExecutionService(logger grpclog.LoggerV2) *ExecutionServiceImpl {
	return &ExecutionServiceImpl{
		logger: logger,
	}
}

func (s *ExecutionServiceImpl) Execute(ctx context.Context, params *dto.ExecutionParams) (*dto.ExecutionResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Create directory
	dirnamePattern := fmt.Sprintf("*-execution-%s", params.EID)
	dir, err := os.MkdirTemp("tmp", dirnamePattern)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dir)

	// Create execution file
	file := filepath.Join(dir, "main.go")
	if err := os.WriteFile(file, []byte(params.Code), 0444); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Execute code
	cmd := exec.Command("go", "run", file)
	cmd.Stdin = strings.NewReader(params.Stdin)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	var result *dto.ExecutionResult

	if err = cmd.Run(); err != nil {
		result = &dto.ExecutionResult{
			Stdout:   stdout.String(),
			ExitCode: uint32(cmd.ProcessState.ExitCode()),
		}

		return result, status.Error(codes.InvalidArgument, stderr.String())
	}

	result = &dto.ExecutionResult{
		Stdout:   stdout.String(),
		ExitCode: uint32(cmd.ProcessState.ExitCode()),
	}

	return result, nil
}
