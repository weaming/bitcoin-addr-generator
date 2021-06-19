use bitcoin::blockdata::opcodes::{all as opcodes, All};
use bitcoin::blockdata::script::Builder;
use bitcoin::hashes::Hash;
use bitcoin::network::constants::Network;
use bitcoin::util::address::{Address, Payload};
use std::str::FromStr;
use tide::prelude::*;
use tide::{Request, Response};

#[derive(Debug, Deserialize)]
struct RequestMultisigP2SHAddr {
    pub_keys: Vec<String>,
    n: u32,
}

lazy_static! {
    static ref OP_ARR: Vec<All> = vec!(
        opcodes::OP_RESERVED, // do not use it
        opcodes::OP_PUSHNUM_1,
        opcodes::OP_PUSHNUM_2,
        opcodes::OP_PUSHNUM_3,
        opcodes::OP_PUSHNUM_4,
        opcodes::OP_PUSHNUM_5,
        opcodes::OP_PUSHNUM_6,
        opcodes::OP_PUSHNUM_7,
        opcodes::OP_PUSHNUM_8,
        opcodes::OP_PUSHNUM_9,
        opcodes::OP_PUSHNUM_10,
        opcodes::OP_PUSHNUM_11,
        opcodes::OP_PUSHNUM_12,
        opcodes::OP_PUSHNUM_13,
        opcodes::OP_PUSHNUM_14,
        opcodes::OP_PUSHNUM_15,
        opcodes::OP_PUSHNUM_16,
    );
}

pub async fn handle_multisig_p2sh_addr(mut req: Request<()>) -> tide::Result {
    let RequestMultisigP2SHAddr { pub_keys, n } = req.body_json().await?;
    let pk_len = pub_keys.len();
    if pk_len <= 1 {
        return Ok(Response::from(json!({
            "error": "need at least 2 public keys"
        })));
    }
    if pk_len > 16 {
        return Ok(Response::from(json!({
            "error": "at most 16 public keys"
        })));
    }
    if !(1 <= n && n <= 16) {
        return Ok(Response::from(json!({
            "error": "invalid n, it should be in range [2, 16]"
        })));
    }
    if n > pk_len as u32 {
        return Ok(Response::from(json!({
            "error": "n should not be greater than the number of public keys"
        })));
    }

    let mut bd = Builder::new().push_opcode(OP_ARR[n as usize]);
    for x in &pub_keys {
        match Address::from_str(x) {
            Ok(a) => {
                if let Payload::PubkeyHash(hash) = a.payload {
                    bd = bd.push_slice(&hash.into_inner());
                } else {
                    return Ok(Response::from(json!({
                        "error": format!("invalid public address: {}", x)
                    })));
                }
            }
            Err(e) => {
                return Ok(Response::from(json!({
                    "error": format!("invalid public address: {}: {}", x, e)
                })));
            }
        }
    }
    bd = bd
        .push_opcode(OP_ARR[pub_keys.len() as usize])
        .push_opcode(opcodes::OP_CHECKMULTISIG);
    let redeem_script = bd.into_script();
    let a = Address::p2sh(&redeem_script, Network::Bitcoin);

    Ok(Response::from(json!({
        "address": a,
    })))
}
