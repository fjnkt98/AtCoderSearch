use crate::atcoder::{AtCoderProblemsClient, Difficulty};
use anyhow::Context;
use itertools::Itertools;
use sqlx::{self, postgres::Postgres, Acquire, Pool, QueryBuilder};
use tracing::instrument;

#[instrument(skip(client, pool))]
pub async fn crawl_difficulties(
    client: &AtCoderProblemsClient,
    pool: &Pool<Postgres>,
) -> anyhow::Result<()> {
    tracing::info!("start to crawl difficulties");

    let difficulties = client
        .fetch_difficulties()
        .await
        .with_context(|| "crawl difficulties")?;

    let mut tx = pool
        .begin()
        .await
        .with_context(|| "begin transaction to save difficulties")?;

    let mut count = 0;
    for difficulties in difficulties.into_iter().collect_vec().chunks(1000) {
        count += insert_difficulties(&mut tx, difficulties)
            .await
            .with_context(|| "insert difficulties")?;
    }

    tx.commit()
        .await
        .with_context(|| "commit transaction to save difficulties")?;

    tracing::info!("saved {} difficulties successfully", count);
    Ok(())
}

#[instrument(skip(db, difficulties))]
async fn insert_difficulties<'a, A>(
    db: A,
    difficulties: &[(String, Difficulty)],
) -> anyhow::Result<u64>
where
    A: Acquire<'a, Database = Postgres>,
{
    if difficulties.is_empty() {
        return Ok(0);
    }

    let mut builder: QueryBuilder<Postgres> = sqlx::QueryBuilder::new(
        r#"
INSERT INTO "difficulties"
    (
        "problem_id",
        "slope",
        "intercept",
        "variance",
        "difficulty",
        "discrimination",
        "irt_loglikelihood",
        "irt_users",
        "is_experimental",
        "updated_at"
    )
"#,
    );
    builder.push_values(difficulties.iter(), |mut s, d| {
        s.push_bind(&d.0)
            .push_bind(&d.1.slope)
            .push_bind(&d.1.intercept)
            .push_bind(&d.1.variance)
            .push_bind(&d.1.difficulty)
            .push_bind(&d.1.discrimination)
            .push_bind(&d.1.irt_loglikelihood)
            .push_bind(&d.1.irt_users)
            .push_bind(&d.1.is_experimental)
            .push("NOW()");
    });
    builder.push(
        r#"
ON CONFLICT ("problem_id") DO
UPDATE
SET
    "problem_id" = EXCLUDED."problem_id",
    "slope" = EXCLUDED."slope",
    "intercept" = EXCLUDED."intercept",
    "variance" = EXCLUDED."variance",
    "difficulty" = EXCLUDED."difficulty",
    "discrimination" = EXCLUDED."discrimination",
    "irt_loglikelihood" = EXCLUDED."irt_loglikelihood",
    "irt_users" = EXCLUDED."irt_users",
    "is_experimental" = EXCLUDED."is_experimental",
    "updated_at" = NOW();
"#,
    );

    let mut conn = db.acquire().await.with_context(|| "acquire connection")?;
    let res = builder
        .build()
        .execute(&mut *conn)
        .await
        .with_context(|| "execute insert difficulties query")?;

    Ok(res.rows_affected())
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_insert_difficulties() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // test insert empty
        let difficulties = vec![];
        let count = insert_difficulties(&pool, &difficulties).await.unwrap();
        assert_eq!(count, 0);

        // test insert single difficulty
        let difficulties = vec![(
            String::from("abc073_b"),
            Difficulty {
                slope: Some(-0.0007632375406559968),
                intercept: Some(6.20973431279239),
                variance: Some(0.41039685374181584),
                difficulty: Some(-179),
                discrimination: Some(0.004479398673070138),
                irt_loglikelihood: Some(-126.67033990479806),
                irt_users: Some(770.0),
                is_experimental: Some(false),
            },
        )];
        let count = insert_difficulties(&pool, &difficulties).await.unwrap();
        assert_eq!(count, 1);

        // test insert multiple difficulties
        let difficulties = vec![
            (
                String::from("abc118_d"),
                Difficulty {
                    slope: Some(-0.0006619775680720775),
                    intercept: Some(8.881759153638702),
                    variance: Some(0.30752713797776526),
                    difficulty: Some(1657),
                    discrimination: Some(0.004479398673070138),
                    irt_loglikelihood: Some(-491.8630322466751),
                    irt_users: Some(2442.0),
                    is_experimental: Some(false),
                },
            ),
            (
                String::from("agc026_d"),
                Difficulty {
                    slope: Some(-0.0004027506918277324),
                    intercept: Some(9.274529080920633),
                    variance: Some(0.12135365008788429),
                    difficulty: Some(2746),
                    discrimination: Some(0.004479398673070138),
                    irt_loglikelihood: Some(-145.66848869773756),
                    irt_users: Some(1799.0),
                    is_experimental: Some(false),
                },
            ),
        ];
        let count = insert_difficulties(&pool, &difficulties).await.unwrap();
        assert_eq!(count, 2);
    }
}
