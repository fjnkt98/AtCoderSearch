use crate::models::errors::GeneratingError;
use ego_tree::NodeRef;
use scraper::node::Node;
use scraper::{Html, Selector};

type Result<T> = std::result::Result<T, GeneratingError>;

/// HTMLから問題文を取得する構造体
pub struct FullTextExtractor {
    span_ja: Selector,
    span_en: Selector,
    section: Selector,
    h3: Selector,
}

impl FullTextExtractor {
    pub fn new() -> Result<Self> {
        let span_ja =
            Selector::parse("span.lang-ja").map_err(|e| GeneratingError::SelectorError(e))?;
        let span_en =
            Selector::parse("span.lang-en").map_err(|e| GeneratingError::SelectorError(e))?;
        let section = Selector::parse("section").map_err(|e| GeneratingError::SelectorError(e))?;
        let h3 = Selector::parse("h3").map_err(|e| GeneratingError::SelectorError(e))?;

        Ok(FullTextExtractor {
            span_ja,
            span_en,
            section,
            h3,
        })
    }

    fn dfs(&self, element: &NodeRef<Node>) -> String {
        let mut result = Vec::new();

        for child in element.children() {
            match child.value() {
                Node::Element(e) => match e.name() {
                    // preタグ(入力例などのコードブロックのタグ)やh3タグは収集範囲外とする
                    "pre" | "h3" => continue,
                    "var" => {
                        // varタグの値の周りには空白を空ける
                        result.push(format!(" {} ", self.dfs(&child)));
                    }
                    _ => {
                        result.push(self.dfs(&child));
                    }
                },
                Node::Text(text) => {
                    result.push(text.trim().to_string());
                }
                _ => continue,
            };
        }

        result.join("")
    }

    /// HTML本文から問題文を取得するメソッド
    pub fn extract(&self, html: &str) -> Result<(Vec<String>, Vec<String>)> {
        let html = Html::parse_document(html);

        let mut text_ja: Vec<String> = Vec::new();
        let mut text_en: Vec<String> = Vec::new();

        // 日本語版の問題文は<span class="lang-ja">タグ内に定義されている。そのため日本語版の問題文を取得したい場合はこのタグの子要素を探しにいけばよい。
        // 英語版の問題文が用意されていない問題はこのタグが存在しないので、その場合はsectionタグを走査し、ボディが「問題文」であるh3タグを持つsectionをパースする。
        if let Some(ja) = html.select(&self.span_ja).next() {
            for section in ja.select(&self.section) {
                let Some(h3) = section.select(&self.h3).next() else {continue};
                let Some(h3) = h3.text().next() else {continue};

                // ボディに「問題文」を含むh3タグだった場合にその本文を取得する
                // 単に等価比較していないのはどっかの問題で「問題分」と誤字っている問題があった気がしたのと、両端に空白が含まれている場合でも対応するため。
                if h3.contains("問題") {
                    text_ja.push(self.dfs(&section));
                }
            }
        } else {
            for section in html.select(&self.section) {
                let Some(h3) = section.select(&self.h3).next() else {continue};
                let Some(h3) = h3.text().next() else {continue};

                // ボディに「問題文」を含むh3タグだった場合にその本文を取得する
                // 単に等価比較していないのはどっかの問題で「問題分」と誤字っている問題があった気がしたのと、両端に空白が含まれている場合でも対応するため。
                if h3.contains("問題") {
                    text_ja.push(self.dfs(&section));
                }
            }
        }

        // 英語版の問題分を取得する
        if let Some(en) = html.select(&self.span_en).next() {
            for section in en.select(&self.section) {
                let Some(h3) = section.select(&self.h3).next() else {continue};
                let Some(h3) = h3.text().next() else {continue};

                if h3.contains("Statement") {
                    text_en.push(self.dfs(&section));
                }
            }
        }

        Ok((text_ja, text_en))
    }
}
