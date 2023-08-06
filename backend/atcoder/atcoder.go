package atcoder

import (
	"bytes"
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
	client http.Client
}

func NewAtCoderClient(username string, password string) (AtCoderClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return AtCoderClient{}, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest(http.MethodGet, "https://atcoder.jp/login", nil)
	if err != nil {
		return AtCoderClient{}, fmt.Errorf("failed to create new request to get login csrf token: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return AtCoderClient{}, fmt.Errorf("failed to request to get login csrf token: %w", err)
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

	req, err = http.NewRequest(http.MethodPost, "https://atcoder.jp/login", strings.NewReader(params.Encode()))
	if err != nil {
		return AtCoderClient{}, fmt.Errorf("failed to create new request to login to atcoder: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = client.Do(req)
	if err != nil || (res.StatusCode/100) != 2 {
		return AtCoderClient{}, fmt.Errorf("login authentication to atcoder failed: %w", err)
	}

	return AtCoderClient{client: client}, nil
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

func (c *AtCoderClient) FetchSubmissionList(contestID string, page uint) (SubmissionList, error) {
	p, err := url.JoinPath("https://atcoder.jp", "contests", contestID, "submissions")
	if err != nil {
		return SubmissionList{}, fmt.Errorf("invalid contestID `%s` was given: %w", contestID, err)
	}
	u, err := url.Parse(p)
	if err != nil {
		return SubmissionList{}, fmt.Errorf("invalid contestID `%s` was given: %w", contestID, err)
	}
	u.RawQuery = fmt.Sprintf("page=%d", page)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return SubmissionList{}, fmt.Errorf("failed to create new request to get submissions html at page %d of the contest `%s`: %w", page, contestID, err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return SubmissionList{}, fmt.Errorf("failed to get submissions html at page %d of the contest `%s`: %w", page, contestID, err)
	}
	if res.StatusCode == http.StatusNotFound {
		return SubmissionList{
			MaxPage:     0,
			Submissions: nil,
		}, nil
	}
	if res.StatusCode != http.StatusOK {
		return SubmissionList{}, fmt.Errorf("non-ok status returned at page %d of the contest `%s`: %w", page, contestID, err)
	}

	defer res.Body.Close()
	list, err := scrapeSubmissions(res.Body)
	if err != nil {
		return SubmissionList{}, fmt.Errorf("failed to scrape submission list from html at page %d of the contest `%s`: %w", page, contestID, err)
	}
	submissions := make([]Submission, len(list.Submissions))
	for i, s := range list.Submissions {
		s.ContestID = contestID
		submissions[i] = s
	}
	list.Submissions = submissions

	return list, nil
}

func scrapeSubmissions(html io.Reader) (SubmissionList, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return SubmissionList{}, fmt.Errorf("failed to read html from reader: %w", err)
	}

	var maxPage uint
	pagePattern, _ := regexp.Compile(`page=\d+$`)
	doc.Find("a").Each(func(_ int, a *goquery.Selection) {
		if href, ok := a.Attr("href"); ok {
			if pagePattern.MatchString(href) {
				pageString := a.Text()
				page, err := strconv.Atoi(pageString)
				if err != nil {
					return
				}

				if page := uint(page); maxPage < uint(page) {
					maxPage = page
				}
			}
		}
	})

	tbody := doc.Find("tbody")
	if tbody == nil {
		return SubmissionList{}, fmt.Errorf("failed to read html from reader: %w", err)
	}
	submissions := make([]Submission, 0)
	tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		var s Submission
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				dt, _ := time.Parse("2006-01-02 15:04:05Z0700", td.Text())
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
				s.Point, _ = strconv.ParseFloat(td.Text(), 64)
			case 5:
				lengthStr := td.Text()
				lengthStr = strings.TrimSuffix(lengthStr, "Byte")
				lengthStr = strings.TrimSpace(lengthStr)
				length, _ := strconv.Atoi(lengthStr)
				s.Length = uint64(length)
			case 6:
				s.Result = td.Text()
			case 7, 9:
				a := td.Find("td > a")
				if href, ok := a.Attr("href"); ok {
					s.ID, _ = strconv.ParseInt(path.Base(href), 10, 64)
				} else {
					t := td.Text()
					t = strings.TrimSuffix(t, "ms")
					t = strings.TrimSpace(t)
					et, _ := strconv.ParseUint(t, 10, 64)
					s.ExecutionTime = &et

				}
			}
		})
		submissions = append(submissions, s)
	})

	return SubmissionList{
		MaxPage:     maxPage,
		Submissions: submissions,
	}, nil
}
