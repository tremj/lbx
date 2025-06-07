package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return api.StartServer()
	},
}
