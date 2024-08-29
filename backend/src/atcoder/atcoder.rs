use std::{
    borrow::Cow,
    collections::HashMap,
    sync::{Arc, LazyLock},
    time::Duration,
};

use anyhow::Context as _;
use regex::Regex;
use reqwest::{cookie::Jar, header::CONTENT_TYPE, Client, Url};

static CSRF_PATTERN: LazyLock<Regex> =
    LazyLock::new(|| Regex::new(r#"var csrfToken = "(.+)""#).unwrap());

static RANK_PATTERN: LazyLock<Regex> = LazyLock::new(|| Regex::new(r#"\((\d+)\)"#).unwrap());

pub struct AtCoderClient {
    client: Client,
}

impl AtCoderClient {
    pub fn new() -> anyhow::Result<Self> {
        let client = Client::builder()
            .cookie_store(true)
            .timeout(Duration::from_secs(30))
            .build()
            .with_context(|| "create http client")?;

        return Ok(Self { client });
    }

    pub async fn login<'a>(&self, username: &'a str, password: &'a str) -> anyhow::Result<()> {
        let url = Url::parse("https://atcoder.jp/login").with_context(|| "parse login url")?;

        let res = self
            .client
            .get(url.clone())
            .send()
            .await
            .with_context(|| "request to login url")?;
        let body = res.text().await.with_context(|| "get response body")?;

        let token = extract_csrf_token(&body).with_context(|| "extract csrf token")?;

        let form = HashMap::from([
            ("username", username),
            ("password", password),
            ("csrf_token", &token),
        ]);
        let res = self
            .client
            .post(url)
            .form(&form)
            .send()
            .await
            .with_context(|| "request to login")?;

        match res.error_for_status() {
            Ok(_) => Ok(()),
            Err(e) => Err(anyhow::anyhow!(
                "login authentication failed with status: {:?}",
                e.status()
            )),
        }
    }

    pub async fn fetch_problem_html<'a>(
        &self,
        contest_id: &'a str,
        problem_id: &'a str,
    ) -> anyhow::Result<String> {
        let url = Url::parse(
            format!(
                "https://atcoder.jp/contests/{}/tasks/{}",
                contest_id, problem_id
            )
            .as_str(),
        )
        .with_context(|| "parse problem url")?;

        let res = self
            .client
            .get(url)
            .send()
            .await
            .with_context(|| "request to get problem html")?;

        let body = res.text().await.with_context(|| "get response body")?;

        Ok(body)
    }

    pub async fn fetch_users(&self, page: i64) -> anyhow::Result<Vec<User>> {
        let url =
            Url::parse("https://atcoder.jp/ranking/all").with_context(|| "parse users url")?;

        let params = HashMap::from([
            ("contestType", String::from("algo")),
            ("page", page.to_string()),
        ]);

        let res = self
            .client
            .get(url)
            .query(&params)
            .send()
            .await
            .with_context(|| "request to fetch users")?;

        let body = res.text().await.with_context(|| "get response body")?;

        let users = scrape_users(&body)?;

        Ok(users)
    }

    pub async fn fetch_submissions<'a>(
        &self,
        contest_id: &'a str,
        page: i64,
    ) -> anyhow::Result<Vec<Submission>> {
        let url =
            Url::parse(format!("https://atcoder.jp/contests/{}/submissions", contest_id).as_str())
                .with_context(|| "parse submissions url")?;

        let res = self
            .client
            .get(url)
            .query(&[("page", page)])
            .send()
            .await
            .with_context(|| "request to fetch submissions")?;

        let body = res.text().await.with_context(|| "get response body")?;

        let submissions = scrape_submissions(&body)?;

        Ok(submissions)
    }
}

fn extract_csrf_token(body: &str) -> Option<String> {
    CSRF_PATTERN
        .captures(body)
        .and_then(|caps| caps.get(1))
        .and_then(|m| Some(m.as_str().to_owned()))
}

fn scrape_submissions<'a>(html: &'a str) -> anyhow::Result<Vec<Submission>> {
    todo!()
}

fn scrape_users<'a>(html: &'a str) -> anyhow::Result<Vec<User>> {
    todo!()
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct User {
    pub user_id: String,
    pub rating: i64,
    pub highest_rating: i64,
    pub affiliation: Option<String>,
    pub birth_year: Option<i64>,
    pub country: Option<String>,
    pub crown: Option<String>,
    pub join_count: i64,
    pub rank: i64,
    pub active_rank: Option<i64>,
    pub wins: i64,
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct Submission {
    pub id: i64,
    pub epoch_second: i64,
    pub problem_id: String,
    pub contest_id: String,
    pub user_id: String,
    pub language: String,
    pub point: f64,
    pub length: i64,
    pub result: String,
    pub execution_time: Option<i64>,
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::PathBuf;

    use super::*;
    use rstest::rstest;

    #[rstest]
    fn test_extract_csrf_token(#[files("testdata/atcoder/login.html")] path: PathBuf) {
        let html = fs::read_to_string(path).unwrap();
        let token = extract_csrf_token(&html).unwrap();

        assert_eq!(token, "KrVShPadRMxPBKM9LmjWJHaQvjC7ALXz6DXgHOCL1LQ=");
    }

    #[rstest]
    fn test_scrape_submissions(#[files("testdata/atcoder/submissions.html")] path: PathBuf) {
        let html = fs::read_to_string(path).unwrap();

        let want = vec![
            Submission {
                id: 48852107,
                epoch_second: 1703553569,
                problem_id: String::from("abc300_a"),
                user_id: String::from("Orkhon2010"),
                contest_id: String::from("abc300"),
                language: String::from("C++ 20 (gcc 12.2)"),
                point: 100.0,
                length: 259,
                result: String::from("AC"),
                execution_time: Some(1),
            },
            Submission {
                id: 48852073,
                epoch_second: 1703553403,
                problem_id: String::from("abc300_f"),
                user_id: String::from("ecsmtlir"),
                contest_id: String::from("abc300"),
                language: String::from("C++ 20 (gcc 12.2)"),
                point: 500.0,
                length: 14721,
                result: String::from("AC"),
                execution_time: Some(11),
            },
        ];

        let actual = scrape_submissions(&html).unwrap();

        assert_eq!(want, actual)
    }

    #[rstest]
    fn test_scrape_users(#[files("testdata/atcoder/users.html")] path: PathBuf) {
        let html = fs::read_to_string(path).unwrap();

        let want = vec![
            User {
                user_id: String::from("tourist"),
                rating: 3863,
                highest_rating: 4229,
                affiliation: Some(String::from("ITMO University")),
                birth_year: Some(1994),
                country: Some(String::from("BY")),
                crown: Some(String::from("crown_champion")),
                join_count: 59,
                rank: 1,
                active_rank: Some(1),
                wins: 22,
            },
            User {
                user_id: String::from("w4yneb0t"),
                rating: 3710,
                highest_rating: 3802,
                affiliation: Some(String::from("ETH Zurich")),
                birth_year: None,
                country: Some(String::from("CH")),
                crown: None,
                join_count: 21,
                rank: 2,
                active_rank: None,
                wins: 2,
            },
            User {
                user_id: String::from("ksun48"),
                rating: 3681,
                highest_rating: 3802,
                affiliation: Some(String::from("MIT")),
                birth_year: Some(1998),
                country: Some(String::from("CA")),
                crown: Some(String::from("crown_gold")),
                join_count: 58,
                rank: 3,
                active_rank: Some(2),
                wins: 5,
            },
            User {
                user_id: String::from("ecnerwala"),
                rating: 3663,
                highest_rating: 3814,
                affiliation: Some(String::from("MIT")),
                birth_year: Some(1997),
                country: Some(String::from("US")),
                crown: Some(String::from("crown_gold")),
                join_count: 36,
                rank: 4,
                active_rank: Some(3),
                wins: 2,
            },
            User {
                user_id: String::from("Benq"),
                rating: 3658,
                highest_rating: 3683,
                affiliation: Some(String::from("MIT")),
                birth_year: Some(2001),
                country: Some(String::from("US")),
                crown: None,
                join_count: 48,
                rank: 5,
                active_rank: None,
                wins: 0,
            },
            User {
                user_id: String::from("cospleermusora"),
                rating: 3606,
                highest_rating: 3783,
                affiliation: None,
                birth_year: None,
                country: Some(String::from("RU")),
                crown: None,
                join_count: 25,
                rank: 5,
                active_rank: None,
                wins: 3,
            },
            User {
                user_id: String::from("apiad"),
                rating: 3600,
                highest_rating: 3852,
                affiliation: None,
                birth_year: Some(1997),
                country: Some(String::from("CN")),
                crown: Some(String::from("crown_gold")),
                join_count: 51,
                rank: 7,
                active_rank: Some(4),
                wins: 6,
            },
            User {
                user_id: String::from("Um_nik"),
                rating: 3571,
                highest_rating: 3948,
                affiliation: None,
                birth_year: Some(1996),
                country: Some(String::from("UA")),
                crown: Some(String::from("crown_gold")),
                join_count: 60,
                rank: 8,
                active_rank: Some(5),
                wins: 7,
            },
            User {
                user_id: String::from("mnbvmar"),
                rating: 3555,
                highest_rating: 3736,
                affiliation: Some(String::from("University of Warsaw")),
                birth_year: Some(1996),
                country: Some(String::from("PL")),
                crown: Some(String::from("crown_gold")),
                join_count: 22,
                rank: 9,
                active_rank: Some(6),
                wins: 1,
            },
            User {
                user_id: String::from("Stonefeang"),
                rating: 3554,
                highest_rating: 3658,
                affiliation: Some(String::from("University of Warsaw")),
                birth_year: Some(1997),
                country: Some(String::from("PL")),
                crown: Some(String::from("crown_gold")),
                join_count: 37,
                rank: 10,
                active_rank: Some(7),
                wins: 2,
            },
        ];

        let actual = scrape_users(&html).unwrap();
        assert_eq!(want, actual)
    }

    #[test]
    fn test_new_atcoder_client() {
        AtCoderClient::new().unwrap();
    }
}
