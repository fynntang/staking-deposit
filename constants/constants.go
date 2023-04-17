package constants

var (
	ZeroBytes32 [32]byte

	// DomainDeposit Execution-spec constants taken from https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md
	DomainDeposit               = [4]byte{0x03, 0x00, 0x00, 0x00}
	BlsWithdrawalPrefix         = []byte{0x00}
	Eth1AddressWithdrawalPrefix = []byte{0x01}
)

const (
	ETH2GWei         = 1 * 1e9
	MinDepositAmount = (1 << 0) * ETH2GWei
	MaxDepositAmount = (1 << 5) * ETH2GWei

	// DefaultValidatorKeysFolderName File/folder constants
	DefaultValidatorKeysFolderName = "validator_keys"
)
