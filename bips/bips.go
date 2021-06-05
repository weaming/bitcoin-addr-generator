package bips

import "encoding/hex"

const (
	BitSize128 = 128
	BitSize256 = 256
	Compresss  = true
)

type AddressRoot struct {
	Mnemonic   string
	Passphrase string
	Seed       string
	RootKey    string
}

type Address struct {
	BIP          string `json:",omitempty"`
	Path         string `json:",omitempty"`
	WIF          string `json:",omitempty"`
	Address      string `json:",omitempty"`
	SegwitBech32 string `json:",omitempty"`
	SegwitNested string `json:",omitempty"`
}

func BIPCommon(mnemonic, passphrase string) (*KeyManager, *AddressRoot, error) {
	km, e := NewKeyManager(BitSize256, passphrase, mnemonic)
	if e != nil {
		return nil, nil, e
	}
	masterKey, e := km.GetMasterKey()
	if e != nil {
		return km, nil, e
	}
	rootKey := masterKey.B58Serialize()
	return km, &AddressRoot{km.GetMnemonic(), km.GetPassphrase(), hex.EncodeToString(km.GetSeed()), rootKey}, nil
}

func BIP44(mnemonic, passphrase string, index uint32) (*AddressRoot, *Address, error) {
	km, addrRoot, e := BIPCommon(mnemonic, passphrase)
	if e != nil {
		return nil, nil, e
	}

	key, e := km.GetKey(PurposeBIP44, CoinTypeBTC, 0, 0, index)
	if e != nil {
		return addrRoot, nil, e
	}
	wif, address, _, _, e := key.Encode(Compresss)
	if e != nil {
		return addrRoot, nil, e
	}
	return addrRoot, &Address{
		BIP:  "BIP44",
		Path: key.GetPath(), WIF: wif,
		Address:      address,
		SegwitBech32: "",
		SegwitNested: "",
	}, nil
}

func BIP49(mnemonic, passphrase string, index uint32) (*AddressRoot, *Address, error) {
	km, addrRoot, e := BIPCommon(mnemonic, passphrase)
	if e != nil {
		return nil, nil, e
	}

	key, e := km.GetKey(PurposeBIP49, CoinTypeBTC, 0, 0, index)
	if e != nil {
		return nil, nil, e
	}
	wif, _, _, segwitNested, e := key.Encode(Compresss)
	if e != nil {
		return nil, nil, e
	}
	return addrRoot, &Address{
		BIP:  "BIP49",
		Path: key.GetPath(), WIF: wif,
		Address:      "",
		SegwitBech32: "",
		SegwitNested: segwitNested,
	}, nil
}

func BIP84(mnemonic, passphrase string, index uint32) (*AddressRoot, *Address, error) {
	km, addrRoot, e := BIPCommon(mnemonic, passphrase)
	if e != nil {
		return nil, nil, e
	}

	key, e := km.GetKey(PurposeBIP84, CoinTypeBTC, 0, 0, index)
	if e != nil {
		return nil, nil, e
	}
	wif, _, segwitBech32, _, e := key.Encode(Compresss)
	if e != nil {
		return nil, nil, e
	}
	return addrRoot, &Address{
		BIP:  "BIP84",
		Path: key.GetPath(), WIF: wif,
		Address:      "",
		SegwitBech32: segwitBech32,
		SegwitNested: "",
	}, nil
}
