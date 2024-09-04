package atcoder

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
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

func (c *AtCoderClient) Login(username, password string) error {
	panic("not implemented")
}

func (c *AtCoderClient) FetchProblemHTML(contestID, problemID string) (string, error) {
	panic("not implemented")
}

func (c *AtCoderClient) FetchUsers(page int) ([]User, error) {
	panic("not implemented")
}

func (c *AtCoderClient) FetchSubmissions(contestID string, page int) ([]Submission, error) {
	panic("not implemented")
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

func extractCSRFToken(html string) string {
	panic("not implemented")
}

func scrapeSubmissions(html io.Reader) ([]Submission, error) {
	panic("not implemented")
}

func scrapeUsers(html io.Reader) ([]User, error) {
	panic("not implemented")
}
