use crate::utils::models::*;
use sqlx::postgres::Postgres;
use sqlx::Pool;
use tokio::macros::support::Pin;
use tokio_stream::Stream;

type Result<T> = std::result::Result<T, IndexingError>;

pub struct RecordReader<'a> {
    pool: &'a Pool<Postgres>,
}

impl<'a> RecordReader<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> Self {
        RecordReader { pool: pool }
    }

    pub async fn read_rows(
        &self,
    ) -> Result<Pin<Box<dyn Stream<Item = std::result::Result<Record, sqlx::Error>> + Send + 'a>>>
    {
        let stream = sqlx::query_as::<_, Record>(
            "
            SELECT
                problems.id AS problem_id,
                problems.title AS problem_title,
                problems.url AS problem_url,
                contests.id AS contest_id,
                contests.title AS contest_title,
                problems.difficulty AS difficulty,
                contests.start_epoch_second AS start_at,
                contests.duration_second AS duration,
                contests.rate_change AS rate_change,
                contests.category AS category,
                problems.html AS html
            FROM
                problems
                LEFT JOIN contests ON problems.contest_id = contests.id;
            ",
        )
        .fetch(self.pool);

        Ok(stream)
    }
}
