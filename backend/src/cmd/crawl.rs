use std::time::Duration;

use crate::{
    atcoder::{AtCoderClient, AtCoderProblemsClient},
    crawl::{contest, difficulty, problem::ProblemCrawler},
};
use clap::Subcommand;
use sqlx::postgres::PgPoolOptions;
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
        let problems_client = AtCoderProblemsClient::new()?;
        let pool = PgPoolOptions::new()
            .max_connections(8)
            .connect(&args.database_url)
            .await?;

        match self {
            Self::Problem { duration, all } => {
                contest::crawl_contests(&problems_client, &pool).await?;
                difficulty::crawl_difficulties(&problems_client, &pool).await?;

                let atcoder_client = AtCoderClient::new()?;
                let crawler = ProblemCrawler::new(
                    &pool,
                    &atcoder_client,
                    &problems_client,
                    Duration::from_secs(*duration as u64),
                );
                crawler.crawl(token.clone(), *all).await?;

                Ok(())
            }
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
