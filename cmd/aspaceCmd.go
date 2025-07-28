package cmd

import (
	"github.com/nyudlts/rename-workflow-tool/lib"
	"github.com/spf13/cobra"
)

func init() {
	aspaceCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(aspaceCmd)
}

var aspaceCmd = &cobra.Command{
	Use: "aspace",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("ASpace command executed")
	},
}

var checkCmd = &cobra.Command{
	Use: "check",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("ASpace check command executed")
		if err := lib.CheckAspace(); err != nil {
			panic(err)
		}
	},
}
