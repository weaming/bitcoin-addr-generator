package base58check

import (
	"crypto/sha256"

	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
)

const (
	// compressMagic is the magic byte used to identify a WIF encoding for
	// an address created from a compressed serialized public key.
	compressMagic byte = 0x01
)

// https://bitcoin.stackexchange.com/a/3839
func WIF(key *bip32.Key, compress bool) string {
	a := []byte{0x80}
	a = append(a, key.Key...)
	if compress {
		a = append(a, compressMagic)
	}
	s1 := sha256.Sum256(a)
	s2 := sha256.Sum256(s1[:])
	a = append(a, s2[:4]...)
	return base58.Encode(a)
}
