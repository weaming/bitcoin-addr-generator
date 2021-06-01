package bip44

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip32"
)

var hardenedOffset uint64 = 0x01 << 31

func BIP44Addr(extendedKey *bip32.Key, keyPath string) (*bip32.Key, error) {
	// path format: m / purpose' / coin_type' / account' / change / address_index
	// Apostrophe in the path indicates that BIP32 hardened derivation is used.
	xs := strings.Split(keyPath, "/")
	if len(xs) < 6 {
		return nil, fmt.Errorf("invalid path: depth not enough")
	}
	if xs[0] != "m" {
		return nil, fmt.Errorf("invalid path: first part is not 'm'")
	}
	for i, s := range xs[1:] {
		depth := i + 2 // human readable index in the path
		useHardendedAddr := false
		if 2 <= depth && depth <= 4 {
			if !strings.HasSuffix(s, "'") {
				return nil, fmt.Errorf(`invalid path: path(depth %d) must ends with "'"`, depth)
			}
			s = strings.Replace(s, "'", "", 1)
			useHardendedAddr = true
		}

		i2, e := strconv.ParseUint(s, 10, 32)
		if e != nil {
			return nil, fmt.Errorf(`invalid path: path(depth %d) parsed fail, error: %v`, depth, e)
		}

		if 2 == depth {
			if i2 != 44 {
				return nil, fmt.Errorf(`invalid path: path (%s)(depth %d) should be 44'`, s, depth)
			}
		}
		if 3 == depth {
			if i2 != 0 && i2 != 1 {
				return nil, fmt.Errorf(`invalid path: path (%s)(depth %d) represents unknown coin type`, s, depth)
			}
		}
		if 5 == depth {
			if i2 != 0 && i2 != 1 {
				return nil, fmt.Errorf(`invalid path: path (%s)(depth %d) represents unknown chain type`, s, depth)
			}
		}
		if 6 == depth {
			if strings.HasSuffix(s, "'") {
				s = strings.Replace(s, "'", "", 1)
				useHardendedAddr = true
			}
		}

		keyIndex := uint64(i2)
		if useHardendedAddr {
			keyIndex += hardenedOffset
		}
		extendedKey, e = extendedKey.NewChildKey(uint32(keyIndex))
		if e != nil {
			return nil, fmt.Errorf(`generate child key fail: depth %d, error: %v`, depth, e)
		}
	}
	return extendedKey, nil
}
