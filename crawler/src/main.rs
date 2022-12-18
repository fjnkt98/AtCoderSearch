mod contest;
mod problem;

use anyhow::Result;
use contest::crawler::ContestCrawler;
use contest::models::Contest;
use dotenvy::dotenv;
use problem::crawler::ProblemCrawler;
use problem::models::{Problem, ProblemJson};

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let database_url: String =
        std::env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await?;

    let crawler = ContestCrawler::new(&pool);
    let contests: Vec<Contest> = crawler.crawl().await.expect("Failed to get contests.");

    crawler
        .save(&contests)
        .await
        .expect("Failed to save contests.");

    let crawler = ProblemCrawler::new(&pool);
    let target: Vec<ProblemJson> = crawler.detect_diff().await?;
    println!("{} Unknown problems detected", target.len());
    let problems: Vec<Problem> = crawler
        .crawl(&target)
        .await
        .expect("Failed to get problems");

    crawler
        .save(&problems)
        .await
        .expect("Failed to save problems");

    Ok(())
}
