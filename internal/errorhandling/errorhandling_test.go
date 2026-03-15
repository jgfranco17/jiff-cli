package errorhandling

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			base := &ExitError{
				Err:      errors.New("base error"),
				ExitCode: ExitCodeConfigurationError,
			}
			e := base.WithSolution(tc.solution)

			assert.Equal(t, tc.solution, e.Solution)

			rendered := e.String()
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
			code:          ExitCodeOperationFailed,
			wantFragments: []string{"[ERROR]", "disk full", fmt.Sprintf("code %d", ExitCodeOperationFailed)},
		},
		{
			name:          "render includes solution when set",
			err:           errors.New("not found"),
			code:          ExitCodeInvalidInput,
			solution:      "check the file path",
			wantFragments: []string{"[ERROR]", "not found", "[SOLUTION]", "check the file path"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base := &ExitError{
				Err:      tc.err,
				ExitCode: tc.code,
			}
			e := base.WithSolution(tc.solution)
			rendered := e.String()
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
			wrapped:   ErrFailOnDiff,
			targetErr: ErrFailOnDiff,
			wantIs:    true,
		},
		{
			name:      "errors.Is returns false for unrelated sentinel",
			wrapped:   errors.New("unrelated"),
			targetErr: ErrFailOnDiff,
			wantIs:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := ExitError{
				Err:      tc.wrapped,
				ExitCode: ExitCodeGenericError,
			}
			assert.Equal(t, tc.wantIs, errors.Is(e, tc.targetErr))
		})
	}
}
