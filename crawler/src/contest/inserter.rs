use crate::contest::models::Contest;
use anyhow::{Error, Result};
use sqlx::postgres::Postgres;
use sqlx::Pool;

pub async fn insert(pool: &Pool<Postgres>, contests: &Vec<Contest>) -> Result<()> {
    let mut tx = pool.begin().await?;

    for contest in contests.iter() {
        let result = sqlx::query(
            r#"
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
            "#).bind(&contest.id).bind(contest.start_epoch_second).bind(contest.duration_second).bind(&contest.title).bind(&contest.rate_change).bind(&contest.category)
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
