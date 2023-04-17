package credential

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fynntang/staking-deposit/account_manager/verifier"
	"github.com/fynntang/staking-deposit/constants"
	"github.com/fynntang/staking-deposit/crypto/bls"
	"github.com/fynntang/staking-deposit/crypto/bls/blst"
	"github.com/fynntang/staking-deposit/crypto/key_derivation"
	"github.com/fynntang/staking-deposit/crypto/keystore"
	_mnemonic "github.com/fynntang/staking-deposit/crypto/mnemonic"
	"github.com/fynntang/staking-deposit/crypto/ssz"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	consensus "github.com/umbracle/go-eth-consensus"
	"math/big"
	"strconv"
)

type WithdrawalType int

const (
	BlsWithdrawal WithdrawalType = iota
	Eth1AddressWithdrawal
)

type Credential struct {
	SigningSk             *bls.PrivateKey
	WithdrawalSk          *bls.PrivateKey
	Amount                *big.Int
	Chain                 constants.Chain
	Eth1WithdrawalAddress common.Address
	SigningKeyPath        string
}

// NewCredential
//
// Set path as EIP-2334 format
// https://eips.ethereum.org/EIPS/eip-2334
func NewCredential(mnemonic, mnemonicPassword string, index int, amount *big.Int, chain constants.Chain, eth1WithdrawalAddress string) (*Credential, error) {
	purpose := 12381
	coinType := 3600
	account := strconv.Itoa(index)

	withdrawalKeyPath := fmt.Sprintf("m/%d/%d/%s/0", purpose, coinType, account)
	signingKeyPath := fmt.Sprintf("%s/0", withdrawalKeyPath)

	seed := _mnemonic.NewSeed(mnemonic, mnemonicPassword)

	withdrawalSK, err := key_derivation.MnemonicAndPathToKey(seed, withdrawalKeyPath)
	if err != nil {
		return nil, err
	}
	signingSK, err := key_derivation.MnemonicAndPathToKey(seed, signingKeyPath)
	if err != nil {
		return nil, err
	}

	return &Credential{
		SigningSk:             signingSK,
		WithdrawalSk:          withdrawalSK,
		Amount:                amount,
		Chain:                 chain,
		Eth1WithdrawalAddress: common.HexToAddress(eth1WithdrawalAddress),
		SigningKeyPath:        signingKeyPath,
	}, nil
}

func (c *Credential) SigningPK() blst.PublicKey    { return c.SigningSk.PublicKey() }
func (c *Credential) WithdrawalPK() blst.PublicKey { return c.WithdrawalSk.PublicKey() }
func (c *Credential) IsEmptyEth1WithdrawalAddress() bool {
	return c.Eth1WithdrawalAddress == (common.Address{})
}
func (c *Credential) WithdrawalPrefix() []byte {
	if !c.IsEmptyEth1WithdrawalAddress() {
		return constants.Eth1AddressWithdrawalPrefix

	}
	return constants.BlsWithdrawalPrefix
}
func (c *Credential) WithdrawalType() WithdrawalType {
	return WithdrawalType(bytes.Compare(c.WithdrawalPrefix(), constants.BlsWithdrawalPrefix))
}
func (c *Credential) WithdrawalCredentials() ([]byte, error) {

	var withdrawalCredentials []byte
	switch c.WithdrawalType() {
	case BlsWithdrawal:
		withdrawalCredentials = constants.BlsWithdrawalPrefix
		sum := sha256.New().Sum(c.WithdrawalPK().Marshal()) // todo: 取值待确认
		withdrawalCredentials = append(withdrawalCredentials, sum...)
	case Eth1AddressWithdrawal:
		if !c.IsEmptyEth1WithdrawalAddress() {
			withdrawalCredentials = constants.Eth1AddressWithdrawalPrefix
			withdrawalCredentials = append(withdrawalCredentials, make([]byte, 11)...)
			withdrawalCredentials = append(withdrawalCredentials, c.Eth1WithdrawalAddress.Bytes()...)
		}
	}
	return withdrawalCredentials, nil
}

func (c *Credential) DepositMessage() (*consensus.DepositMessage, error) {
	minAmount := decimal.NewFromFloat(constants.MinDepositAmount)
	amount := decimal.NewFromBigInt(c.Amount, 0)
	maxAmount := decimal.NewFromFloat(constants.MaxDepositAmount)

	if minAmount.Cmp(amount) == 1 && amount.Cmp(maxAmount) == 1 {
		return nil, fmt.Errorf("%s ETH deposits are not within the bounds of this cli", amount.Div(decimal.NewFromFloat(constants.ETH2GWei)).String())
	}

	withdrawalCredentials, err := c.WithdrawalCredentials()
	if err != nil {
		return nil, err
	}

	return &consensus.DepositMessage{
		Pubkey:                [48]byte(c.SigningPK().Marshal()),
		WithdrawalCredentials: [32]byte(withdrawalCredentials),
		Amount:                c.Amount.Uint64(),
	}, nil
}
func (c *Credential) SignedDeposit() (*consensus.DepositData, error) {
	depositMessage, err := c.DepositMessage()
	if err != nil {
		return nil, err
	}

	// deposit domain
	domain, err := ssz.ComputeDepositDomain(c.Chain.GenesisForkVersion)
	if err != nil {
		return nil, err
	}

	// sign
	messageToSign, err := ssz.ComputeSigningRoot(depositMessage, domain)
	if err != nil {
		return nil, err
	}

	signature := c.SigningSk.Sign(messageToSign)

	return &consensus.DepositData{
		Pubkey:                depositMessage.Pubkey,
		WithdrawalCredentials: depositMessage.WithdrawalCredentials,
		Amount:                depositMessage.Amount,
		Signature:             consensus.Signature(signature.Marshal()),
	}, nil
}

func (c *Credential) SigningDepositData() (*verifier.DepositData, error) {
	signedDeposit, err := c.SignedDeposit()
	if err != nil {
		return nil, err
	}

	depositMessage, err := c.DepositMessage()
	if err != nil {
		return nil, err
	}

	depositMessageRoot, err := depositMessage.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	signedDepositRoot, err := signedDeposit.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	return &verifier.DepositData{
		PublicKey:             hex.EncodeToString(signedDeposit.Pubkey[:]),
		WithdrawalCredentials: hex.EncodeToString(signedDeposit.WithdrawalCredentials[:]),
		Amount:                signedDeposit.Amount,
		Signature:             hex.EncodeToString(signedDeposit.Signature[:]),
		DepositMessageRoot:    hex.EncodeToString(depositMessageRoot[:]),
		DepositDataRoot:       hex.EncodeToString(signedDepositRoot[:]),
		ForkVersion:           hex.EncodeToString(c.Chain.GenesisForkVersion[:]),
		NetworkName:           string(c.Chain.NetworkName),
		DepositCliVersion:     constants.DepositCliVersion,
	}, nil
}
func (c *Credential) VerifyDepositData(data *verifier.DepositData) error {
	pubkey, err := hex.DecodeString(data.PublicKey)
	if err != nil {
		return err
	}
	pk, err := bls.PublicKeyFromBytes(pubkey)
	if err != nil {
		return err
	}

	withdrawalCredentials, err := hex.DecodeString(data.WithdrawalCredentials)
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(data.Signature)
	if err != nil {
		return err
	}
	sign, err := bls.SignatureFromBytes(signature)
	if err != nil {
		return err
	}
	depositMessageRoot, err := hex.DecodeString(data.DepositDataRoot)
	if err != nil {
		return err
	}

	forkVersion, err := hex.DecodeString(data.ForkVersion)
	if err != nil {
		return err
	}

	if !bytes.Equal(pk.Marshal(), c.SigningPK().Marshal()) {
		return fmt.Errorf("public key not match")
	}

	if len(withdrawalCredentials) != 32 {
		return fmt.Errorf("withdrawal credentials length not match")
	}
	if bytes.Equal(withdrawalCredentials[:1], constants.BlsWithdrawalPrefix) &&
		bytes.Equal(withdrawalCredentials[:1], c.WithdrawalPrefix()) {
		if !bytes.Equal(withdrawalCredentials[1:], key_derivation.SHA256(c.SigningPK().Marshal())[1:]) {
			return fmt.Errorf("withdrawal credentials not match")
		}
	} else if bytes.Equal(withdrawalCredentials[:1], constants.Eth1AddressWithdrawalPrefix) &&
		bytes.Equal(withdrawalCredentials[:1], c.WithdrawalPrefix()) {
		if !bytes.Equal(withdrawalCredentials[1:12], make([]byte, 11)) {
			return fmt.Errorf("withdrawal credentials not match")
		}
		if c.IsEmptyEth1WithdrawalAddress() {
			return fmt.Errorf("eth1 withdrawal address is empty")
		}

		if !bytes.Equal(withdrawalCredentials[12:], c.Eth1WithdrawalAddress.Bytes()) {
			return fmt.Errorf("withdrawal credentials not match")
		}
	} else {
		return fmt.Errorf("withdrawal credentials not match")
	}

	minAmount := decimal.NewFromFloat(constants.MinDepositAmount)
	amount := decimal.NewFromBigInt(c.Amount, 0)
	maxAmount := decimal.NewFromFloat(constants.MaxDepositAmount)

	if minAmount.Cmp(amount) == 1 && amount.Cmp(maxAmount) == 1 {
		return fmt.Errorf("%s ETH deposits are not within the bounds of this cli", amount.Div(decimal.NewFromFloat(constants.ETH2GWei)).String())
	}

	domain, err := ssz.ComputeDepositDomain([4]byte(forkVersion))
	if err != nil {
		return err
	}

	depositMessage := &consensus.DepositMessage{Pubkey: [48]byte(pk.Marshal()), WithdrawalCredentials: [32]byte(withdrawalCredentials), Amount: data.Amount}

	signingRoot, err := ssz.ComputeSigningRoot(depositMessage, domain)
	if err != nil {
		return err
	}

	if !sign.Verify(signingRoot, pk) {
		return fmt.Errorf("signature Verify failed")
	}
	depositData := &consensus.DepositData{
		Pubkey:                [48]byte(pk.Marshal()),
		WithdrawalCredentials: [32]byte(withdrawalCredentials),
		Amount:                data.Amount,
		Signature:             consensus.Signature(sign.Marshal()),
	}

	depositDataRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return err
	}

	if !bytes.Equal(depositDataRoot[:], depositMessageRoot) {
		return fmt.Errorf("deposit data root not match")
	}

	return nil
}

func (c *Credential) SigningKeystore(password string) (*verifier.Keystore, error) {
	encryptor := keystore.NewScryptKeystoreEncryptor()
	cryptoFields, err := encryptor.Encrypt(c.SigningSk.Marshal(), password)
	if err != nil {
		return nil, err
	}

	return &verifier.Keystore{
		Description: encryptor.Name(),
		Version:     encryptor.Version(),
		UUID:        uuid.New().String(),
		Path:        c.SigningKeyPath,
		Pubkey:      hex.EncodeToString(c.SigningPK().Marshal()),
		Crypto:      cryptoFields,
	}, nil
}
func (c *Credential) VerifyKeystore(cryptoFields map[string]interface{}, password string) error {
	decryptor := keystore.NewScryptKeystoreEncryptor()
	decrypt, err := decryptor.Decrypt(cryptoFields, password)
	if err != nil {
		return err
	}

	if !bytes.Equal(decrypt, c.SigningSk.Marshal()) {
		return fmt.Errorf("keystore password is incorrect")
	}

	return nil
}
