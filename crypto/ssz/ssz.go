package ssz

import (
	"errors"
	"fmt"
	ssz "github.com/ferranbt/fastssz"
	"github.com/fynntang/staking-deposit/constants"
	consensus "github.com/umbracle/go-eth-consensus"
)

// ComputeForkDataRoot Return the appropriate ForkData root for a given deposit version.
func ComputeForkDataRoot(currentVersion [4]byte, genesisValidatorsRoot [32]byte) ([]byte, error) {
	fd := consensus.ForkData{
		CurrentVersion:        currentVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}

	if err := fd.HashTreeRootWith(ssz.DefaultHasherPool.Get()); err != nil {
		return nil, err
	}

	forkRoot, err := fd.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return forkRoot[:], nil
}

// ComputeDepositDomain Deposit-only `compute_domain`
func ComputeDepositDomain(forkVersion [4]byte) ([]byte, error) {
	domainType := constants.DomainDeposit

	forkDataRoot, err := ComputeDepositForkDataRoot(forkVersion)
	if err != nil {
		return nil, err
	}

	return append(domainType[:], forkDataRoot[:28]...), nil
}

// ComputeDepositForkDataRoot Return the appropriate ForkData root for a given deposit version.
func ComputeDepositForkDataRoot(currentVersion [4]byte) ([]byte, error) {

	genesisValidatorsRoot := constants.ZeroBytes32 // For deposit, it's fixed value

	return ComputeForkDataRoot(currentVersion, genesisValidatorsRoot)
}

// ComputeSigningRoot
// Return the signing root of an object by calculating the root of the object-domain tree.
// The root is the hash tree root of:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signingdata
func ComputeSigningRoot(sszObject ssz.HashRoot, domain []byte) ([]byte, error) {

	if len(domain) != 32 {
		return nil, errors.New(fmt.Sprintf("Domain should be in 32 bytes. Got %d.", len(domain)))
	}

	signingRoot := new(consensus.SigningData)
	copy(signingRoot.Domain[:], domain)
	objectRoot, err := sszObject.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	signingRoot.ObjectRoot = objectRoot

	messageToSign, err := signingRoot.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return messageToSign[:], nil
}
