package solr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/goark/errs"
)

var ErrNonOKResponse = errs.New("non-ok status code returned")

type SolrCore struct {
	host   url.URL
	name   string
	client *http.Client
}

func NewSolrCore(host string, name string) (*SolrCore, error) {
	parsedHost, err := url.Parse(host)
	if err != nil {
		return nil, errs.New(
			"invalid host or core name was given",
			errs.WithCause(err),
			errs.WithContext("core name", name),
			errs.WithContext("host", host),
		)
	}

	return &SolrCore{
		host:   url.URL{Scheme: parsedHost.Scheme, Host: parsedHost.Host},
		name:   name,
		client: &http.Client{},
	}, nil
}

func MustNewSolrCore(host string, name string) *SolrCore {
	core, err := NewSolrCore(host, name)
	if err != nil {
		panic(err)
	}

	return core
}

func (c *SolrCore) Name() string {
	return c.name
}

func (c *SolrCore) Ping(ctx context.Context) (*PingResponse, error) {
	u := c.host.JoinPath("solr", c.name, "admin", "ping")
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare ping request",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute ping request",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	defer res.Body.Close()
	var body PingResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errs.New(
			"failed to decode ping response",
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}

	if res.StatusCode != http.StatusOK {
		return &body, errs.Wrap(
			ErrNonOKResponse,
			errs.WithCause(err),
			errs.WithContext("uri", u.String()),
		)
	}
	return &body, nil
}

func (c *SolrCore) Status(ctx context.Context) (*CoreStatus, error) {
	v := url.Values{}
	v.Set("action", "STATUS")
	v.Set("core", c.name)
	u := c.host.JoinPath("solr", "admin", "cores")
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare status request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute status request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	defer res.Body.Close()
	var body CoreStatuses
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errs.New(
			"failed to decode status response",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	if res.StatusCode != http.StatusOK {
		return nil, errs.Wrap(
			ErrNonOKResponse,
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	status, ok := body.Status[c.name]
	if !ok {
		return nil, errs.New(
			"core not found",
			errs.WithContext("statues", body),
		)
	}
	return &status, nil
}

func (c *SolrCore) Reload(ctx context.Context) (*SimpleResponse, error) {
	v := url.Values{}
	v.Set("action", "RELOAD")
	v.Set("core", c.name)
	u := c.host.JoinPath("solr", "admin", "cores")
	u.RawQuery = v.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare reload request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute reload request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	defer res.Body.Close()
	var body SimpleResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errs.New(
			"failed to decode reload response",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	if res.StatusCode != http.StatusOK {
		return &body, errs.Wrap(
			ErrNonOKResponse,
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	return &body, err
}

func (c *SolrCore) NewSelect() *SelectQuery {
	return NewSelectQuery(c.client, c.host.JoinPath("solr", c.name, "select"))
}

func (c *SolrCore) Post(ctx context.Context, src io.Reader, contentType string) (*SimpleResponse, error) {
	u := c.host.JoinPath("solr", c.name, "update")

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), src)
	if err != nil {
		return nil, errs.New(
			"failed to prepare post request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	req.Header.Add("Content-Type", contentType)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute post request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	defer res.Body.Close()
	var body SimpleResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errs.New(
			"failed to decode post response",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	if res.StatusCode != http.StatusOK {
		return &body, errs.Wrap(
			ErrNonOKResponse,
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	return &body, nil
}

func (c *SolrCore) Commit(ctx context.Context) (*SimpleResponse, error) {
	src := strings.NewReader(`{"commit": {}}`)
	if res, err := c.Post(ctx, src, "application/json"); err != nil {
		return res, errs.New(
			"failed to commit",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}

func (c *SolrCore) Optimize(ctx context.Context) (*SimpleResponse, error) {
	src := strings.NewReader(`{"optimize": {}}`)
	if res, err := c.Post(ctx, src, "application/json"); err != nil {
		return res, errs.New(
			"failed to optimize",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}

func (c *SolrCore) Rollback() (*SimpleResponse, error) {
	src := strings.NewReader(`{"rollback": {}}`)
	if res, err := c.Post(context.Background(), src, "application/json"); err != nil {
		return res, errs.New(
			"failed to rollback",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}

func (c *SolrCore) Delete(ctx context.Context) (*SimpleResponse, error) {
	src := strings.NewReader(`{"delete":{"query":"*:*"}}`)
	if res, err := c.Post(ctx, src, "application/json"); err != nil {
		return res, errs.New(
			"failed to delete",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}
