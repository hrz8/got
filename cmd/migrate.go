package cmd

import (
	"github.com/hrz8/got/internal/storage"
	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate the database",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(NewMigrateUpCommand())
	cmd.AddCommand(NewMigrateDownCommand())

	return cmd
}

func NewMigrateUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "migrate the database to the most recent version available",
		RunE:  migrateUp,
		Args:  cobra.NoArgs,
	}

	return cmd
}

func NewMigrateDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "roll back the migration for the current database",
		RunE:  migrateDown,
		Args:  cobra.NoArgs,
	}

	return cmd
}

func migrateUp(cmd *cobra.Command, args []string) error {
	return storage.MigrateUp()
}

func migrateDown(cmd *cobra.Command, args []string) error {
	return nil
}
