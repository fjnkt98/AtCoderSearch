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

// ErrNonOKResponse is the error returned when Solr returns a status code other than 200.
var ErrNonOKResponse = errs.New("non-ok status code returned")

// SolrCore is a struct that abstracts the core of Solr.
type SolrCore struct {
	host   url.URL
	name   string
	client *http.Client
}

// NewSolrCore creates a new SolrCore instance.
// It returns an error when the given host url is invalid.
// It doesn't return an error to specify the name of a core that does not exist.
func NewSolrCore(host string, name string) (*SolrCore, error) {
	parsedHost, err := url.Parse(host)
	if err != nil {
		return nil, errs.New(
			"invalid host or core name was given",
			errs.WithCause(err),
			errs.WithContext("core", name),
			errs.WithContext("host", host),
		)
	}

	return &SolrCore{
		host:   url.URL{Scheme: parsedHost.Scheme, Host: parsedHost.Host},
		name:   name,
		client: &http.Client{},
	}, nil
}

// MustNewSolrCore creates a new SolrCore instance.
// Panic if an error occurs internally.
func MustNewSolrCore(host string, name string) *SolrCore {
	core, err := NewSolrCore(host, name)
	if err != nil {
		panic(err)
	}

	return core
}

// Name returns the name of the Solr core.
func (c *SolrCore) Name() string {
	return c.name
}

// Ping pings to the Solr core.
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

// Status gets the status of the Solr core.
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

// Reload reloads the Solr core.
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

// NewSelect creates a new *SelectQuery instance.
func (c *SolrCore) NewSelect() *SelectQuery {
	return newSelectQuery(c.client, c.host.JoinPath("solr", c.name, "select"))
}

// NewMoreLikeThis creates a new *MoreLikeThisQuery instance.
func (c *SolrCore) NewMoreLikeThis() *MoreLikeThisQuery {
	return newMoreLikeThisQuery(c.client, c.host.JoinPath("solr", c.name, "select"))
}

// Post posts the data into Solr core.
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

// Commit commits the Solr core.
func (c *SolrCore) Commit(ctx context.Context) (*SimpleResponse, error) {
	src := strings.NewReader(`{"commit":{}}`)
	if res, err := c.Post(ctx, src, "application/json"); err != nil {
		return res, errs.New(
			"failed to commit",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}

// Optimize commits the Solr core with optimize.
func (c *SolrCore) Optimize(ctx context.Context) (*SimpleResponse, error) {
	src := strings.NewReader(`{"optimize":{}}`)
	if res, err := c.Post(ctx, src, "application/json"); err != nil {
		return res, errs.New(
			"failed to optimize",
			errs.WithCause(err),
		)
	} else {
		return res, nil
	}
}

// Rollback rollbacks the Solr core.
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

// Delete deletes all documents in the Solr core.
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
