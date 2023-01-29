use crate::models::errors::UploadingError;
use solr_client::clients::core::SolrCore;
use std::path::{Path, PathBuf};
use tokio::fs;
use tokio::fs::File;
use tokio::io::AsyncReadExt;

type Result<T> = std::result::Result<T, UploadingError>;

pub struct DocumentUploader<'a> {
    savedir: PathBuf,
    core: &'a SolrCore,
}

impl<'a> DocumentUploader<'a> {
    pub fn new(savedir: &Path, core: &'a SolrCore) -> Self {
        Self {
            savedir: savedir.to_path_buf(),
            core: core,
        }
    }

    pub async fn upload(&self, optimize: bool) -> Result<()> {
        self.core.reload().await?;

        let mut files = fs::read_dir(&self.savedir)
            .await
            .map_err(|e| UploadingError::FileOperationError(e))?;
        let mut target = Vec::new();
        while let Some(file) = files
            .next_entry()
            .await
            .map_err(|e| UploadingError::FileOperationError(e))?
        {
            let path = file.path();
            if let Some(extension) = path.extension() {
                if extension == "json" {
                    target.push(path);
                }
            }
        }

        for filepath in target.into_iter() {
            let mut file = File::open(filepath).await?;
            let mut buffer: Vec<u8> = Vec::new();
            file.read_to_end(&mut buffer).await?;

            self.core.post(buffer).await?;
        }

        self.core.commit(optimize).await?;

        Ok(())
    }
}
