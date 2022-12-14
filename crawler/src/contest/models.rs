use regex::Regex;
use serde::Deserialize;
use sqlx::{FromRow, Type};

#[derive(PartialEq, Debug)]
pub enum RatedTargetType {
    ALL,
    UNRATED,
    UPPERBOUND(i64),
    LOWERBOUND(i64),
}

const AGC001_STARTED_AT: i64 = 1468670400;

#[derive(Deserialize)]
pub struct ContestJson {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
}

#[derive(FromRow, Type)]
pub struct Contest {
    pub id: String,
    pub start_epoch_second: i64,
    pub duration_second: i64,
    pub title: String,
    pub rate_change: String,
    pub category: String,
}

impl ContestJson {
    fn rated_target(&self) -> RatedTargetType {
        if self.start_epoch_second < AGC001_STARTED_AT {
            return RatedTargetType::UNRATED;
        }

        match self.rate_change.as_str() {
            "-" => RatedTargetType::UNRATED,
            "All" => RatedTargetType::ALL,
            _ => {
                let range = self
                    .rate_change
                    .split('~')
                    .map(|word| word.trim())
                    .collect::<Vec<&str>>();
                if range.len() != 2 {
                    return RatedTargetType::UNRATED;
                }

                if let Ok(lower_bound) = range[0].parse::<i64>() {
                    return RatedTargetType::LOWERBOUND(lower_bound);
                }
                if let Ok(upper_bound) = range[1].parse::<i64>() {
                    return RatedTargetType::UPPERBOUND(upper_bound);
                }

                return RatedTargetType::UNRATED;
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
            RatedTargetType::ALL => {
                return String::from("AGC-Like");
            }
            RatedTargetType::UPPERBOUND(_) => {
                return String::from("ABC-Like");
            }
            RatedTargetType::LOWERBOUND(_) => {
                return String::from("ARC-Like");
            }
            RatedTargetType::UNRATED => {
                if self.id.starts_with("past") {
                    return String::from("PAST");
                }
                if self.id.starts_with("joi") {
                    return String::from("JOI");
                }
                if Regex::new("^(jag|JAG)").unwrap().is_match(&self.id) {
                    return String::from("JAG");
                }
                if Regex::new(
                "(^Chokudai self|ハーフマラソン|^HACK TO THE FUTURE|Asprova|Heuristics Contest)",
            )
            .unwrap()
            .is_match(&self.title)
                || Regex::new("(^future-meets-you-self|^hokudai-hitachi)")
                    .unwrap()
                    .is_match(&self.id)
                || [
                    "genocon2021",
                    "stage0-2021",
                    "caddi2019",
                    "pakencamp-2019-day2",
                    "kuronekoyamato-self2019",
                    "wn2017_1",
                ]
                .contains(&&self.id.as_str())
            {
                return String::from("Marathon");
            }

                let pattern1 = Regex::new("ドワンゴ|^Mujin|SoundHound|^codeFlyer|^COLOCON|みんなのプロコン|CODE THANKS FESTIVAL").unwrap();
                let pattern2 = Regex::new("(CODE FESTIVAL|^DISCO|日本最強プログラマー学生選手権|全国統一プログラミング王|Indeed)").unwrap();
                let pattern3 = Regex::new(
                    "(^Donuts|^dwango|^DigitalArts|^Code Formula|天下一プログラマーコンテスト)",
                )
                .unwrap();

                if pattern1.is_match(&self.title)
                    || pattern2.is_match(&self.title)
                    || pattern3.is_match(&self.title)
                {
                    return String::from("Other Sponsored");
                }

                return String::from("Other Contests");
            }
        };
    }
}

#[cfg(test)]
mod tests {
    use super::ContestJson;
    use super::*;

    #[test]
    fn when_is_unrated() {
        let contest = ContestJson {
            id: String::from("abc001"),
            start_epoch_second: 1381579200,
            duration_second: 7200,
            title: String::from("AtCoder Beginner Contest 001"),
            rate_change: String::from("-"),
        };

        assert_eq!(contest.rated_target(), RatedTargetType::UNRATED);
    }

    #[test]
    fn when_abc() {
        let contest = ContestJson {
            id: String::from("abc042"),
            start_epoch_second: 1469275200,
            duration_second: 6000,
            title: String::from("AtCoder Beginner Contest 042"),
            rate_change: String::from(" ~ 1199"),
        };

        assert_eq!(contest.categorize(), String::from("ABC"));
    }

    #[test]
    fn when_abc_like() {
        let contest = ContestJson {
            id: String::from("zone2021"),
            start_epoch_second: 1619870400,
            duration_second: 6000,
            title: String::from("ZONeエナジー プログラミングコンテスト  “HELLO SPACE”"),
            rate_change: String::from(" ~ 1999"),
        };

        assert_eq!(contest.categorize(), String::from("ABC-Like"));
    }

    #[test]
    fn when_other_sponsored() {
        let contest = ContestJson {
            id: String::from("jsc2019-final"),
            start_epoch_second: 1569728700,
            duration_second: 10800,
            title: String::from("第一回日本最強プログラマー学生選手権決勝"),
            rate_change: String::from("-"),
        };

        assert_eq!(contest.categorize(), String::from("Other Sponsored"));
    }

    #[test]
    fn when_other_contests() {
        let contest = ContestJson {
            id: String::from("ttpc2019"),
            start_epoch_second: 1567224300,
            duration_second: 18000,
            title: String::from("東京工業大学プログラミングコンテスト2019"),
            rate_change: String::from("-"),
        };

        assert_eq!(contest.categorize(), String::from("Other Contests"));
    }
}
