use serde::Deserialize;

/// AtCoderProblemsから取得できる問題情報のJSONスキーマ
///
/// - id: 問題ID。URIに使用されている文字列
/// - contest_id: 問題が所属するコンテストのID
/// - problem_index: 問題番号。e.g. A, B, Ex
/// - name: 問題名
/// - title: 問題番号と問題名をつなげた文字列。
#[derive(Deserialize, Clone)]
pub struct ProblemJson {
    pub id: String,
    pub contest_id: String,
    pub problem_index: String,
    pub name: String,
    pub title: String,
}

/// 問題の難易度情報
/// `https://kenkoooo.com/atcoder/resources/problem-models.json`から得られるJSONスキーマ
///
/// 問題によっては難易度情報が無いことがあるので、Option型でフィールドを定義している
#[derive(Deserialize)]
pub struct ProblemDifficulty {
    pub slope: Option<f64>,
    pub intercept: Option<f64>,
    pub variance: Option<f64>,
    pub difficulty: Option<i64>,
    pub discrimination: Option<f64>,
    pub irt_loglikelihood: Option<f64>,
    pub irt_users: Option<i64>,
    pub is_experimental: Option<bool>,
}
