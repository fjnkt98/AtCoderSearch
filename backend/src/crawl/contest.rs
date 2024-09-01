use crate::atcoder::{AtCoderProblemsClient, Contest};
use anyhow::Context;
use sqlx::{self, postgres::Postgres, Acquire, Pool, QueryBuilder};

pub async fn crawl_contests(
    client: &AtCoderProblemsClient,
    pool: &Pool<Postgres>,
) -> anyhow::Result<()> {
    let contests = client
        .fetch_contests()
        .await
        .with_context(|| "crawl contest")?;

    let mut tx = pool.begin().await.with_context(|| "begin transaction")?;

    let mut _count = 0;
    for chunk in contests.chunks(1000) {
        _count += insert_contests(&mut tx, chunk)
            .await
            .with_context(|| "insert contests")?;
    }

    tx.commit().await.with_context(|| "commit transaction")?;

    Ok(())
}

async fn insert_contests<'a, A>(db: A, contests: &[Contest]) -> anyhow::Result<u64>
where
    A: Acquire<'a, Database = Postgres>,
{
    if contests.is_empty() {
        return Ok(0);
    }

    let mut builder: QueryBuilder<Postgres> = sqlx::QueryBuilder::new(
        r#"INSERT INTO "contests" ("contest_id", "start_epoch_second", "duration_second", "title", "rate_change", "category", "updated_at") "#,
    );
    builder.push_values(contests.iter(), |mut separated, c| {
        separated
            .push_bind(&c.id)
            .push_bind(&c.start_epoch_second)
            .push_bind(&c.duration_second)
            .push_bind(&c.title)
            .push_bind(&c.rate_change)
            .push_bind(c.categorize())
            .push("NOW()");
    });
    builder.push(r#" ON CONFLICT ("contest_id") DO UPDATE SET "contest_id" = EXCLUDED."contest_id", "start_epoch_second" = EXCLUDED."start_epoch_second", "duration_second" = EXCLUDED."duration_second", "title" = EXCLUDED."title", "rate_change" = EXCLUDED."rate_change", "category" = EXCLUDED."category", "updated_at" = NOW()"#);

    let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

    let res = builder
        .build()
        .execute(&mut *conn)
        .await
        .with_context(|| "execute insert contests query")?;

    Ok(res.rows_affected())
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_insert_contests() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // test insert empty
        let contests = vec![];
        let count = insert_contests(&pool, &contests).await.unwrap();
        assert_eq!(count, 0);

        // test insert single contest
        let contests = vec![Contest {
            id: String::from("abc001"),
            start_epoch_second: 1468670400,
            duration_second: 6000,
            title: String::from("abc001"),
            rate_change: String::from("-"),
        }];
        let count = insert_contests(&pool, &contests).await.unwrap();
        assert_eq!(count, 1);

        // test insert multiple contests
        let contests = vec![
            Contest {
                id: String::from("abc001"),
                start_epoch_second: 1468670400,
                duration_second: 6000,
                title: String::from("abc001"),
                rate_change: String::from("-"),
            },
            Contest {
                id: String::from("abc002"),
                start_epoch_second: 1468670400,
                duration_second: 6000,
                title: String::from("abc002"),
                rate_change: String::from("-"),
            },
        ];
        let count = insert_contests(&pool, &contests).await.unwrap();
        assert_eq!(count, 2);
    }
}
