use serde::Serialize;
use std::fmt::{Display, Formatter};
use std::ops;

//// Solrへのクエリに付加するクエリストリングを生成するビルダのトレイト
pub trait QueryBuilder {
    type Key: Serialize;
    type Value: Serialize;

    fn build(self) -> Vec<(Self::Key, Self::Value)>;
}

/// スタンダードなクエリビルダ
pub struct StandardQueryBuilder {
    q: Option<String>,
    start: Option<u32>,
    rows: Option<u32>,
    fq: Option<String>,
    fl: Option<String>,
    sort: Option<String>,
}

impl StandardQueryBuilder {
    pub fn new() -> Self {
        Self {
            q: None,
            start: None,
            rows: None,
            fq: None,
            fl: None,
            sort: None,
        }
    }

    pub fn q(mut self, q: &impl SolrQueryOperand) -> Self {
        self.q = Some(q.to_string());
        self
    }

    pub fn start(mut self, start: u32) -> Self {
        self.start = Some(start);
        self
    }

    pub fn rows(mut self, rows: u32) -> Self {
        self.rows = Some(rows);
        self
    }

    pub fn fq(mut self, fq: &impl SolrQueryOperand) -> Self {
        self.fq = Some(fq.to_string());
        self
    }

    pub fn fl(mut self, fl: String) -> Self {
        self.fl = Some(fl);
        self
    }

    pub fn sort(mut self, sort: &SortOrderBuilder) -> Self {
        self.sort = Some(sort.build());
        self
    }
}

impl QueryBuilder for StandardQueryBuilder {
    type Key = &'static str;
    type Value = String;
    fn build(self) -> Vec<(&'static str, String)> {
        let mut result: Vec<(&'static str, String)> = Vec::new();
        match self.q {
            Some(q) => result.push(("q", q)),
            None => result.push(("q", String::from("*:*"))),
        };

        match self.start {
            Some(start) => result.push(("start", start.to_string())),
            None => result.push(("start", String::from("0"))),
        };

        match self.rows {
            Some(rows) => result.push(("rows", rows.to_string())),
            None => result.push(("rows", 10.to_string())),
        };

        if let Some(fq) = self.fq {
            result.push(("fq", fq));
        }

        if let Some(fl) = self.fl {
            result.push(("fl", fl));
        }

        match self.sort {
            Some(sort) => result.push(("sort", sort)),
            None => result.push(("sort", String::from("score desc"))),
        };

        result
    }
}

pub trait SolrQueryOperand: Display {}

pub trait SingleSolrQueryOperand: SolrQueryOperand {}

pub enum QueryOperandKind<'a> {
    Standard(StandardQueryOperand),
    Range(RangeQueryOperand),
    Expression(QueryExpression<'a>),
}

pub enum TermModifiers {
    Normal,
    Fuzzy(u32),
    Proximity(u32),
    Boost(f64),
    Constant(f64),
    Phrase,
}

pub struct StandardQueryOperand {
    field: String,
    word: String,
    modifier: TermModifiers,
}

impl SolrQueryOperand for StandardQueryOperand {}
impl SingleSolrQueryOperand for StandardQueryOperand {}

impl StandardQueryOperand {
    pub fn new(field: &str, word: &str) -> Self {
        Self {
            field: field.to_string(),
            word: word.to_string(),
            modifier: TermModifiers::Normal,
        }
    }

    pub fn option(mut self, option: TermModifiers) -> Self {
        self.modifier = option;
        self
    }
}

impl Display for StandardQueryOperand {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        match self.modifier {
            TermModifiers::Normal => {
                write!(f, "{}:{}", self.field, self.word)?;
            }
            TermModifiers::Fuzzy(fuzzy) => {
                todo!();
            }
            TermModifiers::Proximity(proximity) => {
                todo!();
            }
            TermModifiers::Boost(boost) => {
                todo!();
            }
            TermModifiers::Constant(constant) => {
                todo!();
            }
            TermModifiers::Phrase => {
                todo!();
            }
        }

        Ok(())
    }
}

pub struct RangeQueryOperand {
    field: String,
    start: String,
    end: String,
    left_open: bool,
    right_open: bool,
}

impl SolrQueryOperand for RangeQueryOperand {}

impl RangeQueryOperand {
    pub fn new(field: &str) -> Self {
        Self {
            field: field.to_string(),
            start: "*".to_string(),
            end: "*".to_string(),
            left_open: false,
            right_open: true,
        }
    }
}

impl Display for RangeQueryOperand {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        todo!();
    }
}

pub struct QueryExpression<'a> {
    pub operator: String,
    pub operands: Vec<&'a QueryOperandKind<'a>>,
}

impl<'a> Display for QueryExpression<'a> {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        let s = self
            .operands
            .iter()
            .map(|op| match op {
                QueryOperandKind::Standard(op) => op.to_string(),
                QueryOperandKind::Range(op) => op.to_string(),
                QueryOperandKind::Expression(e) => e.to_string(),
            })
            .collect::<Vec<String>>()
            .join(&self.operator);
        write!(f, "{}", s)?;

        Ok(())
    }
}

pub struct SortOrderBuilder {
    order: Vec<String>,
}

impl SortOrderBuilder {
    pub fn build(&self) -> String {
        self.order.join(",")
    }

    pub fn asc(mut self, field: String) -> Self {
        self.order.push(format!("{} ASC", field));
        self
    }

    pub fn desc(mut self, field: String) -> Self {
        self.order.push(format!("{} DESC", field));
        self
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_with_no_params() {
        let builder = StandardQueryBuilder::new();

        assert_eq!(vec![("q", String::from("*:*"))], builder.build());
    }

    #[test]
    fn test_with_q() {
        let q = StandardQueryOperand::new("text_ja", "hoge");
        let builder = StandardQueryBuilder::new().q(&q);

        assert_eq!(vec![("q", String::from("text_ja:hoge"))], builder.build());
    }

    #[test]
    fn test_with_start() {
        let builder = StandardQueryBuilder::new().start(10);

        assert_eq!(
            vec![("q", String::from("*:*")), ("start", 10.to_string())],
            builder.build()
        );
    }

    #[test]
    fn test_with_rows() {
        let builder = StandardQueryBuilder::new().rows(20);

        assert_eq!(
            vec![("q", String::from("*:*")), ("rows", 20.to_string())],
            builder.build()
        );
    }

    #[test]
    fn test_with_fq() {
        let builder = StandardQueryBuilder::new().fq(&StandardQueryOperand::new("name", "alice"));

        assert_eq!(
            vec![
                ("q", String::from("*:*")),
                ("fq", String::from("name:alice"))
            ],
            builder.build()
        );
    }

    #[test]
    fn test_with_multiple_fq() {
        let builder = StandardQueryBuilder::new()
            .fq(&StandardQueryOperand::new("name", "alice"))
            .fq(&StandardQueryOperand::new("age", "24"));

        assert_eq!(
            vec![
                ("q", String::from("*:*")),
                ("fq", String::from("name:alice")),
                ("fq", String::from("age:24"))
            ],
            builder.build()
        );
    }

    #[test]
    fn test_with_fl() {
        let builder = StandardQueryBuilder::new().fl(String::from("id,name"));

        assert_eq!(
            vec![("q", String::from("*:*")), ("fl", String::from("id,name")),],
            builder.build()
        );
    }

    #[test]
    fn test_with_sort() {}
}
