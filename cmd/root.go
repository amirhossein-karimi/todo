package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo application CLI",
	Long:  "A simple Todo application.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Todo CLI")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
