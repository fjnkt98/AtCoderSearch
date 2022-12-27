use anyhow::Result;
use ego_tree::NodeRef;
use regex::Regex;
use scraper::node::Node;
use scraper::{Html, Selector};

pub struct FullTextExtractor<'a> {
    html: &'a str,
}

impl<'a> FullTextExtractor<'a> {
    pub fn new(html: &'a str) -> Self {
        FullTextExtractor { html: html }
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

    pub fn extract(&self) -> Result<(Vec<String>, Vec<String>)> {
        let html = Html::parse_document(self.html);

        let div_part = Selector::parse("div.part").unwrap();
        let section = Selector::parse("section").unwrap();
        let h3 = Selector::parse("h3").unwrap();

        let ascii = Regex::new("^[\x20-\x7E].*$").unwrap();

        let mut text_ja: Vec<String> = Vec::new();
        let mut text_en: Vec<String> = Vec::new();
        for part in html.select(&div_part) {
            let section = part.select(&section).next().unwrap();

            let title = section.select(&h3).next().unwrap().text().next().unwrap();

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

            if ascii.is_match(&title) {
                text_en.push(full_text.join(""));
            } else {
                text_ja.push(full_text.join(" "));
            }
        }

        Ok((text_ja, text_en))
    }
}
