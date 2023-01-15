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
    op: Option<String>,
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
            op: None,
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

        if let Some(start) = self.start {
            result.push(("start", start.to_string()));
        }

        if let Some(rows) = self.rows {
            result.push(("rows", rows.to_string()));
        }

        if let Some(fq) = self.fq {
            result.push(("fq", fq));
        }

        if let Some(fl) = self.fl {
            result.push(("fl", fl));
        }

        if let Some(sort) = self.sort {
            result.push(("sort", sort));
        }

        result
    }
}

pub trait SolrQueryOperand: Display {}

pub trait SingleSolrQueryOperand: SolrQueryOperand {}

pub enum QueryOperandKind {
    Standard(StandardQueryOperand),
    Range(RangeQueryOperand),
    Expression(QueryExpression),
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

    pub fn start(mut self, start: &str) -> Self {
        self.start = start.to_string();
        self
    }

    pub fn end(mut self, end: &str) -> Self {
        self.end = end.to_string();
        self
    }

    pub fn left_open(mut self) -> Self {
        self.left_open = true;
        self
    }

    pub fn left_close(mut self) -> Self {
        self.left_open = false;
        self
    }

    pub fn right_open(mut self) -> Self {
        self.right_open = true;
        self
    }

    pub fn right_close(mut self) -> Self {
        self.right_open = false;
        self
    }
}

impl Display for RangeQueryOperand {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        todo!();
    }
}

/// StandardQueryOperandを左辺に取るクエリ演算子同士の加算演算子のオーバーロードの定義
impl<'a, T> ops::Add<T> for StandardQueryOperand
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn add(self, rhs: T) -> QueryExpression {
        let op = &rhs;
        todo!();
    }
}

/// StandardQueryOperandを左辺に取るクエリ演算子同士の乗算演算子のオーバーロードの定義
impl<T> ops::Mul<T> for StandardQueryOperand
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn mul(self, rhs: T) -> QueryExpression {
        todo!();
    }
}

/// StandardQueryOperandを左辺に取るクエリ式との加算演算子のオーバーロードの定義
impl<'a> ops::Add<QueryExpression> for StandardQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

/// StandardQueryOperandを左辺に取るクエリ式との乗算演算子のオーバーロードの定義
impl<'a> ops::Mul<QueryExpression> for StandardQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

/// RangeQueryOperandを左辺に取るクエリ演算子同士の加算演算子のオーバーロードの定義
impl<T> ops::Add<T> for RangeQueryOperand
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn add(self, rhs: T) -> QueryExpression {
        todo!();
    }
}

/// RangeQueryOperandを左辺に取るクエリ演算子同士の乗算演算子のオーバーロードの定義
impl<T> ops::Mul<T> for RangeQueryOperand
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn mul(self, rhs: T) -> QueryExpression {
        todo!();
    }
}

/// RangeQueryOperandを左辺に取るクエリ式との加算演算子のオーバーロードの定義
impl<'a> ops::Add<QueryExpression> for RangeQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

/// RangeQueryOperandを左辺に取るクエリ式との乗算演算子のオーバーロードの定義
impl<'a> ops::Mul<QueryExpression> for RangeQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

pub struct QueryExpression {
    pub operator: String,
    pub operands: Vec<QueryOperandKind>,
}

impl<'a> Display for QueryExpression {
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

impl<'a> ops::Add<QueryExpression> for QueryExpression {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

impl<'a, T> ops::Add<T> for QueryExpression
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn add(self, rhs: T) -> QueryExpression {
        todo!();
    }
}

impl<'a> ops::Mul<QueryExpression> for QueryExpression {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        todo!();
    }
}

impl<'a, T> ops::Mul<T> for QueryExpression
where
    T: SolrQueryOperand,
{
    type Output = QueryExpression;

    fn mul(self, rhs: T) -> QueryExpression {
        todo!();
    }
}

pub struct SortOrderBuilder {
    order: Vec<String>,
}

impl SortOrderBuilder {
    pub fn new() -> Self {
        Self { order: Vec::new() }
    }

    pub fn build(&self) -> String {
        self.order.join(",")
    }

    pub fn asc(mut self, field: &str) -> Self {
        self.order.push(format!("{} ASC", field));
        self
    }

    pub fn desc(mut self, field: &str) -> Self {
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
    fn test_build_sort_order() {
        let sort = SortOrderBuilder::new().desc("score").asc("name").build();

        assert_eq!(String::from("score desc,name asc"), sort);
    }

    #[test]
    fn test_query_operand_representation() {
        let q = StandardQueryOperand::new("name", "alice");
        assert_eq!(String::from("name:alice"), q.to_string());
    }

    #[test]
    fn test_fuzzy_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Fuzzy(1));
        assert_eq!(String::from("name:alice~1"), q.to_string());
    }

    #[test]
    fn test_proximity_query_operand() {
        let q =
            StandardQueryOperand::new("name", "alice wonder").option(TermModifiers::Proximity(2));
        assert_eq!(String::from(r#"name:"alice"~1"#), q.to_string());
    }

    #[test]
    fn test_boost_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Boost(10.0));
        assert_eq!(String::from("name:alice^10"), q.to_string());
    }

    #[test]
    fn test_constant_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Constant(0.0));
        assert_eq!(String::from("name:alice^=10"), q.to_string());
    }

    #[test]
    fn test_phrase_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Phrase);
        assert_eq!(String::from(r#"name:"alice""#), q.to_string());
    }

    #[test]
    fn test_default_range_query() {
        let q = RangeQueryOperand::new("age");

        assert_eq!(String::from("age:[* TO *}"), q.to_string())
    }

    #[test]
    fn test_left_close_right_open_range_query() {
        let q = RangeQueryOperand::new("age")
            .start("10")
            .end("20")
            .left_close()
            .right_open();
    }

    #[test]
    fn test_add_operands() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("age", "24");

        let q = op1 + op2;

        assert_eq!(String::from("name:alice OR age:24"), q.to_string())
    }

    #[test]
    fn test_mul_operands() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("age", "24");

        let q = op1 * op2;

        assert_eq!(String::from("name:alice AND age:24"), q.to_string())
    }

    #[test]
    fn test_add_operand_to_expression() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("age", "24");

        let q = (op1 * op2) + op3;

        assert_eq!(
            String::from("(name:alice AND name:bob) OR age:24"),
            q.to_string()
        )
    }

    #[test]
    fn test_add_expression_to_operand() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("age", "24");

        let q = op3 * (op1 + op2);

        assert_eq!(
            String::from("(name:alice OR name:bob) AND age:24"),
            q.to_string()
        )
    }

    #[test]
    fn test_add_expression_to_expression() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("age", "24");
        let op3 = StandardQueryOperand::new("name", "bob");
        let op4 = StandardQueryOperand::new("age", "32");

        let q = (op1 * op2) + (op3 * op4);

        assert_eq!(
            String::from("(name:alice AND age:24) OR (name:bob AND age:32)"),
            q.to_string()
        )
    }

    #[test]
    fn test_mul_expression_to_expression() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("age", "24");
        let op4 = StandardQueryOperand::new("age", "32");

        let q = (op1 + op2) * (op3 + op4);

        assert_eq!(
            String::from("(name:alice OR name:bob) OR (age:24 AND age:32)"),
            q.to_string()
        )
    }

    #[test]
    fn test_extend_expression_with_add() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("name", "charles");

        let q = op1 + op2 + op3;

        assert_eq!(
            String::from("name:alice OR name:bob OR name:charles"),
            q.to_string()
        )
    }

    #[test]
    fn test_extend_expression_with_mul() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("name", "charles");

        let q = op1 * op2 * op3;

        assert_eq!(
            String::from("name:alice AND name:bob AND name:charles"),
            q.to_string()
        )
    }
}
