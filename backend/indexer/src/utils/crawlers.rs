use crate::models::{
    contest::ContestJson,
    errors::CrawlingError,
    problem::{ProblemDifficulty, ProblemJson},
    tables::Contest,
};
use minify_html::{minify, Cfg};
use reqwest::{header::ACCEPT_ENCODING, Client};
use sqlx::{
    postgres::{PgRow, Postgres},
    Pool, Row,
};
use std::collections::{HashMap, HashSet};
use tokio::time::{self, Duration};

type Result<T> = std::result::Result<T, CrawlingError>;

pub struct ContestCrawler<'a> {
    url: String,
    pool: &'a Pool<Postgres>,
}

impl<'a> ContestCrawler<'a> {
    pub fn new(pool: &'a Pool<Postgres>) -> Self {
        ContestCrawler {
            url: String::from("https://kenkoooo.com/atcoder/resources/contests.json"),
            pool: pool,
        }
    }

    /// AtCoderProblemsからコンテスト情報を取得するメソッド
    pub async fn get_contest_list(&self) -> Result<Vec<ContestJson>> {
        tracing::info!("Start to retrieve contests information from AtCoder Problems.");
        let client = reqwest::Client::new();

        let response = client
            .get(&self.url)
            // AtCoderProblemsはAccept-Encodingにgzipを指定しないとダウンロードできない(https://twitter.com/kenkoooo/status/1147352545133645824)
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let json = response
            .text()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let contests: Vec<ContestJson> =
            serde_json::from_str(&json).map_err(|e| CrawlingError::DeserializeError(e))?;
        tracing::info!(
            "{} contests information successfully retrieved.",
            contests.len()
        );

        Ok(contests)
    }

    /// AtCoderProblemsから取得したコンテスト情報からデータベースへ格納する用のモデルを作って返すメソッド
    pub async fn crawl(&self) -> Result<Vec<Contest>> {
        tracing::info!("Start to crawl contests information.");
        let contests: Vec<Contest> = self
            .get_contest_list()
            .await?
            .iter()
            .map(|contest| Contest {
                id: contest.id.clone(),
                start_epoch_second: contest.start_epoch_second.clone(),
                duration_second: contest.duration_second.clone(),
                title: contest.title.clone(),
                rate_change: contest.rate_change.clone(),
                category: contest.categorize(),
            })
            .collect();
        tracing::info!(
            "{} contests information successfully crawled.",
            contests.len()
        );

        Ok(contests)
    }

    /// コンテスト情報をデータベースへ保存するメソッド
    ///
    /// データの保存にMERGE INTO文(PostgreSQL 15から)を使用している
    /// コンテスト情報の存在判定にIDを使用し、IDが存在すればUPDATE、IDが存在しなければINSERTを実行する
    /// UPDATE時はすべての情報をUPDATEするようにしている
    pub async fn save(&self, contests: &Vec<Contest>) -> Result<()> {
        tracing::info!("Start to save contests information.");
        // トランザクション開始
        let mut tx = self.pool.begin().await?;

        // 各コンテスト情報を一つずつ処理する
        for contest in contests.iter() {
            tracing::info!("Start to save {} into database...", contest.id);
            let result = sqlx::query("
                MERGE INTO contests
                USING
                    (VALUES($1, $2, $3, $4, $5, $6)) AS contest(id, start_epoch_second, duration_second, title, rate_change, category)
                ON
                    contests.id = contest.id
                WHEN MATCHED THEN
                    UPDATE SET (id, start_epoch_second, duration_second, title, rate_change, category) = (contest.id, contest.start_epoch_second, contest.duration_second, contest.title, contest.rate_change, contest.category)
                WHEN NOT MATCHED THEN
                    INSERT (id, start_epoch_second, duration_second, title, rate_change, category)
                    VALUES (contest.id, contest.start_epoch_second, contest.duration_second, contest.title, contest.rate_change, contest.category);
                ")
                .bind(&contest.id)
                .bind(&contest.start_epoch_second)
                .bind(&contest.duration_second)
                .bind(&contest.title)
                .bind(&contest.rate_change)
                .bind(&contest.category)
                .execute(&mut tx)
                .await;

            // エラーが発生したらトランザクションをロールバックしてエラーを早期リターンする
            if let Err(e) = result {
                tracing::error!("An error occurred at saving {:?}.", contest);
                tx.rollback().await?;
                return Err(CrawlingError::SqlExecutionError(e));
            }

            tracing::info!("Contest {} was saved.", contest.id);
        }

        tx.commit().await?;
        tracing::info!("{} contests successfully saved.", contests.len());

        Ok(())
    }

    /// コンテスト情報の取得からデータベースへの保存までの一連の処理を行うメソッド
    pub async fn run(&self) -> Result<()> {
        let contests = self.crawl().await?;
        self.save(&contests).await?;

        Ok(())
    }
}
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

        tracing::info!("Attempting to get problem list from AtCoder Problems...");
        let response = client
            .get(&self.url)
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let json = response
            .text()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let problems: Vec<ProblemJson> =
            serde_json::from_str(&json).map_err(|e| CrawlingError::DeserializeError(e))?;

        tracing::info!("{} problems collected.", problems.len());

        Ok(problems)
    }

    /// 問題ページをクロールしてHTML情報を取得するメソッド
    ///
    /// クロール間隔は300msにしてある。
    ///
    /// - target: クロール対象の問題のリスト
    pub async fn crawl(&self, url: &str, client: &Client, config: &Cfg) -> Result<String> {
        tracing::info!("Crawl {}", url);
        let html = client
            .get(url)
            .send()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?
            .bytes()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;
        let html =
            String::from_utf8(minify(&html, config)).map_err(|e| CrawlingError::ParseError(e))?;

        Ok(html)
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
            .await
            .map_err(|e| CrawlingError::SqlExecutionError(e))?
            .iter()
            .cloned(),
        );

        let target: Vec<ProblemJson> = self
            .get_problem_list()
            .await?
            .into_iter()
            .filter(|problem| !exists_problems.contains(&problem.id))
            .collect();

        tracing::info!("{} problems are now target for collection.", target.len());

        Ok(target)
    }

    /// 問題の難易度情報を取得してハッシュマップとして返すメソッド
    async fn get_difficulties(&self) -> Result<HashMap<String, ProblemDifficulty>> {
        let client = reqwest::Client::new();

        tracing::info!("Attempting to get difficulties from AtCoder Problems...");

        let response = client
            .get("https://kenkoooo.com/atcoder/resources/problem-models.json")
            .header(ACCEPT_ENCODING, "gzip")
            .send()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let json = response
            .text()
            .await
            .map_err(|e| CrawlingError::RequestError(e))?;

        let difficulties: HashMap<String, ProblemDifficulty> =
            serde_json::from_str(&json).map_err(|e| CrawlingError::DeserializeError(e))?;

        Ok(difficulties)
    }

    /// 問題データをデータベースに格納するメソッド
    pub async fn save(&self, targets: &Vec<ProblemJson>, duration: Duration) -> Result<()> {
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
        let client = Client::new();

        let difficulties = self.get_difficulties().await?;

        for problem in targets.iter() {
            let mut tx = self.pool.begin().await?;

            let difficulty = difficulties
                .get(&problem.id)
                .and_then(|difficulty| difficulty.difficulty);
            let url = format!(
                "https://atcoder.jp/contests/{}/tasks/{}",
                problem.contest_id, problem.id
            );
            let html = self.crawl(&url, &client, &config).await?;

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
                .bind(&url)
                .bind(html)
                .bind(difficulty)
                .execute(&mut tx)
                .await;

            match result {
                Ok(_) => {
                    tracing::info!("Problem {} was saved.", problem.id);
                    tx.commit().await?;
                }
                Err(e) => {
                    tracing::error!("An error occurred at {:?}", problem.id);
                    tx.rollback().await?;
                    return Err(CrawlingError::SqlExecutionError(e));
                }
            }

            time::sleep(duration).await;
        }

        Ok(())
    }

    /// 問題情報の取得からデータベースへの保存までの一連の処理を行うメソッド
    ///
    /// - allがtrueのときはすべての問題を対象にクロールを行う
    /// - allがfalseのときは差分取得のみを行う
    pub async fn run(&self, all: bool, duration: Duration) -> Result<()> {
        let targets = if all {
            self.get_problem_list().await?
        } else {
            self.detect_diff().await?
        };

        self.save(&targets, duration).await?;

        Ok(())
    }
}
