use anyhow::Context;
use chrono::{DateTime, FixedOffset};
use serde_json::Value;
use sqlx::{Acquire, FromRow, Pool, Postgres};
use std::{error, fmt};

#[derive(Debug, Clone, PartialEq, PartialOrd)]
pub enum HistoryStatus {
    Working,
    Completed,
    Aborted,
}

impl TryFrom<&str> for HistoryStatus {
    type Error = UnknownStatusError;

    fn try_from(value: &str) -> Result<Self, Self::Error> {
        match value {
            "working" => Ok(Self::Working),
            "completed" => Ok(Self::Completed),
            "aborted" => Ok(Self::Aborted),
            _ => Err(UnknownStatusError::new(value)),
        }
    }
}

impl TryFrom<String> for HistoryStatus {
    type Error = UnknownStatusError;

    fn try_from(value: String) -> Result<Self, Self::Error> {
        Self::try_from(value.as_ref())
    }
}

impl AsRef<str> for HistoryStatus {
    fn as_ref(&self) -> &str {
        match self {
            Self::Working => "working",
            Self::Completed => "completed",
            Self::Aborted => "aborted",
        }
    }
}

#[derive(Debug, Clone)]
pub struct UnknownStatusError {
    value: String,
}

impl fmt::Display for UnknownStatusError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "unknown status literal: {}", self.value)
    }
}

impl error::Error for UnknownStatusError {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}

impl UnknownStatusError {
    pub fn new(value: &str) -> Self {
        Self {
            value: String::from(value),
        }
    }
}
#[derive(Debug, Clone)]
pub struct BatchHistory {
    pool: Pool<Postgres>,
    with_transaction: bool,
    pub data: BatchHistoryRow,
}

#[derive(Debug, Clone, FromRow, PartialEq)]
pub struct BatchHistoryRow {
    pub id: i64,
    pub name: String,
    pub started_at: DateTime<FixedOffset>,
    pub finished_at: Option<DateTime<FixedOffset>>,
    #[sqlx(try_from = "String")]
    pub status: HistoryStatus,
    pub options: Option<Value>,
}

impl BatchHistory {
    pub async fn new(pool: Pool<Postgres>, name: &str, options: Value) -> anyhow::Result<Self> {
        let data = BatchHistoryRow::new(&pool, name, options)
            .await
            .with_context(|| "create new history row")?;

        Ok(Self {
            pool,
            with_transaction: false,
            data,
        })
    }

    pub async fn with_transaction(
        pool: Pool<Postgres>,
        name: &str,
        options: Value,
    ) -> anyhow::Result<Self> {
        let mut tx = pool
            .begin()
            .await
            .with_context(|| "begin transaction to create batch crawl history")?;
        let data = BatchHistoryRow::new(&mut tx, name, options)
            .await
            .with_context(|| "create new history row")?;
        tx.commit()
            .await
            .with_context(|| "commit transaction to create batch crawl history")?;

        Ok(Self {
            pool,
            with_transaction: true,
            data,
        })
    }

    pub async fn complete(&mut self) -> anyhow::Result<()> {
        if self.with_transaction {
            let mut tx = self
                .pool
                .begin()
                .await
                .with_context(|| "begin transaction to update batch crawl history")?;

            self.data
                .update(&mut tx, HistoryStatus::Completed)
                .await
                .with_context(|| "update history row")?;

            tx.commit()
                .await
                .with_context(|| "commit transaction to update batch crawl history")?;
        } else {
            self.data
                .update(&self.pool, HistoryStatus::Completed)
                .await
                .with_context(|| "update history row")?;
        };

        Ok(())
    }

    pub async fn abort(&mut self) -> anyhow::Result<()> {
        if self.data.status != HistoryStatus::Working {
            return Ok(());
        }

        if self.with_transaction {
            let mut tx = self
                .pool
                .begin()
                .await
                .with_context(|| "begin transaction to update batch crawl history")?;

            self.data
                .update(&mut tx, HistoryStatus::Aborted)
                .await
                .with_context(|| "update history row")?;

            tx.commit()
                .await
                .with_context(|| "commit transaction to update batch crawl history")?;
        } else {
            self.data
                .update(&self.pool, HistoryStatus::Aborted)
                .await
                .with_context(|| "update history row")?;
        };

        Ok(())
    }

    pub async fn fetch_latest(
        pool: &Pool<Postgres>,
        name: &str,
    ) -> anyhow::Result<Option<BatchHistoryRow>> {
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
            AND "status" = $2
        ORDER BY
            "started_at" DESC
        LIMIT
            1;
        "#;

        sqlx::query_as(sql)
            .bind(name)
            .bind(HistoryStatus::Completed.as_ref())
            .fetch_one(pool)
            .await
            .and_then(|h| Ok(Some(h)))
            .or_else(|err| match err {
                sqlx::Error::RowNotFound => Ok(None),
                _ => Err(err).with_context(|| "exec query")?,
            })
    }
}

impl BatchHistoryRow {
    pub(crate) async fn new<'a, A>(db: A, name: &str, options: Value) -> anyhow::Result<Self>
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
        let history: BatchHistoryRow = sqlx::query_as(sql)
            .bind(name)
            .bind(&options)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "exec query")?;

        Ok(history)
    }

    pub(crate) async fn update<'a, A>(&mut self, db: A, status: HistoryStatus) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "batch_histories"
        SET
            "status" = $2,
            "finished_at" = NOW()
        WHERE
            "id" = $1
        RETURNING
            *;
        "#;

        let history: BatchHistoryRow = sqlx::query_as(sql)
            .bind(self.id)
            .bind(status.as_ref())
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "exec query")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }
}

#[derive(Debug, Clone)]
pub struct SubmissionCrawlHistory {
    pool: Pool<Postgres>,
    with_transaction: bool,
    pub data: SubmissionCrawlHistoryRow,
}

#[derive(Debug, Clone, FromRow, PartialEq, PartialOrd)]
pub struct SubmissionCrawlHistoryRow {
    pub id: i64,
    pub contest_id: String,
    pub started_at: DateTime<FixedOffset>,
    #[sqlx(try_from = "String")]
    pub status: HistoryStatus,
    pub finished_at: Option<DateTime<FixedOffset>>,
}

impl SubmissionCrawlHistory {
    pub async fn new(pool: Pool<Postgres>, contest_id: &str) -> anyhow::Result<Self> {
        let data = SubmissionCrawlHistoryRow::new(&pool, contest_id)
            .await
            .with_context(|| "create new history row")?;

        Ok(Self {
            pool,
            with_transaction: false,
            data,
        })
    }

    pub async fn with_transaction(pool: Pool<Postgres>, contest_id: &str) -> anyhow::Result<Self> {
        let mut tx = pool
            .begin()
            .await
            .with_context(|| "begin transaction to create submission crawl history")?;

        let data = SubmissionCrawlHistoryRow::new(&mut tx, contest_id)
            .await
            .with_context(|| "create new history row")?;

        tx.commit()
            .await
            .with_context(|| "commit transaction to create submission crawl history")?;

        Ok(Self {
            pool,
            with_transaction: true,
            data,
        })
    }

    pub async fn complete(&mut self) -> anyhow::Result<()> {
        if self.with_transaction {
            let mut tx = self
                .pool
                .begin()
                .await
                .with_context(|| "begin transaction to update submission crawl history")?;

            self.data
                .update(&mut tx, HistoryStatus::Completed)
                .await
                .with_context(|| "update history row")?;

            tx.commit()
                .await
                .with_context(|| "commit transaction to update submission crawl history")?;
        } else {
            self.data
                .update(&self.pool, HistoryStatus::Completed)
                .await
                .with_context(|| "update history row")?;
        }

        Ok(())
    }

    pub async fn abort(&mut self) -> anyhow::Result<()> {
        if self.data.status != HistoryStatus::Working {
            return Ok(());
        }

        if self.with_transaction {
            let mut tx = self
                .pool
                .begin()
                .await
                .with_context(|| "begin transaction to update submission crawl history")?;

            self.data
                .update(&mut tx, HistoryStatus::Aborted)
                .await
                .with_context(|| "update history row")?;

            tx.commit()
                .await
                .with_context(|| "commit transaction to update submission crawl history")?;
        } else {
            self.data
                .update(&self.pool, HistoryStatus::Aborted)
                .await
                .with_context(|| "update history row")?;
        }

        Ok(())
    }

    pub async fn fetch_latest(
        pool: &Pool<Postgres>,
        contest_id: &str,
    ) -> anyhow::Result<Option<SubmissionCrawlHistoryRow>> {
        let sql = r#"
        SELECT
            *
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

        sqlx::query_as(sql)
            .bind(contest_id)
            .fetch_one(pool)
            .await
            .and_then(|h| Ok(Some(h)))
            .or_else(|err| match err {
                sqlx::Error::RowNotFound => Ok(None),
                _ => Err(err).with_context(|| "exec query")?,
            })
    }
}

impl SubmissionCrawlHistoryRow {
    pub(crate) async fn new<'a, A>(db: A, contest_id: &str) -> anyhow::Result<Self>
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

        let history: SubmissionCrawlHistoryRow = sqlx::query_as(sql)
            .bind(contest_id)
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "exec query")?;

        Ok(history)
    }

    pub async fn update<'a, A>(&mut self, db: A, status: HistoryStatus) -> anyhow::Result<()>
    where
        A: Acquire<'a, Database = Postgres>,
    {
        let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

        let sql = r#"
        UPDATE "submission_crawl_histories"
        SET
            "status" = $2,
            "finished_at" = NOW()
        WHERE
            "id" = $1
        RETURNING
            *;
        "#;

        let history: SubmissionCrawlHistoryRow = sqlx::query_as(sql)
            .bind(self.id)
            .bind(status.as_ref())
            .fetch_one(&mut *conn)
            .await
            .with_context(|| "exec query")?;

        self.status = history.status;
        self.finished_at = history.finished_at;

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_create_and_complete_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.name, "TestBatch");
        assert_eq!(history.data.finished_at, None);
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.complete().await?;

        assert!(history.data.finished_at.is_some());
        assert_eq!(history.data.status, HistoryStatus::Completed);

        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_complete_batch_history_with_transaction() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::with_transaction(pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.name, "TestBatch");
        assert_eq!(history.data.finished_at, None);
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.complete().await?;

        assert!(history.data.finished_at.is_some());
        assert_eq!(history.data.status, HistoryStatus::Completed);

        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_abort_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.name, "TestBatch");
        assert_eq!(history.data.finished_at, None);
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.abort().await?;

        assert!(history.data.finished_at.is_some());
        assert_eq!(history.data.status, HistoryStatus::Aborted);

        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_abort_batch_history_with_transaction() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::with_transaction(pool, "TestBatch", Value::Null).await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.name, "TestBatch");
        assert_eq!(history.data.finished_at, None);
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.abort().await?;

        assert!(history.data.finished_at.is_some());
        assert_eq!(history.data.status, HistoryStatus::Aborted);

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_completed_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = BatchHistory::new(pool, "TestBatch", Value::Null).await?;
        history.complete().await?;

        history.abort().await?;
        assert_eq!(history.data.status, HistoryStatus::Completed);

        Ok(())
    }

    #[tokio::test]
    async fn test_fetch_latest_batch_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let latest = BatchHistory::fetch_latest(&pool, "TestBatch").await?;
        assert_eq!(latest, None);

        let mut history1 = BatchHistory::new(pool.clone(), "TestBatch", Value::Null).await?;
        history1.complete().await?;

        let _history2 = BatchHistory::new(pool.clone(), "TestBatch", Value::Null).await?;

        let mut history3 = BatchHistory::new(pool.clone(), "TestBatch", Value::Null).await?;
        history3.abort().await?;

        let latest = BatchHistory::fetch_latest(&pool, "TestBatch").await?;

        assert_eq!(latest, Some(history1.data));
        Ok(())
    }

    #[tokio::test]
    async fn test_create_and_update_submission_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(pool, "abc001").await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.contest_id, "abc001");
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.complete().await?;

        assert_eq!(history.data.status, HistoryStatus::Completed);

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(pool, "abc001").await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.contest_id, "abc001");
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.abort().await?;

        assert_eq!(history.data.status, HistoryStatus::Aborted);

        Ok(())
    }

    #[tokio::test]
    async fn test_abort_completed_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        let mut history = SubmissionCrawlHistory::new(pool, "abc001").await?;

        assert_eq!(history.data.id, 1);
        assert_eq!(history.data.contest_id, "abc001");
        assert_eq!(history.data.status, HistoryStatus::Working);

        history.complete().await?;

        history.abort().await?;

        assert_eq!(history.data.status, HistoryStatus::Completed);

        Ok(())
    }

    #[tokio::test]
    async fn test_fetch_latest_submission_crawl_history() -> anyhow::Result<()> {
        let container = create_container()?.start().await?;
        let pool = create_pool_from_container(&container).await?;

        // history not found
        let latest = SubmissionCrawlHistory::fetch_latest(&pool, "abc001").await?;
        assert_eq!(latest, None);

        let mut history1 = SubmissionCrawlHistory::new(pool.clone(), "abc001").await?;
        history1.complete().await.unwrap();

        let _history2 = SubmissionCrawlHistory::new(pool.clone(), "abc001").await?;

        let mut history3 = SubmissionCrawlHistory::new(pool.clone(), "abc001").await?;
        history3.abort().await.unwrap();

        let latest = SubmissionCrawlHistory::fetch_latest(&pool, "abc001").await?;

        assert_eq!(latest, Some(history1.data));

        Ok(())
    }
}
