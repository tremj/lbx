package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lbx",
	Short: "Central command for all functionalities within the load balancer",
}

func Execute() error {
	return rootCmd.Execute()
}
