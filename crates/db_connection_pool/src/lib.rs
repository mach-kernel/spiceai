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

use crate::dbconnection::DbConnection;
use async_trait::async_trait;
use spicepod::component::dataset::acceleration;

#[cfg(feature = "clickhouse")]
pub mod clickhousepool;
pub mod dbconnection;
#[cfg(feature = "duckdb")]
pub mod duckdbpool;
#[cfg(feature = "mysql")]
pub mod mysqlpool;
#[cfg(feature = "postgres")]
pub mod postgrespool;
#[cfg(feature = "sqlite")]
pub mod sqlitepool;

pub type Error = Box<dyn std::error::Error + Send + Sync>;
type Result<T, E = Error> = std::result::Result<T, E>;

#[async_trait]
pub trait DbConnectionPool<T, P: 'static> {
    async fn connect(&self) -> Result<Box<dyn DbConnection<T, P>>>;
}

#[derive(Default)]
pub enum Mode {
    #[default]
    Memory,
    File,
}

impl From<&str> for Mode {
    fn from(m: &str) -> Self {
        match m {
            "file" => Mode::File,
            "memory" => Mode::Memory,
            _ => Mode::default(),
        }
    }
}

impl From<acceleration::Mode> for Mode {
    fn from(m: acceleration::Mode) -> Self {
        match m {
            acceleration::Mode::File => Mode::File,
            acceleration::Mode::Memory => Mode::Memory,
        }
    }
}
