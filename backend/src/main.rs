mod atcoder;
mod cmd;
mod crawl;

use clap::Parser;
use cmd::App;

#[tokio::main]
async fn main() {
    let app = App::parse();

    if let Err(err) = app.run().await {
        println!("command failed: {:#}", err)
    }
}
