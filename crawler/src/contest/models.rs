use serde::Deserialize;
use sqlx::{FromRow, Type};

#[derive(Deserialize)]
pub struct ContestJson {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
}

#[derive(FromRow, Type)]
pub struct Contest {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
    pub category: String,
}
