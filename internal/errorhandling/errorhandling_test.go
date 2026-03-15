package errorhandling_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		code         int
		wantMessage  string
		wantExitCode int
	}{
		{
			name:         "generic error code",
			err:          errors.New("something went wrong"),
			code:         errorhandling.ExitCodeGenericError,
			wantMessage:  "something went wrong",
			wantExitCode: errorhandling.ExitCodeGenericError,
		},
		{
			name:         "invalid input error code",
			err:          errors.New("bad input"),
			code:         errorhandling.ExitCodeInvalidInput,
			wantMessage:  "bad input",
			wantExitCode: errorhandling.ExitCodeInvalidInput,
		},
		{
			name:         "operation failed error code",
			err:          errors.New("operation failed"),
			code:         errorhandling.ExitCodeOperationFailed,
			wantMessage:  "operation failed",
			wantExitCode: errorhandling.ExitCodeOperationFailed,
		},
		{
			name:         "configuration error code",
			err:          errors.New("bad config"),
			code:         errorhandling.ExitCodeConfigurationError,
			wantMessage:  "bad config",
			wantExitCode: errorhandling.ExitCodeConfigurationError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := errorhandling.New(tc.err, tc.code)
			require.NotNil(t, e)
			assert.Equal(t, tc.wantMessage, e.Error())
			assert.Equal(t, tc.wantExitCode, e.ExitCode)
		})
	}
}

func TestWithSolution(t *testing.T) {
	tests := []struct {
		name             string
		solution         string
		wantInRender     string
		wantAbsentRender string
	}{
		{
			name:         "solution is stored and appears in Render",
			solution:     "try running with --verbose",
			wantInRender: "[SOLUTION]",
		},
		{
			name:             "empty solution omitted from Render",
			solution:         "",
			wantAbsentRender: "[SOLUTION]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := errorhandling.New(errors.New("base error"), errorhandling.ExitCodeGenericError).
				WithSolution(tc.solution)

			assert.Equal(t, tc.solution, e.Solution)

			rendered := e.Render()
			if tc.wantInRender != "" {
				assert.Contains(t, rendered, tc.wantInRender)
			}
			if tc.wantAbsentRender != "" {
				assert.NotContains(t, rendered, tc.wantAbsentRender)
			}
		})
	}
}

func TestExitError_Render(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		code          int
		solution      string
		wantFragments []string
	}{
		{
			name:          "render includes error prefix and message",
			err:           errors.New("disk full"),
			code:          errorhandling.ExitCodeOperationFailed,
			wantFragments: []string{"[ERROR]", "disk full", fmt.Sprintf("code %d", errorhandling.ExitCodeOperationFailed)},
		},
		{
			name:          "render includes solution when set",
			err:           errors.New("not found"),
			code:          errorhandling.ExitCodeInvalidInput,
			solution:      "check the file path",
			wantFragments: []string{"[ERROR]", "not found", "[SOLUTION]", "check the file path"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := errorhandling.New(tc.err, tc.code).WithSolution(tc.solution)
			rendered := e.Render()
			for _, fragment := range tc.wantFragments {
				assert.True(t, strings.Contains(rendered, fragment),
					"expected Render() to contain %q, got:\n%s", fragment, rendered)
			}
		})
	}
}

func TestExitError_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		wrapped   error
		targetErr error
		wantIs    bool
	}{
		{
			name:      "errors.Is resolves sentinel through ExitError",
			wrapped:   errorhandling.ErrFailOnDiff,
			targetErr: errorhandling.ErrFailOnDiff,
			wantIs:    true,
		},
		{
			name:      "errors.Is returns false for unrelated sentinel",
			wrapped:   errors.New("unrelated"),
			targetErr: errorhandling.ErrFailOnDiff,
			wantIs:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := errorhandling.New(tc.wrapped, errorhandling.ExitCodeGenericError)
			assert.Equal(t, tc.wantIs, errors.Is(e, tc.targetErr))
		})
	}
}
