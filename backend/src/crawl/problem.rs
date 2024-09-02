use anyhow::Context;
use itertools::Itertools;
use sqlx::{Acquire, Pool, Postgres};
use std::{collections::BTreeSet, time::Duration};
use tracing::instrument;

use crate::atcoder::{AtCoderClient, AtCoderProblemsClient, Problem};

#[instrument(skip(atcoder_client, problems_client, pool))]
pub async fn crawl_problems<'a>(
    atcoder_client: &'a AtCoderClient,
    problems_client: &'a AtCoderProblemsClient,
    pool: &'a Pool<Postgres>,
    all: bool,
    duration: Duration,
) -> anyhow::Result<()> {
    let mut targets = problems_client.fetch_problems().await?;
    if !all {
        targets = detect_diff(pool, &targets).await?;
    }

    tracing::info!("start to crawl {} problems", targets.len());

    for target in targets.iter() {
        let html = atcoder_client
            .fetch_problem_html(&target.contest_id, &target.id)
            .await?;

        let mut tx = pool.begin().await?;
        insert_problem(&mut tx, target, &html).await?;

        tx.commit().await?;

        tracing::info!("saved problem {} successfully", target.id);
        tokio::time::sleep(duration).await;
    }

    Ok(())
}

async fn detect_diff<'a, A>(db: A, problems: &[Problem]) -> anyhow::Result<Vec<Problem>>
where
    A: Acquire<'a, Database = Postgres>,
{
    let mut conn = db.acquire().await?;

    let rows: Vec<(String,)> = sqlx::query_as(r#"SELECT "problem_id" FROM "problems";"#)
        .fetch_all(&mut *conn)
        .await
        .with_context(|| "fetch problem ids from database")?;
    let exists = BTreeSet::from_iter(rows.iter().map(|r| r.0.clone()));

    let result = problems
        .iter()
        .filter(|p| !exists.contains(&p.id))
        .cloned()
        .collect_vec();

    Ok(result)
}

#[instrument(skip(db, problem, html))]
async fn insert_problem<'a, A>(db: A, problem: &Problem, html: &str) -> anyhow::Result<u64>
where
    A: Acquire<'a, Database = Postgres>,
{
    let sql = r#"
INSERT INTO
    "problems" (
        "problem_id",
        "contest_id",
        "problem_index",
        "name",
        "title",
        "url",
        "html",
        "updated_at"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW())
ON CONFLICT ("problem_id") DO
UPDATE
SET
    "contest_id" = EXCLUDED."contest_id",
    "problem_index" = EXCLUDED."problem_index",
    "name" = EXCLUDED."name",
    "title" = EXCLUDED."title",
    "url" = EXCLUDED."url",
    "html" = EXCLUDED."html",
    "updated_at" = NOW();
    "#;

    let mut conn = db.acquire().await?;
    let res = sqlx::query(sql)
        .bind(&problem.id)
        .bind(&problem.contest_id)
        .bind(&problem.problem_index)
        .bind(&problem.name)
        .bind(&problem.title)
        .bind(&problem.url())
        .bind(html)
        .execute(&mut *conn)
        .await
        .with_context(|| "execute insert problem query")?;

    Ok(res.rows_affected())
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use sqlx::Row;
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_insert_problem() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        let problem = Problem {
            id: String::from("abc001_a"),
            contest_id: String::from("abc001"),
            problem_index: String::from("A"),
            title: String::from("title"),
            name: String::from("A. title"),
        };
        let html = "html";

        let count = insert_problem(&pool, &problem, &html).await.unwrap();
        assert_eq!(count, 1);

        let sql = r#"SELECT * FROM "problems""#;
        let rows = sqlx::query(sql).fetch_all(&pool).await.unwrap();
        assert_eq!(rows.len(), 1);
        let row = &rows[0];

        assert_eq!(
            problem.id,
            row.try_get::<String, &str>("problem_id").unwrap()
        );
        assert_eq!(
            problem.contest_id,
            row.try_get::<String, &str>("contest_id").unwrap()
        );
        assert_eq!(
            problem.problem_index,
            row.try_get::<String, &str>("problem_index").unwrap()
        );
        assert_eq!(problem.name, row.try_get::<String, &str>("name").unwrap());
        assert_eq!(problem.title, row.try_get::<String, &str>("title").unwrap());
        assert_eq!("html", row.try_get::<String, &str>("html").unwrap());
        assert_eq!(
            "https://atcoder.jp/contests/abc001/tasks/abc001_a",
            row.try_get::<String, &str>("url").unwrap()
        );
    }

    #[tokio::test]
    async fn test_detect_diff() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        let sql = r#"
INSERT INTO "problems" ("problem_id", "contest_id", "problem_index", "name", "title", "url", "html")
VALUES
    ('abc001_a', 'abc001', 'A', 'test problem 1', 'A. test problem 1', 'url', 'html');
        "#;
        let res = sqlx::query(sql).execute(&pool).await.unwrap();
        assert_eq!(res.rows_affected(), 1);

        let problems = vec![
            Problem {
                id: String::from("abc001_a"),
                contest_id: String::from("abc001"),
                problem_index: String::from("A"),
                name: String::from("test problem 1"),
                title: String::from("A. test problem 1"),
            },
            Problem {
                id: String::from("abc001_b"),
                contest_id: String::from("abc001"),
                problem_index: String::from("B"),
                name: String::from("test problem 2"),
                title: String::from("B. test problem 2"),
            },
        ];

        let diff = detect_diff(&pool, &problems).await.unwrap();
        let want = vec![problems[1].clone()];
        assert_eq!(diff, want);
    }
}
