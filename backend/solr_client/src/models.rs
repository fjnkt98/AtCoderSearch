/// 命名規則
///
/// - Solrから返ってくるレスポンス本体のモデル -> SolrXXXXResponse
/// - レスポンスの一部 -> SolrXXXX(Header|Body|Info|e.t.c)
///
use chrono::{DateTime, FixedOffset, Utc};
use itertools::Itertools;
use serde::de::{Error, Unexpected};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use serde_json::Value;
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrResponseHeader {
    pub status: u32,
    #[serde(alias = "QTime")]
    pub qtime: u32,
    pub params: Option<HashMap<String, Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrErrorInfo {
    pub metadata: Vec<String>,
    pub msg: String,
    pub code: u32,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct LuceneInfo {
    #[serde(alias = "solr-spec-version")]
    pub solr_spec_version: String,
    #[serde(alias = "solr-impl-version")]
    pub solr_impl_version: String,
    #[serde(alias = "lucene-spec-version")]
    pub lucene_spec_version: String,
    #[serde(alias = "lucene-impl-version")]
    pub lucene_impl_version: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSystemInfo {
    #[serde(alias = "responseHeader")]
    pub header: SolrResponseHeader,
    pub mode: String,
    pub solr_home: String,
    pub core_root: String,
    pub lucene: LuceneInfo,
    pub jvm: Value,
    pub security: Value,
    pub system: Value,
    pub error: Option<SolrErrorInfo>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct IndexInfo {
    #[serde(alias = "numDocs")]
    pub num_docs: u64,
    #[serde(alias = "maxDoc")]
    pub max_doc: u64,
    #[serde(alias = "deletedDocs")]
    pub deleted_docs: u64,
    pub version: u64,
    #[serde(alias = "segmentCount")]
    pub segment_count: u64,
    pub current: bool,
    #[serde(alias = "hasDeletions")]
    pub has_deletions: bool,
    pub directory: String,
    #[serde(alias = "segmentsFile")]
    pub segments_file: String,
    #[serde(alias = "segmentsFileSizeInBytes")]
    pub segments_file_size_in_bytes: u64,
    #[serde(alias = "userData")]
    pub user_data: Value,
    #[serde(alias = "sizeInBytes")]
    pub size_in_bytes: u64,
    pub size: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SolrCoreStatus {
    pub name: String,
    #[serde(alias = "instanceDir")]
    pub instance_dir: String,
    #[serde(alias = "dataDir")]
    pub data_dir: String,
    pub config: String,
    pub schema: String,
    #[serde(alias = "startTime")]
    pub start_time: String,
    pub uptime: u64,
    pub index: IndexInfo,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrCoreList {
    #[serde(alias = "responseHeader")]
    pub header: SolrResponseHeader,
    #[serde(alias = "initFailures")]
    pub init_failures: Value,
    pub status: Option<HashMap<String, SolrCoreStatus>>,
    pub error: Option<SolrErrorInfo>,
}

impl SolrCoreList {
    pub fn as_vec(&self) -> Option<Vec<String>> {
        if let Some(cores) = &self.status {
            Some(cores.keys().cloned().collect())
        } else {
            return None;
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSimpleResponse {
    #[serde(alias = "responseHeader")]
    pub header: SolrResponseHeader,
    pub error: Option<SolrErrorInfo>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSelectResponse {
    #[serde(alias = "responseHeader")]
    pub header: SolrResponseHeader,
    pub response: SolrSelectBody,
    pub facet_counts: Option<SolrFacetBody>,
    pub error: Option<SolrErrorInfo>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrSelectBody {
    #[serde(alias = "numFound")]
    pub num_found: u32,
    pub start: u32,
    #[serde(alias = "numFoundExact")]
    pub num_found_exact: bool,
    // TODO: ジェネリクス化
    pub docs: Value,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrFacetBody {
    pub facet_queries: Value,
    #[serde(deserialize_with = "deserialize_facet_fields")]
    pub facet_fields: HashMap<String, Vec<(String, u32)>>,
    #[serde(deserialize_with = "deserialize_facet_ranges")]
    pub facet_ranges: HashMap<String, RangeFacetKind>,
    pub facet_intervals: Value,
    pub facet_heatmaps: Value,
}

fn deserialize_facet_fields<'de, D>(
    deserializer: D,
) -> Result<HashMap<String, Vec<(String, u32)>>, D::Error>
where
    D: Deserializer<'de>,
{
    let value: HashMap<String, Vec<Value>> = Deserialize::deserialize(deserializer)?;
    let value: HashMap<String, Vec<(String, u32)>> = value
        .iter()
        .map(|(k, v)| {
            (
                k.to_string(),
                v.iter()
                    .tuples()
                    .map(|(v1, v2)| {
                        (
                            v1.as_str().unwrap_or("").to_string(),
                            v2.as_u64().unwrap_or(0) as u32,
                        )
                    })
                    .collect::<Vec<(String, u32)>>(),
            )
        })
        .collect();

    Ok(value)
}

/// レンジファセットの結果をデシリアライズする関数
///
/// レンジファセットの結果の型はフィールドの型依存なのでアドホックに場合分けする必要があった
fn deserialize_facet_ranges<'de, D>(
    deserializer: D,
) -> Result<HashMap<String, RangeFacetKind>, D::Error>
where
    D: Deserializer<'de>,
{
    let value: HashMap<String, Value> = Deserialize::deserialize(deserializer)?;
    let mut result: HashMap<String, RangeFacetKind> = HashMap::new();
    for (field, value) in value.iter() {
        match &value["start"] {
            Value::Number(start) => {
                if start.is_i64() {
                    let value: IntegerRangeFacet =
                        serde_json::from_value(value.clone()).map_err(|e| {
                            D::Error::custom(format!(
                                "Failed to parse integer range facet result. [{}]",
                                e.to_string()
                            ))
                        })?;
                    result.insert(field.to_string(), RangeFacetKind::Integer(value));
                } else {
                    let value: FloatRangeFacet =
                        serde_json::from_value(value.clone()).map_err(|e| {
                            D::Error::custom(format!(
                                "Failed to parse float range facet result. [{}]",
                                e.to_string()
                            ))
                        })?;
                    result.insert(field.to_string(), RangeFacetKind::Float(value));
                }
            }
            Value::String(start) => {
                if DateTime::parse_from_rfc3339(&start.replace("Z", "+00:00")).is_ok() {
                    let value: DateTimeRangeFacet =
                        serde_json::from_value(value.clone()).map_err(|e| {
                            D::Error::custom(format!(
                                "Failed to parse datetime range facet result. [{}]",
                                e.to_string()
                            ))
                        })?;
                    result.insert(field.to_string(), RangeFacetKind::DateTime(value));
                } else {
                    // TODO; 数値、日付型以外のレンジファセットがあったら処理を追加する
                    return Err(D::Error::custom("Unexpected range facet value type."));
                }
            }
            _ => {
                return Err(D::Error::custom("Mismatched range facet value type."));
            }
        }
    }
    Ok(result)
}

#[derive(Serialize, Deserialize, Debug)]
pub enum RangeFacetKind {
    Integer(IntegerRangeFacet),
    Float(FloatRangeFacet),
    DateTime(DateTimeRangeFacet),
}

#[derive(Serialize, Deserialize, Debug)]
pub struct IntegerRangeFacet {
    #[serde(deserialize_with = "deserialize_range_facet_counts")]
    pub counts: Vec<(String, u32)>,
    pub gap: i64,
    pub start: i64,
    pub end: i64,
    pub before: Option<i64>,
    pub after: Option<i64>,
    pub between: Option<i64>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct FloatRangeFacet {
    #[serde(deserialize_with = "deserialize_range_facet_counts")]
    pub counts: Vec<(String, u32)>,
    pub gap: f64,
    pub start: f64,
    pub end: f64,
    pub before: Option<f64>,
    pub after: Option<f64>,
    pub between: Option<f64>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DateTimeRangeFacet {
    #[serde(deserialize_with = "deserialize_range_facet_counts")]
    pub counts: Vec<(String, u32)>,
    #[serde(
        serialize_with = "serialize_datetime",
        deserialize_with = "deserialize_datetime"
    )]
    pub gap: DateTime<FixedOffset>,
    #[serde(
        serialize_with = "serialize_datetime",
        deserialize_with = "deserialize_datetime"
    )]
    pub start: DateTime<FixedOffset>,
    #[serde(
        serialize_with = "serialize_datetime",
        deserialize_with = "deserialize_datetime"
    )]
    pub end: DateTime<FixedOffset>,
    #[serde(
        serialize_with = "serialize_optional_datetime",
        deserialize_with = "deserialize_optional_datetime"
    )]
    pub before: Option<DateTime<FixedOffset>>,
    #[serde(
        serialize_with = "serialize_optional_datetime",
        deserialize_with = "deserialize_optional_datetime"
    )]
    pub after: Option<DateTime<FixedOffset>>,
    #[serde(
        serialize_with = "serialize_optional_datetime",
        deserialize_with = "deserialize_optional_datetime"
    )]
    pub between: Option<DateTime<FixedOffset>>,
}

/// DateTime型をSolrが扱える文字列形式にシリアライズする関数
fn serialize_datetime<S>(timestamp: &DateTime<FixedOffset>, s: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    s.serialize_str(
        &timestamp
            .with_timezone(&Utc)
            .to_rfc3339()
            .replace("+00:00", "Z"),
    )
}

/// DateTime型をSolrが扱える文字列形式にシリアライズする関数(Option版)
/// TODO: コードの重複が多いのでスマートに解決する方法を探す
fn serialize_optional_datetime<S>(
    timestamp: &Option<DateTime<FixedOffset>>,
    s: S,
) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    match timestamp {
        Some(timestamp) => s.serialize_str(
            &timestamp
                .with_timezone(&Utc)
                .to_rfc3339()
                .replace("+00:00", "Z"),
        ),
        None => s.serialize_none(),
    }
}

/// Solrのタイムスタンプ文字列表現をDateTime型にデシリアライズする関数
fn deserialize_datetime<'de, D>(deserializer: D) -> Result<DateTime<FixedOffset>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    if let Ok(timestamp) = DateTime::parse_from_rfc3339(&value.replace("Z", "+00:00")) {
        return Ok(timestamp);
    } else {
        return Err(Error::invalid_value(
            Unexpected::Str(&value),
            &"Invalid timestamp string",
        ));
    }
}

/// Solrのタイムスタンプ文字列表現をDateTime型にデシリアライズする関数(Option版)
/// TODO: コードの重複が多いのでスマートに解決する方法を探す
fn deserialize_optional_datetime<'de, D>(
    deserializer: D,
) -> Result<Option<DateTime<FixedOffset>>, D::Error>
where
    D: Deserializer<'de>,
{
    let value = String::deserialize(deserializer)?;
    if let Ok(timestamp) = DateTime::parse_from_rfc3339(&value.replace("Z", "+00:00")) {
        return Ok(Some(timestamp));
    } else {
        return Ok(None);
    }
}

/// ファセットの結果の配列をRustが扱える配列にデシリアライズする関数
///
/// Solrのファセットの結果は「文字列、数値」が交互に格納された配列で返ってくる。Rustは型が混じった配列を扱えないので、タプルのリストに変換する。
fn deserialize_range_facet_counts<'de, D>(deserializer: D) -> Result<Vec<(String, u32)>, D::Error>
where
    D: Deserializer<'de>,
{
    let value: Vec<Value> = Deserialize::deserialize(deserializer)?;
    let value: Vec<(String, u32)> = value
        .iter()
        .tuples()
        .map(|(v1, v2)| {
            (
                v1.as_str().unwrap_or("").to_string(),
                v2.as_u64().unwrap_or(0) as u32,
            )
        })
        .collect();

    Ok(value)
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisBody {
    pub field_types: HashMap<String, SolrAnalysisField>,
    pub field_names: HashMap<String, SolrAnalysisField>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisField {
    pub index: Option<Vec<Value>>,
    pub query: Option<Vec<Value>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SolrAnalysisResponse {
    #[serde(alias = "responseHeader")]
    pub header: SolrResponseHeader,
    pub analysis: SolrAnalysisBody,
    pub error: Option<SolrErrorInfo>,
}

#[cfg(test)]
mod test {
    use super::*;
    // TODO: テスト書く
}
