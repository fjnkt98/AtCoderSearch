use crate::models::errors::UploadingError;
use futures::stream::FuturesUnordered;
use futures::StreamExt;
use solrust::client::core::SolrCore;
use std::path::{Path, PathBuf};
use tokio::fs;

type Result<T> = std::result::Result<T, UploadingError>;

pub struct DocumentUploader {
    savedir: PathBuf,
    core: SolrCore,
}

impl DocumentUploader {
    pub fn new(savedir: &Path, core: &SolrCore) -> Self {
        Self {
            savedir: savedir.to_path_buf(),
            core: core.clone(),
        }
    }

    pub async fn upload(&self, optimize: bool) -> Result<()> {
        self.core.reload().await?;

        let mut files = fs::read_dir(&self.savedir)
            .await
            .map_err(|e| UploadingError::FileOperationError(e))?;

        let mut tasks = FuturesUnordered::new();
        while let Some(file) = files
            .next_entry()
            .await
            .map_err(|e| UploadingError::FileOperationError(e))?
        {
            let path = file.path();
            let core = self.core.clone();
            if let Some(extension) = path.extension() {
                if extension == "json" {
                    let metadata = file
                        .metadata()
                        .await
                        .map_err(|e| UploadingError::FileOperationError(e))?;
                    tracing::info!(
                        "Processing file: {}, size: {}kB",
                        path.display(),
                        metadata.len() / 1024
                    );
                    let task = tokio::spawn(async move {
                        let content = fs::read(path).await.unwrap();
                        core.post(content).await.unwrap();
                    });
                    tasks.push(task);
                }
            }
        }
        while let Some(task) = tasks.next().await {
            if let Err(e) = task {
                self.core.rollback().await?;
                return Err(UploadingError::UnexpectedError(e));
            }
        }

        self.core.commit(optimize).await?;

        Ok(())
    }
}
