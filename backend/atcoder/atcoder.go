package atcoder

import (
	"bytes"
	"context"
	"fjnkt98/atcodersearch/acs"
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
	"github.com/morikuni/failure"
)

type AtCoderClient struct {
	client http.Client
}

func NewAtCoderClient(ctx context.Context, username string, password string) (AtCoderClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return AtCoderClient{}, failure.Translate(err, acs.CookieInitializeError, failure.Message("failed to create cookie jar"))
	}

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://atcoder.jp/login", nil)
	if err != nil {
		return AtCoderClient{}, failure.Translate(err, acs.LoginError, failure.Context{"url": "https://atcoder.jp/login"}, failure.Message("failed to create new request to get login csrf token"))
	}

	res, err := client.Do(req)
	if err != nil {
		return AtCoderClient{}, failure.Translate(err, acs.LoginError, failure.Context{"url": "https://atcoder.jp/login"}, failure.Message("failed to request to get login csrf token"))
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

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, "https://atcoder.jp/login", strings.NewReader(params.Encode()))
	if err != nil {
		return AtCoderClient{}, failure.Translate(err, acs.LoginError, failure.Context{"url": "https://atcoder.jp/login"}, failure.Message("failed to create new request to login to atcoder"))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = client.Do(req)
	if err != nil || (res.StatusCode/100) != 2 {
		return AtCoderClient{}, failure.Translate(err, acs.LoginError, failure.Context{"url": "https://atcoder.jp/login"}, failure.Message("login authentication to atcoder failed"))
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

func (c *AtCoderClient) FetchSubmissionList(ctx context.Context, contestID string, page int) (SubmissionList, error) {
	p, err := url.JoinPath("https://atcoder.jp", "contests", contestID, "submissions")
	if err != nil {
		return SubmissionList{}, failure.Translate(err, acs.InvalidURL, failure.Context{"contestID": contestID}, failure.Message("invalid contestID was given"))
	}
	u, err := url.Parse(p)
	if err != nil {
		return SubmissionList{}, failure.Translate(err, acs.InvalidURL, failure.Context{"contestID": contestID}, failure.Message("invalid contestID was given"))
	}
	u.RawQuery = fmt.Sprintf("page=%d", page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return SubmissionList{}, failure.Translate(err, acs.RequestError, failure.Context{"url": u.String(), "page": strconv.Itoa(int(page)), "contestID": contestID}, failure.Message("failed to create new request to get submissions html"))
	}

	res, err := c.client.Do(req)
	if err != nil {
		return SubmissionList{}, failure.Translate(err, acs.RequestError, failure.Context{"url": u.String(), "page": strconv.Itoa(int(page)), "contestID": contestID}, failure.Messagef("failed to get submissions html at page %d of the contest `%s`", page, contestID))
	}
	if res.StatusCode == http.StatusNotFound {
		return SubmissionList{
			MaxPage:     0,
			Submissions: nil,
		}, nil
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		return SubmissionList{}, failure.New(acs.RequestError, failure.Context{"url": u.String(), "page": strconv.Itoa(int(page)), "contestID": contestID, "response": buf.String()}, failure.Messagef("non-ok status returned at page %d of the contest `%s`", page, contestID))
	}

	defer res.Body.Close()
	list, err := scrapeSubmissions(res.Body)
	if err != nil {
		return SubmissionList{}, failure.Translate(err, acs.ScrapeError, failure.Context{}, failure.Messagef("failed to scrape submission list from html at page %d of the contest `%s`", page, contestID))
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
		return SubmissionList{}, failure.Translate(err, acs.ReadError, failure.Message("failed to read html from reader"))
	}

	var maxPage int
	pagePattern, _ := regexp.Compile(`page=\d+$`)
	doc.Find("a").Each(func(_ int, a *goquery.Selection) {
		if href, ok := a.Attr("href"); ok {
			if pagePattern.MatchString(href) {
				pageString := a.Text()
				page, err := strconv.Atoi(pageString)
				if err != nil {
					return
				}

				if page := int(page); maxPage < int(page) {
					maxPage = page
				}
			}
		}
	})

	submissions := make([]Submission, 0)
	doc.Find("tbody").Find("tr").Each(func(_ int, tr *goquery.Selection) {
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
				s.Length = int64(length)
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
					et, _ := strconv.ParseInt(t, 10, 64)
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

func (c *AtCoderClient) FetchSubmissionResult(ctx context.Context, contestID string, submissionID int64) (string, error) {
	p, err := url.JoinPath("https://atcoder.jp/contests", contestID, "submissions", strconv.Itoa(int(submissionID)))
	if err != nil {
		return "", failure.Translate(err, acs.InvalidURL, failure.Context{"contestID": contestID, "submissionID": strconv.Itoa(int(submissionID))}, failure.Message("invalid contestID or submissionID was given"))
	}
	u, err := url.Parse(p)
	if err != nil {
		return "", failure.Translate(err, acs.InvalidURL, failure.Context{"contestID": contestID, "submissionID": strconv.Itoa(int(submissionID))}, failure.Message("invalid contestID or submissionID was given"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", failure.Translate(err, acs.RequestError, failure.Context{"url": u.String()}, failure.Message("failed to create new request to get submission html"))
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", failure.Translate(err, acs.RequestError, failure.Context{"url": u.String()}, failure.Message("failed to get submission html"))
	}
	if res.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		return "", failure.New(acs.RequestError, failure.Context{"url": u.String(), "response": buf.String()}, failure.Message("non-ok status returned"))
	}

	defer res.Body.Close()
	result, err := scrapeSubmission(res.Body)
	if err != nil {
		return "", failure.Translate(err, acs.ScrapeError, failure.Context{"url": u.String()}, failure.Message("failed to scrape submission from html"))
	}

	return result, nil
}

func scrapeSubmission(html io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return "", failure.Translate(err, acs.ReadError, failure.Message("failed to read html from reader"))
	}

	result := strings.TrimSpace(doc.Find("#judge-status").Find("span").Text())
	return result, nil
}
