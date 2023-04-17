package utils

import "github.com/fynntang/staking-deposit/constants"

func GetMnemonicLanguageFromName(name string) constants.MnemonicLanguage {

	for language, s := range constants.MnemonicLanguagesName {
		if s == name {
			return language
		}
	}

	return ""
}
