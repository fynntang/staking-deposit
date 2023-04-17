package mnemonic

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"strings"
	"testing"
)

func TestNewMnemonic(t *testing.T) {

	tests := []struct {
		name     string
		wordlist []string
		mnemonic string
		password string
		want     string
	}{
		{
			name:     "Test Case 1",
			wordlist: wordlists.ChineseSimplified,
			mnemonic: "迎 插 督 弃 谁 备 眉 章 炒 案 土 嫩 此 拖 腰 参 客 渗 极 壁 亦 领 靠 败",
			password: "MyMnemonicPassword",
			want:     "bca1c28513a463e31088e5640f714de21220bd56172b31c79cf71eab1297a46922fe8d98334c2fc5bdec0068fd8d51126ab81d8bcc94c3ce1c937674e03a5c4c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnemonic, err := NewMnemonic(tt.wordlist)
			assert.NoError(t, err)
			t.Logf("mnemonic: %s", mnemonic)
			assert.Equal(t, 24, len(strings.Split(mnemonic, " ")), "Mnemonic length is not equal")

			seed := bip39.NewSeed(tt.mnemonic, tt.password)
			assert.Equal(t, 64, len(seed), "Seed length is not equal")

			assert.Equal(t, tt.want, hex.EncodeToString(seed), "Seed is not equal")
		})
	}
}
