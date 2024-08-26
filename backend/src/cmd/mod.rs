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

#[derive(Subcommand)]
enum CrawlCommand {
    Problem {
        #[arg(long, default_value_t = 1, help = "crawl duration in sec.")]
        duration: i64,
        #[arg(short, long, help = "if true, crawl all problems.")]
        all: bool,
    },
    User {
        #[arg(long, default_value_t = 1)]
        duration: i64,
    },
    Submission {
        #[arg(long, default_value_t = 1)]
        duration: i64,
        #[arg(long, default_value_t = 0)]
        retry: i64,
        #[arg(long)]
        target: String,
        #[arg(long, env, hide_env_values = true)]
        atcoder_username: String,
        #[arg(long, env, hide_env_values = true)]
        atcoder_password: String,
    },
}

#[derive(Subcommand)]
enum UpdateCommands {
    Problem,
    User,
}

impl App {
    pub async fn run(&self) -> anyhow::Result<()> {
        let token = CancellationToken::new();

        match &self.command {
            Command::Crawl { args, command } => match command {
                CrawlCommand::Problem { duration, all } => {
                    println!("crawl problem with duration: {}, all: {}", duration, all)
                }
                CrawlCommand::User { duration } => {
                    println!("crawl user with duration: {}", duration);
                }
                CrawlCommand::Submission {
                    duration,
                    retry,
                    target,
                    atcoder_username,
                    atcoder_password,
                } => {
                    println!("crawl submission with duration: {}", duration);
                }
            },
            Command::Update { args, command } => {}
            Command::Serve { args, port } => {}
        }

        Ok(())
    }
}
