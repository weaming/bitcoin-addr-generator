#![deny(unused_imports)]

use tide::log;
#[macro_use]
extern crate lazy_static;

mod handle_segwit_addr;
use handle_segwit_addr::handle_segwit_addr;
mod handle_multisig_p2sh_addr;
use handle_multisig_p2sh_addr::handle_multisig_p2sh_addr;

#[async_std::main]
async fn main() -> tide::Result<()> {
    log::start();

    let mut app = tide::new();
    app.at("/api/hd-segwit-address").post(handle_segwit_addr);
    app.at("/api/multisig-p2sh-address")
        .post(handle_multisig_p2sh_addr);

    app.listen("0.0.0.0:8080").await?;
    Ok(())
}
