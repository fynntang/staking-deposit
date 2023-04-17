package blst

// Signature is an interface for signatures.
type Signature interface {
	Verify(msg []byte, pub PublicKey) bool
	VerifyAggregate(msgs [][]byte, pubKeys []PublicKey) bool
	VerifyAggregateCommon(msg []byte, pubKeys []PublicKey) bool
	Marshal() []byte
}
