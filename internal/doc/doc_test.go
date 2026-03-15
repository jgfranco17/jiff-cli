package doc_test

import (
	"testing"

	"github.com/jgfranco17/jiff-cli/internal/doc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMinimalCmd(name, version, short string) *cobra.Command {
	return &cobra.Command{
		Use:     name,
		Version: version,
		Short:   short,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
}

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name          string
		buildCmd      func() *cobra.Command
		wantFragments []string
		wantErr       bool
	}{
		{
			name: "header contains CLI name, version and description",
			buildCmd: func() *cobra.Command {
				return newMinimalCmd("mytool", "1.2.3", "Does something useful")
			},
			wantFragments: []string{
				"# mytool CLI Documentation",
				"**Version:** 1.2.3",
				"**Description:** Does something useful",
			},
		},
		{
			name: "usage section is present",
			buildCmd: func() *cobra.Command {
				return newMinimalCmd("jiff", "0.0.1", "JSON diff tool")
			},
			wantFragments: []string{
				"## Usage",
				"jiff [command] [flags] [arguments]",
			},
		},
		{
			name: "global flags section lists flag names",
			buildCmd: func() *cobra.Command {
				cmd := newMinimalCmd("jiff", "0.0.1", "JSON diff tool")
				cmd.PersistentFlags().BoolP("dry-run", "d", false, "Perform a dry run")
				cmd.PersistentFlags().StringP("output", "o", "", "Output format")
				return cmd
			},
			wantFragments: []string{
				"## Global Flags",
				"--dry-run",
				"--output",
				"-d",
				"-o",
			},
		},
		{
			name: "commands section is present",
			buildCmd: func() *cobra.Command {
				root := newMinimalCmd("jiff", "0.0.1", "JSON diff tool")
				sub := &cobra.Command{
					Use:   "compare",
					Short: "Compare two files",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				root.AddCommand(sub)
				return root
			},
			wantFragments: []string{
				"## Commands",
				"### compare",
				"Compare two files",
			},
		},
		{
			name: "no global flags section when there are no persistent flags",
			buildCmd: func() *cobra.Command {
				return newMinimalCmd("bare", "0.1.0", "A bare command")
			},
			wantFragments: []string{
				"## Commands",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := doc.GenerateMarkdown(tc.buildCmd())

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			for _, fragment := range tc.wantFragments {
				assert.Contains(t, output, fragment,
					"expected GenerateMarkdown output to contain %q", fragment)
			}
		})
	}
}
