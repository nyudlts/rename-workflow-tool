package cmd

import (
	"os"

	lib "github.com/nyudlts/rename-workflow-tool/lib"
	"github.com/spf13/cobra"
)

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
		if err := lib.TransferSourcePackage(); err != nil {
			cmd.Println("Error transferring source package:", err)
			os.Exit(1)
		}
	},
}

var sourceSizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("source size command executed")
		if err := lib.PrintSourceSize(); err != nil {
			cmd.Println("Error getting source size:", err)
			os.Exit(1)
		}
	},
}
