package constants

import "encoding/hex"

type NetworkName string

const (
	ChainMAINNET  NetworkName = "mainnet"
	ChainGOERLI   NetworkName = "goerli"
	ChainPRATER   NetworkName = "prater"
	ChainSEPOLIA  NetworkName = "sepolia"
	ChainZHEJIANG NetworkName = "zhejiang"
)

type Chain struct {
	NetworkName           NetworkName
	GenesisForkVersion    [4]byte
	GenesisValidatorsRoot [32]byte
}

var NetworkNames = []NetworkName{
	ChainMAINNET,
	ChainGOERLI,
	ChainPRATER,
	ChainSEPOLIA,
	ChainZHEJIANG,
}

var chains = map[NetworkName]Chain{
	ChainMAINNET:  newChain(ChainMAINNET, "00000000", "4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95"),
	ChainGOERLI:   newChain(ChainGOERLI, "00001020", "043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb"),
	ChainSEPOLIA:  newChain(ChainSEPOLIA, "90000069", "d8ea171f3c94aea21ebc42a1ed61052acf3f9209c00e4efbaaddac09ed9b8078"),
	ChainZHEJIANG: newChain(ChainZHEJIANG, "00000069", "53a92d8f2bb1d85f62d16a156e6ebcd1bcaba652d0900b2c2f387826f3481f6f"),
}

func GetChain(chain NetworkName) Chain {
	return chains[chain]
}

func newChain(networkName NetworkName, genesisForkVersion, genesisValidatorsRoot string) Chain {
	gfv, _ := hex.DecodeString(genesisForkVersion)
	gvr, _ := hex.DecodeString(genesisValidatorsRoot)

	return Chain{NetworkName: networkName, GenesisForkVersion: [4]byte(gfv), GenesisValidatorsRoot: [32]byte(gvr)}
}

func (n NetworkName) ToString() string {
	return string(n)
}
