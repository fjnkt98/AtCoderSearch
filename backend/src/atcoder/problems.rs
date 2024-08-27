use anyhow::Context as _;
use regex::Regex;
use std::{sync::LazyLock, time::Duration};

use reqwest::Client;

static AGC001_STARTED_AT: i64 = 1468670400;
static JOI_PATTERN: LazyLock<Regex> = LazyLock::new(|| Regex::new(r#"^(jag|JAG)"#).unwrap());
static MARATHON_PATTERN_1: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(r#"(^Chokudai self|ハーフマラソン|^HACK TO THE FUTURE|Asprova|Heuristics Contest)"#)
        .unwrap()
});
static MARATHON_PATTERN_2: LazyLock<Regex> =
    LazyLock::new(|| Regex::new(r#"(^future-meets-you-self|^hokudai-hitachi)"#).unwrap());
static MARATHON_PATTERN_3: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(r#"^(genocon2021|stage0-2021|caddi2019|pakencamp-2019-day2|kuronekoyamato-self2019|wn2017_1)$"#).unwrap()
});
static SPONSORED_PATTERN_1: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(
        r#"ドワンゴ|^Mujin|SoundHound|^codeFlyer|^COLOCON|みんなのプロコン|CODE THANKS FESTIVAL"#,
    )
    .unwrap()
});
static SPONSORED_PATTERN_2: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(
        r#"(CODE FESTIVAL|^DISCO|日本最強プログラマー学生選手権|全国統一プログラミング王|Indeed)"#,
    )
    .unwrap()
});
static SPONSORED_PATTERN_3: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(r#"(^Donuts|^dwango|^DigitalArts|^Code Formula|天下一プログラマーコンテスト)"#)
        .unwrap()
});

#[derive(Debug, PartialEq, PartialOrd)]
pub enum RatedType {
    Unrated,
    All,
    Range,
    UpperBound,
    LowerBound,
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct RatedTarget {
    typ: RatedType,
    from: Option<i64>,
    to: Option<i64>,
}

pub struct AtCoderProblemsClient {
    client: Client,
}

impl AtCoderProblemsClient {
    pub fn new() -> anyhow::Result<Self> {
        let client = Client::builder()
            .timeout(Duration::from_secs(30))
            .build()
            .with_context(|| "create http client")?;

        return Ok(Self { client });
    }

    pub async fn fetch_contests(&self) -> anyhow::Result<Vec<Contest>> {
        todo!();
    }

    pub async fn fetch_problems(&self) -> anyhow::Result<Vec<Problem>> {
        todo!();
    }

    pub async fn fetch_difficulties(&self) -> anyhow::Result<Vec<Difficulty>> {
        todo!();
    }
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct Contest {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
}

impl Contest {
    pub fn rated_target(&self) -> RatedTarget {
        todo!();
    }

    pub fn categorize(&self) -> String {
        todo!();
    }
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct Problem {
    pub id: String,
    pub contest_id: String,
    pub problem_index: String,
    pub name: String,
    pub title: String,
}

#[derive(Debug, PartialEq, PartialOrd)]
pub struct Difficulty {
    pub slope: Option<f64>,
    pub intercept: Option<f64>,
    pub variance: Option<f64>,
    pub difficulty: Option<i64>,
    pub discrimination: Option<f64>,
    pub irt_loglikelihood: Option<f64>,
    pub irt_users: Option<f64>,
    pub is_experimental: bool,
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::PathBuf;

    use super::*;
    use rstest::rstest;

    #[test]
    fn test_contest_rated_target() {
        let cases = [(
            Contest {
                id: String::from("test"),
                start_epoch_second: 1468670399,
                duration_second: 0,
                title: String::from("test"),
                rate_change: String::from("-"),
            },
            RatedTarget {
                typ: RatedType::Unrated,
                from: None,
                to: None,
            },
        )];

        for (c, t) in cases {
            assert_eq!(t, c.rated_target())
        }
    }
}
