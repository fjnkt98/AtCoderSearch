use crate::querybuilders::common::SolrCommonQueryBuilder;
use crate::querybuilders::facet::FacetBuilder;
use crate::querybuilders::q::{Operator, SolrQueryExpression};
use crate::querybuilders::sort::SortOrderBuilder;
use std::collections::HashMap;

pub trait SolrStandardQueryBuilder: SolrCommonQueryBuilder {
    fn q(self, q: &impl SolrQueryExpression) -> Self;
    fn df(self, df: &str) -> Self;
}

pub struct StandardQueryBuilder {
    params: HashMap<String, String>,
    multi_params: HashMap<String, Vec<String>>,
}

impl StandardQueryBuilder {
    pub fn new() -> Self {
        Self {
            params: HashMap::new(),
            multi_params: HashMap::new(),
        }
    }
}

impl SolrStandardQueryBuilder for StandardQueryBuilder {
    fn q(mut self, q: &impl SolrQueryExpression) -> Self {
        self.params.insert("q".to_string(), q.to_string());
        self
    }
    fn df(mut self, df: &str) -> Self {
        self.params.insert("df".to_string(), df.to_string());
        self
    }
}

impl SolrCommonQueryBuilder for StandardQueryBuilder {
    fn sort(mut self, sort: &SortOrderBuilder) -> Self {
        self.params.insert("sort".to_string(), sort.build());
        self
    }

    fn start(mut self, start: u32) -> Self {
        self.params.insert("start".to_string(), start.to_string());
        self
    }

    fn rows(mut self, rows: u32) -> Self {
        self.params.insert("rows".to_string(), rows.to_string());
        self
    }

    fn fq(mut self, fq: &impl SolrQueryExpression) -> Self {
        self.multi_params
            .entry("fq".to_string())
            .or_default()
            .push(fq.to_string());
        self
    }

    fn fl(mut self, fl: String) -> Self {
        self.params.insert("fl".to_string(), fl);
        self
    }

    fn debug(mut self) -> Self {
        self.params.insert("debug".to_string(), "all".to_string());
        self.params
            .insert("debug.explain.structured".to_string(), "true".to_string());
        self
    }
    fn wt(mut self, wt: &str) -> Self {
        self.params.insert("wt".to_string(), wt.to_string());
        self
    }
    fn facet(mut self, facet: &impl FacetBuilder) -> Self {
        self.params.insert("facet".to_string(), "true".to_string());
        for (key, value) in facet.build() {
            self.params.insert(key, value);
        }
        self
    }

    fn op(mut self, op: Operator) -> Self {
        match op {
            Operator::AND => {
                self.params.insert("q.op".to_string(), "AND".to_string());
            }
            Operator::OR => {
                self.params.insert("q.op".to_string(), "OR".to_string());
            }
        }
        self
    }

    fn build(self) -> Vec<(String, String)> {
        let mut params = Vec::new();

        params.extend(self.params.into_iter());
        for (key, values) in self.multi_params.into_iter() {
            params.extend(values.into_iter().map(|param| (key.clone(), param)));
        }

        params
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::querybuilders::facet::{FieldFacetBuilder, RangeFacetBuilder};
    use crate::querybuilders::q::{QueryOperand, StandardQueryOperand};

    #[test]
    fn test_with_q() {
        let q = QueryOperand::from("text_ja:hoge");
        let builder = StandardQueryBuilder::new().q(&q);

        assert_eq!(
            vec![("q".to_string(), "text_ja:hoge".to_string())],
            builder.build()
        );
    }

    #[test]
    fn test_sample_query() {
        let q = QueryOperand::from(StandardQueryOperand::new("text_ja", "高橋?"));
        let sort = SortOrderBuilder::new().desc("score").desc("difficulty");
        let facet1 = FieldFacetBuilder::new("category");
        let facet2 = RangeFacetBuilder::new(
            "difficulty",
            0.to_string(),
            2000.to_string(),
            400.to_string(),
        );
        let builder = StandardQueryBuilder::new()
            .q(&q)
            .sort(&sort)
            .facet(&facet1)
            .facet(&facet2);

        let mut expected = vec![
            ("q".to_string(), r#"text_ja:高橋\?"#.to_string()),
            ("sort".to_string(), "score desc,difficulty desc".to_string()),
            ("facet".to_string(), "true".to_string()),
            ("facet.field".to_string(), "category".to_string()),
            ("facet.range".to_string(), "difficulty".to_string()),
            (
                "f.difficulty.facet.range.start".to_string(),
                "0".to_string(),
            ),
            (
                "f.difficulty.facet.range.end".to_string(),
                "2000".to_string(),
            ),
            (
                "f.difficulty.facet.range.gap".to_string(),
                "400".to_string(),
            ),
        ];
        expected.sort();
        let mut actual = builder.build();
        actual.sort();

        assert_eq!(actual, expected);
    }
}
