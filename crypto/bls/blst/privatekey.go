package blst

// PrivateKey is a private key in Ethereum 2.
type PrivateKey interface {
	PublicKey() PublicKey
	Sign(msg []byte) Signature
	Marshal() []byte
}
