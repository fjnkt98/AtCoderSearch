use serde::Serialize;
use sqlx::{prelude::FromRow, Pool, Postgres};
use std::pin::Pin;
use tokio_stream::Stream;

use super::{ReadRows, ToDocument};

pub struct UserRowReader<'a> {
    pool: &'a Pool<Postgres>,
}

impl<'a> UserRowReader<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> Self {
        Self { pool }
    }
}

#[derive(Debug, Clone, PartialEq, PartialOrd, FromRow)]
pub struct User {
    user_id: String,
    rating: i32,
    highest_rating: i32,
    affiliation: Option<String>,
    birth_year: Option<i32>,
    country: Option<String>,
    crown: Option<String>,
    join_count: i32,
    rank: i32,
    active_rank: Option<i32>,
    wins: i32,
}

impl ToDocument for User {
    type Document = UserDoc;

    fn to_document(self) -> anyhow::Result<Self::Document> {
        Ok(UserDoc {
            user_id: self.user_id,
            rating: self.rating,
            highest_rating: self.highest_rating,
            affiliation: self.affiliation,
            birth_year: self.birth_year,
            country: self.country,
            crown: self.crown,
            join_count: self.join_count,
            rank: self.rank,
            active_rank: self.active_rank,
            wins: self.wins,
            color: rate_to_color(self.rating),
            highest_color: rate_to_color(self.highest_rating),
        })
    }
}

#[derive(Debug, Serialize, Clone, PartialEq, PartialOrd)]
#[serde(rename_all = "camelCase")]
pub struct UserDoc {
    user_id: String,
    rating: i32,
    highest_rating: i32,
    affiliation: Option<String>,
    birth_year: Option<i32>,
    country: Option<String>,
    crown: Option<String>,
    join_count: i32,
    rank: i32,
    active_rank: Option<i32>,
    wins: i32,
    color: String,
    highest_color: String,
}

fn rate_to_color(rate: i32) -> String {
    match rate {
        0..=399 => "gray",
        400..=799 => "brown",
        800..=1199 => "green",
        1200..=1599 => "cyan",
        1600..=1999 => "blue",
        2000..=2399 => "yellow",
        2400..=2799 => "orange",
        2800..=3199 => "red",
        3200..=3599 => "silver",
        _ => "gold",
    }
    .to_string()
}

impl<'a> ReadRows<'a> for UserRowReader<'a> {
    type Row = User;

    async fn read_rows(
        &'a self,
    ) -> anyhow::Result<Pin<Box<dyn Stream<Item = Result<Self::Row, sqlx::Error>> + Send + 'a>>>
    {
        let sql = r#"
    SELECT
        user_id,
        rating,
        highest_rating,
        affiliation,
        birth_year,
        country,
        crown,
        join_count,
        rank,
        active_rank,
        wins
    FROM
        "users"
        "#;
        let stream = sqlx::query_as(sql).fetch(self.pool);

        Ok(stream)
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::testutil::{create_container, create_pool_from_container};
    use futures::TryStreamExt;
    use testcontainers::runners::AsyncRunner;

    #[tokio::test]
    async fn test_read_user_rows() {
        let container = create_container().unwrap().start().await.unwrap();
        let pool = create_pool_from_container(&container).await.unwrap();

        sqlx::query(
            r#"
INSERT INTO "users" ("user_id", "rating", "highest_rating", "affiliation", "birth_year", "country", "crown", "join_count", "rank", "active_rank", "wins") VALUES
('user01', 1000, 1100, 'AtCoder', 1998, 'JP', 'gold', 10, 3000, 2500, 0),
('user02', 1000, 1100, NULL, NULL, NULL, NULL, 10, 3000, NULL, 0);
"#,
        )
        .execute(&pool)
        .await
        .unwrap();

        let reader = UserRowReader { pool: &pool };
        let stream = reader.read_rows().await.unwrap();

        let users: Vec<User> = stream.try_collect().await.unwrap();
        let want = vec![
            User {
                user_id: String::from("user01"),
                rating: 1000,
                highest_rating: 1100,
                affiliation: Some(String::from("AtCoder")),
                birth_year: Some(1998),
                country: Some(String::from("JP")),
                crown: Some(String::from("gold")),
                join_count: 10,
                rank: 3000,
                active_rank: Some(2500),
                wins: 0,
            },
            User {
                user_id: String::from("user02"),
                rating: 1000,
                highest_rating: 1100,
                affiliation: None,
                birth_year: None,
                country: None,
                crown: None,
                join_count: 10,
                rank: 3000,
                active_rank: None,
                wins: 0,
            },
        ];

        assert_eq!(users, want);
    }

    #[test]
    fn test_serialize_user_doc() {
        let users = vec![
            UserDoc {
                user_id: String::from("user01"),
                rating: 1000,
                highest_rating: 1100,
                affiliation: Some(String::from("AtCoder")),
                birth_year: Some(1998),
                country: Some(String::from("JP")),
                crown: Some(String::from("gold")),
                join_count: 10,
                rank: 3000,
                active_rank: Some(2500),
                wins: 0,
                color: String::from("green"),
                highest_color: String::from("green"),
            },
            UserDoc {
                user_id: String::from("user02"),
                rating: 1000,
                highest_rating: 1100,
                affiliation: None,
                birth_year: None,
                country: None,
                crown: None,
                join_count: 10,
                rank: 3000,
                active_rank: None,
                wins: 0,
                color: String::from("green"),
                highest_color: String::from("green"),
            },
        ];

        let result = serde_json::to_string(&users).unwrap();
        let want = String::from(
            r#"[{"userId":"user01","rating":1000,"highestRating":1100,"affiliation":"AtCoder","birthYear":1998,"country":"JP","crown":"gold","joinCount":10,"rank":3000,"activeRank":2500,"wins":0,"color":"green","highestColor":"green"},{"userId":"user02","rating":1000,"highestRating":1100,"affiliation":null,"birthYear":null,"country":null,"crown":null,"joinCount":10,"rank":3000,"activeRank":null,"wins":0,"color":"green","highestColor":"green"}]"#,
        );
        assert_eq!(result, want);
    }
}
