package cmd

import (
	"os"

	lib "github.com/nyudlts/rename-workflow-tool/lib"
	"github.com/spf13/cobra"
)

func init() {
	sipCmd.AddCommand(sipSizeCmd)
	sipCmd.AddCommand(sipValidateCmd)
	rootCmd.AddCommand(sipCmd)
}

var sipCmd = &cobra.Command{
	Use: "sip",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("sip command executed")
	},
}

var sipValidateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("sip validate command executed")
		if err := lib.ValidateSip(); err != nil {
			cmd.Printf("  [ERROR] %s", err.Error())
			os.Exit(1)
		}
		cmd.Println("SIP validation completed successfully")
	},
}

var sipSizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("sip size command executed")
		if err := lib.GetSipSize(); err != nil {
			cmd.Println("Error getting SIP size")
			os.Exit(1)
		}
	},
}
