package commandline

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jgfranco17/dev-tooling-go/logging"
	"github.com/jgfranco17/jiff-cli/internal/diffs"
	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CLI represents the command-line interface instance of the application.
// Allows for configuration of global flags and holds the root command.
type CLI struct {
	rootCmd *cobra.Command

	verbosity  int
	failOnDiff bool
}

// NewCommand creates a new instance of Command.
func New(name string, description string, version string) *CLI {
	var verbosity int
	var failOnDiff bool
	ctx := context.Background()

	root := &cobra.Command{
		Use:               name,
		Version:           version,
		Short:             description,
		Example:           fmt.Sprintf("%s file1.json file2.json", name),
		Args:              cobra.ExactArgs(2),
		PersistentPreRunE: preRunFunc(ctx),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.FromContext(cmd.Context())

			sourceFile, targetFile := args[0], args[1]
			sourceData, err := loadJsonFile(sourceFile)
			if err != nil {
				return errorhandling.ExitError{
					Err:      fmt.Errorf("failed to read source file: %w", err),
					ExitCode: errorhandling.ExitCodeInvalidInput,
					Solution: "Please ensure the source file exists and is accessible.",
				}
			}
			targetData, err := loadJsonFile(targetFile)
			if err != nil {
				return errorhandling.ExitError{
					Err:      fmt.Errorf("failed to read target file: %w", err),
					ExitCode: errorhandling.ExitCodeInvalidInput,
					Solution: "Please ensure the target file exists and is accessible.",
				}
			}
			logger.WithFields(logrus.Fields{
				"source": sourceFile,
				"target": targetFile,
			}).Debug("Comparing files")

			diffs, err := diffs.CompareJSON(sourceData, targetData)
			if err != nil {
				return errorhandling.ExitError{
					Err:      fmt.Errorf("failed to compare files: %w", err),
					ExitCode: errorhandling.ExitCodeOperationFailed,
				}
			}

			if !diffs.IsEmpty() {
				logger.WithFields(logrus.Fields{
					"added":   len(diffs.Added),
					"removed": len(diffs.Removed),
					"changed": len(diffs.Changed),
				}).Info("Differences found")

				diffs.Render(cmd.OutOrStdout())
				if failOnDiff {
					return errorhandling.ExitError{
						Err:      errorhandling.ErrFailOnDiff,
						ExitCode: errorhandling.ExitCodeOperationFailed,
					}
				}
				return nil
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity (-v or -vv)")
	root.PersistentFlags().BoolVarP(&failOnDiff, "fail", "x", false, "Return nonzero if differences are found")
	return &CLI{
		rootCmd:    root,
		verbosity:  verbosity,
		failOnDiff: failOnDiff,
	}
}

// Execute executes the root command.
func (cr *CLI) Execute() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errorhandling.ExitError{
				Err:      fmt.Errorf("unexpected internal error: %v", r),
				ExitCode: errorhandling.ExitCodeCrashError,
				Solution: "Please report this issue to the developers with steps to reproduce.",
			}
		}
	}()
	err = cr.rootCmd.Execute()
	return err
}

type cobraRunFunc func(cmd *cobra.Command, args []string) error

func preRunFunc(ctx context.Context) cobraRunFunc {
	return func(cmd *cobra.Command, args []string) error {
		verbosity, _ := cmd.Flags().GetCount("verbose")
		var level logrus.Level
		switch verbosity {
		case 1:
			level = logrus.InfoLevel
		case 2:
			level = logrus.DebugLevel
		case 3:
			level = logrus.TraceLevel
		default:
			level = logrus.WarnLevel
		}

		logger := logging.New(cmd.ErrOrStderr(), level)
		ctx = logging.WithContext(cmd.Context(), logger)

		ctx, cancel := context.WithCancel(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			select {
			case <-c:
				logger.Warn("Received shutdown signal, exiting...")
				cancel()
			case <-ctx.Done():
			}
		}()

		cmd.SetContext(ctx)
		return nil
	}

}

func loadJsonFile(filePathInput string) (io.Reader, error) {
	if filepath.Ext(filePathInput) != ".json" {
		return nil, fmt.Errorf("invalid file type: %s", filePathInput)
	}
	data, err := os.Open(filePathInput)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil

}
