package cli

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fynntang/staking-deposit/constants"
	"github.com/fynntang/staking-deposit/credential"
	_mnemonic "github.com/fynntang/staking-deposit/crypto/mnemonic"
	"github.com/fynntang/staking-deposit/utils"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var newMnemonicCmd = &cobra.Command{
	Use:   "new_mnemonic",
	Short: "Generate a new mnemonic and keys",
	RunE:  runNewMnemonicCmd,
}

func init() {
	rootCmd.AddCommand(newMnemonicCmd)
}

type Prompt survey.Prompt

var (
	numValidatorsPrompt                Prompt = &survey.Input{Message: "How many validators would you like to create?", Default: "1"}
	chainPrompt                        Prompt = &survey.Select{Message: "Please choose your chain:", Options: chainNetworkNames(), Default: constants.ChainMAINNET.ToString()}
	cliLanguagePrompt                  Prompt = &survey.Select{Message: "Please choose your language:", Options: languageOptionsKeys(), Default: constants.MnemonicLanguageEnglish.ToName()}
	eth1WithdrawalAddressPrompt        Prompt = &survey.Input{Message: "Please enter your execution address:", Default: ""}
	confirmEth1WithdrawalAddressPrompt Prompt = &survey.Input{Message: "Please confirm your execution address:", Default: ""}
	mnemonicLanguagePrompt             Prompt = &survey.Select{Message: "Please choose the language of the mnemonic word list:", Options: languageOptionsKeys(), Default: constants.MnemonicLanguageEnglish.ToName()}
	keystorePasswordPrompt             Prompt = &survey.Password{Message: "Please create a password that secures your validator keystore(s).\n  You will need to re-enter this to decrypt them when you setup your Ethereum validators:"}
	confirmKeystorePasswordPrompt      Prompt = &survey.Password{Message: "Please confirm your password:"}
	confirmSaveMnemonicPrompt          Prompt = &survey.Confirm{Message: "\033[31mPlease confirm you have written down your mnemonic.\033[0;39m", Default: false}
	confirmMnemonicPrompt              Prompt = &survey.Input{Message: "Please enter your mnemonic:"}
)

func runNewMnemonicCmd(cmd *cobra.Command, args []string) error {
	var numValidators uint64
	if err := survey.AskOne(numValidatorsPrompt, &numValidators, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var chain string
	if err := survey.AskOne(chainPrompt, &chain, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	//var cliLanguage string
	//if err := survey.AskOne(cliLanguagePrompt, &cliLanguage, survey.WithValidator(survey.Required)); err != nil {
	//	return err
	//}

	var eth1WithdrawalAddress string
	if err := survey.AskOne(eth1WithdrawalAddressPrompt, &eth1WithdrawalAddress, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if !common.IsHexAddress(eth1WithdrawalAddress) {
		return errors.New("invalid eth1 withdrawal address")
	}

	var confirmEth1WithdrawalAddress string
	if err := survey.AskOne(confirmEth1WithdrawalAddressPrompt, &confirmEth1WithdrawalAddress, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	if strings.Compare(eth1WithdrawalAddress, confirmEth1WithdrawalAddress) != 0 {
		return errors.New("execution address do not match")
	}

	var mnemonicLanguage string
	if err := survey.AskOne(mnemonicLanguagePrompt, &mnemonicLanguage, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var keystorePassword string
	if err := survey.AskOne(keystorePasswordPrompt, &keystorePassword, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var confirmKeystorePassword string
	if err := survey.AskOne(confirmKeystorePasswordPrompt, &confirmKeystorePassword, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	if strings.Compare(keystorePassword, confirmKeystorePassword) != 0 {
		return errors.New("keystore Password do not match")
	}

	mnemonic, err := _mnemonic.NewMnemonic(utils.GetMnemonicLanguageFromName(mnemonicLanguage).ToWordList())
	if err != nil {
		return err
	}

	fmt.Printf("This is your mnemonic (seed phrase). "+
		"Write it down and store it safely. "+
		"It is the ONLY way to retrieve your deposit.\n\n\033[32m%s\033[0;39m\n\n", mnemonic)

	var confirmSaveMnemonic bool
	if err := survey.AskOne(confirmSaveMnemonicPrompt, &confirmSaveMnemonic, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if !confirmSaveMnemonic {
		return errors.New("you must save your mnemonic to continue")
	}

	fmt.Print("\033[H\033[2J")

	var confirmMnemonic string
	if err := survey.AskOne(confirmMnemonicPrompt, &confirmMnemonic, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if strings.Compare(mnemonic, confirmMnemonic) != 0 {
		return errors.New("mnemonic do not match")
	}

	printLife()
	fmt.Print("Creating your keys.\n\n")

	amounts := make([]int64, 0)
	for i := 0; i < int(numValidators); i++ {
		amounts = append(amounts, int64(constants.MaxDepositAmount*numValidators))
	}

	folder := utils.ValidatorKeysFolderPath()

	credentials, err := credential.NewCredentialsList(
		mnemonic,
		keystorePassword,
		int(numValidators),
		amounts,
		constants.GetChain(constants.NetworkName(chain)),
		0,
		eth1WithdrawalAddress)
	if err != nil {
		return err
	}

	var depositsFile string
	var keystoreFiles []string

	depositsFile, err = credentials.ExportDepositData(folder)
	if err != nil {
		return err
	}

	keystoreFiles, err = credentials.ExportKeystores(keystorePassword, folder)
	if err != nil {
		return err
	}

	if err := credentials.VerifyDepositData(depositsFile); err != nil {
		return err
	}
	fmt.Printf("Verifying your DepositData: \033[32m%s\033[0;39m\n", "Success!")

	if err := credentials.VerifyKeystores(keystoreFiles, keystorePassword); err != nil {
		return err
	}
	fmt.Printf("Verifying your Keystore(s): \033[32m%s\033[0;39m\n", "Success!")

	fmt.Printf("Your keys can be found at: %s\n\n", folder)
	return nil
}

func chainNetworkNames() []string {
	var names []string
	for _, network := range constants.NetworkNames {
		names = append(names, network.ToString())
	}
	return names
}

func languageOptionsKeys() []string {
	var keys []string
	for _, language := range constants.MnemonicLanguagesName {
		keys = append(keys, language)
	}
	return keys
}

func printLife() {
	l := utils.NewLife(40, 15)
	for i := 0; i < 300; i++ {
		l.Step()
		fmt.Print("\x0c", l)
		time.Sleep(time.Millisecond * 50)
	}
}
