use crate::querybuilders::common::SolrCommonQueryBuilder;
use crate::querybuilders::facet::FacetBuilder;
use crate::querybuilders::q::{Operator, SolrQueryExpression};
use crate::querybuilders::sort::SortOrderBuilder;
use std::collections::HashMap;

pub trait SolrDisMaxQueryBuilder: SolrCommonQueryBuilder {
    fn q(self, q: String) -> Self;
    fn qf(self, qf: &str) -> Self;
    fn qs(self, qs: u32) -> Self;
    fn pf(self, pf: &str) -> Self;
    fn ps(self, ps: u32) -> Self;
    fn mm(self, mm: &str) -> Self;
    fn q_alt(self, q: &impl SolrQueryExpression) -> Self;
    fn tie(self, tie: f64) -> Self;
    fn bq(self, bq: &impl SolrQueryExpression) -> Self;
    fn bf(self, bf: &str) -> Self;
}

pub struct DisMaxQueryBuilder {
    params: HashMap<String, String>,
    multi_params: HashMap<String, Vec<String>>,
}

impl DisMaxQueryBuilder {
    pub fn new() -> Self {
        Self {
            params: HashMap::new(),
            multi_params: HashMap::new(),
        }
    }
}

impl SolrCommonQueryBuilder for DisMaxQueryBuilder {
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
        params.push(("defType".to_string(), "dismax".to_string()));

        params.extend(self.params.into_iter());
        for (key, values) in self.multi_params.into_iter() {
            params.extend(values.into_iter().map(|param| (key.clone(), param)));
        }

        params
    }
}

impl SolrDisMaxQueryBuilder for DisMaxQueryBuilder {
    fn q(mut self, q: String) -> Self {
        // TODO: 引数の型を抽象モデルに変更する
        self.params.insert("q".to_string(), q.to_string());
        self
    }
    fn qf(mut self, qf: &str) -> Self {
        self.params.insert("qf".to_string(), qf.to_string());
        self
    }
    fn qs(mut self, qs: u32) -> Self {
        self.params.insert("qs".to_string(), qs.to_string());
        self
    }
    fn pf(mut self, pf: &str) -> Self {
        self.params.insert("pf".to_string(), pf.to_string());
        self
    }
    fn ps(mut self, ps: u32) -> Self {
        self.params.insert("ps".to_string(), ps.to_string());
        self
    }
    fn mm(mut self, mm: &str) -> Self {
        self.params.insert("mm".to_string(), mm.to_string());
        self
    }
    fn q_alt(mut self, q: &impl SolrQueryExpression) -> Self {
        self.params.insert("q.alt".to_string(), q.to_string());
        self
    }
    fn tie(mut self, tie: f64) -> Self {
        self.params.insert("tie".to_string(), tie.to_string());
        self
    }
    fn bq(mut self, bq: &impl SolrQueryExpression) -> Self {
        self.multi_params
            .entry("bq".to_string())
            .or_default()
            .push(bq.to_string());
        self
    }
    fn bf(mut self, bf: &str) -> Self {
        self.multi_params
            .entry("bf".to_string())
            .or_default()
            .push(bf.to_string());
        self
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use crate::querybuilders::q::QueryOperand;

    #[test]
    fn test_q() {
        let q = QueryOperand::from("プログラミング Rust");
        let builder = DisMaxQueryBuilder::new().q(q.to_string());

        let mut expected = vec![
            ("defType".to_string(), "dismax".to_string()),
            ("q".to_string(), "プログラミング Rust".to_string()),
        ];
        let mut actual = builder.build();
        expected.sort();
        actual.sort();
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_qf() {
        let q = QueryOperand::from("プログラミング Rust");
        let builder = DisMaxQueryBuilder::new().q(q.to_string()).qf("title text");

        let mut expected = vec![
            ("defType".to_string(), "dismax".to_string()),
            ("q".to_string(), "プログラミング Rust".to_string()),
            ("qf".to_string(), "title text".to_string()),
        ];
        let mut actual = builder.build();
        expected.sort();
        actual.sort();
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_sample_query() {
        let q = QueryOperand::from("*:*");
        let sort = SortOrderBuilder::new().desc("score").asc("start_at");
        let builder = DisMaxQueryBuilder::new()
            .q("すぬけ 耳".to_string())
            .qf("text_ja")
            .op(Operator::AND)
            .wt("json")
            .debug()
            .q_alt(&q)
            .sort(&sort)
            .fl("problem_title".to_string());

        let mut expected = vec![
            ("defType".to_string(), "dismax".to_string()),
            ("q".to_string(), "すぬけ 耳".to_string()),
            ("qf".to_string(), "text_ja".to_string()),
            ("q.op".to_string(), "AND".to_string()),
            ("wt".to_string(), "json".to_string()),
            ("debug".to_string(), "all".to_string()),
            ("debug.explain.structured".to_string(), "true".to_string()),
            ("q.alt".to_string(), "*:*".to_string()),
            ("sort".to_string(), "score desc,start_at asc".to_string()),
            ("fl".to_string(), "problem_title".to_string()),
        ];
        let mut actual = builder.build();
        expected.sort();
        actual.sort();
        assert_eq!(actual, expected);
    }
}
