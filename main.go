package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/weaming/bitcoin-addr-generator/random"
)

func main() {
	// load our rand seed
	random.InitRandSeed()

	// set mnenomic words list based on environment LANG
	lang := os.Getenv("LANG")
	for l, wl := range map[string][]string{
		"en_US": wordlists.English,
		"zh_CN": wordlists.ChineseSimplified,
	} {
		if strings.Contains(lang, l) {
			bip39.SetWordList(wl)
			break
		}
	}

	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "secret passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
}
