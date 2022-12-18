use crate::contest::models::{Contest, ContestJson};
use anyhow::{Context, Error, Result};
use reqwest::header::ACCEPT_ENCODING;
use sqlx::postgres::Postgres;
use sqlx::Pool;

pub struct ContestCrawler<'a> {
    url: String,
    pool: &'a Pool<Postgres>,
}

impl<'a> ContestCrawler<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> ContestCrawler {
        ContestCrawler {
            url: String::from("https://kenkoooo.com/atcoder/resources/contests.json"),
            pool: pool,
        }
    }

    /// AtCoderProblemsからコンテスト情報を取得するメソッド
    pub async fn get_contest_list(&self) -> Result<Vec<ContestJson>> {
        let client = reqwest::Client::new();

        let response = client
            .get(&self.url)
            // AtCoderProblemsはAccept-Encodingにgzipを指定しないとダウンロードできない(https://twitter.com/kenkoooo/status/1147352545133645824)
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .context("Failed to get contest information from AtCoder Problems.")?;

        let json = response
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let contests: Vec<ContestJson> =
            serde_json::from_str(&json).context("Failed to parse JSON body.")?;

        Ok(contests)
    }

    /// AtCoderProblemsから取得したコンテスト情報からデータベースへ格納する用のモデルを作って返すメソッド
    pub async fn crawl(&self) -> Result<Vec<Contest>> {
        let contests: Vec<Contest> = self
            .get_contest_list()
            .await?
            .iter()
            .map(|contest| Contest {
                id: contest.id.clone(),
                start_epoch_second: contest.start_epoch_second.clone(),
                duration_second: contest.duration_second.clone(),
                title: contest.title.clone(),
                rate_change: contest.rate_change.clone(),
                category: contest.categorize(),
            })
            .collect();

        Ok(contests)
    }

    pub async fn save(&self, contests: &Vec<Contest>) -> Result<()> {
        let mut tx = self.pool.begin().await?;

        for contest in contests.iter() {
            let result = sqlx::query("
                MERGE INTO contests
                USING
                    (VALUES($1, $2, $3, $4, $5, $6)) AS contest(id, start_epoch_second, duration_second, title, rate_change, category)
                ON
                    contests.id = contest.id
                WHEN MATCHED THEN
                    UPDATE SET (id, start_epoch_second, duration_second, title, rate_change, category) = (contest.id, contest.start_epoch_second, contest.duration_second, contest.title, contest.rate_change, contest.category)
                WHEN NOT MATCHED THEN
                    INSERT (id, start_epoch_second, duration_second, title, rate_change, category)
                    VALUES (contest.id, contest.start_epoch_second, contest.duration_second, contest.title, contest.rate_change, contest.category);
                ")
                .bind(&contest.id)
                .bind(&contest.start_epoch_second)
                .bind(&contest.duration_second)
                .bind(&contest.title)
                .bind(&contest.rate_change)
                .bind(&contest.category)
                .execute(&mut tx)
                .await;

            if let Err(e) = result {
                tx.rollback().await?;
                return Err(Error::new(e));
            }

            tracing::debug!("Contest {} was saved.", contest.id);
        }

        tx.commit().await?;

        Ok(())
    }
}
