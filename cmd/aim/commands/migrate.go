package commands

import (
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Automatic Migration",
	RunE: func(cmd *cobra.Command, args []string) error {

		// db start
		err := entity.Start(logger.With("lib", "ent"), config)
		if err != nil {
			return err
		}
		db, err := entity.Client(cmd.Context())
		if err != nil {
			return err
		}
		err = db.Debug().Schema.Create(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
}
