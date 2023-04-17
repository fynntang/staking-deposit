package key_derivation

import (
	"fmt"
	mnemonic "github.com/fynntang/staking-deposit/crypto/mnemonic"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestMnemonicAndPathToKey(t *testing.T) {
	type args struct {
		mnemonic string
		path     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Nil",
			args: args{
				mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
				path:     "m/0",
				password: "TREZOR",
			},
			want: _bigInt("20397789859736650942317412262472558107875392172444076792671091975210932703118"),
			wantErr: func(t assert.TestingT, err error, v ...interface{}) bool {
				return assert.NoError(t, err, v...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := mnemonic.NewSeed(tt.args.mnemonic, tt.args.password)
			got, err := MnemonicAndPathToKey(seed, tt.args.path)
			if !tt.wantErr(t, err, fmt.Sprintf("MnemonicAndPathToKey(%v, %v, %v)", tt.args.mnemonic, tt.args.path, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "MnemonicAndPathToKey(%v, %v, %v)", tt.args.mnemonic, tt.args.path, tt.args.password)
		})
	}
}
