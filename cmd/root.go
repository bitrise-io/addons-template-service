package cmd

import (
	"os"

	"github.com/bitrise-team/addons-template-service/models"
	"github.com/bitrise-team/addons-template-service/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "addons-template-service",
	Short: "A brief description of your application",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := os.Getenv("PORT")
		dbConn := os.Getenv("DATABASE_URL")

		if err := models.SetupDB(dbConn); err != nil {
			return err
		}

		s := server.NewServer(server.Settings{Port: port})
		logrus.Infof("Listening on port %s...", port)

		return s.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
