package credential

import (
	"encoding/json"
	"github.com/fynntang/staking-deposit/constants"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func getCredential() *Credential {
	credential, _ := NewCredential(
		"迎 插 督 弃 谁 备 眉 章 炒 案 土 嫩 此 拖 腰 参 客 渗 极 壁 亦 领 靠 败",
		"MyMnemonicPassword",
		0,
		big.NewInt(constants.MaxDepositAmount),
		constants.GetChain(constants.ChainGOERLI),
		"0xc91e61C36b47B13B24C58535ff543Aa218ef61C4",
	)
	return credential
}

func TestCredential_ToDepositData(t *testing.T) {
	data, err := getCredential().SigningDepositData()
	assert.NoError(t, err)

	marshal, err := json.Marshal(data)
	assert.NoError(t, err)

	t.Log(string(marshal))
}

func TestBLS12318(t *testing.T) {
	sk := getCredential().SigningSk
	t.Logf("sk: %x", sk.Marshal())
	pk := sk.PublicKey()

	t.Logf("pk: %x", pk.Marshal())

	msg := []byte("hello world")

	t.Logf("msg: %s", msg)

	signature := sk.Sign(msg)

	assert.Equal(t, 96, len(signature.Marshal()))
	assert.True(t, signature.Verify(msg, pk))
}

func TestSigningKeystore(t *testing.T) {
	credential := getCredential()

	keystore, err := credential.SigningKeystore("MyMnemonicPassword")
	assert.NoError(t, err)

	marshal, err := json.Marshal(keystore)
	assert.NoError(t, err)

	t.Logf("keystore: %s", string(marshal))
}
