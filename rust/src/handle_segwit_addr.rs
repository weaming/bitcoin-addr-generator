use bitcoin::network::constants::Network;
use bitcoin::util::{
    address::Address as Addr,
    bip32::{ChildNumber, ExtendedPubKey},
};
use bitcoin_wallet::account::{
    AccountAddressType, MasterAccount, MasterKeyEntropy, Seed, Unlocker,
};
use bitcoin_wallet::mnemonic::Mnemonic;
use hex;
use regex::Regex;
use std::str::FromStr;
use tide::log;
use tide::prelude::*;
use tide::{Request, Response};

lazy_static! {
    // coin types include 2 to allow bad test case
    static ref BIP44_PATH: Regex = Regex::new(r"^m/(44|49|84)'/([012])'/(\d+)'/([01])/(\d+)$").unwrap();
}

#[derive(Debug, Deserialize)]
struct RequestHDSegwitAddr {
    #[serde(default)]
    mnemonic: String,

    #[serde(default)]
    passphrase: String,

    path: String,
}

#[derive(Debug, Serialize)]
struct AddressRoot {
    mnemonic: String,

    #[serde(default)]
    passphrase: String,

    seed: String,
    root_priv_key: String,
    root_pub_key: String,
}

#[derive(Debug, Serialize)]
struct Address {
    bip: String,
    path: String,
    wif: String,

    #[serde(skip_serializing_if = "Option::is_none")]
    address: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    segwit_nested: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    segwit_bech32: Option<String>,
}

pub async fn handle_segwit_addr(mut req: Request<()>) -> tide::Result {
    let RequestHDSegwitAddr {
        mut mnemonic,
        passphrase,
        path,
    } = req.body_json().await?;
    if mnemonic == "" {
        if let Some(_) = option_env!("TEST") {
            mnemonic = "scale unfair later desert panda boost clap van census advice liar bomb manual subway cruise swing virtual access pig topple midnight double vague expect".to_owned();
        }
    }

    if let Some(indexes) = BIP44_PATH.captures(&path) {
        log::info!("indexes {:?}", indexes);

        let coin_type = u32::from_str(&indexes[2]).unwrap();
        let account = u32::from_str(&indexes[3]).unwrap();
        let change = u32::from_str(&indexes[4]).unwrap();
        let index = u32::from_str(&indexes[5]).unwrap();

        let _mnemonic = {
            if mnemonic == "" {
                Mnemonic::new_random(MasterKeyEntropy::Paranoid).unwrap()
            } else {
                Mnemonic::from_str(&mnemonic).unwrap()
            }
        };
        let mnemonic = _mnemonic.to_string();
        let seed: Seed = _mnemonic.to_seed(None);

        let network = match coin_type {
            0 => Network::Bitcoin,
            1 => Network::Testnet,
            // 2 is allowed in BIP44_PATH regexp
            _ => {
                return Ok(Response::from(json!({
                    "error": format!("invalid coin type: {}", coin_type)
                })));
            }
        };
        let master = MasterAccount::from_seed(&seed, 0, network, &passphrase).unwrap();

        // create an unlocker that is able to decrypt the encrypted mnemonic and then calculate private keys
        let mut unlocker = Unlocker::new_for_master(&master, &passphrase).unwrap();

        let address_root = AddressRoot {
            mnemonic: mnemonic,
            passphrase: passphrase,
            seed: hex::encode(&seed.0),
            root_priv_key: unlocker.master_private().to_string(),
            root_pub_key: master.master_public().to_string(),
        };

        if let Ok(purpose) = u32::from_str(&indexes[1]) {
            match purpose {
                44u32 => {
                    let address_type = AccountAddressType::P2PKH; // pay-to-public-key-hash (legacy)
                    let by_change = unlocker.sub_account_key(address_type, account, change)?;
                    let ext_priv_key = unlocker
                        .context()
                        .private_child(&by_change, ChildNumber::Normal { index })
                        .unwrap();
                    let ext_pub_key: ExtendedPubKey = unlocker
                        .context()
                        .extended_public_from_private(&ext_priv_key);
                    let ext_pub_addr = Addr::p2pkh(&ext_pub_key.public_key, ext_pub_key.network);

                    let address = Address {
                        bip: "BIP44".to_string(),
                        path: path,
                        wif: ext_priv_key.private_key.to_wif(),
                        address: Some(ext_pub_addr.to_string()),
                        segwit_nested: None,
                        segwit_bech32: None,
                    };
                    return Ok(Response::from(json!({
                        "AddressRoot": address_root,
                        "Address": address,
                    })));
                }
                49u32 => {
                    let address_type = AccountAddressType::P2SHWPKH; // pay-to-script-hash-witness-public-key-hash (transitional single key segwit)
                    let by_change = unlocker.sub_account_key(address_type, account, change)?;
                    let ext_priv_key = unlocker
                        .context()
                        .private_child(&by_change, ChildNumber::Normal { index })
                        .unwrap();
                    let ext_pub_key: ExtendedPubKey = unlocker
                        .context()
                        .extended_public_from_private(&ext_priv_key);
                    let ext_pub_addr = Addr::p2shwpkh(&ext_pub_key.public_key, ext_pub_key.network);

                    let address = Address {
                        bip: "BIP49".to_string(),
                        path: path,
                        wif: ext_priv_key.private_key.to_wif(),
                        address: None,
                        segwit_nested: Some(ext_pub_addr.to_string()),
                        segwit_bech32: None,
                    };
                    return Ok(Response::from(json!({
                        "AddressRoot": address_root,
                        "Address": address,
                    })));
                }
                84u32 => {
                    let address_type = AccountAddressType::P2WPKH; // pay-to-witness-public-key-hash (native single key segwit)
                    let by_change = unlocker.sub_account_key(address_type, account, change)?;
                    let ext_priv_key = unlocker
                        .context()
                        .private_child(&by_change, ChildNumber::Normal { index })
                        .unwrap();
                    let ext_pub_key: ExtendedPubKey = unlocker
                        .context()
                        .extended_public_from_private(&ext_priv_key);
                    let ext_pub_addr = Addr::p2wpkh(&ext_pub_key.public_key, ext_pub_key.network);

                    let address = Address {
                        bip: "BIP84".to_string(),
                        path: path,
                        wif: ext_priv_key.private_key.to_wif(),
                        address: None,
                        segwit_nested: None,
                        segwit_bech32: Some(ext_pub_addr.to_string()),
                    };
                    return Ok(Response::from(json!({
                        "AddressRoot": address_root,
                        "Address": address,
                    })));
                }
                _ => {}
            }
        };
    } else {
        return Ok(Response::from(json!({
            "error": format!("invalid path: {}", path)
        })));
    }
    Ok(Response::from(json!(
        {
            "error": "should never happen"
        }
    )))
}
