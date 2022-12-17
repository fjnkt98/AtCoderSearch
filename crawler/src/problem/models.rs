use serde::Deserialize;
use sqlx::{FromRow, Type};

#[derive(Deserialize, Clone)]
pub struct ProblemJson {
    pub id: String,
    pub contest_id: String,
    pub problem_index: String,
    pub name: String,
    pub title: String,
}

#[derive(FromRow, Type, Clone)]
pub struct Problem {
    pub id: String,
    pub contest_id: String,
    pub problem_index: String,
    pub name: String,
    pub title: String,
    pub url: String,
    pub html: String,
}
