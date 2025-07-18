package cmd

import "github.com/spf13/cobra"

func init() {
	aipCmd.AddCommand(aipValidateCmd)
	rootCmd.AddCommand(aipCmd)
}

var aipCmd = &cobra.Command{
	Use: "aip",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip command executed")
	},
}

var aipValidateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip validate command executed")
	},
}

var aipTransferCmd = &cobra.Command{
	Use: "transfer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip transfer command executed")
	},
}
