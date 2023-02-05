use crate::models::errors::GeneratingError;
use crate::utils::extractor::FullTextExtractor;
use chrono::{DateTime, NaiveDateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_with::serde_as;
use solr_client::models::datetime::SolrDateTime;

type Result<T> = std::result::Result<T, GeneratingError>;

#[derive(sqlx::FromRow)]
pub struct Record {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub difficulty: i32,
    pub start_at: i64,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
    pub html: String,
}

impl Record {
    pub fn to_document(self, extractor: &FullTextExtractor) -> Result<Document> {
        let start_at = DateTime::<Utc>::from_utc(
            NaiveDateTime::from_timestamp_opt(self.start_at, 0).unwrap(),
            Utc,
        );

        let (text_ja, text_en) = extractor.extract(&self.html)?;

        let contest_url: String = format!("https://atcoder.jp/contests/{}", self.contest_id);

        let document = Document {
            problem_id: self.problem_id,
            problem_title: self.problem_title,
            problem_url: self.problem_url,
            contest_id: self.contest_id,
            contest_title: self.contest_title,
            contest_url: contest_url,
            difficulty: self.difficulty,
            start_at: start_at,
            duration: self.duration,
            rate_change: self.rate_change,
            category: self.category,
            text_ja: text_ja,
            text_en: text_en,
        };

        Ok(document)

        // todo!();
    }
}

#[serde_as]
#[derive(Serialize, Deserialize, Debug)]
pub struct Document {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub contest_url: String,
    pub difficulty: i32,
    #[serde(default)]
    #[serde_as(as = "SolrDateTime")]
    pub start_at: DateTime<Utc>,
    pub duration: i64,
    pub rate_change: String,
    pub category: String,
    pub text_ja: Vec<String>,
    pub text_en: Vec<String>,
}
