package bls

import (
	_bls "github.com/herumi/bls-eth-go-binary/bls"
)

// InitBLS initialises the BLS library with the appropriate curve and parameters for Ethereum 2.
func InitBLS() error {
	if err := _bls.Init(_bls.BLS12_381); err != nil {
		return err
	}
	return _bls.SetETHmode(_bls.EthModeDraft07)
}
