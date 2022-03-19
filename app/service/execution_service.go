package service

import (
  "os"
  "os/exec"
  "bytes"
  "strings"
  "path/filepath"
  "fmt"
  "context"
)

type ExecutionService interface {
  Execute(ctx context.Context, eid string, code string, stdin string) (stdout string, exitCode uint32, err error)
} 

type executionService struct {}

func NewExecutionService() *executionService {
  return &executionService{}
}

func (s *executionService) Execute(ctx context.Context, eid string, code string, stdin string) (stdout string, exitCode uint32, err error) {
  var out bytes.Buffer

  // Create directory
  dirnamePattern := fmt.Sprintf("*-execution-%s", eid)
  dir, err := os.MkdirTemp("tmp", dirnamePattern)
  defer os.RemoveAll(dir)

  // Create execution file
  file := filepath.Join(dir, "main.go")
  if err := os.WriteFile(file, []byte(code), 0444); err != nil {
    return stdout, exitCode, err
	}

  // Execute code
  cmd := exec.Command("go", "run", file)
  cmd.Stdin = strings.NewReader(stdin)
  cmd.Stdout = &out

  err = cmd.Run()

  // FIXME: Find out how to get exitCode
  return out.String(), exitCode, nil
}

