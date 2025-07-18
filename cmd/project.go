package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectInitCmd)
	rootCmd.AddCommand(projectCmd)
}

var projectCmd = &cobra.Command{
	Use: "project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project command executed")
	},
}

var projectArchiveCmd = &cobra.Command{
	Use: "archive",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for initializing a project goes here
		cmd.Println("project archive command executed")
	},
}

var projectInitCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for initializing a project goes here
		cmd.Println("project init command executed")
	},
}
