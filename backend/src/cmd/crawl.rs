use crate::{
    atcoder::AtCoderProblemsClient,
    crawl::contest::{self, crawl_contest},
};
use clap::Subcommand;
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};
use tokio_util::sync::CancellationToken;

use super::CommonArgs;

#[derive(Subcommand)]
pub(crate) enum CrawlCommand {
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

impl CrawlCommand {
    pub async fn exec(&self, token: CancellationToken, args: &CommonArgs) -> anyhow::Result<()> {
        let client = AtCoderProblemsClient::new()?;
        let pool = PgPoolOptions::new()
            .max_connections(8)
            .connect(&args.database_url)
            .await?;

        match self {
            Self::Problem { duration, all } => crawl_contest(&client, &pool).await,
            Self::User { duration } => {
                todo!();
            }
            Self::Submission {
                duration,
                retry,
                target,
                atcoder_username,
                atcoder_password,
            } => {
                todo!();
            }
        }
    }
}
