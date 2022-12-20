package helper

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func SetCommandArgs(t *testing.T, args []string) {
	t.Helper()

	t.Cleanup(func() {
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	})

	commandArgs := os.Args

	os.Args = append([]string{os.Args[0]}, args...)

	t.Cleanup(func() {
		os.Args = commandArgs
	})
}
