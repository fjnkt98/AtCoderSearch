package atcoder

import (
	"bytes"
	"context"
	"errors"
	"fjnkt98/atcodersearch/pkg/ptr"
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

type AtCoderClient struct {
	client *http.Client
}

func NewAtCoderClient() (*AtCoderClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: time.Duration(30) * time.Second,
	}

	return &AtCoderClient{client}, nil
}

var ErrLoginAuth = errors.New("login authentication failed")

func (c *AtCoderClient) Login(ctx context.Context, username, password string) error {
	uri := "https://atcoder.jp/login"

	token, err := func() (string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return "", fmt.Errorf("create new request to get login csrf token: %w", err)
		}

		res, err := c.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("do request to get login csrf token: %w", err)
		}
		defer func() {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(res.Body); err != nil {
			return "", fmt.Errorf("read response body: %w", err)
		}

		return extractCSRFToken(buf.String())
	}()
	if err != nil {
		return err
	}

	params := url.Values{
		"username":   {username},
		"password":   {password},
		"csrf_token": {token},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("create new request to login: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request to login: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("login failed with code %s: %w", res.Status, ErrLoginAuth)
	}

	return nil
}

func (c *AtCoderClient) FetchProblemHTML(ctx context.Context, contestID, problemID string) (string, error) {
	u, err := url.Parse(fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", contestID, problemID))
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create new request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return buf.String(), nil
}

func (c *AtCoderClient) FetchUsers(ctx context.Context, page int) ([]User, error) {
	u, err := url.Parse(fmt.Sprintf("https://atcoder.jp/ranking/all?contestType=algo&page=%d", page))
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	users, err := scrapeUsers(res.Body)
	if err != nil {
		return nil, fmt.Errorf("scrape users: %w", err)
	}

	return users, nil
}

func (c *AtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]Submission, error) {
	u, err := url.Parse(fmt.Sprintf("https://atcoder.jp/contests/%s/submissions?page=%d", contestID, page))
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	submissions, err := scrapeSubmissions(res.Body)
	if err != nil {
		return nil, fmt.Errorf("scrape submissions: %w", err)
	}

	return submissions, err
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

var ErrCSRFTokenNotFound = errors.New("csrf token not found")
var CSRFPattern = regexp.MustCompile(`var csrfToken = "(.+)"`)

func extractCSRFToken(html string) (string, error) {
	m := CSRFPattern.FindStringSubmatch(html)

	if len(m) > 1 {
		return m[1], nil
	} else {
		return "", ErrCSRFTokenNotFound
	}
}

var ProblemPathPattern = regexp.MustCompile(`^/contests/([^/]+)/tasks/([^/]+)$`)
var ErrProblemIDOrContestIDNotFound = errors.New("problem id or contest id not match")
var ErrUserIDNotFound = errors.New("user id not found")

func scrapeSubmissions(html io.Reader) ([]Submission, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0)

	submissions := make([]Submission, 0, 20)
	doc.Find("tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var s Submission
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				if dt, err := time.Parse("2006-01-02 15:04:05Z0700", strings.TrimSpace(td.Text())); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					s.EpochSecond = dt.Unix()
				}
			case 1:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					m := ProblemPathPattern.FindStringSubmatch(href)
					if len(m) != 3 {
						errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, ErrProblemIDOrContestIDNotFound))
					} else {
						s.ContestID = m[1]
						s.ProblemID = m[2]
					}
				} else {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, ErrProblemIDOrContestIDNotFound))
				}
			case 2:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					s.UserID = path.Base(href)
				} else {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, ErrUserIDNotFound))
				}
			case 3:
				s.Language = td.Text()

			case 4:
				s.Point, err = strconv.ParseFloat(td.Text(), 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				}
			case 5:
				if length, err := strconv.Atoi(strings.TrimSuffix(td.Text(), " Byte")); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					s.Length = int32(length)
				}
			case 6:
				s.Result = td.Text()
			case 7, 9:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					if id, err := strconv.Atoi(path.Base(href)); err != nil {
						errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
					} else {
						s.ID = int64(id)
					}
				} else {
					if et, err := strconv.Atoi(strings.TrimSuffix(td.Text(), " ms")); err != nil {
						errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
					} else {
						s.ExecutionTime = ptr.To(int32(et))
					}
				}
			}
		})

		submissions = append(submissions, s)
	})

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return submissions, nil
}

var RankPattern = regexp.MustCompile(`\((\d+)\)`)
var ErrRankNotMatch = errors.New("rank not match")

func scrapeUsers(html io.Reader) ([]User, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0)

	users := make([]User, 0, 100)
	doc.Find(".table > tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var user User
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				m := RankPattern.FindStringSubmatch(td.Find("span").Text())
				if len(m) != 2 {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, ErrRankNotMatch))
				} else {
					if rank, err := strconv.Atoi(m[1]); err != nil {
						errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
					} else {
						user.Rank = int32(rank)
					}
				}
				if rank, err := strconv.Atoi(strings.TrimSpace(td.Contents().Not("span").Text())); err == nil {
					user.ActiveRank = ptr.To(int32(rank))
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
				if year, err := strconv.Atoi(td.Text()); err == nil {
					user.BirthYear = ptr.To(int32(year))
				}
			case 3:
				if rating, err := strconv.Atoi(td.Text()); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					user.Rating = int32(rating)
				}
			case 4:
				if highestRating, err := strconv.Atoi(td.Text()); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					user.HighestRating = int32(highestRating)
				}
			case 5:
				if joinCount, err := strconv.Atoi(td.Text()); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					user.JoinCount = int32(joinCount)
				}
			case 6:
				if wins, err := strconv.Atoi(td.Text()); err != nil {
					errs = append(errs, fmt.Errorf("at row %d col %d: %w", i, j, err))
				} else {
					user.Wins = int32(wins)
				}
			}
		})
		users = append(users, user)
	})

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return users, nil
}
