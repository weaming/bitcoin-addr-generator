package bips

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func MarshalJSON(in interface{}) string {
	bs, e := json.Marshal(in)
	if e != nil {
		log.Println("marshal json fail:", e)
		return "{}"
	}
	return string(bs)
}

func MarshalJSONIndent(in interface{}) string {
	bs, e := json.MarshalIndent(in, "", "  ")
	if e != nil {
		log.Println("marshal json with indent fail:", e)
		return "{}"
	}
	return string(bs)
}

func CheckHDPath(keyPath string) ([]uint32, error) {
	indexes := []uint32{}
	// path format: m / purpose' / coin_type' / account' / change / address_index
	// Apostrophe in the path indicates that BIP32 hardened derivation is used.
	xs := strings.Split(keyPath, "/")
	if len(xs) < 6 {
		return indexes, fmt.Errorf("invalid path: depth not enough")
	}
	if xs[0] != "m" {
		return indexes, fmt.Errorf("invalid path: first part is not 'm'")
	}

	for i, s := range xs[1:] {
		depth := i + 2 // human readable index in the path
		useHardendedAddr := false
		if 2 <= depth && depth <= 4 {
			if !strings.HasSuffix(s, "'") {
				return indexes, fmt.Errorf(`invalid path: path(depth %d) must ends with "'"`, depth)
			}
			s = strings.Replace(s, "'", "", 1)
			useHardendedAddr = true
		}
		if 6 == depth {
			if strings.HasSuffix(s, "'") {
				s = strings.Replace(s, "'", "", 1)
				useHardendedAddr = true
			}
		}

		i2, e := strconv.ParseUint(s, 10, 32)
		if e != nil {
			return indexes, fmt.Errorf(`invalid path: path(depth %d) parsed fail, error: %v`, depth, e)
		}

		if 2 == depth {
			if i2 != 44 {
				return indexes, fmt.Errorf(`invalid path: path (%s)(depth %d) should be 44'`, s, depth)
			}
		}
		if 3 == depth {
			if i2 != 0 && i2 != 1 {
				return indexes, fmt.Errorf(`invalid path: path (%s)(depth %d) represents unknown coin type`, s, depth)
			}
		}
		if 5 == depth {
			if i2 != 0 && i2 != 1 {
				return indexes, fmt.Errorf(`invalid path: path (%s)(depth %d) represents unknown chain type`, s, depth)
			}
		}

		keyIndex := uint32(i2)
		if useHardendedAddr {
			keyIndex += Apostrophe
		}
		indexes = append(indexes, keyIndex)
	}
	return indexes, nil
}
