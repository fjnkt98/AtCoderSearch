use std::{
    sync::{Arc, OnceLock},
    time::Duration,
};

use anyhow::Context as _;
use regex::Regex;
use reqwest::{cookie::Jar, Client, Url};

static CSRF_PATTERN: OnceLock<Regex> = OnceLock::new();

fn extract_csrf_token(body: &str) -> Option<String> {
    let pattern = CSRF_PATTERN.get_or_init(|| Regex::new(r#"var csrfToken = "(.+)""#).unwrap());

    pattern
        .captures(body)
        .and_then(|caps| caps.get(0))
        .and_then(|m| Some(m.as_str().to_owned()))
}

pub struct AtCoderClient {
    client: Client,
}

pub struct User {
    user_id: String,
    rating: i64,
    highest_rating: i64,
    affiliation: Option<String>,
    birth_year: Option<String>,
    country: Option<String>,
    crown: Option<String>,
    join_count: i64,
    rank: i64,
    active_rank: i64,
    wins: i64,
}

pub struct Submission {
    id: i64,
    epoch_second: i64,
    problem_id: String,
    contest_id: String,
    user_id: String,
    language: String,
    point: f64,
    length: i64,
    result: String,
    execution_time: Option<i64>,
}

impl AtCoderClient {
    pub fn new() -> anyhow::Result<Self> {
        let jar = Arc::new(Jar::default());
        let client = Client::builder()
            .cookie_provider(jar)
            .timeout(Duration::from_secs(30))
            .build()
            .with_context(|| "create http client")?;

        return Ok(Self { client });
    }

    pub async fn login(
        &self,
        username: impl AsRef<str>,
        password: impl AsRef<str>,
    ) -> anyhow::Result<()> {
        todo!()
    }

    pub async fn fetch_problem_html(
        &self,
        contest_id: impl AsRef<str>,
        problem_id: impl AsRef<str>,
    ) -> anyhow::Result<String> {
        todo!();
    }

    pub async fn fetch_users(&self, page: i64) -> anyhow::Result<Vec<User>> {
        todo!();
    }

    pub async fn fetch_submissions(
        &self,
        contest_id: impl AsRef<str>,
        page: i64,
    ) -> anyhow::Result<Vec<Submission>> {
        todo!();
    }
}
