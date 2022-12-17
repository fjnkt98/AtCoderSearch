use crate::contest::models::Contest;
use anyhow::{Error, Result};
use sqlx::postgres::Postgres;
use sqlx::Pool;

pub async fn insert(pool: &Pool<Postgres>, contests: &Vec<Contest>) -> Result<()> {
    let mut tx = pool.begin().await?;

    let result = sqlx::query!("TRUNCATE contest CASCADE;")
        .execute(&mut tx)
        .await;

    if let Err(e) = result {
        tx.rollback().await?;
        return Err(Error::new(e));
    }

    for contest in contests.iter() {
        let result = sqlx::query!(
            r#"
            INSERT INTO contest
            (id, start_epoch_second, duration_second, title, rate_change, category)
            VALUES ($1, $2, $3, $4, $5, $6);
            "#,
            contest.id,
            contest.start_epoch_second,
            contest.duration_second,
            contest.title,
            contest.rate_change,
            contest.category
        )
        .execute(&mut tx)
        .await;

        if let Err(e) = result {
            tx.rollback().await?;
            return Err(Error::new(e));
        }
    }

    tx.commit().await?;

    Ok(())
}
