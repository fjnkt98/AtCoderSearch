use crate::models::SearchParams;
use chrono::{Local, TimeZone, Utc};
use itertools::Itertools;
use serde::ser::Error;
use serde::{Deserialize, Serialize, Serializer};
use solrust::types::response::{SolrFacetBody, SolrRangeFacetKind};

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
    pub params: SearchParams,
    pub facet: FacetResult,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FacetResult {
    pub category: FieldFacetResultPart,
    pub difficulty: RangeFacetResultPart,
}

impl FacetResult {
    pub fn empty() -> Self {
        Self {
            category: FieldFacetResultPart::empty(),
            difficulty: RangeFacetResultPart::empty(),
        }
    }
}

impl From<Option<SolrFacetBody>> for FacetResult {
    fn from(facet: Option<SolrFacetBody>) -> FacetResult {
        let result = if let Some(facet) = facet {
            let category = if let Some(counts) = facet.facet_fields.get("category") {
                FieldFacetResultPart {
                    counts: counts
                        .iter()
                        .cloned()
                        .map(|(key, count)| FieldFacetCount { key, count })
                        .collect_vec(),
                }
            } else {
                FieldFacetResultPart::empty()
            };
            let difficulty = if let Some(ranges) = facet.facet_ranges.get("difficulty") {
                match ranges {
                    SolrRangeFacetKind::Integer(range) => {
                        let mut counts = range.counts.clone();
                        counts.push((range.end.to_string(), 0));

                        RangeFacetResultPart {
                            counts: counts
                                .iter()
                                .cloned()
                                .tuple_windows()
                                .map(|(begin, end)| RangeFacetCount {
                                    begin: begin.0,
                                    end: end.0,
                                    count: begin.1,
                                })
                                .collect(),
                            start: range.start.to_string(),
                            end: range.end.to_string(),
                            gap: range.gap.to_string(),
                            before: range.before.and_then(|before| Some(before.to_string())),
                            after: range.after.and_then(|after| Some(after.to_string())),
                            between: range.after.and_then(|between| Some(between.to_string())),
                        }
                    }
                    _ => RangeFacetResultPart::empty(),
                }
            } else {
                RangeFacetResultPart::empty()
            };
            FacetResult {
                category,
                difficulty,
            }
        } else {
            FacetResult::empty()
        };

        result
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FieldFacetResultPart {
    pub counts: Vec<FieldFacetCount>,
}

impl FieldFacetResultPart {
    pub fn empty() -> Self {
        Self { counts: Vec::new() }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FieldFacetCount {
    pub key: String,
    pub count: u32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RangeFacetResultPart {
    pub counts: Vec<RangeFacetCount>,
    pub start: String,
    pub end: String,
    pub gap: String,
    pub before: Option<String>,
    pub after: Option<String>,
    pub between: Option<String>,
}

impl RangeFacetResultPart {
    pub fn empty() -> Self {
        Self {
            counts: Vec::new(),
            start: String::from(""),
            end: String::from(""),
            gap: String::from(""),
            before: None,
            after: None,
            between: None,
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RangeFacetCount {
    pub begin: String,
    pub end: String,
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
