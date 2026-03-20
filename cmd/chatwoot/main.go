package main

import (
	"os"

	"github.com/chatwoot/chatwoot-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
