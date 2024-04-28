package atcoder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goark/errs"
)

type Submission struct {
	ID            int64
	EpochSecond   int64
	ProblemID     string
	ContestID     string
	UserID        string
	Language      string
	Point         float64
	Length        int32
	Result        string
	ExecutionTime *int32
}

type User struct {
	UserName      string
	Rating        int32
	HighestRating int32
	Affiliation   *string
	BirthYear     *int32
	Country       *string
	Crown         *string
	JoinCount     int32
	Rank          int32
	ActiveRank    *int32
	Wins          int32
}

type AtCoderClient struct {
	client *http.Client
}

func NewAtCoderClient() (*AtCoderClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, errs.New(
			"failed to create cookie jar",
			errs.WithCause(err),
		)
	}

	client := http.Client{
		Jar:     jar,
		Timeout: time.Duration(30) * time.Second,
	}

	return &AtCoderClient{client: &client}, nil
}

func extractCSRFToken(body string) string {
	pattern, _ := regexp.Compile(`var csrfToken = "(.+)"`)

	m := pattern.FindStringSubmatch(body)
	var token string
	if len(m) > 1 {
		token = m[1]
	}

	return token
}

func (c *AtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]Submission, error) {
	p, err := url.JoinPath("https://atcoder.jp", "contests", contestID, "submissions")
	if err != nil {
		return nil, errs.New(
			"invalid contest id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
		)
	}
	u, err := url.Parse(p)
	if err != nil {
		return nil, errs.New(
			"invalid contest id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
		)
	}
	u.RawQuery = fmt.Sprintf("page=%d", page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to create new request to get submissions html",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
			errs.WithContext("page", page),
			errs.WithContext("contest id", contestID),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to get submissions html",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
			errs.WithContext("page", page),
			errs.WithContext("contest id", contestID),
		)
	}
	if res.StatusCode == http.StatusNotFound {
		return make([]Submission, 0), nil
	}
	if res.StatusCode != http.StatusOK {
		return nil, errs.New(
			"non-ok status returned",
			errs.WithContext("uri", u.String()),
			errs.WithContext("page", page),
			errs.WithContext("contest id", contestID),
		)
	}

	defer res.Body.Close()
	submissions, err := scrapeSubmissions(res.Body)
	if err != nil {
		return nil, errs.New(
			"failed to scrape submissions",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
			errs.WithContext("page", page),
			errs.WithContext("contest id", contestID),
		)
	}
	for i := 0; i < len(submissions); i++ {
		submissions[i].ContestID = contestID
	}
	return submissions, nil
}

func scrapeSubmissions(html io.Reader) ([]Submission, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, errs.New(
			"failed to read html from reader",
			errs.WithCause(err),
		)
	}

	errors := make([]error, 0)

	submissions := make([]Submission, 0)
	doc.Find("tbody").Find("tr").Each(func(_ int, tr *goquery.Selection) {
		var s Submission
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				dt, err := time.Parse("2006-01-02 15:04:05Z0700", strings.TrimSpace(td.Text()))
				if err != nil {
					errors = append(errors, err)
				}
				s.EpochSecond = dt.Unix()
			case 1:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					s.ProblemID = path.Base(href)
				}
			case 2:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					s.UserID = path.Base(href)
				}
			case 3:
				s.Language = td.Text()
			case 4:
				s.Point, err = strconv.ParseFloat(td.Text(), 64)
				if err != nil {
					errors = append(errors, err)
				}
			case 5:
				lengthStr := td.Text()
				lengthStr = strings.TrimSuffix(lengthStr, "Byte")
				lengthStr = strings.TrimSpace(lengthStr)
				length, err := strconv.Atoi(lengthStr)
				if err != nil {
					errors = append(errors, err)
				}
				s.Length = int32(length)
			case 6:
				s.Result = td.Text()
			case 7, 9:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					id, err := strconv.Atoi(path.Base(href))
					if err != nil {
						errors = append(errors, err)
					}
					s.ID = int64(id)
				} else {
					t := td.Text()
					t = strings.TrimSuffix(t, "ms")
					t = strings.TrimSpace(t)
					et, err := strconv.Atoi(t)
					if err != nil {
						errors = append(errors, err)
					}
					executionTime := int32(et)
					s.ExecutionTime = &executionTime
				}
			}
		})
		submissions = append(submissions, s)
	})

	if len(errors) > 0 {
		return nil, errs.Join(errors...)
	}

	return submissions, nil
}

func (c *AtCoderClient) FetchSubmissionResult(ctx context.Context, contestID string, submissionID int64) (string, error) {
	p, err := url.JoinPath("https://atcoder.jp/contests", contestID, "submissions", strconv.Itoa(int(submissionID)))
	if err != nil {
		return "", errs.New(
			"invalid contest id or submission id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
			errs.WithContext("submission id", submissionID),
		)
	}
	u, err := url.Parse(p)
	if err != nil {
		return "", errs.New(
			"invalid contest id or submission id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
			errs.WithContext("submission id", submissionID),
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", errs.New(
			"failed to create new request to get submission html",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", errs.New(
			"failed to get submission html",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	if res.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if res.StatusCode != http.StatusOK {
		return "", errs.New(
			"non-ok status returned",
			errs.WithContext("uri", u.String()),
		)
	}

	defer res.Body.Close()
	result, err := scrapeSubmissionResult(res.Body)
	if err != nil {
		return "", errs.New(
			"failed to scrape submission from html",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	return result, nil
}

func scrapeSubmissionResult(html io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", errs.New(
			"failed to read html from reader",
			errs.WithCause(err),
		)
	}

	result := strings.TrimSpace(doc.Find("#judge-status").Find("span").Text())
	return result, nil
}

func (c *AtCoderClient) FetchProblem(ctx context.Context, contestID string, problemID string) (string, error) {
	p, err := url.JoinPath("https://atcoder.jp/contests", contestID, "tasks", problemID)
	if err != nil {
		return "", errs.New(
			"invalid contest id or problem id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
			errs.WithContext("problem id", problemID),
		)
	}
	u, err := url.Parse(p)
	if err != nil {
		return "", errs.New(
			"invalid contest id or problem id was given",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
			errs.WithContext("problem id", problemID),
		)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", errs.New(
			"failed to create request",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return "", errs.New(
			"request failed",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", errs.New(
			"failed to read body",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	return buf.String(), nil
}

func (c *AtCoderClient) FetchUsers(ctx context.Context, page int) ([]User, error) {
	u, _ := url.Parse("https://atcoder.jp/ranking/all")
	v := url.Values{}
	v.Set("contestType", "algo")
	v.Set("page", strconv.Itoa(page))
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to create request",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"request failed",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	defer res.Body.Close()

	users, err := scrapeUsers(res.Body)
	if err != nil {
		return nil, errs.New(
			"failed to scrape user page",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	return users, nil
}

var RANK_RE = regexp.MustCompile(`\((\d+)\)`)

func scrapeUsers(html io.Reader) ([]User, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, errs.New(
			"failed to read html from reader",
			errs.WithCause(err),
		)
	}

	errors := make([]error, 0)

	users := make([]User, 0)
	doc.Find(".table > tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var user User
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				m := RANK_RE.FindStringSubmatch(td.Find("span").Text())
				if m == nil || len(m) < 1 {
					errors = append(errors, errs.New("rank text is not match the regexp", errs.WithContext("row number", i), errs.WithContext("col number", j)))
					break
				}
				if rank, err := strconv.Atoi(m[1]); err != nil {
					errors = append(errors, err)
				} else {
					user.Rank = int32(rank)
				}

				activeRankStr := strings.TrimSpace(td.Contents().Not("span").Text())
				if r, err := strconv.Atoi(activeRankStr); err != nil {
					user.ActiveRank = nil
				} else {
					rank := int32(r)
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
				if year, err := strconv.Atoi(td.Text()); err != nil {
					user.BirthYear = nil
				} else {
					year := int32(year)
					user.BirthYear = &year
				}
			case 3:
				if rating, err := strconv.Atoi(td.Text()); err != nil {
					errors = append(errors, errs.New("failed to convert the rating", errs.WithCause(err), errs.WithContext("row number", i), errs.WithContext("col number", j)))
				} else {
					user.Rating = int32(rating)
				}
			case 4:
				if highestRating, err := strconv.Atoi(td.Text()); err != nil {
					errors = append(errors, errs.New("failed to convert the highest rating", errs.WithCause(err), errs.WithContext("row number", i), errs.WithContext("col number", j)))
				} else {
					user.HighestRating = int32(highestRating)
				}
			case 5:
				if joinCount, err := strconv.Atoi(td.Text()); err != nil {
					errors = append(errors, errs.New("failed to convert the join count", errs.WithCause(err), errs.WithContext("row number", i), errs.WithContext("col number", j)))
				} else {
					user.JoinCount = int32(joinCount)
				}
			case 6:
				if wins, err := strconv.Atoi(td.Text()); err != nil {
					errors = append(errors, errs.New("failed to convert the win count", errs.WithCause(err), errs.WithContext("row number", i), errs.WithContext("col number", j)))
				} else {
					user.Wins = int32(wins)
				}
			}
		})
		users = append(users, user)
	})

	if len(errors) > 0 {
		return nil, errs.Join(errors...)
	}

	return users, nil
}

func (c *AtCoderClient) Login(ctx context.Context, username, password string) error {
	uri := "https://atcoder.jp/login"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return errs.New(
			"failed to create new request to get login csrf token",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return errs.New(
			"failed to execute request to get login csrf token",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()

	token := extractCSRFToken(body)

	params := url.Values{}
	params.Set("username", username)
	params.Set("password", password)
	params.Set("csrf_token", token)

	req, err = http.NewRequest(http.MethodPost, uri, strings.NewReader(params.Encode()))
	if err != nil {
		return errs.New(
			"failed to create new request to login to atcoder",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
			errs.WithContext("token", token),
		)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.client.Do(req)
	if err != nil || (res.StatusCode/100) != 2 {
		return errs.New(
			"login authentication to atcoder failed",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
			errs.WithContext("token", token),
		)
	}

	return nil
}
