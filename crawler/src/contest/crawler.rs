use crate::contest::models::ContestJson;
use anyhow::{Context, Result};
use reqwest::header::ACCEPT_ENCODING;

pub struct ContestCrawler {
    pub url: String,
}

impl ContestCrawler {
    pub fn new() -> ContestCrawler {
        ContestCrawler {
            url: String::from("https://kenkoooo.com/atcoder/resources/contests.json"),
        }
    }

    pub fn get(self) -> Result<Vec<ContestJson>> {
        let client = reqwest::blocking::Client::new();

        let response = client
            .get(&self.url)
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .context("Failed to get contest information from AtCoder Problems.")?;

        let json = response
            .text()
            .context("Failed to get JSON body from response.")?;

        let contests: Vec<ContestJson> =
            serde_json::from_str(&json).context("Failed to parse JSON body.")?;

        Ok(contests)
    }
}
