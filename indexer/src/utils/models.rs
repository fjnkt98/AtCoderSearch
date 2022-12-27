use anyhow::{anyhow, Result};
use serde::{Deserialize, Serialize};

#[derive(sqlx::FromRow)]
pub struct Record {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub difficulty: i32,
    pub start_at: i32,
    pub duration: i32,
    pub rate_change: String,
    pub category: String,
    pub html: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Document {
    pub problem_id: String,
    pub problem_title: String,
    pub problem_url: String,
    pub contest_id: String,
    pub contest_title: String,
    pub contest_url: String,
    pub difficulty: i32,
    pub start_at: String,
    pub duration: i32,
    pub rate_change: String,
    pub category: String,
    pub text_ja: Vec<String>,
    pub text_en: Vec<String>,
}

pub struct DocumentBuilder {
    problem_id: Option<String>,
    problem_title: Option<String>,
    problem_url: Option<String>,
    contest_id: Option<String>,
    contest_title: Option<String>,
    contest_url: Option<String>,
    difficulty: Option<i32>,
    start_at: Option<String>,
    duration: Option<i32>,
    rate_change: Option<String>,
    category: Option<String>,
    text_ja: Option<Vec<String>>,
    text_en: Option<Vec<String>>,
}

impl DocumentBuilder {
    pub fn new() -> Self {
        DocumentBuilder {
            problem_id: None,
            problem_title: None,
            problem_url: None,
            contest_id: None,
            contest_title: None,
            contest_url: None,
            difficulty: None,
            start_at: None,
            duration: None,
            rate_change: None,
            category: None,
            text_ja: None,
            text_en: None,
        }
    }

    pub fn build(self) -> Result<Document> {
        let problem_id = self
            .problem_id
            .ok_or_else(|| anyhow!("problem_id is mandatory."))?;
        let problem_title = self
            .problem_title
            .ok_or_else(|| anyhow!("problem_title is mandatory."))?;
        let problem_url = self
            .problem_url
            .ok_or_else(|| anyhow!("problem_url is mandatory."))?;
        let contest_id = self
            .contest_id
            .ok_or_else(|| anyhow!("contest_id is mandatory."))?;
        let contest_title = self
            .contest_title
            .ok_or_else(|| anyhow!("contest_title is mandatory."))?;
        let contest_url = self
            .contest_url
            .ok_or_else(|| anyhow!("contest_url is mandatory."))?;
        let difficulty = self
            .difficulty
            .ok_or_else(|| anyhow!("difficulty is mandatory."))?;
        let start_at = self
            .start_at
            .ok_or_else(|| anyhow!("start_at is mandatory."))?;
        let duration = self
            .duration
            .ok_or_else(|| anyhow!("duration is mandatory."))?;
        let rate_change = self
            .rate_change
            .ok_or_else(|| anyhow!("rate_change is mandatory."))?;
        let category = self
            .category
            .ok_or_else(|| anyhow!("category is mandatory."))?;
        let text_ja = self
            .text_ja
            .ok_or_else(|| anyhow!("text_ja is mandatory."))?;
        let text_en = self
            .text_en
            .ok_or_else(|| anyhow!("text_en is mandatory."))?;

        Ok(Document {
            problem_id: problem_id,
            problem_title: problem_title,
            problem_url: problem_url,
            contest_id: contest_id,
            contest_title: contest_title,
            contest_url: contest_url,
            difficulty: difficulty,
            start_at: start_at,
            duration: duration,
            rate_change: rate_change,
            category: category,
            text_ja: text_ja,
            text_en: text_en,
        })
    }

    pub fn problem_id(mut self, problem_id: String) -> Self {
        self.problem_id = Some(problem_id);
        self
    }

    pub fn problem_title(mut self, problem_title: String) -> Self {
        self.problem_title = Some(problem_title);
        self
    }

    pub fn problem_url(mut self, problem_url: String) -> Self {
        self.problem_url = Some(problem_url);
        self
    }

    pub fn contest_id(mut self, contest_id: String) -> Self {
        self.contest_id = Some(contest_id);
        self
    }

    pub fn contest_title(mut self, contest_title: String) -> Self {
        self.contest_title = Some(contest_title);
        self
    }

    pub fn contest_url(mut self, contest_url: String) -> Self {
        self.contest_url = Some(contest_url);
        self
    }

    pub fn difficulty(mut self, difficulty: i32) -> Self {
        self.difficulty = Some(difficulty);
        self
    }

    pub fn start_at(mut self, start_at: String) -> Self {
        self.start_at = Some(start_at);
        self
    }

    pub fn duration(mut self, duration: i32) -> Self {
        self.duration = Some(duration);
        self
    }

    pub fn rate_change(mut self, rate_change: String) -> Self {
        self.rate_change = Some(rate_change);
        self
    }

    pub fn category(mut self, category: String) -> Self {
        self.category = Some(category);
        self
    }

    pub fn text_ja(mut self, text_ja: Vec<String>) -> Self {
        self.text_ja = Some(text_ja);
        self
    }

    pub fn text_en(mut self, text_en: Vec<String>) -> Self {
        self.text_en = Some(text_en);
        self
    }
}
