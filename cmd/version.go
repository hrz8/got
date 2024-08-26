package cmd

import (
	"fmt"

	"github.com/hrz8/got/config"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "prints the application version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s\n", config.Version)
			return nil
		},
	}
	return cmd
}
