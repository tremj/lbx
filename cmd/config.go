package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/core"
	"github.com/tremj/lbx/internal/utils"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save load balancer configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fileContent, name, err := utils.RetrieveSaveCmdInfo(cmd)
		if err != nil {
			return err
		}
		return core.SaveConfig(cmd.Context(), name, fileContent)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete load balancer configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := utils.RetrieveDeleteCmdInfo(cmd)
		if err != nil {
			return err
		}
		return core.DeleteConfig(cmd.Context(), name)
	},
}

func init() {
	saveCmd.Flags().StringP("filepath", "f", "", "Path to YAML config file")
	saveCmd.Flags().StringP("name", "n", "", "Name of the load balancer configuration")
	saveCmd.MarkFlagRequired("filepath")
	saveCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(saveCmd)

	deleteCmd.Flags().StringP("name", "n", "", "Name of the load balancer configuration")
	deleteCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(deleteCmd)
}
