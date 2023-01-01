use anyhow::Result;
use clap::Parser;
use crawler::crawlers::{ContestCrawler, ProblemCrawler};
use crawler::models::*;
use dotenvy::dotenv;
use tracing_subscriber::filter::EnvFilter;
use tracing_subscriber::fmt;

#[derive(Parser)]
#[clap(
    name = "AtCoder Search Crawler",
    author = "fjnkt98",
    version = "0.1.0",
    about = "Crawl contest and problem information",
    long_about = "Crawl contest information and problem data from AtCoder and AtCoder Problems."
)]
struct Args {
    #[clap(short, long)]
    all: bool,
    #[clap(short, long)]
    verbose: bool,
}

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let args = Args::parse();

    std::env::set_var("RUST_LOG", "info");

    let log_level = match args.verbose {
        true => "debug",
        false => "info",
    };

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
    let target: Vec<ProblemJson> = match args.all {
        true => crawler.get_problem_list().await?,
        false => crawler.detect_diff().await?,
    };
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
