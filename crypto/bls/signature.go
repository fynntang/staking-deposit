package bls

import (
	"fmt"
	"github.com/fynntang/staking-deposit/crypto/bls/blst"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// Signature is a BLS signature.
type Signature struct {
	sig *bls.Sign
}

// Verify a bls signature given a public key and a message.
func (s *Signature) Verify(msg []byte, pubKey blst.PublicKey) bool {
	return s.sig.VerifyByte(pubKey.(*PublicKey).key, msg)
}

// VerifyAggregate verifies each public key against its respective message.
// Note: this is vulnerable to a rogue public-key attack.
func (s *Signature) VerifyAggregate(msgs [][]byte, pubKeys []blst.PublicKey) bool {
	if len(pubKeys) == 0 {
		return false
	}
	keys := make([]bls.PublicKey, len(pubKeys))
	for i, v := range pubKeys {
		keys[i] = *v.(*PublicKey).key
	}
	return s.sig.VerifyAggregateHashes(keys, msgs)
}

// VerifyAggregateCommon verifies each public key against a single message.
// Note: this is vulnerable to a rogue public-key attack.
func (s *Signature) VerifyAggregateCommon(msg []byte, pubKeys []blst.PublicKey) bool {
	if len(pubKeys) == 0 {
		return false
	}
	keys := make([]bls.PublicKey, len(pubKeys))
	for i, v := range pubKeys {
		keys[i] = *v.(*PublicKey).key
	}
	return s.sig.FastAggregateVerify(keys, msg)
}

// Marshal a signature into a byte slice.
func (s *Signature) Marshal() []byte {
	return s.sig.Serialize()
}

// SignatureFromBytes creates a BLS signature from a byte slice.
func SignatureFromBytes(data []byte) (blst.Signature, error) {
	var sig bls.Sign
	if err := sig.Deserialize(data); err != nil {
		return nil, fmt.Errorf("failed to deserialize signature: %v", err)
	}
	return &Signature{sig: &sig}, nil
}

// SignatureFromSig creates a BLS signature from an existing signature.
func SignatureFromSig(sig bls.Sign) (blst.Signature, error) {
	return &Signature{sig: &sig}, nil
}

// AggregateSignatures aggregates signatures.
func AggregateSignatures(sigs []blst.Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	aggSig := &bls.Sign{}
	//#nosec G104
	_ = aggSig.Deserialize(sigs[0].(*Signature).Marshal())
	for _, sig := range sigs[1:] {
		aggSig.Add(sig.(*Signature).sig)
	}
	return &Signature{sig: aggSig}
}
