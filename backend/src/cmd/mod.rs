mod crawl;
mod update;

use crate::cmd::{crawl::CrawlCommand, update::UpdateCommands};
use clap::{Args, Parser, Subcommand};

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
    #[arg(long, env, hide_env_values = true)]
    engine_master_key: String,
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
}

impl App {
    #[tokio::main]
    pub async fn run(&self) -> anyhow::Result<()> {
        match &self.command {
            Command::Crawl { args, command } => command.exec(args).await,
            Command::Update { args, command } => command.exec(args).await,
        }
    }
}
