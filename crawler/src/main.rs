use anyhow::Result;
use crawler::crawlers::{ContestCrawler, ProblemCrawler};
use crawler::models::*;
use dotenvy::dotenv;
use tracing_subscriber::filter::EnvFilter;
use tracing_subscriber::fmt;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let log_level = std::env::var("RUST_LOG").unwrap_or(String::from("info"));
    std::env::set_var("RUST_LOG", "info");

    let filter =
        EnvFilter::try_from_default_env()?.add_directive(format!("crawler={}", log_level).parse()?);
    fmt().with_env_filter(filter).init();

    let database_url: String =
        std::env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await?;

    tracing::info!("Start to crawl contest information.");
    let crawler = ContestCrawler::new(&pool);
    let contests: Vec<Contest> = crawler.crawl().await.expect("Failed to get contests.");
    tracing::info!("Finish to crawl contest information.");
    tracing::debug!("{} contests acquired.", contests.len());

    tracing::info!("Start to save contest information.");
    crawler
        .save(&contests)
        .await
        .expect("Failed to save contests.");
    tracing::info!("Finish to save contest information.");

    tracing::info!("Start to crawl problem information.");
    let crawler = ProblemCrawler::new(&pool);
    let target: Vec<ProblemJson> = crawler.detect_diff().await?;
    tracing::debug!("{} problems are now target for collection.", target.len());

    let problems: Vec<Problem> = crawler
        .crawl(&target)
        .await
        .expect("Failed to get problems");
    tracing::info!("Finish to crawl problem information.");

    tracing::info!("Start to save problem information.");
    crawler
        .save(&problems)
        .await
        .expect("Failed to save problems");
    tracing::info!("Finish to save problem information.");

    Ok(())
}
