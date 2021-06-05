package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// path format: m / purpose' / coin_type' / account' / change / address_index
// Apostrophe in the path indicates that BIP32 hardened derivation is used.
var PathPat = regexp.MustCompile(`^m/(44|49|84)'/([01])'/(\d+)'/([01])/(\d+)$`)

func CheckHDPath(keyPath string) ([]uint32, error) {
	indexes := []uint32{}
	xs := PathPat.FindStringSubmatch(keyPath)
	if len(xs) == 0 {
		return indexes, fmt.Errorf("invalid path: %v", keyPath)
	}

	for i, s := range xs[1:] {
		depth := i + 2 // human readable index in the path
		i32, e := strconv.ParseUint(s, 10, 32)
		if e != nil {
			return indexes, fmt.Errorf(`invalid path: path(depth %d) parsed fail, error: %v`, depth, e)
		}
		indexes = append(indexes, uint32(i32))
	}
	return indexes, nil
}
