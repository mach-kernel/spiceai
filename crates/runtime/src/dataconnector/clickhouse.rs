/*
Copyright 2024 The Spice.ai OSS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

use async_trait::async_trait;
use clickhouse_rs::ClientHandle;
use data_components::clickhouse::ClickhouseTableFactory;
use data_components::Read;
use datafusion::datasource::TableProvider;
use db_connection_pool::clickhousepool::ClickhouseConnectionPool;
use db_connection_pool::DbConnectionPool;
use secrets::Secret;
use snafu::prelude::*;
use spicepod::component::dataset::Dataset;
use std::any::Any;
use std::pin::Pin;
use std::sync::Arc;
use std::{collections::HashMap, future::Future};

use super::{DataConnector, DataConnectorFactory};

#[derive(Debug, Snafu)]
pub enum Error {
    #[snafu(display("Unable to create Clickhouse connection pool: {source}"))]
    UnableToCreateClickhouseConnectionPool { source: db_connection_pool::Error },

    #[snafu(display("{source}"))]
    UnableToGetReadProvider {
        source: Box<dyn std::error::Error + Send + Sync>,
    },
}

pub type Result<T, E = Error> = std::result::Result<T, E>;

pub struct Clickhouse {
    clickhouse_factory: ClickhouseTableFactory,
}

impl DataConnectorFactory for Clickhouse {
    fn create(
        secret: Option<Secret>,
        params: Arc<Option<HashMap<String, String>>>,
    ) -> Pin<Box<dyn Future<Output = super::NewDataConnectorResult> + Send>> {
        Box::pin(async move {
            let pool: Arc<dyn DbConnectionPool<ClientHandle, &'static (dyn Sync)> + Send + Sync> =
                Arc::new(
                    ClickhouseConnectionPool::new(params, secret)
                        .await
                        .context(UnableToCreateClickhouseConnectionPoolSnafu)?,
                );

            let clickhouse_factory = ClickhouseTableFactory::new(pool);

            Ok(Arc::new(Self { clickhouse_factory }) as Arc<dyn DataConnector>)
        })
    }
}

#[async_trait]
impl DataConnector for Clickhouse {
    fn as_any(&self) -> &dyn Any {
        self
    }

    async fn read_provider(
        &self,
        dataset: &Dataset,
    ) -> super::AnyErrorResult<Arc<dyn TableProvider>> {
        Ok(
            Read::table_provider(&self.clickhouse_factory, dataset.path().into())
                .await
                .context(UnableToGetReadProviderSnafu)?,
        )
    }
}
