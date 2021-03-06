This is the note on the learning path and organized by 3 questions in its' order.

- Git clone this repo and run `go run ./cmd/bitcoin-addr-generator` to run the HTTP API server.
- See the [test-api.sh](./test-api.sh) for examples to call the HTTP API.

## 1. How to provide the seed onto this server

1. Upload seed/mnemonic via the HTTP request
    - Pros: random seed provided by the client
    - Cons: users should not use their largely used seed
2. Generate from random number and return the seed to the client.
    - Pros: random seed provided by the server
    - Cons: server should not use PRNG without a randomness source, else will be more predictable by external attackers

Both need the trust of client to the server.

### Generate random number with unpredictable generator

Library [fortuna](https://github.com/seehuhn/fortuna) implemented the Fortuna algorithm. It accept updates of randomness from the environment.

### References

- [Random number generation - Wikipedia](https://en.wikipedia.org/wiki/Random_number_generation)
- The PRNG algorithm [Mersenne Twister](https://en.wikipedia.org/wiki/Mersenne_Twister) is used in many languages and libraries.
- [bitaddress.org](https://www.bitaddress.org)
- Fortuna [random(4)](https://www.freebsd.org/cgi/man.cgi?query=random&apropos=0&sektion=4&manpath=FreeBSD+11.0-RELEASE+and+Ports&arch=default&format=html)

## 2. Generate HD SegWit Bitcoin Address

For BTC, there is 3 types of address:

1. Legacy (P2PKH): addresses start with a 1
2. Nested SegWit (P2SH): addresses start with a 3
3. Native SegWit (bech32): addresses start with bc1

### HD Wallet

- ECDSA
- secp256k1
- chain code
- CDK funtions
    - parent private key -> child private key: ok
    - parent public key -> child public key: only for non-hardended child keys
    - parent private key -> child public key:
        1. parent private key -> child private key -> extended public key + extended "neutered" version private key: ok
        2. parent private key -> extended public key + extended "neutered" version private key ?-> child private key: only for non-hardended child keys
    - parent public key -> child private key: not possible
- normal/hardended child key: With non-hardened keys, you can prove a child public key is linked to a parent public key using just the public keys. You can also derive public child keys from a public parent key, which enables watch-only wallets. With hardened child keys, you cannot prove that a child public key is linked to a parent public key.
- CKDpriv; CKDpub(only for non-hardened child keys)
- Sha256, Base58, HMAC-SHA512


#### Graph for CKD functions

```
-- stands for works always
???? stands for only for non-hardended key

 parent pub ????????????????CKDpub????????????>|-----------|
                                 |           |
                                 | child pub |
                                 |           |
 ext-priv + ext-pub ??????CKDpub??????>|-----------|
          ^
          |
          N
          |
parent priv -----CKDpriv--------> child priv ---N--> ext-pub + ext-priv
```

#### References

- Mnemonic code for generating deterministic keys [bips/bip-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- Hierarchical Deterministic Wallets [bips/bip-0032](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- Multi-Account Hierarchy for Deterministic Wallets [bips/bip-0044](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
- Registered coin types for BIP-0044 [slips/slip-0044.md](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
- [???????????? HD ??????????????? BIP32???BIP44???BIP39](https://learnblockchain.cn/2018/09/28/hdwallet/)
- [BIP39 - Mnemonic Code](https://iancoleman.io/bip39/#english) online tool
- [Bitcoin Private Key - BitcoinWiki](https://en.bitcoinwiki.org/wiki/Private_key)
- [Base58Check encoding - Bitcoin Wiki](https://en.bitcoin.it/wiki/Base58Check_encoding)
- [4. Keys, Addresses, Wallets - Mastering Bitcoin](https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch04.html)

### SegWit Address

- SegWit is the process by which the block size limit on a blockchain is increased by removing signature data from bitcoin transactions. When certain parts of a transaction are removed, this frees up space or capacity to add more transactions to the chain. Segregate means to separate, and witnesses are the transaction signatures. Hence, segregated witness, in short, means to separate transaction signatures.
- hard/soft forks
- Merkle Tree Root
- Block
    - Coinbase
    - transaction
        - TXID
        - version
            - set of rules
        - inputs 
            - outpoint
                - TXID
                - vout
            - signature script
        - outputs
            - outpoint
                - TXID (32 bytes)
                - Output index number (vout, 4 bytes)
            - pubkey script
            - UTXOs (cannot be spent for at least 100 blocks)
            - Full Public Key
        - wtxid: the double SHA256 of the serialization of all witness data of the transaction
        - Locktime
- script languagae: stack based, stateless and not Turing complete (See book *Mastering Bitcoin*)
    - locking script
        - scriptPubKey
        - witness script
        - cryptographic puzzle
    - unlocking script
        - scriptSig
        - witness
        - redeemScript
    - witnessScript
- all similar words found in bips repo: `ag 'P2[A-Z]*[KH]' -o | grep ':' | awk -F: '{print $3}' | sort | uniq`
    - P2PK: pay to public key
    - P2PKH: pay to public key hash
    - P2SH: pay to script hash
    - P2WPKH: pay to witness public key hashs (see BIP84 for address format)
    - P2WPKH nested in BIP16 P2SH (see BIP49 for address format)
    - P2WSH: pay to witness script hash
    - P2WSH nested in BIP16 P2SH
- struct of scriptPubKey (see BIP-0141)
    - version byte (1 byte)
    - witness program (2~40 bytes)
        - two types
            1. native witness program: a version byte + a witness program
            2. P2SH witness program: a P2SH script and the content of scriptSig is a version byte + a witness program
        - versions
            - version 0: the current version
                - witness program length is 20 bytes: P2WPKH
                - witness program length is 32 bytes: P2WSH
                - witness program length is neither 20 nor 32 bytes: fail
            - version 1~16: reserved for future extensions

#### References

- [Cryptocurrency standards - Trezor Wiki](https://wiki.trezor.io/Cryptocurrency_standards)
- [Developer Glossary - Bitcoin](https://btcinformation.org/en/developer-glossary)
- [Developer Guide - Bitcoin](https://btcinformation.org/en/developer-guide)
- [Segregated Witness - Bitcoin Wiki](https://en.bitcoin.it/wiki/Segregated_Witness)
    - [Payment channels - Bitcoin Wiki](https://en.bitcoin.it/wiki/Payment_channels)
    - [Lightning Network - Bitcoin Wiki](https://en.bitcoin.it/wiki/Lightning_Network)
    - [??????????????????????????????????????????????????????](https://www.528btc.com/bk/2019111158642.html)
- Segregated Witness [bips/bip-0141](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki)
    - P2WPKH
        - Derivation scheme for P2WPKH-nested-in-P2SH based accounts [bips/bip-0049](https://github.com/bitcoin/bips/blob/master/bip-0049.mediawiki)
        - Derivation scheme for P2WPKH based accounts [bips/bip-0084](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki)
            - Base32 address format for native v0-16 witness outputs [bips/bip-0173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki)
    - P2WSH
        - Transaction Signature Verification for Version 0 Witness Program [bips/bip-0143](https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki)
- [Script - Bitcoin Wiki](https://en.bitcoin.it/wiki/Script)
    - OP_HASH160: The input is hashed twice: first with SHA-256 and then with RIPEMD-160.
    - OP_HASH256: The input is hashed two times with SHA-256.
- [Difference between SegWit and Legacy address](https://help.crypto.com/en/articles/4056348-send-and-receive-btc-ltc-difference-between-segwit-and-legacy-address)
- [modood/btckeygen: A very simple and easy to use bitcoin(btc) key/wallet generator.](https://github.com/modood/btckeygen)

## 3. Generate n-out-of-m Multisig P2SH Bitcoin Address

How does P2SH works:

- fund to
    - redeem hash
- spend with
    - redeem script
    - signature script (sigScript)

### Format of the example script

- Locking Script `HASH160 <20-byte hash of redeem script> EQUAL`
- Unlocking Script `0 Sig1 Sig2 <redeem script>`
    - Redeem Script `2 PubKey1 PubKey2 PubKey3 PubKey4 PubKey5 5 CHECKMULTISIG`

#### References

- [Pay to script hash - Bitcoin Wiki](https://en.bitcoin.it/wiki/Pay_to_script_hash)
    - [bips/bip-0013.mediawiki at master ?? bitcoin/bips](https://github.com/bitcoin/bips/blob/master/bip-0013.mediawiki)
    - Pay to Script Hash [bips/bip-0016](https://github.com/bitcoin/bips/blob/master/bip-0016.mediawiki)
- [Multi-signature - Bitcoin Wiki](https://en.bitcoin.it/wiki/Multi-signature)
- [Create Raw Multi-Sig P2SH Bitcoin Transaction in Golang](https://medium.com/coinmonks/build-p2sh-address-and-spend-its-fund-in-golang-1a03a4131512)