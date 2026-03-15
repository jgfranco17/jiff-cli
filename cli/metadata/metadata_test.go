package metadata

import (
	"errors"
	"strings"
	"testing"

	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		wantCode   int
		wantOutput ProjectInfo
	}{
		{
			name:  "valid full metadata",
			input: `{"author":"Jane Doe","name":"mytool","description":"A useful tool","version":"1.2.3","repository":"https://github.com/example/mytool"}`,
			wantOutput: ProjectInfo{
				Author:      "Jane Doe",
				Name:        "mytool",
				Description: "A useful tool",
				Version:     "1.2.3",
				Repository:  "https://github.com/example/mytool",
			},
		},
		{
			name:       "partial metadata with missing optional fields",
			input:      `{"name":"jiff","version":"0.0.1"}`,
			wantOutput: ProjectInfo{Name: "jiff", Version: "0.0.1"},
		},
		{
			name:       "empty object produces zero-value struct",
			input:      `{}`,
			wantOutput: ProjectInfo{},
		},
		{
			name:     "invalid JSON returns configuration error",
			input:    `not-json`,
			wantErr:  true,
			wantCode: errorhandling.ExitCodeConfigurationError,
		},
		{
			name:     "malformed JSON returns configuration error",
			input:    `{"name":"jiff"`,
			wantErr:  true,
			wantCode: errorhandling.ExitCodeConfigurationError,
		},
		{
			name:     "empty input returns configuration error",
			input:    ``,
			wantErr:  true,
			wantCode: errorhandling.ExitCodeConfigurationError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Load(strings.NewReader(tc.input))

			if tc.wantErr {
				require.Error(t, err)
				var exitErr *errorhandling.ExitError
				require.True(t, errors.As(err, &exitErr),
					"expected error to be *errorhandling.ExitError, got %T", err)
				assert.Equal(t, tc.wantCode, exitErr.ExitCode)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantOutput, got)
		})
	}
}
