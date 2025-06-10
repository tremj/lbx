package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/api"
	"github.com/tremj/lbx/internal/storage"
)

var serveCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		storage.InitRedis()
		return api.StartServer()
	},
}
