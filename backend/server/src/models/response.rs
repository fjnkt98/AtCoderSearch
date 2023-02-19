use chrono::{DateTime, Local};
use serde::de::{Error, Unexpected};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use std::collections::HashMap;

/// レスポンスのボディに乗せるJSONのスキーマ
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultResponse {
    pub stats: SearchResultStats,
    pub items: SearchResultBody,
}

/// 検索結果の統計情報
/// 総ヒット数、表示開始位置、表示ドキュメント数、ファセット情報、処理時間等。
/// エラー発生時のエラーメッセージもここに含まれる。
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultStats {
    pub time: u32,
    pub message: Option<String>,
    pub total: u32,
    pub offset: u32,
    pub amount: u32,
    pub facet: HashMap<String, FacetResult>,
}

/// 検索にヒットしたドキュメント
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultBody {
    pub docs: Vec<Document>,
}

/// ファセット結果を格納するフィールドのスキーマ
/// フィールドファセットもレンジファセットも同じスキーマにしている
/// (なのでフィールドファセットの場合startやend等のフィールドは完全に無駄になる。どうにかしたい)
#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResult {
    pub counts: Vec<FacetCount>,
    pub start: Option<String>,
    pub end: Option<String>,
    pub gap: Option<String>,
    pub before: Option<String>,
    pub after: Option<String>,
    pub between: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FacetCount {
    pub key: String,
    pub count: u32,
}

/// 検索結果として返すドキュメントのスキーマ
#[derive(Debug, Serialize, Deserialize)]
pub struct Document {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub contest_url: String,
    pub difficulty: i32,
    #[serde(serialize_with = "serialize", deserialize_with = "deserialize")]
    pub start_at: DateTime<Local>,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
}

fn serialize<S>(value: &DateTime<Local>, serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    serializer.serialize_str(&value.to_rfc3339())
}

fn deserialize<'de, D>(deserializer: D) -> Result<DateTime<Local>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    if let Ok(timestamp) = DateTime::parse_from_rfc3339(&value) {
        return Ok(timestamp.with_timezone(&Local));
    } else {
        return Err(Error::invalid_value(
            Unexpected::Str(&value),
            &"Invalid timestamp string",
        ));
    }
}
