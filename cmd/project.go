package cmd

import (
	"fmt"
	"os"

	lib "github.com/nyudlts/rename-workflow-tool/lib"
	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(projectArchiveCmd)
	projectInitCmd.Flags().StringVarP(&collectionCode, "collection-code", "c", "", "the collection code to use for rename project")
	projectInitCmd.Flags().StringVarP(&sourceLocation, "source-location", "s", "", "the source location for the collection")

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
		if err := lib.InitializeProject(collectionCode, sourceLocation); err != nil {
			cmd.Println("Error initializing project:", err)
			os.Exit(1)
		}
	},
}
