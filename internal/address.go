package internal

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"math/big"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 21
	// AddressLengthBase58 is the expected length of the address in base58format
	AddressLengthBase58 = 34
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

// Address represents the 21 byte address of an Tron account.
type Address []byte

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex get bytes from address in string
func (a Address) Hex() string {
	return ToHex(a[:])
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address {
	id := b.Bytes()
	base := bytes.Repeat([]byte{0}, AddressLength-len(id))
	return append(base, id...)
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address {
	addr, err := FromHex(s)
	if err != nil {
		return nil
	}
	return addr
}

// Base58ToAddress returns Address with byte values of s.
func Base58ToAddress(s string) (Address, error) {
	addr := DecodeCheck(s)
	if addr == nil {
		return nil, fmt.Errorf("decode")
	}
	return addr, nil
}

// Base64ToAddress returns Address with byte values of s.
func Base64ToAddress(s string) (Address, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return Address(decoded), nil
}

// String implements fmt.Stringer.
func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}

	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}
	return EncodeCheck(a.Bytes())
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	*a = Address(srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return []byte(a), nil
}

func parsePrivateKey(private string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return nil
	}
	return privateKey
}

func PrivateKeyToPublicKey(pk string) []byte {
	privateKey := parsePrivateKey(pk)
	if privateKey == nil {
		return nil
	}
	return elliptic.Marshal(privateKey.Curve, privateKey.X, privateKey.Y)
}

func PublicKeyToTronAddress(publicKey []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:])
	address := hash.Sum(nil)
	address = address[len(address)-20:]
	tronAddress := make([]byte, 21)
	tronAddress[0] = 0x41
	copy(tronAddress[1:], address)
	return tronAddress
}
