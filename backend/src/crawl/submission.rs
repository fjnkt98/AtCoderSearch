use anyhow::Context;
use sqlx::{Acquire, Pool, Postgres, QueryBuilder};
use std::{collections::VecDeque, time::Duration};
use tracing::instrument;

use crate::{
    atcoder::{AtCoderClient, Submission},
    history::SubmissionCrawlHistory,
};

pub struct SubmissionCrawler<'a> {
    client: &'a AtCoderClient,
    pool: &'a Pool<Postgres>,
    duration: Duration,
    retry: i64,
    targets: Vec<String>,
}

impl<'a> SubmissionCrawler<'a> {
    pub fn new(
        client: &'a AtCoderClient,
        pool: &'a Pool<Postgres>,
        duration: Duration,
        retry: i64,
        targets: Vec<String>,
    ) -> Self {
        Self {
            client,
            pool,
            duration,
            retry,
            targets,
        }
    }

    pub async fn crawl(&self) -> anyhow::Result<()> {
        let targets = fetch_contest_id(self.pool, &self.targets).await?;

        for contest_id in targets.iter() {
            self.crawl_contest(contest_id).await?;

            tokio::time::sleep(self.duration).await;
        }
        todo!();
    }

    #[instrument(skip(self))]
    async fn crawl_contest(&self, contest_id: &str) -> anyhow::Result<()> {
        let last_crawled =
            SubmissionCrawlHistory::fetch_last_crawled(self.pool, contest_id).await?;
        tracing::info!(
            "start to crawl submissions of {}. last crawled at {:?}",
            contest_id,
            last_crawled
        );
        let last_crawled = last_crawled.and_then(|t| Some(t.timestamp())).unwrap_or(0);

        let mut tx = self
            .pool
            .begin()
            .await
            .with_context(|| "begin transaction to create submission crawl history")?;
        let mut history = SubmissionCrawlHistory::new(&mut tx, contest_id)
            .await
            .with_context(|| "create submission crawl history")?;
        tx.commit()
            .await
            .with_context(|| "commit transaction to create submission crawl history")?;

        let mut submissions: Vec<Submission> = Vec::with_capacity(20);

        let mut queue = VecDeque::from([1]);
        let mut remain = self.retry;
        while let Some(page) = queue.pop_front() {
            tracing::info!("fetch submissions at page {}", page);

            match self.client.fetch_submissions(contest_id, page).await {
                Ok(s) => {
                    if s.is_empty() {
                        tracing::info!("there is no submissions in {}", contest_id);
                        break;
                    }

                    let head = s.get(0).unwrap().clone();
                    submissions.extend(s);

                    if head.epoch_second <= last_crawled {
                        tracing::info!("break crawling since all submissions of {} after page {} have been crawled", contest_id, page);
                        break;
                    }

                    queue.push_back(page + 1);
                    remain = self.retry;
                }
                Err(e) => {
                    if remain <= 0 {
                        return Err(e).with_context(|| "fetch submissions");
                    }

                    tracing::error!("failed to crawl submissions of {} cause: {:#}. retry to crawl after 60 seconds...", contest_id, e);

                    queue.push_back(page);
                    remain -= 1;
                    tokio::time::sleep(Duration::from_secs(60)).await;
                }
            }

            tokio::time::sleep(self.duration).await;
        }

        let mut tx = self
            .pool
            .begin()
            .await
            .with_context(|| "begin transaction to save submissions")?;
        let mut count = 0;
        for chunk in submissions.chunks(1000) {
            count += insert_submissions(&mut tx, &chunk)
                .await
                .with_context(|| "insert submissions")?;
        }
        tx.commit()
            .await
            .with_context(|| "commit transaction to save submissions")?;
        tracing::info!("saved {} submissions successfully", count);

        let mut tx = self
            .pool
            .begin()
            .await
            .with_context(|| "begin transaction to update submission crawl history")?;
        history
            .finish(&mut tx)
            .await
            .with_context(|| "finish submission crawl history")?;
        tx.commit()
            .await
            .with_context(|| "commit transaction to update submission crawl history")?;

        Ok(())
    }
}

#[instrument(skip(db))]
async fn fetch_contest_id<'a, A>(db: A, categories: &[String]) -> anyhow::Result<Vec<String>>
where
    A: Acquire<'a, Database = Postgres>,
{
    let mut builder: QueryBuilder<Postgres> = sqlx::QueryBuilder::new(
        r#"
        SELECT
            "contest_id"
        FROM
            "contests"
"#,
    );
    if !categories.is_empty() {
        builder.push(r#"WHERE "category" IN ("#);
        let mut s = builder.separated(", ");
        for c in categories.iter() {
            s.push_bind(c);
        }
        s.push_unseparated(")");
    }
    builder.push(r#"ORDER BY "start_epoch_second" DESC"#);

    let mut conn = db.acquire().await?;

    let res: Vec<(String,)> = builder
        .build_query_as()
        .fetch_all(&mut *conn)
        .await
        .with_context(|| "execute select contest_id query")?;

    Ok(res.into_iter().map(|(c,)| c).collect())
}

#[instrument(skip(db, submissions))]
async fn insert_submissions<'a, A>(db: A, submissions: &[Submission]) -> anyhow::Result<u64>
where
    A: Acquire<'a, Database = Postgres>,
{
    if submissions.is_empty() {
        return Ok(0);
    }

    let mut builder: QueryBuilder<Postgres> = sqlx::QueryBuilder::new(
        r#"
INSERT INTO
    "submissions" (
        "id",
        "epoch_second",
        "problem_id",
        "contest_id",
        "user_id",
        "language",
        "point",
        "length",
        "result",
        "execution_time",
        "updated_at"
    )
"#,
    );
    builder.push_values(submissions.iter(), |mut separated, s| {
        separated
            .push_bind(&s.id)
            .push_bind(&s.epoch_second)
            .push_bind(&s.problem_id)
            .push_bind(&s.contest_id)
            .push_bind(&s.user_id)
            .push_bind(&s.language)
            .push_bind(&s.point)
            .push_bind(&s.length)
            .push_bind(&s.result)
            .push_bind(&s.execution_time)
            .push("NOW()");
    });
    builder.push(
        r#"
ON CONFLICT ("id", "epoch_second") DO
UPDATE
SET
    "id" = excluded."id",
    "epoch_second" = excluded."epoch_second",
    "problem_id" = excluded."problem_id",
    "contest_id" = excluded."contest_id",
    "user_id" = excluded."user_id",
    "language" = excluded."language",
    "point" = excluded."point",
    "length" = excluded."length",
    "result" = excluded."result",
    "execution_time" = excluded."execution_time",
    "updated_at" = NOW();
"#,
    );

    let mut conn = db.acquire().await?;

    let res = builder
        .build()
        .execute(&mut *conn)
        .await
        .with_context(|| "execute insert submissions query")?;

    Ok(res.rows_affected())
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_insert_submissions() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // test insert empty
        let submissions = vec![];
        let count = insert_submissions(&pool, &submissions).await.unwrap();
        assert_eq!(count, 0);

        // test insert single submission
        let submissions = vec![Submission {
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
        }];
        let count = insert_submissions(&pool, &submissions).await.unwrap();
        assert_eq!(count, 1);

        let inserted: Vec<(i64, i64, String, Option<String>, Option<String>, Option<String>, Option<f64>, Option<i32>, Option<String>, Option<i32>,)> = sqlx::query_as(r#"SELECT "id", "epoch_second", "problem_id", "contest_id", "user_id", "language", "point", "length", "result", "execution_time" FROM "submissions""#).fetch_all(&pool).await.unwrap();
        assert_eq!(
            inserted,
            vec![(
                48852107,
                1703553569,
                String::from("abc300_a"),
                Some(String::from("abc300")),
                Some(String::from("Orkhon2010")),
                Some(String::from("C++ 20 (gcc 12.2)")),
                Some(100.0),
                Some(259),
                Some(String::from("AC")),
                Some(1),
            )]
        );

        // test insert multiple submissions
        let submissions = vec![
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
        let count = insert_submissions(&pool, &submissions).await.unwrap();
        assert_eq!(count, 2);
    }

    #[tokio::test]
    async fn test_fetch_contest_id() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // test data
        sqlx::query(
            r#"
INSERT INTO "contests" ("contest_id", "start_epoch_second", "duration_second", "title", "rate_change", "category") VALUES
('abc001', 0, 0, '', '-', 'ABC'),
('abc002', 0, 0, '', '-', 'ABC'),
('arc001', 0, 0, '', '-', 'ARC');
"#,
        )
        .execute(&pool)
        .await
        .unwrap();

        // test fetch all contest
        let categories = vec![];
        let result = fetch_contest_id(&pool, &categories).await.unwrap();
        assert_eq!(result.len(), 3);

        // test fetch specific category
        let categories = vec![String::from("ABC")];
        let result = fetch_contest_id(&pool, &categories).await.unwrap();
        assert_eq!(result.len(), 2);

        let categories = vec![String::from("ABC"), String::from("ARC")];
        let result = fetch_contest_id(&pool, &categories).await.unwrap();
        assert_eq!(result.len(), 3);
    }
}
