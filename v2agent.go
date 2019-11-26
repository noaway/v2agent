package main

import (
	"github.com/noaway/v2-agent/cmd"
	"github.com/noaway/v2-agent/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "v2agent", Version: version.Version()}
	cmd.Commands(root)
	root.Execute()
}
