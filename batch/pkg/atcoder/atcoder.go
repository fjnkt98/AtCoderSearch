package atcoder

import (
	"bytes"
	"context"
	"errors"
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
)

var CSRFTokenPattern = regexp.MustCompile(`var csrfToken = "(.+)"`)

var ErrNotOK = errors.New("non-ok status returned")
var ErrUnmatchedRank = errors.New("rank text is not match the regexp")
var ErrNoCSRFToken = errors.New("there is no csrf token")
var ErrAuthFailed = errors.New("login authentication to atcoder failed")

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
	UserID        string
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
	client := &http.Client{}
	return NewAtCoderClientWithHTTPClient(client)
}

func NewAtCoderClientWithHTTPClient(client *http.Client) (*AtCoderClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	client.Jar = jar

	return &AtCoderClient{client: client}, nil
}

func extractCSRFToken(body string) (string, error) {
	m := CSRFTokenPattern.FindStringSubmatch(body)
	if m == nil {
		return "", ErrNoCSRFToken
	}

	return m[1], nil
}

func (c *AtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]Submission, error) {
	p, err := url.JoinPath("https://atcoder.jp", "contests", contestID, "submissions")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(p)
	if err != nil {
		return nil, fmt.Errorf("invalid contest id: %w", err)
	}
	u.RawQuery = fmt.Sprintf("page=%d", page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to fetch submissions: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		io.Copy(io.Discard, res.Body)
		return make([]Submission, 0), nil
	}
	if res.StatusCode != http.StatusOK {
		io.Copy(io.Discard, res.Body)
		return nil, fmt.Errorf("not ok from `%s`: %w", u, ErrNotOK)
	}

	submissions, err := scrapeSubmissions(res.Body)
	if err != nil {
		return nil, fmt.Errorf("scrape `%s`: %w", u, err)
	}
	for i := 0; i < len(submissions); i++ {
		submissions[i].ContestID = contestID
	}
	return submissions, nil
}

func scrapeSubmissions(html io.Reader) ([]Submission, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	occurredErrors := make([]error, 0)

	submissions := make([]Submission, 0)
	doc.Find("tbody").Find("tr").Each(func(_ int, tr *goquery.Selection) {
		var s Submission
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				dt, err := time.Parse("2006-01-02 15:04:05Z0700", strings.TrimSpace(td.Text()))
				if err != nil {
					occurredErrors = append(occurredErrors, err)
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
					occurredErrors = append(occurredErrors, err)
				}
			case 5:
				lengthStr := td.Text()
				lengthStr = strings.TrimSuffix(lengthStr, "Byte")
				lengthStr = strings.TrimSpace(lengthStr)
				length, err := strconv.Atoi(lengthStr)
				if err != nil {
					occurredErrors = append(occurredErrors, err)
				}
				s.Length = int32(length)
			case 6:
				s.Result = td.Text()
			case 7, 9:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					id, err := strconv.Atoi(path.Base(href))
					if err != nil {
						occurredErrors = append(occurredErrors, err)
					}
					s.ID = int64(id)
				} else {
					t := td.Text()
					t = strings.TrimSuffix(t, "ms")
					t = strings.TrimSpace(t)
					et, err := strconv.Atoi(t)
					if err != nil {
						occurredErrors = append(occurredErrors, err)
					}
					executionTime := int32(et)
					s.ExecutionTime = &executionTime
				}
			}
		})
		submissions = append(submissions, s)
	})

	if len(occurredErrors) > 0 {
		return nil, errors.Join(occurredErrors...)
	}

	return submissions, nil
}

func (c *AtCoderClient) FetchSubmissionResult(ctx context.Context, contestID string, submissionID int64) (string, error) {
	p, err := url.JoinPath("https://atcoder.jp/contests", contestID, "submissions", strconv.Itoa(int(submissionID)))
	if err != nil {
		return "", err
	}

	u, err := url.Parse(p)
	if err != nil {
		return "", fmt.Errorf("invalid contest id or submission id: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create request to `%s`: %w", u, err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to fetch submission result: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		io.Copy(io.Discard, res.Body)
		return "", nil
	}
	if res.StatusCode != http.StatusOK {
		io.Copy(io.Discard, res.Body)
		return "", fmt.Errorf("not ok from `%s`: %w", u, ErrNotOK)
	}

	result, err := scrapeSubmissionResult(res.Body)
	if err != nil {
		return "", fmt.Errorf("scrape `%s`: %w", u, err)
	}

	return result, nil
}

func scrapeSubmissionResult(html io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(doc.Find("#judge-status").Find("span").Text())
	return result, nil
}

func (c *AtCoderClient) FetchProblem(ctx context.Context, contestID string, problemID string) (string, error) {
	p, err := url.JoinPath("https://atcoder.jp/contests", contestID, "tasks", problemID)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(p)
	if err != nil {
		return "", fmt.Errorf("invalid contest id or problem id: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to fetch problems: %w", err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", fmt.Errorf("problem `%s`: %w", u, err)
	}

	return buf.String(), nil
}

func (c *AtCoderClient) FetchUsers(ctx context.Context, page int) ([]User, error) {
	u, err := url.Parse("https://atcoder.jp/ranking/all")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("contestType", "algo")
	v.Set("page", strconv.Itoa(page))
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to fetch users: %w", err)
	}
	defer res.Body.Close()

	users, err := scrapeUsers(res.Body)
	if err != nil {
		return nil, fmt.Errorf("scrape `%s`: %w", u, err)
	}

	return users, nil
}

var RankPattern = regexp.MustCompile(`\((\d+)\)`)

func scrapeUsers(html io.Reader) ([]User, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	occurredErrors := make([]error, 0)

	users := make([]User, 0)
	doc.Find(".table > tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var user User
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				m := RankPattern.FindStringSubmatch(td.Find("span").Text())
				if m == nil || len(m) < 1 {
					occurredErrors = append(occurredErrors, ErrUnmatchedRank)
					break
				}
				if rank, err := strconv.Atoi(m[1]); err != nil {
					occurredErrors = append(occurredErrors, err)
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
						user.UserID = a.Find("a > span").Text()
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
					occurredErrors = append(occurredErrors, fmt.Errorf("rating at row `%d`, col `%d`: %w", i, j, err))
				} else {
					user.Rating = int32(rating)
				}
			case 4:
				if highestRating, err := strconv.Atoi(td.Text()); err != nil {
					occurredErrors = append(occurredErrors, fmt.Errorf("highest rating at row `%d`, col `%d`: %w", i, j, err))
				} else {
					user.HighestRating = int32(highestRating)
				}
			case 5:
				if joinCount, err := strconv.Atoi(td.Text()); err != nil {
					occurredErrors = append(occurredErrors, fmt.Errorf("join count at row `%d`, col `%d`: %w", i, j, err))
				} else {
					user.JoinCount = int32(joinCount)
				}
			case 6:
				if wins, err := strconv.Atoi(td.Text()); err != nil {
					occurredErrors = append(occurredErrors, fmt.Errorf("win count at row `%d`, col `%d`: %w", i, j, err))
				} else {
					user.Wins = int32(wins)
				}
			}
		})
		users = append(users, user)
	})

	if len(occurredErrors) > 0 {
		return nil, errors.Join(occurredErrors...)
	}

	return users, nil
}

func (c *AtCoderClient) Login(ctx context.Context, username, password string) error {
	uri := "https://atcoder.jp/login"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return fmt.Errorf("get csrf token: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request to get csrf token: %w", err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return fmt.Errorf("login: %w", err)
	}
	body := buf.String()

	token, err := extractCSRFToken(body)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("username", username)
	params.Set("password", password)
	params.Set("csrf_token", token)

	req, err = http.NewRequest(http.MethodPost, uri, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = c.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	if (res.StatusCode / 100) != 2 {
		return ErrAuthFailed
	}

	return nil
}
