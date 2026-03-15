package commandline

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
)

type Metadata struct {
	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Repository  string `json:"repository"`
}

// LoadMetadata reads the CLI metadata from the provided reader instance
// and returns a Metadata struct.
func LoadMetadata(r io.Reader) (Metadata, error) {
	var cfg Metadata
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return Metadata{}, &errorhandling.ExitError{
			Err:      fmt.Errorf("failed to parse metadata: %w", err),
			ExitCode: errorhandling.ExitCodeConfigurationError,
			Solution: "Ensure the JSON is valid and contains all required fields.",
		}
	}

	return cfg, nil
}
