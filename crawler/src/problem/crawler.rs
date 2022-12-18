use crate::problem::models::{Problem, ProblemDifficulty, ProblemJson};
use anyhow::{Context, Error, Result};
use minify_html::{minify, Cfg};
use reqwest;
use reqwest::header::ACCEPT_ENCODING;
use serde_json;
use sqlx::postgres::{PgRow, Postgres};
use sqlx::{Pool, Row};
use std::collections::{HashMap, HashSet};
use tokio::time::{sleep, Duration};

pub struct ProblemCrawler<'a> {
    url: String,
    pool: &'a Pool<Postgres>,
}

impl<'a> ProblemCrawler<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> Self {
        ProblemCrawler {
            url: String::from("https://kenkoooo.com/atcoder/resources/problems.json"),
            pool: pool,
        }
    }

    /// AtCoder Problemsから問題情報の一覧を取得するメソッド
    pub async fn get_problem_list(&self) -> Result<Vec<ProblemJson>> {
        let client = reqwest::Client::new();

        let response = client
            .get(&self.url)
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .context("Failed to get contest information from AtCoder Problems.")?;

        let json = response
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let problems: Vec<ProblemJson> =
            serde_json::from_str(&json).context("Failed to parse JSON body.")?;

        Ok(problems)
    }

    /// 問題ページをクロールしてHTML情報を取得するメソッド
    ///
    /// クロール間隔は300msにしてある。
    ///
    /// - target: クロール対象の問題のリスト
    pub async fn crawl(&self, target: &Vec<ProblemJson>) -> Result<Vec<Problem>> {
        let client = reqwest::Client::new();

        let config = Cfg {
            do_not_minify_doctype: true,
            ensure_spec_compliant_unquoted_attribute_values: false,
            keep_closing_tags: true,
            keep_html_and_head_opening_tags: false,
            keep_spaces_between_attributes: false,
            keep_comments: false,
            minify_css: true,
            minify_js: true,
            remove_bangs: false,
            remove_processing_instructions: false,
        };

        let mut problems: Vec<Problem> = Vec::new();
        for problem in target.iter() {
            let url = format!(
                "https://atcoder.jp/contests/{}/tasks/{}",
                problem.contest_id, problem.id
            );
            let html = client.get(&url).send().await?.bytes().await?;
            let html =
                String::from_utf8(minify(&html, &config)).context("Failed to minify HTML.")?;

            problems.push(Problem {
                id: problem.id.clone(),
                contest_id: problem.contest_id.clone(),
                problem_index: problem.problem_index.clone(),
                name: problem.name.clone(),
                title: problem.title.clone(),
                url: url,
                html: html,
            });

            tracing::info!("Problem {} is collected.", problem.id);

            sleep(Duration::from_millis(200)).await;
        }

        Ok(problems)
    }

    /// AtCoder Problemsから得た一覧情報とデータベースにある情報を比較し、
    /// 未取得の問題を検出するメソッド
    pub async fn detect_diff(&self) -> Result<Vec<ProblemJson>> {
        let exists_problems: HashSet<String> = HashSet::from_iter(
            sqlx::query(
                r#"
            SELECT id FROM problems;
            "#,
            )
            .map(|row: PgRow| row.get(0))
            .fetch_all(self.pool)
            .await?
            .iter()
            .cloned(),
        );

        let target: Vec<ProblemJson> = self
            .get_problem_list()
            .await?
            .into_iter()
            .filter(|problem| !exists_problems.contains(&problem.id))
            .collect();

        Ok(target)
    }

    async fn get_difficulties(&self) -> Result<HashMap<String, ProblemDifficulty>> {
        let client = reqwest::Client::new();

        let response = client
            .get("https://kenkoooo.com/atcoder/resources/problem-models.json")
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .context("Failed to get contest information from AtCoder Problems.")?;

        let json = response
            .text()
            .await
            .context("Failed to get JSON body from response.")?;

        let difficulties: HashMap<String, ProblemDifficulty> =
            serde_json::from_str(&json).context("Failed to parse JSON body")?;

        Ok(difficulties)
    }

    /// 問題データをデータベースに格納するメソッド
    pub async fn save(&self, problems: &Vec<Problem>) -> Result<()> {
        let mut tx = self.pool.begin().await?;

        let difficulties = self.get_difficulties().await?;

        for problem in problems.iter() {
            let difficulty = if difficulties.contains_key(&problem.id) {
                match difficulties.get(&problem.id).unwrap().difficulty {
                    Some(difficulty) => difficulty,
                    None => 0,
                }
            } else {
                0
            };

            let result = sqlx::query(r"
                MERGE INTO problems
                USING
                    (VALUES($1, $2, $3, $4, $5, $6, $7, $8)) AS problem(id, contest_id, problem_index, name, title, url, html, difficulty)
                ON
                    problems.id = problem.id
                WHEN MATCHED THEN
                    UPDATE SET (id, contest_id, problem_index, name, title, url, html, difficulty) = (problem.id, problem.contest_id, problem.problem_index, problem.name, problem.title, problem.url, problem.html, problem.difficulty)
                WHEN NOT MATCHED THEN
                    INSERT (id, contest_id, problem_index, name, title, url, html, difficulty)
                    VALUES (problem.id, problem.contest_id, problem.problem_index, problem.name, problem.title, problem.url, problem.html, problem.difficulty);
                ")
                .bind(&problem.id)
                .bind(&problem.contest_id)
                .bind(&problem.problem_index)
                .bind(&problem.name)
                .bind(&problem.title)
                .bind(&problem.url)
                .bind(&problem.html)
                .bind(difficulty)
                .execute(&mut tx)
                .await;

            if let Err(e) = result {
                tx.rollback().await?;
                return Err(Error::new(e));
            }

            tracing::info!("Problem {} was saved.", problem.id);
        }

        tx.commit().await?;

        Ok(())
    }
}
