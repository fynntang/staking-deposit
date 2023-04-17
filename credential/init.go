package credential

import (
	"github.com/fynntang/staking-deposit/crypto/bls"
	"log"
)

func init() {
	if err := bls.InitBLS(); err != nil {
		log.Fatalf("failed to init bls: %v", err)
	}
}
