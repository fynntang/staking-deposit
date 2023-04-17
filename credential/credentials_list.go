package credential

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fynntang/staking-deposit/account_manager/verifier"
	"github.com/fynntang/staking-deposit/constants"
	"math/big"
	"path"
	"strings"
	"time"
)

type CredentialsList struct {
	Credentials []*Credential
}

func NewCredentialsList(
	mnemonic, mnemonicPassword string,
	numKeys int, amounts []int64,
	chain constants.Chain,
	startIndex int,
	eth1WithdrawalAddress string,
) (*CredentialsList, error) {
	if len(amounts) != numKeys {
		return nil, fmt.Errorf("the number of keys (%d) doesn't equal to the corresponding deposit amounts (%d)", numKeys, len(amounts))
	}

	if !common.IsHexAddress(eth1WithdrawalAddress) {
		return nil, fmt.Errorf("invalid eth1 withdrawal address: %s", eth1WithdrawalAddress)
	}

	var creds []*Credential
	for i := startIndex; i < (startIndex + numKeys); i++ {
		cred, err := NewCredential(
			mnemonic,
			mnemonicPassword,
			i,
			big.NewInt(amounts[i-startIndex]),
			chain,
			eth1WithdrawalAddress)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}

	return &CredentialsList{creds}, nil
}

func (cl *CredentialsList) ExportDepositData(folder string) (string, error) {
	if len(cl.Credentials) == 0 {
		return "", fmt.Errorf("no credentials to export")
	}
	depositDatas := make([]*verifier.DepositData, 0)
	for _, credential := range cl.Credentials {
		depositData, err := credential.SigningDepositData()
		if err != nil {
			return "", err
		}
		depositDatas = append(depositDatas, depositData)
	}
	fileFolder := path.Join(folder, fmt.Sprintf("deposit_data-%d.json", time.Now().Unix()))

	if err := verifier.SaveDepositDataToFile(depositDatas, fileFolder); err != nil {
		return "", err
	}

	return fileFolder, nil
}
func (cl *CredentialsList) VerifyDepositData(fileFolder string) error {
	depositDatas, err := verifier.LoadFileToDepositData(fileFolder)
	if err != nil {
		return err
	}
	if len(depositDatas) != len(cl.Credentials) {
		return fmt.Errorf("the number of deposit data (%d) doesn't equal to the number of credentials (%d)", len(depositDatas), len(cl.Credentials))
	}

	for i, depositData := range depositDatas {
		credential := cl.Credentials[i]
		if err := credential.VerifyDepositData(depositData); err != nil {
			return err
		}
	}

	return nil
}

func (cl *CredentialsList) ExportKeystores(password, folder string) ([]string, error) {
	if len(cl.Credentials) == 0 {
		return nil, fmt.Errorf("no credentials to export")
	}

	keystoreFiles := make([]string, 0)
	for _, creds := range cl.Credentials {
		data, err := creds.SigningKeystore(password)
		if err != nil {
			return nil, err
		}
		file := path.Join(folder, fmt.Sprintf("keystore-%s-%d.json", strings.Replace(creds.SigningKeyPath, "/", "_", -1), time.Now().Unix()))
		if err := verifier.SaveSigningKeystoreToFile(data, file); err != nil {
			return nil, err
		}
		keystoreFiles = append(keystoreFiles, file)
	}

	return keystoreFiles, nil
}

func (cl *CredentialsList) VerifyKeystores(keystoreFiles []string, password string) error {
	for i, keystoreFile := range keystoreFiles {
		keystore, err := verifier.LoadFileToKeystore(keystoreFile)
		if err != nil {
			return err
		}
		if err := cl.Credentials[i].VerifyKeystore(keystore.Crypto, password); err != nil {
			return err
		}
	}

	return nil
}
