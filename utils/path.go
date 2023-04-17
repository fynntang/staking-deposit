package utils

import (
	"github.com/fynntang/staking-deposit/constants"
	"log"
	"os"
	"path"
)

func ValidatorKeysFolderPath(p ...string) string {
	pwd, _ := os.Getwd()

	folderPath := path.Join(pwd, constants.DefaultValidatorKeysFolderName)
	if len(p) > 0 {
		folderPath = path.Join(folderPath, p[0])
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			log.Fatalf("Failed to create folder %s: %v", folderPath, err)
		}
	}

	return folderPath
}
