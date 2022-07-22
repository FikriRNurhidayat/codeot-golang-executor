package dto

type ExecutionParams struct {
	EID   string
	Code  string
	Stdin string
}

type ExecutionResult struct {
	Stdout   string
	ExitCode uint32
	err      error
}
