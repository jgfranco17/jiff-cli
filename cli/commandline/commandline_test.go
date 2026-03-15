package commandline

import (
	"errors"
	"testing"

	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute_PanicRecovery(t *testing.T) {
	tests := []struct {
		name           string
		panicVal       any
		wantMsgContain string
	}{
		{
			name:           "string panic is wrapped in ExitError",
			panicVal:       "something went terribly wrong",
			wantMsgContain: "something went terribly wrong",
		},
		{
			name:           "error panic is wrapped in ExitError",
			panicVal:       errors.New("nil pointer dereference"),
			wantMsgContain: "nil pointer dereference",
		},
		{
			name:           "integer panic is wrapped in ExitError",
			panicVal:       42,
			wantMsgContain: "42",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cli := CLI{
				rootCmd: &cobra.Command{
					Use:           "test",
					SilenceErrors: true,
					SilenceUsage:  true,
					RunE: func(cmd *cobra.Command, args []string) error {
						panic(tc.panicVal)
					},
				},
			}

			err := cli.Execute()
			require.Error(t, err)

			var exitErr errorhandling.ExitError
			require.ErrorAs(t, err, &exitErr)
			assert.Equal(t, errorhandling.ExitCodeCrashError, exitErr.ExitCode)
			assert.Contains(t, exitErr.Error(), tc.wantMsgContain)
			assert.NotEmpty(t, exitErr.Solution)
		})
	}
}

func TestExecute_NoPanic_ReturnsNil(t *testing.T) {
	root := &cobra.Command{
		Use:           "test",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cli := &CLI{rootCmd: root}
	err := cli.Execute()
	assert.NoError(t, err)
}

func TestExecute_NoPanic_PropagatesRunError(t *testing.T) {
	sentinel := errors.New("expected run error")
	root := &cobra.Command{
		Use:           "test",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sentinel
		},
	}
	cli := &CLI{rootCmd: root}
	err := cli.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, sentinel)
}
