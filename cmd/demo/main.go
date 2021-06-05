package main

import (
	"log"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/weaming/bitcoin-addr-generator/bips"
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
	mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect"

	// the following output could be verified with `btckeygen` command
	// installation: go get -u github.com/modood/btckeygen@HEAD
	// btckeygen -mnemonic 'scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect' -pass 'secret passphrase' -bip39

	addrRoot, addr, e := bips.BIP44(mnemonic, "secret passphrase", 2)
	if e != nil {
		panic(e)
	}
	log.Println("addrRoot:", bips.MarshalJSONIndent(addrRoot))
	log.Println("addr:", bips.MarshalJSONIndent(addr))

	addrRoot, addr, e = bips.BIP49(mnemonic, "secret passphrase", 2)
	if e != nil {
		panic(e)
	}
	log.Println("addrRoot:", bips.MarshalJSONIndent(addrRoot))
	log.Println("addr:", bips.MarshalJSONIndent(addr))

	addrRoot, addr, e = bips.BIP84(mnemonic, "secret passphrase", 2)
	if e != nil {
		panic(e)
	}
	log.Println("addrRoot:", bips.MarshalJSONIndent(addrRoot))
	log.Println("addr:", bips.MarshalJSONIndent(addr))
}
