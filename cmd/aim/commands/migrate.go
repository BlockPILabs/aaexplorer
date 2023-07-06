package commands

import (
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Automatic Migration",
	RunE: func(cmd *cobra.Command, args []string) error {

		// db start
		err := entity.Start(config)
		if err != nil {
			return err
		}
		db, err := entity.Client()
		if err != nil {
			return err
		}
		err = db.Schema.Create(cmd.Context())
		if err != nil {
			return err
		}
		// Run forever.
		select {}
	},
}
