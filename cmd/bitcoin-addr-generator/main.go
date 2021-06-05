package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/weaming/bitcoin-addr-generator/bips"
	"github.com/weaming/bitcoin-addr-generator/random"
)

var mnemonic = ""

type Request struct {
	BIP  string
	Path string
}

type Response struct {
	AddressRoot *bips.AddressRoot
	Address     *bips.Address
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	err := json.NewDecoder(r.Body).Decode(req)
	log.Printf("req: %+v", req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indexes, e := bips.CheckHDPath(req.Path)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	var res *Response
	switch req.BIP {
	case "BIP44", "BIP0044":
		r, a, e := bips.BIP44(mnemonic, "", indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	case "BIP49", "BIP0049":
		r, a, e := bips.BIP49(mnemonic, "", indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	case "BIP84", "BIP0084":
		r, a, e := bips.BIP84(mnemonic, "", indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	default:
		http.Error(w, "missing BIP", http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	out, e := json.Marshal(res)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

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

	if os.Getenv("TEST") != "" {
		mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect"
	}

	log.Println("serve on http://0.0.0.0:8080")
	http.HandleFunc("/", handleHTTP)
	http.ListenAndServe(":8080", nil)
}
