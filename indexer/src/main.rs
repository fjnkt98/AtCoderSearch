use anyhow::{Context, Result};
use chrono::{DateTime, NaiveDateTime, TimeZone, Utc};
use dotenvy::dotenv;
use ego_tree::NodeRef;
use indexer::solr;
use regex::Regex;
use scraper::node::Node;
use scraper::{Html, Selector};
use serde::{Deserialize, Serialize};
use sqlx::postgres::Postgres;
use sqlx::Pool;

#[derive(Serialize, Deserialize, Debug)]
struct Document {
    problem_id: String,
    problem_title: String,
    problem_url: String,
    contest_id: String,
    contest_title: String,
    contest_url: String,
    difficulty: i32,
    start_at: String,
    duration: i32,
    rate_change: String,
    category: String,
    text_ja: Vec<String>,
    text_en: Vec<String>,
}

fn dfs(element: &NodeRef<Node>) -> String {
    let mut result = Vec::new();

    for child in element.children() {
        match child.value() {
            Node::Element(_) => {
                result.push(dfs(&child));
            }
            Node::Text(text) => {
                result.push(text.trim().to_string());
            }
            _ => {
                continue;
            }
        };
    }

    result.join(" ")
}

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

    let row = sqlx::query!(
        "
        SELECT
            problems.id AS problem_id,
            problems.title AS problem_title,
            problems.url AS problem_url,
            contests.id AS contest_id,
            contests.title AS contest_title,
            problems.difficulty AS difficulty,
            contests.start_epoch_second AS start_at,
            contests.duration_second AS duration,
            contests.rate_change AS rate_change,
            contests.category AS category,
            problems.html AS html
        FROM
            problems
        LEFT JOIN contests
        ON problems.contest_id = contests.id
        WHERE problems.id = 'abc280_a';
        ",
    )
    .fetch_one(&pool)
    .await?;

    let start_at = DateTime::<Utc>::from_utc(
        NaiveDateTime::from_timestamp_opt(row.start_at.unwrap(), 0).unwrap(),
        Utc,
    )
    .to_rfc3339()
    .replace("+00:00", "Z");
    println!("{}", start_at);

    let html = Html::parse_document(&row.html);

    let div = Selector::parse("div.part").unwrap();
    let section = Selector::parse("section").unwrap();
    let h3 = Selector::parse("h3").unwrap();

    let ascii = Regex::new(r"^[a-zA-Z0-9 ]*$").unwrap();

    let mut full_text_ja: Vec<String> = Vec::new();
    let mut full_text_en: Vec<String> = Vec::new();
    for part in html.select(&div) {
        let section = part.select(&section).next().unwrap();

        let title = section.select(&h3).next().unwrap().text().next().unwrap();

        let mut full_text: Vec<String> = Vec::new();
        for e in section.children() {
            match e.value() {
                Node::Element(element) => {
                    if element.name() == "h3" {
                        continue;
                    } else {
                        full_text.push(dfs(&e));
                    }
                }
                Node::Text(text) => full_text.push(text.to_string()),
                _ => {
                    continue;
                }
            }
        }

        if ascii.is_match(&title) {
            full_text_en.push(full_text.join(""));
        } else {
            full_text_ja.push(full_text.join(" "));
        }
    }

    let contest_url: String = format!(
        "https://atcoder.jp/contests/{}",
        row.contest_id.clone().unwrap()
    );
    let document = Document {
        problem_id: row.problem_id,
        problem_title: row.problem_title,
        problem_url: row.problem_url,
        contest_id: row.contest_id.unwrap(),
        contest_title: row.contest_title.unwrap(),
        contest_url: contest_url,
        difficulty: row.difficulty.unwrap(),
        start_at: start_at,
        duration: row.duration.unwrap() as i32,
        rate_change: row.rate_change.unwrap(),
        category: row.category.unwrap(),
        text_ja: full_text_ja,
        text_en: full_text_en,
    };

    // println!("{}", serde_json::to_string(&document).unwrap());
    // println!("{}, {}", document.problem_id, document.start_at);

    let solr = solr::client::SolrClient::new("http://localhost", 8983)
        .context("Failed to create solr client.")?;

    let cores = solr.cores().await?;
    println!("{:?}", cores);

    Ok(())
}
