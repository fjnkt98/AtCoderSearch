use std::time::Duration;

use crate::{
    atcoder::{AtCoderClient, AtCoderProblemsClient},
    crawl::{contest, difficulty, problem, submission, user},
};
use anyhow::Context;
use clap::Subcommand;
use itertools::Itertools;
use sqlx::postgres::PgPoolOptions;

use super::CommonArgs;

#[derive(Subcommand)]
pub(crate) enum CrawlCommand {
    Problem {
        #[arg(long, default_value = "1s")]
        duration: humantime::Duration,
        #[arg(short, long, help = "if true, crawl all problems.")]
        all: bool,
    },
    User {
        #[arg(long, default_value_t = 1)]
        duration: u64,
    },
    Submission {
        #[arg(long, default_value = "1s")]
        duration: humantime::Duration,
        #[arg(
            long,
            default_value_t = 0,
            help = "number of retries allowed when crawling failed. zero means no retry."
        )]
        retry: i64,
        #[arg(
            long,
            default_value = "",
            help = "comma separated list of target contest id. if not specified, crawl all contests."
        )]
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
            .connect(&args.database_url)
            .await
            .with_context(|| "connect to database")?;

        match self {
            Self::Problem { duration, all } => {
                contest::crawl_contests(&problems_client, &pool)
                    .await
                    .with_context(|| "crawl contests")?;

                difficulty::crawl_difficulties(&problems_client, &pool)
                    .await
                    .with_context(|| "crawl difficulties")?;

                problem::crawl_problems(
                    &atcoder_client,
                    &problems_client,
                    &pool,
                    *all,
                    (*duration).into(),
                )
                .await
                .with_context(|| "crawl problems")?;

                Ok(())
            }
            Self::User { duration } => {
                user::crawl_users(&atcoder_client, &pool, Duration::from_secs(*duration))
                    .await
                    .with_context(|| "crawl users")?;

                Ok(())
            }
            Self::Submission {
                duration,
                retry,
                target,
                atcoder_username,
                atcoder_password,
            } => {
                atcoder_client
                    .login(atcoder_username, atcoder_password)
                    .await
                    .with_context(|| "login to atcoder")?;

                let targets = target.split(",").map(String::from).collect_vec();
                let crawler = submission::SubmissionCrawler::new(
                    atcoder_client,
                    pool,
                    (*duration).into(),
                    *retry,
                    targets,
                );

                crawler.crawl().await?;
                Ok(())
            }
        }
    }
}
