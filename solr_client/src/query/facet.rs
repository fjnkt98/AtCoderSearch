use serde::Serialize;
use std::fmt::{Display, Formatter};
use std::ops;

pub trait FacetBuilder {
    type Key: Serialize;
    type Value: Serialize;

    fn build(self) -> Vec<(Self::Key, Self::Value)>;
}

pub enum FieldFacetSortOrder {
    Index,
    Count,
}

pub enum FieldFacetMethod {
    Enum,
    Fc,
    Fcs,
}

pub struct FieldFacetBuilder {
    field: String,
    prefix: Option<String>,
    contains: Option<String>,
    ignore_case: Option<bool>,
    sort: Option<String>,
    limit: Option<u32>,
    offset: Option<u32>,
    min_count: Option<u32>,
    missing: Option<bool>,
    method: Option<String>,
    exists: Option<bool>,
}

impl FieldFacetBuilder {
    pub fn new(field: &str) -> Self {
        Self {
            field: field.to_string(),
            prefix: None,
            contains: None,
            ignore_case: None,
            sort: None,
            limit: None,
            offset: None,
            min_count: None,
            missing: None,
            method: None,
            exists: None,
        }
    }

    pub fn prefix(mut self, prefix: &str) -> Self {
        self.prefix = Some(prefix.to_string());
        self
    }

    pub fn contains(mut self, contains: &str) -> Self {
        self.contains = Some(contains.to_string());
        self
    }

    pub fn ignore_case(mut self, ignore_case: bool) -> Self {
        self.ignore_case = Some(ignore_case);
        self
    }

    pub fn sort(mut self, sort: FieldFacetSortOrder) -> Self {
        self.sort = Some(match sort {
            FieldFacetSortOrder::Count => "count".to_string(),
            FieldFacetSortOrder::Index => "index".to_string(),
        });
        self
    }

    pub fn limit(mut self, limit: u32) -> Self {
        self.limit = Some(limit);
        self
    }

    pub fn offset(mut self, offset: u32) -> Self {
        self.offset = Some(offset);
        self
    }

    pub fn min_count(mut self, min_count: u32) -> Self {
        self.min_count = Some(min_count);
        self
    }

    pub fn missing(mut self, missing: bool) -> Self {
        self.missing = Some(missing);
        self
    }

    pub fn method(mut self, method: FieldFacetMethod) -> Self {
        self.method = Some(match method {
            FieldFacetMethod::Enum => "enum".to_string(),
            FieldFacetMethod::Fc => "fc".to_string(),
            FieldFacetMethod::Fcs => "fcs".to_string(),
        });
        self
    }

    pub fn exists(mut self, exists: bool) -> Self {
        self.exists = Some(exists);
        self
    }
}

impl FacetBuilder for FieldFacetBuilder {
    type Key = &'static str;
    type Value = String;
    fn build(self) -> Vec<(&'static str, String)> {
        let mut result: Vec<(&'static str, String)> = Vec::new();

        result.push(("facet.field", self.field));
        // if let Some(prefix) = self.prefix {
        //     result.push((format!("f.{}.facet.prefix", self.field).as_str(), prefix));
        // }
        todo!();
    }
}
