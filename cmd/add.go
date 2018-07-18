package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/esell/kvpm/util"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update an entry",

	Run: func(cmd *cobra.Command, args []string) {

		vaultName := os.Getenv("KVAULT")

		basicClient, err := util.GetBasicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var secParams keyvault.SecretSetParameters
		secParams.Value = &args[1]
		newBundle, err := basicClient.SetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", args[0], secParams)
		if err != nil {
			fmt.Printf("unable to add/update secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("added/updated: " + *newBundle.ID)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
