use solr_client::clients::client::SolrClientError;
use solr_client::clients::core::SolrCoreError;
use std::string::FromUtf8Error;
use thiserror::Error;
use tokio::task::JoinError;
use url::ParseError;

#[derive(Debug, Error)]
pub enum CrawlingError {
    #[error("Failed to get information from AtCoder Problems")]
    RequestError(#[from] reqwest::Error),
    #[error("Failed to deserialize JSON data")]
    DeserializeError(#[from] serde_json::error::Error),
    #[error("Failed to execute SQL query")]
    SqlExecutionError(#[from] sqlx::Error),
    #[error("Failed to parse HTML")]
    ParseError(#[from] FromUtf8Error),
}

#[derive(Debug, Error)]
pub enum GeneratingError {
    #[error("Failed to execute SQL query")]
    SqlExecutionError(#[from] sqlx::Error),
    #[error("Failed to create selector")]
    SelectorError(#[from] scraper::error::SelectorErrorKind<'static>),
    #[error("Failed to create regular expression pattern")]
    RegexError(#[from] regex::Error),
    // #[error("Field value is mandatory")]
    // FieldValueNotConfiguredError,
    #[error("Failed to serialize JSON data")]
    SerializeError(#[from] serde_json::Error),
    #[error("Failed to operate file")]
    FileOperationError(#[from] std::io::Error),
    #[error("Failed to parse url")]
    UrlParseError(#[from] ParseError),
}

#[derive(Debug, Error)]
pub enum UploadingError {
    #[error("Failed to operate Solr client")]
    SolrClientError(#[from] SolrClientError),
    #[error("Failed to operate Solr core")]
    SolrCoreError(#[from] SolrCoreError),
    #[error("Failed to operate file")]
    FileOperationError(#[from] std::io::Error),
    #[error("Unexpected error occurred")]
    UnexpectedError(JoinError),
}
