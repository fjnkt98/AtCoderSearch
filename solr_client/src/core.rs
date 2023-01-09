use crate::models::*;
use reqwest::header::CONTENT_TYPE;
use reqwest::Client;
use serde::Serialize;
use serde_json::Value;

type Result<T> = std::result::Result<T, SolrError>;

pub struct SolrCore {
    pub name: String,
    pub base_url: String,
    pub core_url: String,
    client: Client,
}

impl SolrCore {
    pub fn new(name: &str, base_url: &str) -> Self {
        let core_url = format!("{}/solr/{}", base_url, name);

        SolrCore {
            name: String::from(name),
            base_url: String::from(base_url),
            core_url: core_url,
            client: reqwest::Client::new(),
        }
    }
    pub async fn status(&self) -> Result<CoreStatus> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.base_url, path))
            .query(&[("action", "status"), ("core", &self.name)])
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let response: SolrCoreList =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        if let Some(error) = response.error {
            return Err(SolrError::UnexpectedError((error.code, error.msg)));
        }

        // コアオブジェクトが作成できた時点で
        //
        // 1. レスポンスのJSONに`status`フィールドが存在すること
        // 2. `status`フィールドのキーにこのコアが含まれていること
        //
        // が保証されているので、`unwrap()`を使用している。
        let status = response.status.unwrap().get(&self.name).unwrap().clone();

        Ok(status)
    }

    pub async fn reload(&self) -> Result<u32> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.base_url, path))
            .query(&[("action", "reload"), ("core", &self.name)])
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let response: SolrSimpleResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        if let Some(error) = response.error {
            return Err(SolrError::UnexpectedError((error.code, error.msg)));
        }

        Ok(response.header.status)
    }

    pub async fn select<S, T>(&self, params: &Vec<(S, T)>) -> Result<SolrSelectResponse>
    where
        S: Serialize,
        T: Serialize,
    {
        let response = self
            .client
            .get(format!("{}/select", self.core_url))
            .query(params)
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let response: SolrSelectResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        if let Some(error) = response.error {
            return Err(SolrError::UnexpectedError((error.code, error.msg)));
        }

        Ok(response)
    }

    pub async fn analyze(&self, word: &str, field: &str, analyzer: &str) -> Result<Vec<String>> {
        let params = [("analysis.fieldvalue", word), ("analysis.fieldtype", field)];

        let response = self
            .client
            .get(format!("{}/analysis/field", self.core_url))
            .query(&params)
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrAnalysisResponse =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        let result = result.analysis.field_types.get(field).unwrap();
        let result = match analyzer {
            "index" => result.index.as_ref().unwrap(),
            "query" => result.query.as_ref().unwrap(),
            _ => return Err(SolrError::InvalidValueError),
        };
        let result = result.last().unwrap().clone();

        let result = match result {
            Value::Array(array) => array
                .iter()
                .map(|e| e["text"].to_string().trim_matches('"').to_string())
                .collect::<Vec<String>>(),
            _ => Vec::new(),
        };

        Ok(result)
    }

    pub async fn post(&self, body: Vec<u8>) -> Result<()> {
        let response = self
            .client
            .post(format!("{}/update", self.core_url))
            .header(CONTENT_TYPE, "application/json")
            .body(body)
            .send()
            .await?;

        if response.status().as_u16() != 200 {
            return Err(SolrError::CorePostError);
        }

        Ok(())
    }

    pub async fn commit(&self, optimize: bool) -> Result<()> {
        if optimize {
            self.post(br#"{"optimize": {}}"#.to_vec()).await?;
        } else {
            self.post(br#"{"commit": {}}"#.to_vec()).await?;
        }

        Ok(())
    }

    pub async fn rollback(&self) -> Result<()> {
        self.post(br#"{"rollback": {}}"#.to_vec()).await?;

        Ok(())
    }

    pub async fn truncate(&self) -> Result<()> {
        self.post(br#"{"delete":{"query": "*:*"}}"#.to_vec())
            .await?;

        Ok(())
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use chrono::{DateTime, Utc};

    /// コアのステータス取得メソッドの正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_status() {
        let core = SolrCore::new("example", "http://localhost:8983");
        let status = core.status().await.unwrap();

        assert_eq!(status.name, String::from("example"));
    }

    /// コアのリロードメソッドの正常系テスト
    ///
    /// コアのリロード実行時の時刻と、リロード後のコアのスタートタイムの差が1秒以内なら
    /// リロードが実行されたと判断する。
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_reload() {
        let core = SolrCore::new("example", "http://localhost:8983");

        let before = Utc::now();

        core.reload().await.unwrap();

        let status = core.status().await.unwrap();
        let after = status.start_time.replace("Z", "+00:00");
        let after = DateTime::parse_from_rfc3339(&after)
            .unwrap()
            .with_timezone(&Utc);

        assert!(before < after);

        let duration = (after - before).num_milliseconds();
        assert!(duration.abs() < 1000);
    }

    #[tokio::test]
    #[ignore]
    async fn test_select_in_normal() {
        let core = SolrCore::new("example", "http://localhost:8983");

        let params = vec![("q", "*:*")];
        let response = core.select(&params).await.unwrap();

        assert_eq!(response.header.status, 0);
    }

    #[tokio::test]
    #[ignore]
    async fn test_select_in_non_normal() {
        let core = SolrCore::new("example", "http://localhost:8983");

        let params = vec![("q", "text_hoge:*")];
        let response = core.select(&params).await;

        assert!(response.is_err());
    }

    /// 単語の解析メソッドの正常系テスト
    ///
    /// とりあえずエラーが出ないことを確認する。
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_analyze() {
        let core = SolrCore::new("example", "http://localhost:8983");

        let word = "solr-client";
        let expected = vec![String::from("solr"), String::from("client")];

        let actual = core.analyze(word, "text_en", "index").await.unwrap();

        assert_eq!(expected, actual);
    }
}
