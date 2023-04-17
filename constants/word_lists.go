package constants

import "github.com/tyler-smith/go-bip39/wordlists"

type WordLists []string

var WordListsMap = map[MnemonicLanguage]WordLists{
	MnemonicLanguageChineseSimplified:  wordlists.ChineseSimplified,
	MnemonicLanguageChineseTraditional: wordlists.ChineseTraditional,
	MnemonicLanguageCzech:              wordlists.Czech,
	MnemonicLanguageEnglish:            wordlists.English,
	MnemonicLanguageFrench:             wordlists.French,
	MnemonicLanguageItalian:            wordlists.Italian,
	MnemonicLanguageJapanese:           wordlists.Japanese,
	MnemonicLanguageKorean:             wordlists.Korean,
	MnemonicLanguageSpanish:            wordlists.Spanish,
}
