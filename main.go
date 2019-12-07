package main

import (
	_ "github.com/noaway/v2agent/agent"
	"github.com/noaway/v2agent/cmd"
	"github.com/noaway/v2agent/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "v2agent", Version: version.Version()}
	cmd.Commands(root)
	root.Execute()
}
