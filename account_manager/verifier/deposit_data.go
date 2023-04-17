package verifier

import (
	"encoding/json"
	"os"
)

type DepositData struct {
	PublicKey             string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                uint64 `json:"amount"`
	Signature             string `json:"signature"`
	DepositMessageRoot    string `json:"deposit_message_root"`
	DepositDataRoot       string `json:"deposit_data_root"`
	ForkVersion           string `json:"fork_version"`
	NetworkName           string `json:"network_name"`
	DepositCliVersion     string `json:"deposit_cli_version"`
}

func LoadFileToDepositData(path string) ([]*DepositData, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	depositData := make([]*DepositData, 0)
	if err := json.NewDecoder(jsonFile).Decode(&depositData); err != nil {
		return nil, err
	}

	return depositData, nil
}

func SaveDepositDataToFile(data []*DepositData, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return err
	}

	if os.PathSeparator == '/' {
		if err := os.Chmod(path, 0440); err != nil {
			return err
		}
	}

	return nil
}
