mod models;
mod utils;

use std::ffi::OsString;

use crate::utils::crawlers::{ContestCrawler, ProblemCrawler};
use crate::utils::generator::DocumentGenerator;
use crate::utils::uploader::DocumentUploader;
use anyhow::{bail, Result};
use clap::{Args, Parser, Subcommand};
use dotenvy::dotenv;
use solrust::client::solr::SolrClient;
use sqlx::postgres::Postgres;
use sqlx::Pool;
use std::env;
use std::path::PathBuf;
use tracing_subscriber::filter::EnvFilter;
use tracing_subscriber::fmt;

#[derive(Debug, Parser)]
#[command(name = "indexer")]
#[command(
    about = "All you need to indexing",
    long_about = "Crawl and problems, generate document json, and post json into Solr core."
)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}
#[derive(Debug, Subcommand)]
enum Commands {
    Crawl(CrawlArgs),
    Generate(GenerateArgs),
    Post(PostArgs),
}

#[derive(Debug, Args)]
struct CrawlArgs {
    #[arg(long)]
    all: bool,
}

#[derive(Debug, Args)]
struct GenerateArgs {
    path: Option<OsString>,
}

#[derive(Debug, Args)]
struct PostArgs {
    path: Option<OsString>,
    #[arg(short, long)]
    optimize: bool,
}

#[tokio::main(flavor = "multi_thread", worker_threads = 10)]
async fn main() -> Result<()> {
    dotenv().ok();
    let args = Cli::parse();

    let log_level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    env::set_var("RUST_LOG", "info");

    let filter =
        EnvFilter::try_from_default_env()?.add_directive(format!("indexer={}", log_level).parse()?);
    fmt().with_env_filter(filter).init();

    let database_url: String = env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    let pool: Pool<Postgres> = sqlx::postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await?;

    let solr_host = env::var("SOLR_HOST").unwrap_or(String::from("http://localhost"));
    let solr_port = env::var("SOLR_PORT")
        .map(|v| v.parse::<u32>().unwrap())
        .unwrap_or(8983u32);

    let core_name = env::var("CORE_NAME").expect("CORE_NAME must be configured");

    let solr = SolrClient::new(&solr_host, solr_port).expect("Failed to create solr client.");
    let core = solr
        .core(&core_name)
        .await
        .expect("Failed to create core client");

    match args.command {
        Commands::Crawl(args) => {
            let crawler = ContestCrawler::new(&pool);
            if let Err(e) = crawler.run().await {
                tracing::error!(
                    "Failed to crawl and save contest information [{}]",
                    e.to_string()
                );
                bail!(e.to_string());
            }

            let crawler = ProblemCrawler::new(&pool);
            if let Err(e) = crawler.run(args.all).await {
                tracing::error!(
                    "Failed to crawl and save problem information [{}]",
                    e.to_string()
                );
                bail!(e.to_string());
            }
        }
        Commands::Generate(args) => {
            let savedir = match args.path {
                Some(path) => PathBuf::from(path),
                None => {
                    let path = env::var("DOCUMENT_SAVE_DIRECTORY")
                        .expect("Default save directory does not configured");
                    PathBuf::from(path)
                }
            };
            let generator = DocumentGenerator::new(&pool, &savedir);
            tracing::info!("Delete existing documents");
            if let Err(e) = generator.truncate().await {
                tracing::error!(
                    "Failed to delete existing json documents [{}]",
                    e.to_string()
                );
                bail!(e.to_string());
            }
            match generator.generate(1000).await {
                Ok(()) => {
                    tracing::info!("Successfully generate documents");
                }
                Err(e) => {
                    tracing::error!("Failed to generate documents");
                    bail!(e.to_string())
                }
            }
        }
        Commands::Post(args) => {
            let savedir = match args.path {
                Some(path) => PathBuf::from(path),
                None => {
                    let path = env::var("DOCUMENT_SAVE_DIRECTORY")
                        .expect("Default save directory does not configured");
                    PathBuf::from(path)
                }
            };
            let uploader = DocumentUploader::new(&savedir, &core);
            tracing::info!("Start to post documents");
            match uploader.upload(args.optimize).await {
                Ok(()) => {
                    tracing::info!("Successfully post documents")
                }
                Err(e) => {
                    tracing::error!("Failed to post documents");
                    bail!(e.to_string());
                }
            }
        }
    }

    // let log_level = env::var("RUST_LOG").unwrap_or(String::from("info"));
    // env::set_var("RUST_LOG", "info");

    // let filter =
    //     EnvFilter::try_from_default_env()?.add_directive(format!("indexer={}", log_level).parse()?);
    // fmt().with_env_filter(filter).init();

    // let database_url: String = env::var("DATABASE_URL").expect("DATABASE_URL must be configured.");

    // let pool: Pool<Postgres> = sqlx::postgres::PgPoolOptions::new()
    //     .max_connections(5)
    //     .connect(&database_url)
    //     .await?;

    // let solr_host = env::var("SOLR_HOST").unwrap_or(String::from("http://localhost"));
    // let solr_port = env::var("SOLR_PORT")
    //     .map(|v| v.parse::<u32>().unwrap())
    //     .unwrap_or(8983u32);

    // let core_name = env::var("CORE_NAME").expect("CORE_NAME must be configured");

    // let solr = SolrClient::new(&solr_host, solr_port).expect("Failed to create solr client.");
    // let core = solr
    //     .core(&core_name)
    //     .await
    //     .expect("Failed to create core client");

    // let manager = IndexingManager::new(&pool, core);
    // manager
    //     .write()
    //     .await
    //     .expect("Failed to write JSON document.");

    // manager
    //     .post()
    //     .await
    //     .expect("Failed to post document to solr.");

    Ok(())
}
