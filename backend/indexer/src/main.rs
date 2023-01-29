// mod extractor;
// mod manager;
// mod models;
// mod reader;

use std::ffi::OsString;

use anyhow::Result;
use clap::{Args, Parser, Subcommand};
use dotenvy::dotenv;
// use manager::IndexingManager;
// use sqlx::postgres::Postgres;
// use sqlx::Pool;
// use std::env;
// use tracing_subscriber::filter::EnvFilter;
// use tracing_subscriber::fmt;

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
    #[arg(short, long)]
    optimize: bool,
}

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().ok();

    let args = Cli::parse();
    match args.command {
        Commands::Crawl(args) => {
            if args.all {
                println!("crawl all!");
            } else {
                println!("crawl diff!");
            }
        }
        Commands::Generate(args) => match args.path {
            Some(path) => println!("generate json into {}!", path.into_string().unwrap()),
            None => println!("generate json into default path!"),
        },
        Commands::Post(args) => {
            if args.optimize {
                println!("post document with optimize!");
            } else {
                println!("post document without optimize!");
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
