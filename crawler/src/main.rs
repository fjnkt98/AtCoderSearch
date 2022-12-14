mod contest;

use contest::crawler;

fn main() {
    let contest_crawler = crawler::ContestCrawler::new();

    let contests = contest_crawler.crawl().expect("Failed to get contests.");

    println!("{:?}", contests);
}
