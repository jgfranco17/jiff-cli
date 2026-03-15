package errorhandling

import (
	"fmt"
)

const (
	ExitCodeGenericError       int = 1
	ExitCodeConfigurationError int = 2
	ExitCodeOperationFailed    int = 3
	ExitCodeInvalidInput       int = 4
)

type ExitError struct {
	Err      error
	ExitCode int
	Solution string
}

func (e *ExitError) WithSolution(solution string) *ExitError {
	e.Solution = solution
	return e
}

func (e ExitError) Error() string {
	return e.Err.Error()
}

func (e ExitError) Unwrap() error {
	return e.Err
}

func (e ExitError) Render() string {
	message := fmt.Sprintf("[ERROR] %s", e.Err.Error())

	if e.Solution != "" {
		message += fmt.Sprintf("\n[SOLUTION]: %s", e.Solution)
	}

	message += fmt.Sprintf("Jiff exited with code %d\n", e.ExitCode)
	return message
}
