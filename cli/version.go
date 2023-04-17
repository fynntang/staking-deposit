package cli

import (
	"fmt"
	"github.com/fynntang/staking-deposit/constants"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of deposit-cli",
	Long:  `All software has versions. This is deposit-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("DepositCliVersion: %s\n", constants.DepositCliVersion)
	},
}
