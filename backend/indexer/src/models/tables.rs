use sqlx::{FromRow, Type};

/// データベースに格納するコンテスト情報のモデル
///
/// - category: コンテストのカテゴリ。e.g. ABC, ARG
#[derive(FromRow, Type, Debug)]
pub struct Contest {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
    pub category: String,
}

/// データベースに格納する問題情報のモデル
///
/// - url: 問題のページのURL
/// - html: 問題のページのHTML
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
