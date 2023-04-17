package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "deposit-cli",
	Short: "DepositCli is A tool for creating EIP-2335 format BLS12-381 keystores and a corresponding deposit_data*.json file for Ethereum Staking Launchpad.",
	Long:  `DepositCli is A tool for creating EIP-2335 format BLS12-381 keystores and a corresponding deposit_data*.json file for Ethereum Staking Launchpad.`,
}

func Execute() error {
	return rootCmd.Execute()
}
