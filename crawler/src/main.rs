mod contest;
mod problem;

use anyhow::Result;
use contest::crawler::ContestCrawler;
use contest::inserter;
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

    let contest_crawler = ContestCrawler::new();
    let contests: Vec<Contest> = contest_crawler
        .crawl()
        .await
        .expect("Failed to get contests.");

    inserter::insert(&pool, &contests).await?;

    let crawler = ProblemCrawler::new(&pool);
    let target: Vec<ProblemJson> = crawler.get_problem_list().await?[..5].to_vec();
    let problems: Vec<Problem> = crawler
        .crawl(&target)
        .await
        .expect("Failed to get problems");

    crawler.save(&problems).await?;

    Ok(())
}
