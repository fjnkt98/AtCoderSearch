use crate::atcoder::{AtCoderClient, User};
use crate::errors::CanceledError;
use anyhow::Context;
use sqlx::{self, postgres::Postgres, Acquire, Pool, QueryBuilder};
use std::time::Duration;
use tokio_util::sync::CancellationToken;
use tracing::instrument;

#[instrument(skip(token, client, pool))]
pub async fn crawl_users(
    token: CancellationToken,
    client: &AtCoderClient,
    pool: &Pool<Postgres>,
    duration: Duration,
) -> anyhow::Result<()> {
    tracing::info!("start to crawl users");

    let mut tx = pool.begin().await.with_context(|| "begin transaction")?;

    let mut page = 1;
    let mut count = 0;
    'l: loop {
        if token.is_cancelled() {
            return Err(CanceledError::default().into());
        }

        tracing::info!("fetch users at page {}", page);

        let users = client
            .fetch_users(page)
            .await
            .with_context(|| "crawl users")?;

        if users.is_empty() {
            tracing::info!("no users found at page {}", page);
            break 'l;
        }

        count += insert_users(&mut tx, &users)
            .await
            .with_context(|| "insert users")?;

        page += 1;

        tokio::time::sleep(duration).await;
    }

    tx.commit().await.with_context(|| "commit transaction")?;

    tracing::info!("saved {} users successfully", count);
    Ok(())
}

#[instrument(skip(db, users))]
async fn insert_users<'a, A>(db: A, users: &[User]) -> anyhow::Result<u64>
where
    A: Acquire<'a, Database = Postgres>,
{
    if users.is_empty() {
        return Ok(0);
    }

    let mut builder: QueryBuilder<Postgres> = sqlx::QueryBuilder::new(
        r#"
INSERT INTO
    "users" (
        "user_id",
        "rating",
        "highest_rating",
        "affiliation",
        "birth_year",
        "country",
        "crown",
        "join_count",
        "rank",
        "active_rank",
        "wins",
        "updated_at"
    )
"#,
    );
    builder.push_values(users.iter(), |mut s, u| {
        s.push_bind(&u.user_id)
            .push_bind(&u.rating)
            .push_bind(&u.highest_rating)
            .push_bind(&u.affiliation)
            .push_bind(&u.birth_year)
            .push_bind(&u.country)
            .push_bind(&u.crown)
            .push_bind(&u.join_count)
            .push_bind(&u.rank)
            .push_bind(&u.active_rank)
            .push_bind(&u.wins)
            .push("NOW()");
    });
    builder.push(
        r#"
ON CONFLICT ("user_id") DO
UPDATE
SET
    "user_id" = excluded."user_id",
    "rating" = excluded."rating",
    "highest_rating" = excluded."highest_rating",
    "affiliation" = excluded."affiliation",
    "birth_year" = excluded."birth_year",
    "country" = excluded."country",
    "crown" = excluded."crown",
    "join_count" = excluded."join_count",
    "rank" = excluded."rank",
    "active_rank" = excluded."active_rank",
    "wins" = excluded."wins",
    "updated_at" = NOW();
"#,
    );

    let mut conn = db.acquire().await.with_context(|| "acquire connection")?;

    let res = builder
        .build()
        .execute(&mut *conn)
        .await
        .with_context(|| "execute insert users query")?;

    Ok(res.rows_affected())
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_insert_users() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        // test insert empty
        let users = vec![];
        let count = insert_users(&pool, &users).await.unwrap();
        assert_eq!(count, 0);

        // test insert single contest
        let users = vec![User {
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
        }];
        let count = insert_users(&pool, &users).await.unwrap();
        assert_eq!(count, 1);

        // test insert multiple users
        let users = vec![
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
        ];
        let count = insert_users(&pool, &users).await.unwrap();
        assert_eq!(count, 2);
    }
}
