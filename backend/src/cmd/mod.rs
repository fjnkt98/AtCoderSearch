mod crawl;
mod update;

use crate::cmd::{crawl::CrawlCommand, update::UpdateCommands};
use clap::{Args, Parser, Subcommand};
use tokio_util::sync::CancellationToken;

#[derive(Parser)]
pub struct App {
    #[command(subcommand)]
    command: Command,
}

#[derive(Args)]
pub struct CommonArgs {
    #[arg(long, env, hide_env_values = true)]
    database_url: String,
    #[arg(long, env, hide_env_values = true)]
    engine_url: String,
}

#[derive(Subcommand)]
enum Command {
    Crawl {
        #[command(flatten)]
        args: CommonArgs,
        #[command(subcommand)]
        command: CrawlCommand,
    },
    Update {
        #[command(flatten)]
        args: CommonArgs,
        #[command(subcommand)]
        command: UpdateCommands,
    },
    Serve {
        #[command(flatten)]
        args: CommonArgs,
        #[arg(long, env)]
        port: i64,
    },
}

impl App {
    #[tokio::main]
    pub async fn run(&self) -> anyhow::Result<()> {
        let token = CancellationToken::new();
        let token2 = token.clone();

        tokio::spawn(async move {
            tokio::signal::ctrl_c().await.unwrap();

            token.cancel();
        });

        match &self.command {
            Command::Crawl { args, command } => command.exec(token2, args).await,
            Command::Update { args, command } => command.exec(token2, args),
            Command::Serve { args, port } => {
                todo!();
            }
        }
    }
}
