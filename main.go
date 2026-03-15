package main

import (
	"bytes"
	"errors"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jgfranco17/jiff-cli/cli/commandline"
	"github.com/jgfranco17/jiff-cli/cli/metadata"
	"github.com/jgfranco17/jiff-cli/internal/errorhandling"

	_ "embed" // Required for the //go:embed directive
)

//go:embed specs.json
var embeddedMetadata []byte

func main() {
	embeddedMetadataReader := bytes.NewReader(embeddedMetadata)
	metadata, err := metadata.Load(embeddedMetadataReader)
	if err != nil {
		os.Exit(handleError(err))
	}

	command := commandline.New(metadata.Name, metadata.Description, metadata.Version)
	if err := command.Execute(); err != nil {
		os.Exit(handleError(err))
	}
}

func handleError(err error) int {
	var exitErr *errorhandling.ExitError
	if errors.As(err, &exitErr) {
		logrus.Error(exitErr.String())
		return exitErr.ExitCode
	} else {
		logrus.Error(err.Error())
	}
	return errorhandling.ExitCodeGenericError
}
