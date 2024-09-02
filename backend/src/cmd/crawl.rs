use std::time::Duration;

use crate::{
    atcoder::{AtCoderClient, AtCoderProblemsClient},
    crawl::{contest, difficulty, problem, user},
};
use clap::Subcommand;
use sqlx::postgres::PgPoolOptions;

use super::CommonArgs;

#[derive(Subcommand)]
pub(crate) enum CrawlCommand {
    Problem {
        #[arg(long, default_value_t = 1, help = "crawl duration in sec.")]
        duration: u64,
        #[arg(short, long, help = "if true, crawl all problems.")]
        all: bool,
    },
    User {
        #[arg(long, default_value_t = 1)]
        duration: u64,
    },
    Submission {
        #[arg(long, default_value_t = 1)]
        duration: u64,
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
    pub async fn exec(&self, args: &CommonArgs) -> anyhow::Result<()> {
        let problems_client = AtCoderProblemsClient::new()?;
        let atcoder_client = AtCoderClient::new()?;
        let pool = PgPoolOptions::new()
            .max_connections(8)
            .connect(&args.database_url)
            .await?;

        match self {
            Self::Problem { duration, all } => {
                contest::crawl_contests(&problems_client, &pool).await?;
                difficulty::crawl_difficulties(&problems_client, &pool).await?;
                problem::crawl_problems(
                    &atcoder_client,
                    &problems_client,
                    &pool,
                    *all,
                    Duration::from_secs(*duration),
                )
                .await?;

                Ok(())
            }
            Self::User { duration } => {
                user::crawl_users(&atcoder_client, &pool, Duration::from_secs(*duration)).await
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
