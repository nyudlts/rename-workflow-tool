package cmd

import (
	lib "github.com/nyudlts/rename-workflow-tool/lib"
	"github.com/spf13/cobra"
)

func init() {
	aipCmd.AddCommand(aipPrepCmd)
	aipCmd.AddCommand(aipBagCmd)
	aipCmd.AddCommand(aipUpdateCmd)
	aipCmd.AddCommand(aipValidateCmd)
	aipCmd.AddCommand(aipTransferCmd)
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
		if err := lib.ValidateAIPs(); err != nil {
			panic(err)
		}
		cmd.Println("all aip packages valid")
	},
}

var aipTransferCmd = &cobra.Command{
	Use: "transfer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip transfer command executed")
		if err := lib.TransferAIPs(); err != nil {
			panic(err)
		}
	},
}

var aipPrepCmd = &cobra.Command{
	Use: "prep",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip prep command executed")
		if err := lib.PrepAIPs(); err != nil {
			panic(err)
		}
	},
}

var aipBagCmd = &cobra.Command{
	Use: "bag",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip bag command executed")
		if err := lib.BagAIPs(); err != nil {
			panic(err)
		}
	},
}

var aipUpdateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("aip update command executed")
		if err := lib.UpdateAIP(); err != nil {
			panic(err)
		}
	},
}
