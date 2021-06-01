package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/weaming/bitcoin-addr-generator/base58check"
	"github.com/weaming/bitcoin-addr-generator/bip44"
	"github.com/weaming/bitcoin-addr-generator/random"
)

func main() {
	// load our rand seed
	fortuna := random.NewFortunaWrap("./randomness")
	defer fortuna.Close()
	fortuna.Update()

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
	// for test
	// mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect"
	seed := bip39.NewSeed(mnemonic, "secret passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	fmt.Println("Mnemonic:", mnemonic)
	fmt.Println("Seed:", hex.EncodeToString(seed))
	fmt.Println("Master private key:", masterKey)
	fmt.Println("Master public key:", publicKey)

	k, e := bip44.BIP44Addr(masterKey, "m/44'/0'/0'/0/2")
	if e != nil {
		panic(e)
	}
	// validate it on https://iancoleman.io/bip39
	fmt.Println("BIP44", base58check.WIF(k, true), k)
}
