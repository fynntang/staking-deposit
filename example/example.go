package main

import (
	"github.com/fynntang/staking-deposit/constants"
	"github.com/fynntang/staking-deposit/credential"
	_mnemonic "github.com/fynntang/staking-deposit/crypto/mnemonic"
	"github.com/fynntang/staking-deposit/utils"
	"github.com/tyler-smith/go-bip39/wordlists"
	"log"
)

func main() {

	var DefaultMnemonicPassword = "MyMnemonicPassword"

	mnemonic, err := _mnemonic.NewMnemonic(wordlists.ChineseSimplified)
	if err != nil {
		log.Fatalf("NewMnemonic() error = %v", err)
	}
	log.Printf("Mnemonic = %s\n", mnemonic)

	numValidators := 1
	log.Printf("NumValidators = %d\n", numValidators)

	amounts := make([]int64, 0)
	for i := 0; i < numValidators; i++ {
		amounts = append(amounts, int64(constants.MaxDepositAmount*numValidators))
	}
	log.Printf("Amounts = %d\n", amounts)

	folder := utils.ValidatorKeysFolderPath()
	log.Printf("Folder = %s\n", folder)

	chain := constants.GetChain(constants.ChainGOERLI)

	executionAddress := "0xc91e61C36b47B13B24C58535ff543Aa218ef61C4"
	log.Printf("ExecutionAddress = %s\n", executionAddress)

	credentials, err := credential.NewCredentialsList(
		mnemonic,
		DefaultMnemonicPassword,
		numValidators,
		amounts,
		chain,
		0,
		executionAddress)
	if err != nil {
		log.Fatalf("NewCredentialListFromMnemonic() error = %v", err)
	}

	keystoreFileFolders, err := credentials.ExportKeystores(DefaultMnemonicPassword, folder)
	if err != nil {
		log.Fatalf("ExportKeystores error %v", err)
	}
	log.Printf("KeystoreFileFolders = %v\n", keystoreFileFolders)

	depositsFile, err := credentials.ExportDepositData(folder)
	if err != nil {
		log.Fatalf("ExportDepositDataJson error %v", err)
	}
	log.Printf("DepositsFile = %v\n", depositsFile)

	if err := credentials.VerifyKeystores(keystoreFileFolders, DefaultMnemonicPassword); err != nil {
		log.Fatalf("VerifyKeystores error %v", err)
	}
	log.Printf("VerifyKeystores = success\n")

	if err := credentials.VerifyDepositData(depositsFile); err != nil {
		log.Fatalf("VerifyDepositData error %v", err)
	}
	log.Printf("VerifyDepositData = success\n")

}
