package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/parser"
	redis "github.com/tremj/lbx/internal/redisClient"
	"github.com/tremj/lbx/internal/store"
)

func main() {
	redisClient := redis.NewClient("localhost:6379")

	rootCmd := &cobra.Command{
		Use:   "lbx",
		Short: "Central command for all functionalities within the load balancer",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(cmd.Context(), "redisClient", redisClient)
			cmd.SetContext(ctx)
		},
	}

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse YAML configuration files",
		Run:   parser.Parse,
	}
	parseCmd.Flags().StringP("filepath", "f", "", "Path to YAML config file")
	rootCmd.AddCommand(parseCmd)

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save load balancer configuration",
		Run:   store.Save,
	}
	saveCmd.Flags().StringP("filepath", "f", "", "Path to YAML config file")
	saveCmd.Flags().StringP("name", "n", "", "Name of the load balancer configuration")
	rootCmd.AddCommand(saveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
