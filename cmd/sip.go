package cmd

import "github.com/spf13/cobra"

func init() {
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
	},
}
