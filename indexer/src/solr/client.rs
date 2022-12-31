use crate::solr::core::SolrCore;
use crate::solr::models::*;
use reqwest::Client;
use url::Url;

type Result<T> = std::result::Result<T, SolrError>;

#[derive(Debug)]
pub struct SolrClient {
    url: String,
    client: Client,
}

impl SolrClient {
    pub fn new(url: &str, port: u32) -> Result<Self> {
        let url = Url::parse(url).map_err(|e| SolrError::UrlParseError(e))?;

        let scheme = url.scheme();
        let host = url.host_str().ok_or_else(|| SolrError::InvalidHostError)?;

        Ok(SolrClient {
            url: format!("{}://{}:{}", scheme, host, port),
            client: reqwest::Client::new(),
        })
    }

    pub async fn status(&self) -> Result<SolrSystemInfo> {
        let path = "solr/admin/info/system";

        let response = self
            .client
            .get(format!("{}/{}", self.url, path))
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrSystemInfo =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        Ok(result)
    }

    pub async fn cores(&self) -> Result<SolrCoreList> {
        let path = "solr/admin/cores";

        let response = self
            .client
            .get(format!("{}/{}", self.url, path))
            .send()
            .await
            .map_err(|e| SolrError::RequestError(e))?
            .text()
            .await
            .map_err(|e| SolrError::RequestError(e))?;

        let result: SolrCoreList =
            serde_json::from_str(&response).map_err(|e| SolrError::DeserializeError(e))?;

        Ok(result)
    }

    pub async fn core(&self, name: &str) -> Result<SolrCore> {
        let cores = self
            .cores()
            .await?
            .status
            .ok_or_else(|| SolrError::SpecifiedCoreNotFoundError)?;

        if !cores.contains_key(name) {
            return Err(SolrError::SpecifiedCoreNotFoundError);
        }

        Ok(SolrCore::new(name, &self.url))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    /// 通常系テスト
    #[test]
    fn test_create_solr_client() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();
        assert_eq!(client.url, "http://localhost:8983");
    }

    /// 冗長なURLを与えられたときのテスト
    /// 与えられたURLのスキーマとホストのみを読み取る仕様なので、冗長なURLを与えられても
    /// スキーマとホスト以外の情報は無視される。
    #[test]
    fn test_create_solr_client_with_redundant_url() {
        let client = SolrClient::new("http://localhost:8983/solr", 8983).unwrap();
        assert_eq!(client.url, "http://localhost:8983");
    }

    /// 異常系テスト
    #[test]
    fn test_create_solr_client_with_invalid_url() {
        let client = SolrClient::new("hogehoge", 3000);
        assert!(client.is_err());
    }

    /// 正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_status() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let response = client.status().await.unwrap();
        assert_eq!(response.header.status, 0);
    }

    /// 正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_cores() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let response = client.cores().await.unwrap();
        assert!(response.status.unwrap().contains_key("example"));
    }

    /// 正常系テスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_core() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let core = client.core("example").await.unwrap();
        assert_eq!(core.name, String::from("example"));
    }

    /// 存在しないコアを指定したときのテスト
    ///
    /// 以下のコマンドでDockerコンテナを起動してからテストを実行すること。
    ///
    /// ```ignore
    /// docker run --rm -d -p 8983:8983 solr:9.1.0 solr-precreate example
    /// ```
    #[tokio::test]
    #[ignore]
    async fn test_get_non_existent_core() {
        let client = SolrClient::new("http://localhost", 8983).unwrap();

        let core = client.core("hoge").await;
        assert!(core.is_err());
    }
}
