package cmd

import "github.com/spf13/cobra"

func init() {
	sourceCmd.AddCommand(sourceTransferCmd)
	sourceCmd.AddCommand(sourceSizeCmd)
	rootCmd.AddCommand(sourceCmd)
}

var sourceCmd = &cobra.Command{
	Use: "source",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("source command executed")
	},
}

var sourceTransferCmd = &cobra.Command{
	Use: "transfer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("source transfer command executed")
	},
}

var sourceSizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("source size command executed")
	},
}
