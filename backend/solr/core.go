package solr

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/morikuni/failure"
)

const (
	CoreNotFound          failure.StringCode = "CoreNotFound"
	NotOK                 failure.StringCode = "NotOK"
	InvalidHost           failure.StringCode = "InvalidHost"
	RequestCreationError  failure.StringCode = "RequestCreationError"
	RequestExecutionError failure.StringCode = "RequestExecutionError"
	DecodeError           failure.StringCode = "DecodeError"
)

type SolrCore[D any, F any] struct {
	Name      string
	adminURL  url.URL
	pingURL   url.URL
	postURL   url.URL
	selectURL url.URL
	client    *http.Client
}

func NewSolrCore[D any, F any](coreName string, baseURL string) (SolrCore[D, F], error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return SolrCore[D, F]{}, failure.Translate(err, InvalidHost, failure.Context{"coreName": coreName, "baseURL": baseURL}, failure.Message("failed to create solr core"))
	}
	solrBaseURL := url.URL{Scheme: parsedBaseURL.Scheme, Host: parsedBaseURL.Host}

	return SolrCore[D, F]{
		Name:      coreName,
		adminURL:  *solrBaseURL.JoinPath("solr", "admin", "cores"),
		pingURL:   *solrBaseURL.JoinPath("solr", coreName, "admin", "ping"),
		postURL:   *solrBaseURL.JoinPath("solr", coreName, "update"),
		selectURL: *solrBaseURL.JoinPath("solr", coreName, "select"),
		client:    &http.Client{},
	}, nil
}

func (c *SolrCore[D, F]) Ping() (SolrPingResponse, error) {
	req, err := http.NewRequest("GET", c.pingURL.String(), nil)
	if err != nil {
		return SolrPingResponse{}, failure.Translate(err, RequestCreationError, failure.Context{"url": c.pingURL.String()}, failure.Message("failed to prepare request"))
	}

	if res, err := c.client.Do(req); err != nil {
		return SolrPingResponse{}, failure.Translate(err, RequestExecutionError, failure.Context{"url": c.pingURL.String()}, failure.Message("ping failed"))
	} else {
		if res.StatusCode != http.StatusOK {
			return SolrPingResponse{}, failure.New(NotOK, failure.Context{"url": c.pingURL.String()}, failure.Message("ping failed"))
		}

		var pingResponse SolrPingResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&pingResponse); err != nil {
			return SolrPingResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": c.pingURL.String()}, failure.Message("failed to decode solr ping response"))
		} else {
			return pingResponse, nil
		}
	}
}

func (c *SolrCore[D, F]) Status() (SolrCoreStatus, error) {
	v := url.Values{}
	v.Set("action", "STATUS")
	v.Set("core", c.Name)
	u := c.adminURL
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SolrCoreStatus{}, failure.Translate(err, RequestCreationError, failure.Context{"url": u.String()}, failure.Message("failed to prepare request"))
	}

	if res, err := c.client.Do(req); err != nil {
		return SolrCoreStatus{}, failure.Translate(err, RequestExecutionError, failure.Context{"url": u.String()}, failure.Message("status request failed"))
	} else {
		if res.StatusCode != http.StatusOK {
			return SolrCoreStatus{}, failure.New(NotOK, failure.Context{"url": u.String()}, failure.Message("couldn't get core status"))
		}

		var coreStatus SolrCoreList
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&coreStatus); err != nil {
			return SolrCoreStatus{}, failure.Translate(err, DecodeError, failure.Context{"url": u.String()}, failure.Messagef("failed to decode solr status response"))
		}

		status, ok := coreStatus.Status[c.Name]
		if ok {
			return status, nil
		} else {
			return SolrCoreStatus{}, failure.New(CoreNotFound, failure.Context{"url": u.String()}, failure.Messagef("the core `%s` doesn't exists", c.Name))
		}
	}
}

func (core *SolrCore[D, F]) Reload() (SolrSimpleResponse, error) {
	v := url.Values{}
	v.Set("action", "RELOAD")
	v.Set("core", core.Name)
	u := core.adminURL
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SolrSimpleResponse{}, failure.Translate(err, RequestCreationError, failure.Context{"url": u.String()}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrSimpleResponse{}, failure.Translate(err, RequestExecutionError, failure.Context{"url": u.String()}, failure.Message("reload request failed"))
	} else {
		var reloadResponse SolrSimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&reloadResponse); err != nil {
			return SolrSimpleResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": u.String()}, failure.Message("failed to decode solr status response"))
		}

		if res.StatusCode != http.StatusOK {
			return SolrSimpleResponse{}, failure.New(NotOK, failure.Context{"url": u.String()}, failure.Messagef("failed to reload core: %s", reloadResponse.Error.Msg))
		} else {
			return reloadResponse, nil
		}
	}
}

func (core *SolrCore[D, F]) Select(params url.Values) (SolrSelectResponse[D, F], error) {
	u := core.selectURL
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SolrSelectResponse[D, F]{}, failure.Translate(err, RequestCreationError, failure.Context{"url": u.String()}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrSelectResponse[D, F]{}, failure.Translate(err, RequestExecutionError, failure.Context{"url": u.String()}, failure.Message("failed to select request"))
	} else {
		var selectResponse SolrSelectResponse[D, F]
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&selectResponse); err != nil {
			return SolrSelectResponse[D, F]{}, failure.Translate(err, DecodeError, failure.Context{"url": u.String()}, failure.Message("failed to decode Solr select response"))
		}

		if res.StatusCode != http.StatusOK {
			return SolrSelectResponse[D, F]{}, failure.New(NotOK, failure.Context{"url": u.String()}, failure.Messagef("select request failed: %s", selectResponse.Error.Msg))
		} else {
			return selectResponse, nil
		}
	}
}

func (core *SolrCore[D, F]) Post(body io.Reader, contentType string) (SolrSimpleResponse, error) {
	req, err := http.NewRequest("POST", core.postURL.String(), body)
	if err != nil {
		return SolrSimpleResponse{}, failure.Translate(err, RequestCreationError, failure.Context{"url": core.postURL.String(), "Content-Type": contentType}, failure.Message("failed to prepare request"))
	}
	req.Header.Add("Content-Type", contentType)

	if res, err := core.client.Do(req); err != nil {
		return SolrSimpleResponse{}, failure.Translate(err, RequestExecutionError, failure.Context{"url": core.postURL.String(), "Content-Type": contentType}, failure.Message("failed to post request"))
	} else {
		var response SolrSimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return SolrSimpleResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": core.postURL.String(), "Content-Type": contentType}, failure.Message("failed to decode Solr post response"))
		}

		if res.StatusCode != http.StatusOK {
			return SolrSimpleResponse{}, failure.New(NotOK, failure.Context{"url": core.postURL.String(), "Content-Type": contentType}, failure.Messagef("post failed: %s", response.Error.Msg))
		} else {
			return response, nil
		}
	}
}

func (core *SolrCore[D, F]) Commit() (SolrSimpleResponse, error) {
	body := strings.NewReader(`{"commit": {}}`)
	return core.Post(body, "application/json")
}

func (core *SolrCore[D, F]) Optimize() (SolrSimpleResponse, error) {
	body := strings.NewReader(`{"optimize": {}}`)
	return core.Post(body, "application/json")
}

func (core *SolrCore[D, F]) Rollback() (SolrSimpleResponse, error) {
	body := strings.NewReader(`{"rollback": {}}`)
	return core.Post(body, "application/json")
}

func (core *SolrCore[D, F]) Truncate() (SolrSimpleResponse, error) {
	body := strings.NewReader(`{"delete":{"query": "*:*"}}`)
	return core.Post(body, "application/json")
}
