package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"github.com/weaming/bitcoin-addr-generator/bips"
	"github.com/weaming/bitcoin-addr-generator/random"
)

var fortuna = random.NewFortunaWrap("./randomness")

type RequestHDSegwitAddr struct {
	Mnemonic, Passphrase string
	Path                 string
}

type ResponseHDSegwitAddr struct {
	AddressRoot *bips.AddressRoot
	Address     *bips.Address
}

func handleHDSegwitAddr(w http.ResponseWriter, r *http.Request) {
	// update randomness
	fortuna.Update()

	req := &RequestHDSegwitAddr{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpError(w, "parse request err: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Mnemonic == "" && os.Getenv("TEST") != "" {
		// mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		req.Mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect"
	}
	log.Printf("req: %+v", req)

	indexes, e := CheckHDPath(req.Path)
	if e != nil {
		httpError(w, e.Error(), http.StatusBadRequest)
		return
	}
	log.Println("indexes:", indexes)

	var res *ResponseHDSegwitAddr
	switch indexes[0] {
	case 44:
		r, a, e := bips.BIP44(req.Mnemonic, req.Passphrase, indexes[len(indexes)-1])
		if e != nil {
			httpError(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &ResponseHDSegwitAddr{r, a}
	case 49:
		r, a, e := bips.BIP49(req.Mnemonic, req.Passphrase, indexes[len(indexes)-1])
		if e != nil {
			httpError(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &ResponseHDSegwitAddr{r, a}
	case 84:
		r, a, e := bips.BIP84(req.Mnemonic, req.Passphrase, indexes[len(indexes)-1])
		if e != nil {
			httpError(w, e.Error(), http.StatusBadRequest)
			return
		}
		res = &ResponseHDSegwitAddr{r, a}
	default:
		httpError(w, "unknown BIP", http.StatusBadRequest)
		return
	}
	httpJSON(w, res, 200)
}

type RequestMultisigP2SHAddr struct {
	PubKeys []string
	N       uint32
}

type ResponseMultisigP2SHAddr struct {
	Address string
}

var OpArr = []byte{
	0x00,
	txscript.OP_1,
	txscript.OP_2,
	txscript.OP_3,
	txscript.OP_4,
	txscript.OP_5,
	txscript.OP_6,
	txscript.OP_7,
	txscript.OP_8,
	txscript.OP_9,
	txscript.OP_10,
	txscript.OP_11,
	txscript.OP_12,
	txscript.OP_13,
	txscript.OP_14,
	txscript.OP_15,
	txscript.OP_16,
}

func handleMultisigP2SHAddr(w http.ResponseWriter, r *http.Request) {
	// update randomness
	fortuna.Update()

	req := &RequestMultisigP2SHAddr{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpError(w, "parse request err: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("req: %+v", req)

	if len(req.PubKeys) <= 1 {
		httpError(w, "need at least 2 public keys", http.StatusBadRequest)
		return
	}
	if len(req.PubKeys) > 16 {
		httpError(w, "at most 16 public keys", http.StatusBadRequest)
		return
	}
	if req.N < 1 || 16 < req.N {
		httpError(w, "N must be 2 <= N <= 16", http.StatusBadRequest)
		return
	}
	if req.N > uint32(len(req.PubKeys)) {
		httpError(w, "N should not be greater than the number of public keys", http.StatusBadRequest)
		return
	}

	pks := []btcutil.Address{}
	for _, pk := range req.PubKeys {
		addr, e := btcutil.DecodeAddress(pk, &chaincfg.MainNetParams)
		if e != nil {
			httpError(w, fmt.Sprintf("public address is invalid: %s", pk), http.StatusBadRequest)
			return
		}
		pks = append(pks, addr)
	}

	res := &ResponseMultisigP2SHAddr{}
	builder := txscript.NewScriptBuilder()

	builder.AddOp(OpArr[req.N])
	for _, pk := range pks {
		builder.AddData(pk.ScriptAddress())
	}

	builder.AddOp(OpArr[len(req.PubKeys)])
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	redeemScript, e := builder.Script()
	if e != nil {
		httpError(w, e.Error(), http.StatusInternalServerError)
		return
	}

	redeemHash := btcutil.Hash160(redeemScript)
	addr, e := btcutil.NewAddressScriptHashFromHash(redeemHash, &chaincfg.MainNetParams)
	if e != nil {
		httpError(w, e.Error(), http.StatusInternalServerError)
		return
	}
	res.Address = addr.EncodeAddress()
	httpJSON(w, res, 200)
}

// Write error string which is wrapped in a JSON.
func httpError(w http.ResponseWriter, err string, code int) {
	w.Header().Add("Content-Type", "application/json")
	out, e := json.Marshal(map[string]interface{}{"error": err})
	if e != nil {
		panic(e) // should not happen
	}
	w.Write(out)
}

// Write something as JSON to the client.
func httpJSON(w http.ResponseWriter, v interface{}, code int) {
	out, e := json.Marshal(v)
	if e != nil {
		httpError(w, e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
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
	http.HandleFunc("/api/hd-segwit-address", handleHDSegwitAddr)
	http.HandleFunc("/api/multisig-p2sh-address", handleMultisigP2SHAddr)
	http.ListenAndServe(":8080", nil)
}
