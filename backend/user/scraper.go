package user

import (
	"io"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/morikuni/failure"
)

var RANK_RE, _ = regexp.Compile(`\((\d+)\)`)

func Scrape(html io.Reader) ([]User, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, failure.Translate(err, ReadError, failure.Message("failed to read html from reader"))
	}

	users := make([]User, 0)
	doc.Find(".table > tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var user User
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				m := RANK_RE.FindStringSubmatch(td.Find("span").Text())
				if rank, err := strconv.Atoi(m[1]); err == nil {
					user.Rank = int(rank)
				}

				activeRankStr := strings.TrimSpace(td.Contents().Not("span").Text())
				if r, err := strconv.Atoi(activeRankStr); err == nil {
					rank := int(r)
					user.ActiveRank = &rank
				}
			case 1:
				td.Find("a").Each(func(k int, a *goquery.Selection) {
					switch k {
					case 0:
						img := a.Find("img")
						if src, ok := img.Attr("src"); ok {
							country, _, _ := strings.Cut(path.Base(src), ".")
							user.Country = &country
						}
					case 1:
						user.UserName = a.Find("a > span").Text()
					case 2:
						affiliation := a.Find("a > span").Text()
						if affiliation != "" {
							user.Affiliation = &affiliation
						}
					}
				})

				img := td.Find("td > img")
				if src, ok := img.Attr("src"); ok {
					crown, _, _ := strings.Cut(path.Base(src), ".")
					user.Crown = &crown
				}
			case 2:
				if year, err := strconv.Atoi(td.Text()); err == nil {
					year := int(year)
					user.BirthYear = &year
				}
			case 3:
				if rating, err := strconv.Atoi(td.Text()); err == nil {
					user.Rating = rating
				}
			case 4:
				if highestRating, err := strconv.Atoi(td.Text()); err == nil {
					user.HighestRating = highestRating
				}
			case 5:
				if joinCount, err := strconv.Atoi(td.Text()); err == nil {
					user.JoinCount = int(joinCount)
				}
			case 6:
				if wins, err := strconv.Atoi(td.Text()); err == nil {
					user.Wins = int(wins)
				}
			}
		})
		users = append(users, user)
	})

	return users, nil
}
