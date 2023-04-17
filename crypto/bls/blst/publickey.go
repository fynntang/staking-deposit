package blst

// PublicKey is an interface for public keys
type PublicKey interface {
	Marshal() []byte
	Aggregate(other PublicKey)
	Copy() PublicKey
}
