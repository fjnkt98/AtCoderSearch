pub trait QueryBuilder {
    fn build(self) -> Vec<(String, String)>;
}

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

    pub fn q(mut self, q: String) -> Self {
        self.q = Some(q);
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

    pub fn fq(mut self, fq: String) -> Self {
        self.fq = Some(fq);
        self
    }

    pub fn fl(mut self, fl: String) -> Self {
        self.fl = Some(fl);
        self
    }

    pub fn sort(mut self, sort: String) -> Self {
        self.sort = Some(sort);
        self
    }
}

impl QueryBuilder for StandardQueryBuilder {
    fn build(self) -> Vec<(String, String)> {
        let mut result: Vec<(String, String)> = Vec::new();
        match self.q {
            Some(q) => result.push((String::from("q"), q)),
            None => result.push((String::from("q"), String::from("*:*"))),
        };

        match self.start {
            Some(start) => result.push((String::from("start"), start.to_string())),
            None => result.push((String::from("start"), String::from("0"))),
        };

        match self.rows {
            Some(rows) => result.push((String::from("rows"), rows.to_string())),
            None => result.push((String::from("rows"), String::from("10"))),
        };

        if let Some(fq) = self.fq {
            result.push((String::from("fq"), fq));
        }

        if let Some(fl) = self.fl {
            result.push((String::from("fl"), fl));
        }

        match self.sort {
            Some(sort) => result.push((String::from("sort"), sort)),
            None => result.push((String::from("sort"), String::from("score desc"))),
        };

        result
    }
}
