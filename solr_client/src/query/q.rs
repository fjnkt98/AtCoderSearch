use std::fmt::{Display, Formatter};
use std::ops;

pub trait SolrQueryOperand: Display {}

#[derive(Clone)]
pub enum QueryOperandKind {
    Standard(StandardQueryOperand),
    Range(RangeQueryOperand),
    Expression(QueryExpression),
}

#[derive(Clone)]
pub enum TermModifiers {
    Normal,
    Fuzzy(u32),
    Proximity(u32),
    Boost(f64),
    Constant(f64),
    Phrase,
}

#[derive(Clone)]
pub struct StandardQueryOperand {
    field: String,
    word: String,
    modifier: TermModifiers,
}

impl SolrQueryOperand for StandardQueryOperand {}

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
                write!(f, "{}:{}~{}", self.field, self.word, fuzzy)?;
            }
            TermModifiers::Proximity(proximity) => {
                write!(f, r#"{}:"{}"~{}"#, self.field, self.word, proximity)?;
            }
            TermModifiers::Boost(boost) => {
                write!(f, "{}:{}^{}", self.field, self.word, boost)?;
            }
            TermModifiers::Constant(constant) => {
                write!(f, "{}:{}^={}", self.field, self.word, constant)?;
            }
            TermModifiers::Phrase => {
                write!(f, r#"{}:"{}""#, self.field, self.word)?;
            }
        }

        Ok(())
    }
}

#[derive(Clone)]
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
        let left_parenthesis = match self.left_open {
            false => '[',
            true => '{',
        };

        let right_parenthesis = match self.right_open {
            false => ']',
            true => '}',
        };

        write!(
            f,
            "{}:{}{} TO {}{}",
            self.field, left_parenthesis, self.start, self.end, right_parenthesis
        )?;

        Ok(())
    }
}

impl ops::Add<StandardQueryOperand> for StandardQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: StandardQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::OR,
            vec![
                QueryOperandKind::Standard(self),
                QueryOperandKind::Standard(rhs),
            ],
        )
    }
}

impl ops::Mul<StandardQueryOperand> for StandardQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: StandardQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::AND,
            vec![
                QueryOperandKind::Standard(self),
                QueryOperandKind::Standard(rhs),
            ],
        )
    }
}

impl ops::Add<RangeQueryOperand> for StandardQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: RangeQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::OR,
            vec![
                QueryOperandKind::Standard(self),
                QueryOperandKind::Range(rhs),
            ],
        )
    }
}

impl ops::Mul<RangeQueryOperand> for StandardQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: RangeQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::AND,
            vec![
                QueryOperandKind::Standard(self),
                QueryOperandKind::Range(rhs),
            ],
        )
    }
}

impl ops::Add<QueryExpression> for StandardQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        let operands = match rhs.operator {
            Operator::AND => {
                let operands = vec![
                    QueryOperandKind::Standard(self),
                    QueryOperandKind::Expression(rhs),
                ];
                operands
            }
            Operator::OR => {
                let mut operands = vec![QueryOperandKind::Standard(self)];
                operands.extend(rhs.operands);
                operands
            }
        };

        QueryExpression {
            operator: Operator::OR,
            operands: operands,
        }
    }
}

impl ops::Mul<QueryExpression> for StandardQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        let operands = match rhs.operator {
            Operator::OR => {
                let operands = vec![
                    QueryOperandKind::Standard(self),
                    QueryOperandKind::Expression(rhs),
                ];
                operands
            }
            Operator::AND => {
                let mut operands = vec![QueryOperandKind::Standard(self)];
                operands.extend(rhs.operands);
                operands
            }
        };

        QueryExpression {
            operator: Operator::AND,
            operands: operands,
        }
    }
}

impl ops::Add<RangeQueryOperand> for RangeQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: RangeQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::OR,
            vec![QueryOperandKind::Range(self), QueryOperandKind::Range(rhs)],
        )
    }
}

/// RangeQueryOperandを左辺に取るクエリ演算子同士の乗算演算子のオーバーロードの定義
impl ops::Mul<RangeQueryOperand> for RangeQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: RangeQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::AND,
            vec![QueryOperandKind::Range(self), QueryOperandKind::Range(rhs)],
        )
    }
}

impl ops::Add<StandardQueryOperand> for RangeQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: StandardQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::OR,
            vec![
                QueryOperandKind::Range(self),
                QueryOperandKind::Standard(rhs),
            ],
        )
    }
}

impl ops::Mul<StandardQueryOperand> for RangeQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: StandardQueryOperand) -> QueryExpression {
        QueryExpression::new(
            Operator::AND,
            vec![
                QueryOperandKind::Range(self),
                QueryOperandKind::Standard(rhs),
            ],
        )
    }
}

impl ops::Add<QueryExpression> for RangeQueryOperand {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        let operands = match rhs.operator {
            Operator::OR => {
                let operands = vec![
                    QueryOperandKind::Range(self),
                    QueryOperandKind::Expression(rhs),
                ];
                operands
            }
            Operator::AND => {
                let mut operands = vec![QueryOperandKind::Range(self)];
                operands.extend(rhs.operands);
                operands
            }
        };

        QueryExpression {
            operator: Operator::OR,
            operands: operands,
        }
    }
}

impl ops::Mul<QueryExpression> for RangeQueryOperand {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        let operands = match rhs.operator {
            Operator::OR => {
                let operands = vec![
                    QueryOperandKind::Range(self),
                    QueryOperandKind::Expression(rhs),
                ];
                operands
            }
            Operator::AND => {
                let mut operands = vec![QueryOperandKind::Range(self)];
                operands.extend(rhs.operands);
                operands
            }
        };

        QueryExpression {
            operator: Operator::AND,
            operands: operands,
        }
    }
}

#[derive(Clone, PartialEq, Eq)]
pub enum Operator {
    AND,
    OR,
}

#[derive(Clone)]
pub struct QueryExpression {
    pub operator: Operator,
    pub operands: Vec<QueryOperandKind>,
}

impl QueryExpression {
    pub fn new(operator: Operator, operands: Vec<QueryOperandKind>) -> Self {
        Self { operator, operands }
    }
}

impl Display for QueryExpression {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        let operator = match self.operator {
            Operator::AND => " AND ",
            Operator::OR => " OR ",
        };

        let s = self
            .operands
            .iter()
            .map(|op| match op {
                QueryOperandKind::Standard(op) => op.to_string(),
                QueryOperandKind::Range(op) => op.to_string(),
                QueryOperandKind::Expression(e) => e.to_string(),
            })
            .collect::<Vec<String>>()
            .join(operator);
        write!(f, "({})", s)?;

        Ok(())
    }
}

impl ops::Add<QueryExpression> for QueryExpression {
    type Output = QueryExpression;

    fn add(self, rhs: QueryExpression) -> QueryExpression {
        let operands = if self.operator == Operator::OR && rhs.operator == Operator::OR {
            let mut operands = Vec::new();
            operands.extend(self.operands);
            operands.extend(rhs.operands);
            operands
        } else {
            vec![
                QueryOperandKind::Expression(self),
                QueryOperandKind::Expression(rhs),
            ]
        };

        QueryExpression {
            operator: Operator::OR,
            operands: operands,
        }
    }
}

impl ops::Mul<QueryExpression> for QueryExpression {
    type Output = QueryExpression;

    fn mul(self, rhs: QueryExpression) -> QueryExpression {
        let operands = if self.operator == Operator::AND && rhs.operator == Operator::AND {
            let mut operands = Vec::new();
            operands.extend(self.operands);
            operands.extend(rhs.operands);
            operands
        } else {
            vec![
                QueryOperandKind::Expression(self),
                QueryOperandKind::Expression(rhs),
            ]
        };

        QueryExpression {
            operator: Operator::AND,
            operands: operands,
        }
    }
}

impl ops::Add<StandardQueryOperand> for QueryExpression {
    type Output = QueryExpression;

    fn add(self, rhs: StandardQueryOperand) -> QueryExpression {
        match self.operator {
            Operator::OR => {
                let mut operands = self.operands.clone();
                operands.push(QueryOperandKind::Standard(rhs));
                QueryExpression {
                    operator: Operator::OR,
                    operands: operands,
                }
            }
            Operator::AND => QueryExpression {
                operator: Operator::OR,
                operands: vec![
                    QueryOperandKind::Expression(self),
                    QueryOperandKind::Standard(rhs),
                ],
            },
        }
    }
}

impl ops::Mul<StandardQueryOperand> for QueryExpression {
    type Output = QueryExpression;

    fn mul(self, rhs: StandardQueryOperand) -> QueryExpression {
        match self.operator {
            Operator::AND => {
                let mut operands = self.operands.clone();
                operands.push(QueryOperandKind::Standard(rhs));
                QueryExpression {
                    operator: Operator::AND,
                    operands: operands,
                }
            }
            Operator::OR => QueryExpression {
                operator: Operator::AND,
                operands: vec![
                    QueryOperandKind::Expression(self),
                    QueryOperandKind::Standard(rhs),
                ],
            },
        }
    }
}

impl ops::Add<RangeQueryOperand> for QueryExpression {
    type Output = QueryExpression;

    fn add(self, rhs: RangeQueryOperand) -> QueryExpression {
        match self.operator {
            Operator::OR => {
                let mut operands = self.operands.clone();
                operands.push(QueryOperandKind::Range(rhs));
                QueryExpression {
                    operator: Operator::OR,
                    operands: operands,
                }
            }
            Operator::AND => QueryExpression {
                operator: Operator::OR,
                operands: vec![
                    QueryOperandKind::Expression(self),
                    QueryOperandKind::Range(rhs),
                ],
            },
        }
    }
}

impl ops::Mul<RangeQueryOperand> for QueryExpression {
    type Output = QueryExpression;

    fn mul(self, rhs: RangeQueryOperand) -> QueryExpression {
        match self.operator {
            Operator::AND => {
                let mut operands = self.operands.clone();
                operands.push(QueryOperandKind::Range(rhs));
                QueryExpression {
                    operator: Operator::AND,
                    operands: operands,
                }
            }
            Operator::OR => QueryExpression {
                operator: Operator::AND,
                operands: vec![
                    QueryOperandKind::Expression(self),
                    QueryOperandKind::Range(rhs),
                ],
            },
        }
    }
}

#[cfg(test)]
mod test {
    use super::*;

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
        assert_eq!(String::from(r#"name:"alice wonder"~2"#), q.to_string());
    }

    #[test]
    fn test_boost_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Boost(10.0));
        assert_eq!(String::from("name:alice^10"), q.to_string());
    }

    #[test]
    fn test_constant_query_operand() {
        let q = StandardQueryOperand::new("name", "alice").option(TermModifiers::Constant(0.0));
        assert_eq!(String::from("name:alice^=0"), q.to_string());
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
    fn test_left_close_right_close_range_query() {
        let q = RangeQueryOperand::new("age")
            .start("10")
            .end("20")
            .left_close()
            .right_close();

        assert_eq!(String::from("age:[10 TO 20]"), q.to_string())
    }

    #[test]
    fn test_add_operands() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("age", "24");

        let q = op1 + op2;

        assert_eq!(String::from("(name:alice OR age:24)"), q.to_string())
    }

    #[test]
    fn test_mul_operands() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("age", "24");

        let q = op1 * op2;

        assert_eq!(String::from("(name:alice AND age:24)"), q.to_string())
    }

    #[test]
    fn test_add_operand_to_expression() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("age", "24");

        let q = (op1 * op2) + op3;

        assert_eq!(
            String::from("((name:alice AND name:bob) OR age:24)"),
            q.to_string()
        )
    }

    #[test]
    fn test_add_expression_to_operand() {
        let op1 = StandardQueryOperand::new("name", "alice");
        let op2 = StandardQueryOperand::new("name", "bob");
        let op3 = StandardQueryOperand::new("age", "24");

        let q = op1 * (op2 + op3);

        assert_eq!(
            String::from("(name:alice AND (name:bob OR age:24))"),
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
            String::from("((name:alice AND age:24) OR (name:bob AND age:32))"),
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
            String::from("((name:alice OR name:bob) AND (age:24 OR age:32))"),
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
            String::from("(name:alice OR name:bob OR name:charles)"),
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
            String::from("(name:alice AND name:bob AND name:charles)"),
            q.to_string()
        )
    }
}
