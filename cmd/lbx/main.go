package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/parser"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "lbx",
		Short: "Central command for all functionalities within the load balancer",
	}

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse YAML configuration files",
		Run:   parser.Parse,
	}
	parseCmd.Flags().StringP("filepath", "f", "config.yaml", "Path to YAML config file")
	rootCmd.AddCommand(parseCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
