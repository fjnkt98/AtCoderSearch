mod contest;

use anyhow::Result;
use contest::crawler;
use contest::inserter;
use contest::models::Contest;
use dotenvy::dotenv;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let database_url: String =
        std::env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(10)
        .connect(&database_url)
        .await?;

    let contest_crawler = crawler::ContestCrawler::new();
    let contests: Vec<Contest> = contest_crawler
        .crawl()
        .await
        .expect("Failed to get contests.");

    inserter::insert(&pool, &contests).await?;

    Ok(())
}
