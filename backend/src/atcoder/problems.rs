use anyhow::Context as _;
use itertools::Itertools;
use regex::Regex;
use serde::Deserialize;
use std::{collections::BTreeMap, sync::LazyLock, time::Duration};

use reqwest::{Client, Url};

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

#[derive(Debug, PartialEq, PartialOrd, Clone)]
pub enum RatedTarget {
    Unrated,
    All,
    Range { from: i64, to: i64 },
    UpperBound { to: i64 },
    LowerBound { from: i64 },
}

pub struct AtCoderProblemsClient {
    client: Client,
}

impl AtCoderProblemsClient {
    pub fn new() -> anyhow::Result<Self> {
        let client = Client::builder()
            .timeout(Duration::from_secs(30))
            .gzip(true)
            .build()
            .with_context(|| "create http client")?;

        return Ok(Self { client });
    }

    pub async fn fetch_contests(&self) -> anyhow::Result<Vec<Contest>> {
        let url = Url::parse("https://kenkoooo.com/atcoder/resources/contests.json")
            .with_context(|| "parse contests.json url")?;
        let res = self
            .client
            .get(url)
            .send()
            .await
            .with_context(|| "request to contests.json")?;
        let contests: Vec<Contest> = res
            .json()
            .await
            .with_context(|| "deserialize contests json")?;

        Ok(contests)
    }

    pub async fn fetch_problems(&self) -> anyhow::Result<Vec<Problem>> {
        let url = Url::parse("https://kenkoooo.com/atcoder/resources/problems.json")
            .with_context(|| "parse problems.json url")?;
        let res = self
            .client
            .get(url)
            .send()
            .await
            .with_context(|| "request to problems.json")?;
        let problems: Vec<Problem> = res
            .json()
            .await
            .with_context(|| "deserialize problems json")?;

        Ok(problems)
    }

    pub async fn fetch_difficulties(&self) -> anyhow::Result<BTreeMap<String, Difficulty>> {
        let url = Url::parse("https://kenkoooo.com/atcoder/resources/problem-models.json")
            .with_context(|| "parse problem-models.json url")?;
        let res = self
            .client
            .get(url)
            .send()
            .await
            .with_context(|| "request to problem-models.json")?;
        let difficulties: BTreeMap<String, Difficulty> = res
            .json()
            .await
            .with_context(|| "deserialize problems json")?;

        Ok(difficulties)
    }
}

#[derive(Debug, PartialEq, PartialOrd, Deserialize, Clone)]
pub struct Contest {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
}

impl Contest {
    pub fn rated_target(&self) -> RatedTarget {
        if self.start_epoch_second < AGC001_STARTED_AT {
            return RatedTarget::Unrated;
        }

        match self.rate_change.as_str() {
            "-" => RatedTarget::Unrated,
            "All" => RatedTarget::All,
            _ => {
                if let Some((from, to)) = self.rate_change.split(" ~ ").collect_tuple() {
                    match (from.parse::<i64>(), to.parse::<i64>()) {
                        (Ok(from), Ok(to)) => RatedTarget::Range { from, to },
                        (Ok(from), Err(_)) => RatedTarget::LowerBound { from },
                        (Err(_), Ok(to)) => RatedTarget::UpperBound { to },
                        (Err(_), Err(_)) => RatedTarget::Unrated,
                    }
                } else {
                    RatedTarget::Unrated
                }
            }
        }
    }

    pub fn categorize(&self) -> String {
        if self.id.starts_with("abc") {
            return String::from("ABC");
        }
        if self.id.starts_with("arc") {
            return String::from("ARC");
        }
        if self.id.starts_with("agc") {
            return String::from("AGC");
        }
        if self.id.starts_with("ahc") {
            return String::from("AHC");
        }

        match self.rated_target() {
            RatedTarget::All => String::from("AGC-Like"),
            RatedTarget::UpperBound { .. } => String::from("ABC-Like"),
            RatedTarget::LowerBound { .. } => String::from("ARC-Like"),
            RatedTarget::Unrated => {
                if self.id.starts_with("past") {
                    return String::from("PAST");
                }
                if self.id.starts_with("joi") {
                    return String::from("JOI");
                }
                if JOI_PATTERN.is_match(&self.id) {
                    return String::from("JOI");
                }

                if MARATHON_PATTERN_1.is_match(&self.id)
                    || MARATHON_PATTERN_2.is_match(&self.id)
                    || MARATHON_PATTERN_3.is_match(&self.id)
                {
                    return String::from("Marathon");
                }

                if SPONSORED_PATTERN_1.is_match(&self.title)
                    || SPONSORED_PATTERN_2.is_match(&self.title)
                    || SPONSORED_PATTERN_3.is_match(&self.title)
                {
                    return String::from("Other Sponsored");
                }

                return String::from("Other Contests");
            }
            _ => String::from("Other Contests"),
        }
    }
}

#[derive(Debug, PartialEq, PartialOrd, Deserialize, Clone)]
pub struct Problem {
    pub id: String,
    pub contest_id: String,
    pub problem_index: String,
    pub name: String,
    pub title: String,
}

impl Problem {
    pub fn url(&self) -> String {
        format!(
            "https://atcoder.jp/contests/{}/tasks/{}",
            self.contest_id, self.id
        )
    }
}

#[derive(Debug, PartialEq, PartialOrd, Deserialize, Clone)]
pub struct Difficulty {
    pub slope: Option<f64>,
    pub intercept: Option<f64>,
    pub variance: Option<f64>,
    pub difficulty: Option<i64>,
    pub discrimination: Option<f64>,
    pub irt_loglikelihood: Option<f64>,
    pub irt_users: Option<f64>,
    pub is_experimental: Option<bool>,
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::PathBuf;

    use super::*;
    use rstest::rstest;

    #[test]
    fn test_new_atcoder_problems_client() {
        AtCoderProblemsClient::new().unwrap();
    }

    #[tokio::test]
    async fn test_fetch_contests() {
        let client = AtCoderProblemsClient::new().unwrap();
        client.fetch_contests().await.unwrap();
    }

    #[tokio::test]
    async fn test_fetch_problems() {
        let client = AtCoderProblemsClient::new().unwrap();
        client.fetch_problems().await.unwrap();
    }

    #[tokio::test]
    async fn test_fetch_difficulties() {
        let client = AtCoderProblemsClient::new().unwrap();
        client.fetch_difficulties().await.unwrap();
    }

    #[test]
    fn test_contest_rated_target() {
        let cases = [
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670399,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from(""),
                },
                RatedTarget::Unrated,
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from("-"),
                },
                RatedTarget::Unrated,
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from("All"),
                },
                RatedTarget::All,
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from(" ~ 1199"),
                },
                RatedTarget::UpperBound { to: 1199 },
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from(" ~ 2799"),
                },
                RatedTarget::UpperBound { to: 2799 },
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from("1200 ~ "),
                },
                RatedTarget::LowerBound { from: 1200 },
            ),
            (
                Contest {
                    id: String::from("test"),
                    start_epoch_second: 1468670401,
                    duration_second: 0,
                    title: String::from("test"),
                    rate_change: String::from("1200 ~ 2799"),
                },
                RatedTarget::Range {
                    from: 1200,
                    to: 2799,
                },
            ),
        ];

        for (c, t) in cases {
            assert_eq!(t, c.rated_target())
        }
    }

    #[test]
    fn test_contest_categorize() {
        let cases = [
            (
                Contest {
                    id: String::from("abc042"),
                    start_epoch_second: 1469275200,
                    duration_second: 6000,
                    title: String::from("AtCoder Beginner Contest 042"),
                    rate_change: String::from(" ~ 1199"),
                },
                "ABC",
            ),
            (
                Contest {
                    id: String::from("zone2021"),
                    duration_second: 6000,
                    rate_change: String::from(" ~ 1999"),
                    start_epoch_second: 1619870400,
                    title: String::from("ZONeエナジー プログラミングコンテスト  “HELLO SPACE”"),
                },
                "ABC-Like",
            ),
            (
                Contest {
                    id: String::from("jsc2019-final"),
                    duration_second: 10800,
                    rate_change: String::from("-"),
                    start_epoch_second: 1569728700,
                    title: String::from("第一回日本最強プログラマー学生選手権決勝"),
                },
                "Other Sponsored",
            ),
            (
                Contest {
                    id: String::from("ttpc2019"),
                    duration_second: 18000,
                    rate_change: String::from("-"),
                    start_epoch_second: 1567224300,
                    title: String::from("東京工業大学プログラミングコンテスト2019"),
                },
                "Other Contests",
            ),
        ];

        for (c, want) in cases {
            assert_eq!(want, c.categorize())
        }
    }

    #[test]
    fn test_problem_url() {
        let p = Problem {
            id: String::from("abc001_a"),
            contest_id: String::from("abc001"),
            problem_index: String::from("A"),
            name: String::from("problem name"),
            title: String::from("problem title"),
        };
        assert_eq!("https://atcoder.jp/contests/abc001/tasks/abc001_a", p.url());
    }
}
