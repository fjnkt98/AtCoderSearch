use chrono::{DateTime, Local, TimeZone, Utc};
use itertools::Itertools;
use serde::{
    de::{Error as deError, Unexpected},
    ser::Error as seError,
};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use solrust::types::response::{SolrFacetBody, SolrRangeFacetKind};
use std::collections::BTreeMap;

/// レスポンスのボディに乗せるJSONのスキーマ
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultResponse {
    pub stats: SearchResultStats,
    pub items: Vec<Document>,
    pub message: Option<String>,
}

/// 検索結果の統計情報
/// 総ヒット数、表示開始位置、表示ドキュメント数、ファセット情報、処理時間等。
/// エラー発生時のエラーメッセージもここに含まれる。
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResultStats {
    pub time: u32,
    pub total: u32,
    pub index: u32,
    pub pages: u32,
    pub count: u32,
    pub facet: FacetResult,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResult(BTreeMap<String, FacetResultPart>);

impl From<Option<SolrFacetBody>> for FacetResult {
    fn from(facet: Option<SolrFacetBody>) -> FacetResult {
        let mut result: BTreeMap<String, FacetResultPart> = BTreeMap::new();
        if let Some(facet) = facet {
            for (key, value) in facet.facet_fields.iter() {
                result.insert(
                    key.clone(),
                    FacetResultPart {
                        counts: value
                            .iter()
                            .map(|(key, count)| FacetCount {
                                key: key.clone(),
                                count: count.clone(),
                            })
                            .collect(),
                        range_info: None,
                    },
                );
            }

            for (key, value) in facet.facet_ranges.iter() {
                result.insert(
                    key.clone(),
                    match value {
                        SolrRangeFacetKind::Integer(count) => FacetResultPart {
                            counts: count
                                .counts
                                .iter()
                                .tuple_windows()
                                .map(|(begin, end)| FacetCount {
                                    key: format!("{} ~ {}", begin.0, end.0),
                                    count: begin.1,
                                })
                                .collect(),
                            range_info: Some(RangeFacetInfo {
                                start: count.start.to_string(),
                                end: count.end.to_string(),
                                gap: count.gap.to_string(),
                                before: count.before.and_then(|before| Some(before.to_string())),
                                after: count.after.and_then(|after| Some(after.to_string())),
                                between: count.after.and_then(|between| Some(between.to_string())),
                            }),
                        },
                        SolrRangeFacetKind::Float(count) => FacetResultPart {
                            counts: count
                                .counts
                                .iter()
                                .tuple_windows()
                                .map(|(begin, end)| FacetCount {
                                    key: format!("{} ~ {}", begin.0, end.0),
                                    count: begin.1,
                                })
                                .collect(),
                            range_info: Some(RangeFacetInfo {
                                start: count.start.to_string(),
                                end: count.end.to_string(),
                                gap: count.gap.to_string(),
                                before: count.before.and_then(|before| Some(before.to_string())),
                                after: count.after.and_then(|after| Some(after.to_string())),
                                between: count.after.and_then(|between| Some(between.to_string())),
                            }),
                        },
                        SolrRangeFacetKind::DateTime(count) => FacetResultPart {
                            counts: count
                                .counts
                                .iter()
                                .tuple_windows()
                                .map(|(begin, end)| FacetCount {
                                    key: format!("{} ~ {}", begin.0, end.0),
                                    count: begin.1,
                                })
                                .collect(),
                            range_info: Some(RangeFacetInfo {
                                start: count.start.format("%Y-%m-%dT%H:%M:%S%:z").to_string(),
                                end: count.end.format("%Y-%m-%dT%H:%M:%S%:z").to_string(),
                                gap: count.gap.clone(),
                                before: count.before.and_then(|before| {
                                    Some(before.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                                }),
                                after: count.after.and_then(|after| {
                                    Some(after.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                                }),
                                between: count.after.and_then(|between| {
                                    Some(between.format("%Y-%m-%dT%H:%M:%S%:z").to_string())
                                }),
                            }),
                        },
                    },
                );
            }
        }

        FacetResult(result)
    }
}

/// ファセット結果を格納するフィールドのスキーマ
/// フィールドファセットもレンジファセットも同じスキーマにしている
/// (なのでフィールドファセットの場合startやend等のフィールドは完全に無駄になる。どうにかしたい)
#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResultPart {
    pub counts: Vec<FacetCount>,
    pub range_info: Option<RangeFacetInfo>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RangeFacetInfo {
    pub start: String,
    pub end: String,
    pub gap: String,
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
    #[serde(serialize_with = "serialize")]
    pub start_at: i64,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
}

fn serialize<S>(value: &i64, serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    let value = Utc
        .timestamp_opt(*value, 0)
        .single()
        .and_then(|d| Some(d.with_timezone(&Local)))
        .ok_or_else(|| S::Error::custom("Failed to deserialize Unix epoch time."))?;
    serializer.serialize_str(&value.to_rfc3339())
}

fn deserialize<'de, D>(deserializer: D) -> Result<DateTime<Local>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    let timestamp = value.parse::<i64>().map_err(|e| {
        D::Error::invalid_value(
            Unexpected::Str(&e.to_string()),
            &"Valid Unix epoch timestamp",
        )
    })?;
    let datetime = Utc
        .timestamp_opt(timestamp, 0)
        .single()
        .and_then(|d| Some(d.with_timezone(&Local)))
        .ok_or_else(|| D::Error::custom("Failed to deserialize Unix epoch time."))?;
    Ok(datetime)
}
