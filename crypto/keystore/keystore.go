package keystore

import "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

func NewScryptKeystoreEncryptor() *keystorev4.Encryptor {
	return keystorev4.New(keystorev4.WithCipher("scrypt"))
}
