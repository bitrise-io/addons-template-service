package cmd

import (
	"errors"
	"os"

	"github.com/bitrise-team/addons-template-service/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration task",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(connStr) < 1 {
			return errors.New("--connection-string, -c must be set")
		}
		if err := models.Migrate(connStr); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	},
}

var connStr string

func init() {
	migrateCmd.Flags().StringVarP(&connStr, "connection-string", "c", os.Getenv("DATABASE_URL"), "DB connection string")
	rootCmd.AddCommand(migrateCmd)
}
