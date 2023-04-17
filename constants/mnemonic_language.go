package constants

type MnemonicLanguage string

const (
	MnemonicLanguageChineseSimplified  MnemonicLanguage = "chinese_simplified"
	MnemonicLanguageChineseTraditional MnemonicLanguage = "chinese_traditional"
	MnemonicLanguageCzech              MnemonicLanguage = "czech"
	MnemonicLanguageEnglish            MnemonicLanguage = "english"
	MnemonicLanguageFrench             MnemonicLanguage = "french"
	MnemonicLanguageItalian            MnemonicLanguage = "italian"
	MnemonicLanguageJapanese           MnemonicLanguage = "japanese"
	MnemonicLanguageKorean             MnemonicLanguage = "korean"
	MnemonicLanguageSpanish            MnemonicLanguage = "spanish"
)

var MnemonicLanguagesName = map[MnemonicLanguage]string{
	MnemonicLanguageChineseSimplified:  "🇨🇳 Chinese (Simplified)",
	MnemonicLanguageChineseTraditional: "🇨🇳 Chinese (Traditional)",
	MnemonicLanguageCzech:              "🇨🇿 Czech",
	MnemonicLanguageEnglish:            "🇺🇸 English",
	MnemonicLanguageFrench:             "🇫🇷 French",
	MnemonicLanguageItalian:            "🇮🇹 Italian",
	MnemonicLanguageJapanese:           "🇯🇵 Japanese",
	MnemonicLanguageKorean:             "🇰🇷 Korean",
	MnemonicLanguageSpanish:            "🇪🇸 Spanish",
}

func (l MnemonicLanguage) ToString() string {
	return string(l)
}
func (l MnemonicLanguage) ToName() string {
	return MnemonicLanguagesName[l]
}
func (l MnemonicLanguage) ToWordList() WordLists {
	return WordListsMap[l]
}
