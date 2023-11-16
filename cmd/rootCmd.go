package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitadr",
	Short: "gitadr is a tool to help you manage your adrs the gitops way",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no command is given
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
