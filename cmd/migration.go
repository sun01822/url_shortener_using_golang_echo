package cmd

import (
	"url_shortener/internal/database"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrationCmd)
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.New()
		migrate(db)
	},
}

func migrate(db database.Service) {
	err := db.Migrate()
	panicError(err)
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
