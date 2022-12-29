use crate::utils::models::*;
use ego_tree::NodeRef;
use regex::Regex;
use scraper::node::Node;
use scraper::{Html, Selector};

type Result<T> = std::result::Result<T, IndexingError>;

pub struct FullTextExtractor {
    div: Selector,
    section: Selector,
    h3: Selector,
    ascii: Regex,
}

impl FullTextExtractor {
    pub fn new() -> Result<Self> {
        let div = Selector::parse("div.part").map_err(|e| IndexingError::SelectorError(e))?;
        let section = Selector::parse("section").map_err(|e| IndexingError::SelectorError(e))?;
        let h3 = Selector::parse("h3").map_err(|e| IndexingError::SelectorError(e))?;

        let ascii = Regex::new("^[\x20-\x7E].*$").map_err(|e| IndexingError::RegexError(e))?;

        Ok(FullTextExtractor {
            div: div,
            section: section,
            h3: h3,
            ascii: ascii,
        })
    }

    fn dfs(&self, element: &NodeRef<Node>) -> String {
        let mut result = Vec::new();

        for child in element.children() {
            match child.value() {
                Node::Element(e) => {
                    if e.name() == "pre" {
                        continue;
                    }
                    result.push(self.dfs(&child));
                }
                Node::Text(text) => {
                    result.push(text.trim().to_string());
                }
                _ => {
                    continue;
                }
            };
        }

        result.join(" ")
    }

    pub fn extract(&self, html: &str) -> Result<(Vec<String>, Vec<String>)> {
        let html = Html::parse_document(html);

        let mut text_ja: Vec<String> = Vec::new();
        let mut text_en: Vec<String> = Vec::new();
        for part in html.select(&self.div) {
            let Some(section) = part.select(&self.section).next() else {continue};

            let Some(title) = section.select(&self.h3).next() else {continue};
            let Some(title) = title.text().next() else {continue};

            let mut full_text: Vec<String> = Vec::new();
            for e in section.children() {
                match e.value() {
                    Node::Element(element) => {
                        if element.name() == "h3" {
                            continue;
                        } else {
                            full_text.push(self.dfs(&e));
                        }
                    }
                    Node::Text(text) => full_text.push(text.to_string()),
                    _ => {
                        continue;
                    }
                }
            }

            if self.ascii.is_match(&title) {
                text_en.push(full_text.join(""));
            } else {
                text_ja.push(full_text.join(" "));
            }
        }

        Ok((text_ja, text_en))
    }
}
