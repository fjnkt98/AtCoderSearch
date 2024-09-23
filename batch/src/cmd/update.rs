use clap::Subcommand;
use meilisearch_sdk::client::Client;
use sqlx::postgres::PgPoolOptions;

use super::CommonArgs;
use crate::update::{update_index, user::UserRowReader};

#[derive(Subcommand)]
pub(crate) enum UpdateCommands {
    Problem,
    User {
        #[arg(long, default_value_t = 1000)]
        chunk_size: usize,
    },
}

impl UpdateCommands {
    pub async fn exec(&self, args: &CommonArgs) -> anyhow::Result<()> {
        let pool = PgPoolOptions::new()
            .max_connections(8)
            .connect(&args.database_url)
            .await?;
        let client = Client::new(&args.engine_url, Some(&args.engine_master_key))?;

        match self {
            UpdateCommands::Problem => todo!(),
            UpdateCommands::User { chunk_size } => {
                let reader = UserRowReader::new(pool.clone());

                update_index(reader, client, *chunk_size).await?;
                Ok(())
            }
        }
    }
}
