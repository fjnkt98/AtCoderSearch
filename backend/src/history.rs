use anyhow::Context;
use chrono::{DateTime, FixedOffset};
use serde_json::Value;
use sqlx::{Acquire, FromRow, Postgres};

#[derive(Debug, FromRow, PartialEq)]
pub struct BatchHistory {
    pub id: i64,
    pub name: String,
    pub started_at: DateTime<FixedOffset>,
    pub finished_at: Option<DateTime<FixedOffset>>,
    pub status: String,
    pub options: Option<Value>,
}

impl BatchHistory {
    pub async fn new<'a, A>(db: A, name: &str, options: Value) -> anyhow::Result<Self>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        INSERT INTO
            "batch_histories" ("name", "started_at", "options")
        VALUES
            ($1, NOW(), $2)
        RETURNING
            *;
        "#;
        let history: BatchHistory = sqlx::query_as(sql)
            .bind(name)
            .bind(options)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "create batch history")?;

        Ok(history)
    }

    pub async fn finish<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "batch_histories"
        SET
            "status" = 'finished',
            "finished_at" = NOW()
        WHERE
            "id" = $1
        RETURNING
            *;
        "#;

        let history: BatchHistory = sqlx::query_as(sql)
            .bind(self.id)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "update batch history")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }

    pub async fn fetch_latest<'a, A>(db: A, name: &'a str) -> anyhow::Result<Self>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        SELECT
            "id",
            "name",
            "started_at",
            "finished_at",
            "status",
            "options"
        FROM
            "batch_histories"
        WHERE
            "name" = $1
            AND "status" = 'finished'
        ORDER BY
            "started_at" DESC
        LIMIT
            1;
        "#;

        let history: BatchHistory = sqlx::query_as(sql)
            .bind(name)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "fetch latest batch history")?;

        Ok(history)
    }
}

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct Contest {
//     pub contest_id: String,
//     pub start_epoch_second: i64,
//     pub duration_second: i64,
//     pub title: String,
//     pub rate_change: String,
//     pub category: String,
//     pub updated_at: DateTime<FixedOffset>,
// }

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct Difficulty {
//     pub problem_id: String,
//     pub slope: Option<f64>,
//     pub intercept: Option<f64>,
//     pub variance: Option<f64>,
//     pub difficulty: Option<i64>,
//     pub discrimination: Option<f64>,
//     pub irt_loglikelihood: Option<f64>,
//     pub irt_users: Option<f64>,
//     pub is_experimental: Option<bool>,
//     pub updated_at: DateTime<FixedOffset>,
// }

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct Language {
//     pub language: String,
//     pub group: Option<String>,
// }

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct Problem {
//     pub problem_id: String,
//     pub contest_id: String,
//     pub problem_index: String,
//     pub name: String,
//     pub title: String,
//     pub url: String,
//     pub html: String,
//     pub updated_at: DateTime<FixedOffset>,
// }

#[derive(Debug, FromRow, PartialEq, PartialOrd)]
pub struct SubmissionCrawlHistory {
    pub id: i64,
    pub contest_id: String,
    pub started_at: DateTime<FixedOffset>,
    pub status: String,
    pub finished_at: Option<DateTime<FixedOffset>>,
}

impl SubmissionCrawlHistory {
    pub async fn new<'a, A>(db: A, contest_id: &str) -> anyhow::Result<Self>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        INSERT INTO
            "submission_crawl_histories" ("contest_id")
        VALUES
            ($1)
        RETURNING
            *;
        "#;

        let history: SubmissionCrawlHistory = sqlx::query_as(sql)
            .bind(contest_id)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "create submission crawl history")?;

        Ok(history)
    }

    pub async fn finish<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "submission_crawl_histories"
        SET
            "status" = 'finished',
            "finished_at" = NOW()
        WHERE
            "id" = $1
        RETURNING
            *;
        "#;

        let history: SubmissionCrawlHistory = sqlx::query_as(sql)
            .bind(self.id)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "update submission crawl history")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }

    pub async fn fetch_last_crawled<'a, A>(
        db: A,
        contest_id: &str,
    ) -> anyhow::Result<Option<DateTime<FixedOffset>>>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        SELECT
            "started_at"
        FROM
            "submission_crawl_histories"
        WHERE
            "contest_id" = $1
            AND "status" = 'finished'
        ORDER BY
            "started_at" DESC
        LIMIT
            1;
        "#;

        let result: sqlx::Result<(DateTime<FixedOffset>,)> = sqlx::query_as(sql)
            .bind(contest_id)
            .fetch_one(&mut *conn)
            .await;

        match result {
            Ok((latest,)) => {
                return Ok(Some(latest));
            }
            Err(e) => match e {
                sqlx::Error::RowNotFound => {
                    return Ok(None);
                }
                _ => {
                    return Err(e).with_context(|| "fetch last crawled epoch");
                }
            },
        };
    }
}

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct Submission {
//     pub id: i64,
//     pub epoch_second: i64,
//     pub problem_id: String,
//     pub contest_id: Option<String>,
//     pub user_id: Option<String>,
//     pub language: Option<String>,
//     pub point: Option<f64>,
//     pub length: Option<i64>,
//     pub result: Option<String>,
//     pub execution_time: Option<i64>,
//     pub updated_at: DateTime<FixedOffset>,
// }

// #[derive(Debug, FromRow, PartialEq, PartialOrd)]
// pub struct User {
//     pub user_id: String,
//     pub rating: i64,
//     pub highest_rating: i64,
//     pub affiliation: Option<String>,
//     pub birth_year: Option<i64>,
//     pub country: Option<String>,
//     pub crown: Option<String>,
//     pub join_count: i64,
//     pub rank: i64,
//     pub active_rank: Option<i64>,
//     pub wins: i64,
//     pub updated_at: DateTime<FixedOffset>,
// }

#[cfg(test)]
mod tests {
    use std::time::Duration;

    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use rstest::rstest;
    use testcontainers::runners::AsyncRunner;

    #[rstest]
    #[tokio::test]
    async fn test_create_and_update_batch_history() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        let mut history = BatchHistory::new(&pool, "TestBatch", Value::Null)
            .await
            .expect("create new batch history");

        assert_eq!(history.id, 1);
        assert_eq!(history.name, "TestBatch");
        assert_eq!(history.finished_at, None);
        assert_eq!(history.status, "working");

        history.finish(&pool).await.expect("update batch history");

        assert!(history.finished_at.is_some());
        assert_eq!(history.status, "finished");
    }

    #[rstest]
    #[tokio::test]
    async fn test_fetch_latest_batch_history() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        let mut history1 = BatchHistory::new(&pool, "TestBatch", Value::Null)
            .await
            .unwrap();

        history1.finish(&pool).await.unwrap();

        let _history2 = BatchHistory::new(&pool, "TestBatch", Value::Null)
            .await
            .unwrap();

        let latest = BatchHistory::fetch_latest(&pool, "TestBatch")
            .await
            .unwrap();

        assert_eq!(latest, history1);
    }

    #[rstest]
    #[tokio::test]
    async fn test_create_and_update_submission_crawl_history() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        let mut history = SubmissionCrawlHistory::new(&pool, "abc001").await.unwrap();

        assert_eq!(history.id, 1);
        assert_eq!(history.contest_id, "abc001");
        assert_eq!(history.status, "working");

        history.finish(&pool).await.unwrap();

        assert_eq!(history.status, "finished");
    }

    #[rstest]
    #[tokio::test]
    async fn test_fetch_latest_submission_crawl_history() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // histories not found
        assert_eq!(
            SubmissionCrawlHistory::fetch_last_crawled(&pool, "abc001")
                .await
                .unwrap(),
            None
        );

        let mut history1 = SubmissionCrawlHistory::new(&pool, "abc001").await.unwrap();
        tokio::time::sleep(Duration::from_secs(1)).await;
        let _history2 = SubmissionCrawlHistory::new(&pool, "abc001").await.unwrap();

        history1.finish(&pool).await.unwrap();

        let latest = SubmissionCrawlHistory::fetch_last_crawled(&pool, "abc001")
            .await
            .unwrap()
            .unwrap();

        assert_eq!(history1.started_at, latest);
    }
}
