package mnemonic

import "github.com/tyler-smith/go-bip39"

func NewMnemonic(wordlist []string) (string, error) {
	bip39.SetWordList(wordlist)

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func NewSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}
