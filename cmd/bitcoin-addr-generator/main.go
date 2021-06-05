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

var fortuna = random.NewFortunaWrap("./randomness")

type Request struct {
	Mnemonic, Passphase string
	Path                string
}

type Response struct {
	AddressRoot *bips.AddressRoot
	Address     *bips.Address
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// update randomness
	fortuna.Update()

	req := &Request{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, "parse request err: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Mnemonic == "" && os.Getenv("TEST") != "" {
		// mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		req.Mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect"
	}
	log.Printf("req: %+v", req)

	indexes, e := CheckHDPath(req.Path)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	log.Println("indexes:", indexes)

	var res *Response
	switch indexes[0] {
	case 44:
		r, a, e := bips.BIP44(req.Mnemonic, req.Passphase, indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	case 49:
		r, a, e := bips.BIP49(req.Mnemonic, req.Passphase, indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	case 84:
		r, a, e := bips.BIP84(req.Mnemonic, req.Passphase, indexes[len(indexes)-1])
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &Response{r, a}
	default:
		http.Error(w, "unknown BIP", http.StatusBadRequest)
		return
	}
	out, e := json.Marshal(res)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

func main() {
	// load our rand seed
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

	log.Println("serve on http://0.0.0.0:8080")
	http.HandleFunc("/", handleHTTP)
	http.ListenAndServe(":8080", nil)
}
