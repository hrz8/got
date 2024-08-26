package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "got",
		Short: "golang template",
		Long:  "simple golang template",
	}
	return cmd
}
