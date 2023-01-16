use crate::query::*;
use serde::Serialize;

//// Solrへのクエリに付加するクエリストリングを生成するビルダのトレイト
pub trait QueryBuilder {
    type Key: Serialize;
    type Value: Serialize;

    fn build(self) -> Vec<(Self::Key, Self::Value)>;
}

pub struct StandardQueryBuilder {
    q: Option<String>,
    start: Option<u32>,
    rows: Option<u32>,
    fq: Vec<String>,
    fl: Option<String>,
    sort: Option<String>,
    op: Option<String>,
    facet: Option<String>,
}

impl StandardQueryBuilder {
    pub fn new() -> Self {
        Self {
            q: None,
            start: None,
            rows: None,
            fq: Vec::new(),
            fl: None,
            sort: None,
            op: None,
            facet: None,
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
        self.fq.push(fq.to_string());
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

    pub fn op(mut self, op: &str) -> Self {
        self.op = Some(op.to_string());
        self
    }

    pub fn facet(mut self, facet: &FieldFacetBuilder) -> Self {
        todo!();
    }
}

impl QueryBuilder for StandardQueryBuilder {
    type Key = String;
    type Value = String;
    fn build(self) -> Vec<(String, String)> {
        let mut result: Vec<(String, String)> = Vec::new();
        match self.q {
            Some(q) => result.push((String::from("q"), q)),
            None => result.push((String::from("q"), String::from("*:*"))),
        };

        if let Some(start) = self.start {
            result.push((String::from("start"), start.to_string()));
        }

        if let Some(rows) = self.rows {
            result.push((String::from("rows"), rows.to_string()));
        }

        if !self.fq.is_empty() {
            for fq in self.fq {
                result.push((String::from("fq"), fq));
            }
        }

        if let Some(fl) = self.fl {
            result.push((String::from("fl"), fl));
        }

        if let Some(sort) = self.sort {
            result.push((String::from("sort"), sort));
        }

        if let Some(op) = self.op {
            result.push((String::from("q.op"), op));
        }

        result
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_with_no_params() {
        let builder = StandardQueryBuilder::new();

        assert_eq!(
            vec![(String::from("q"), String::from("*:*"))],
            builder.build()
        );
    }

    #[test]
    fn test_with_q() {
        let q = StandardQueryOperand::new("text_ja", "hoge");
        let builder = StandardQueryBuilder::new().q(&q);

        assert_eq!(
            vec![(String::from("q"), String::from("text_ja:hoge"))],
            builder.build()
        );
    }

    #[test]
    fn test_with_start() {
        let builder = StandardQueryBuilder::new().start(10);

        assert_eq!(
            vec![
                (String::from("q"), String::from("*:*")),
                (String::from("start"), 10.to_string())
            ],
            builder.build()
        );
    }

    #[test]
    fn test_with_rows() {
        let builder = StandardQueryBuilder::new().rows(20);

        assert_eq!(
            vec![
                (String::from("q"), String::from("*:*")),
                (String::from("rows"), 20.to_string())
            ],
            builder.build()
        );
    }

    #[test]
    fn test_with_fq() {
        let builder = StandardQueryBuilder::new().fq(&StandardQueryOperand::new("name", "alice"));

        assert_eq!(
            vec![
                (String::from("q"), String::from("*:*")),
                (String::from("fq"), String::from("name:alice"))
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
                (String::from("q"), String::from("*:*")),
                (String::from("fq"), String::from("name:alice")),
                (String::from("fq"), String::from("age:24"))
            ],
            builder.build()
        );
    }

    #[test]
    fn test_with_fl() {
        let builder = StandardQueryBuilder::new().fl(String::from("id,name"));

        assert_eq!(
            vec![
                (String::from("q"), String::from("*:*")),
                (String::from("fl"), String::from("id,name")),
            ],
            builder.build()
        );
    }

    #[test]
    fn test_q_op() {
        let builder = StandardQueryBuilder::new().op("AND");

        assert_eq!(
            vec![
                (String::from("q"), String::from("*:*")),
                (String::from("q.op"), String::from("AND")),
            ],
            builder.build()
        )
    }
}
