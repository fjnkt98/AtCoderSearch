use anyhow::Result;
use dotenvy::dotenv;
use indexer::utils::manager::IndexingManager;
use sqlx::postgres::Postgres;
use sqlx::Pool;

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let log_level = std::env::var("RUST_LOG").unwrap_or(String::from("info"));
    std::env::set_var("RUST_LOG", log_level);
    tracing_subscriber::fmt::init();

    let database_url: String =
        std::env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    let pool: Pool<Postgres> = sqlx::postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await?;

    let manager = IndexingManager::new(&pool);
    manager
        .write()
        .await
        .expect("Failed to write JSON document.");

    manager
        .post()
        .await
        .expect("Failed to post document to solr.");

    Ok(())
}
