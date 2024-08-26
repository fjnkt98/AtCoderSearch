use sqlx::{self, postgres::Postgres, Pool};

pub async fn crawl_contest(pool: Pool<Postgres>) -> anyhow::Result<()> {
    todo!()
}
