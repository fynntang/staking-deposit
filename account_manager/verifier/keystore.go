package verifier

import (
	"encoding/json"
	"os"
)

// Keystore
//
//	Implement an EIP 2335-compliant keystore. A keystore is a JSON file that
//	stores an encrypted version of a private key under a user-supplied password.
//
//	Ref: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2335.md
type Keystore struct {
	Description string                 `json:"description"`
	Version     uint                   `json:"version"`
	UUID        string                 `json:"uuid"`
	Path        string                 `json:"path"`
	Pubkey      string                 `json:"pubkey"`
	Crypto      map[string]interface{} `json:"crypto"`
}

func LoadFileToKeystore(path string) (*Keystore, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	keystore := &Keystore{}
	if err := json.NewDecoder(jsonFile).Decode(keystore); err != nil {
		return nil, err
	}

	return keystore, nil
}

func SaveSigningKeystoreToFile(data *Keystore, path string) error {
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
