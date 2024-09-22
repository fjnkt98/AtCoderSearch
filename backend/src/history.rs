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
            .bind(&options)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "exec query")?;

        Ok(history)
    }

    pub async fn complete<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "batch_histories"
        SET
            "status" = 'completed',
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
            .with_context(|| "exec query")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }

    pub async fn abort<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        if self.status != "working" {
            return Ok(());
        }

        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "batch_histories"
        SET
            "status" = 'aborted',
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
            .with_context(|| "exec query")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }

    pub async fn fetch_latest<'a, A>(db: A, name: &'a str) -> anyhow::Result<Option<Self>>
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
            AND "status" = 'completed'
        ORDER BY
            "started_at" DESC
        LIMIT
            1;
        "#;

        sqlx::query_as(sql)
            .bind(name)
            .fetch_one(&mut *conn)
            .await
            .and_then(|h| Ok(Some(h)))
            .or_else(|err| match err {
                sqlx::Error::RowNotFound => Ok(None),
                _ => Err(err).with_context(|| "exec query")?,
            })
    }
}

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
            .with_context(|| "exec query")?;

        Ok(history)
    }

    pub async fn complete<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "submission_crawl_histories"
        SET
            "status" = 'completed',
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
            .with_context(|| "exec query")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }

    pub async fn abort<'a, A>(&mut self, db: A) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        if self.status != "working" {
            return Ok(());
        }

        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "submission_crawl_histories"
        SET
            "status" = 'aborted',
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
            .with_context(|| "exec query")?;

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
            AND "status" = 'completed'
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
                    return Err(e).with_context(|| "exec query");
                }
            },
        };
    }
}

#[cfg(test)]
mod tests {
    use std::time::Duration;

    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_create_and_complete_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.id, 1);
        assert_eq!(history.name, "TestBatch");
        assert_eq!(history.finished_at, None);
        assert_eq!(history.status, "working");

        history.complete(&pool).await?;

        assert!(history.finished_at.is_some());
        assert_eq!(history.status, "completed");

        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_abort_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.id, 1);
        assert_eq!(history.name, "TestBatch");
        assert_eq!(history.finished_at, None);
        assert_eq!(history.status, "working");

        history.abort(&pool).await?;

        assert!(history.finished_at.is_some());
        assert_eq!(history.status, "aborted");

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_completed_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;
        history.complete(&pool).await?;

        history.abort(&pool).await?;
        assert_eq!(history.status, "completed");

        Ok(())
    }

    #[tokio::test]
    async fn test_fetch_latest_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let latest = BatchHistory::fetch_latest(&pool, "TestBatch").await?;
        assert_eq!(latest, None);

        let mut history1 = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;
        history1.complete(&pool).await?;

        let _history2 = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;

        let mut history3 = BatchHistory::new(&pool, "TestBatch", Value::Null).await?;
        history3.abort(&pool).await?;

        let latest = BatchHistory::fetch_latest(&pool, "TestBatch").await?;

        assert_eq!(latest, Some(history1));
        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_update_submission_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(&pool, "abc001").await?;

        assert_eq!(history.id, 1);
        assert_eq!(history.contest_id, "abc001");
        assert_eq!(history.status, "working");

        history.complete(&pool).await?;

        assert_eq!(history.status, "completed");

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(&pool, "abc001").await?;

        assert_eq!(history.id, 1);
        assert_eq!(history.contest_id, "abc001");
        assert_eq!(history.status, "working");

        history.abort(&pool).await?;

        assert_eq!(history.status, "aborted");

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_completed_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(&pool, "abc001").await?;

        assert_eq!(history.id, 1);
        assert_eq!(history.contest_id, "abc001");
        assert_eq!(history.status, "working");

        history.complete(&pool).await?;

        history.abort(&pool).await?;

        assert_eq!(history.status, "completed");

        Ok(())
    }

    #[tokio::test]
    async fn test_fetch_latest_submission_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        // history not found
        let latest = SubmissionCrawlHistory::fetch_last_crawled(&pool, "abc001").await?;
        assert_eq!(latest, None);

        let mut history1 = SubmissionCrawlHistory::new(&pool, "abc001").await?;
        tokio::time::sleep(Duration::from_secs(1)).await;
        history1.complete(&pool).await.unwrap();

        let _history2 = SubmissionCrawlHistory::new(&pool, "abc001").await?;

        let mut history3 = SubmissionCrawlHistory::new(&pool, "abc001").await?;
        tokio::time::sleep(Duration::from_secs(1)).await;
        history3.abort(&pool).await.unwrap();

        let latest = SubmissionCrawlHistory::fetch_last_crawled(&pool, "abc001").await?;

        assert_eq!(latest, Some(history1.started_at));

        Ok(())
    }
}
