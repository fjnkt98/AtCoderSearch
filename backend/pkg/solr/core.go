package solr

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/goark/errs"
)

var ErrNonOKResponse = errs.New("non-ok status code returned")

type SolrCore interface {
	Name() string
	Ping(ctx context.Context) (io.ReadCloser, error)
	Status(ctx context.Context) (io.ReadCloser, error)
	Reload(ctx context.Context) (io.ReadCloser, error)
	Select(ctx context.Context, params url.Values) (io.ReadCloser, error)
	Post(ctx context.Context, body io.Reader, contentType string) (io.ReadCloser, error)
}

type core struct {
	host   url.URL
	name   string
	client *http.Client
}

func NewSolrCore(host string, name string) (SolrCore, error) {
	parsedHost, err := url.Parse(host)
	if err != nil {
		return nil, errs.New(
			"invalid host or core name was given",
			errs.WithCause(err),
			errs.WithContext("core name", name),
			errs.WithContext("host", host),
		)
	}

	return &core{
		host:   url.URL{Scheme: parsedHost.Scheme, Host: parsedHost.Host},
		name:   name,
		client: &http.Client{},
	}, nil
}

func MustNewSolrCore(host string, name string) SolrCore {
	core, err := NewSolrCore(host, name)
	if err != nil {
		panic(err)
	}

	return core
}

func (c *core) Name() string {
	return c.name
}

func (c *core) Ping(ctx context.Context) (io.ReadCloser, error) {
	u := c.host.JoinPath("solr", c.name, "admin", "ping")
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare ping request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}

	if res, err := c.client.Do(req); err != nil {
		return nil, errs.New(
			"failed to execute ping request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	} else {
		if res.StatusCode != http.StatusOK {
			return res.Body, errs.Wrap(
				ErrNonOKResponse,
				errs.WithCause(err),
				errs.WithContext("url", u.String()),
			)
		}
		return res.Body, nil
	}
}

func (c *core) Status(ctx context.Context) (io.ReadCloser, error) {
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

	if res, err := c.client.Do(req); err != nil {
		return nil, errs.New(
			"failed to execute status request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	} else {
		if res.StatusCode != http.StatusOK {
			return res.Body, errs.Wrap(
				ErrNonOKResponse,
				errs.WithCause(err),
				errs.WithContext("url", u.String()),
			)
		}
		return res.Body, nil
	}
}

func (c *core) Reload(ctx context.Context) (io.ReadCloser, error) {
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

	if res, err := c.client.Do(req); err != nil {
		return nil, errs.New(
			"failed to execute reload request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	} else {
		if res.StatusCode != http.StatusOK {
			return res.Body, errs.Wrap(
				ErrNonOKResponse,
				errs.WithCause(err),
				errs.WithContext("url", u.String()),
			)
		}
		return res.Body, nil
	}
}

func (c *core) Select(ctx context.Context, params url.Values) (io.ReadCloser, error) {
	u := c.host.JoinPath("solr", c.name, "select")
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare select request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}

	if res, err := c.client.Do(req); err != nil {
		return nil, errs.New(
			"failed to execute select request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	} else {
		if res.StatusCode != http.StatusOK {
			return res.Body, errs.Wrap(
				ErrNonOKResponse,
				errs.WithCause(err),
				errs.WithContext("url", u.String()),
			)
		}
		return res.Body, nil
	}
}

func (c *core) Post(ctx context.Context, body io.Reader, contentType string) (io.ReadCloser, error) {
	u := c.host.JoinPath("solr", c.name, "update")

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, errs.New(
			"failed to prepare post request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	}
	req.Header.Add("Content-Type", contentType)

	if res, err := c.client.Do(req); err != nil {
		return nil, errs.New(
			"failed to execute post request",
			errs.WithCause(err),
			errs.WithContext("url", u.String()),
		)
	} else {
		if res.StatusCode != http.StatusOK {
			return res.Body, errs.Wrap(
				ErrNonOKResponse,
				errs.WithCause(err),
				errs.WithContext("url", u.String()),
			)
		}
		return res.Body, nil
	}
}
