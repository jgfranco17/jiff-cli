package metadata

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jgfranco17/jiff-cli/internal/errorhandling"
)

// ProjectInfo represents the developer metadata for a project.
// This allows for embedding important information about the
// project directly within the binary, and can be used for display
// in help commands.
type ProjectInfo struct {
	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Repository  string `json:"repository"`
}

// Load reads the metadata from the provided reader instance
// and returns a ProjectInfo struct.
func Load(r io.Reader) (ProjectInfo, error) {
	var data ProjectInfo
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return ProjectInfo{}, &errorhandling.ExitError{
			Err:      fmt.Errorf("failed to parse metadata: %w", err),
			ExitCode: errorhandling.ExitCodeConfigurationError,
			Solution: "Ensure the JSON is valid and contains all required fields.",
		}
	}
	return data, nil
}
