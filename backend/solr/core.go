package solr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		return SolrCore[D, F]{}, fmt.Errorf("failed to create Solr core because: %s", err.Error())
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

func (core *SolrCore[D, F]) Ping() (SolrPingResponse, error) {
	req, err := http.NewRequest("GET", core.pingURL.String(), nil)
	if err != nil {
		return SolrPingResponse{}, fmt.Errorf("failed to prepare request: %s", err.Error())
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrPingResponse{}, fmt.Errorf("ping request failed: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			return SolrPingResponse{}, fmt.Errorf("ping returns non-ok response: %s: ", res.Status)
		}

		var pingResponse SolrPingResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&pingResponse); err != nil {
			return SolrPingResponse{}, fmt.Errorf("failed to unmarshal Solr ping response: %s", err.Error())
		} else {
			return pingResponse, nil
		}
	}
}

func (core *SolrCore[D, F]) Status() (SolrCoreStatus, error) {
	v := url.Values{}
	v.Set("action", "STATUS")
	v.Set("core", core.Name)
	u := core.adminURL
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SolrCoreStatus{}, fmt.Errorf("failed to prepare request: %s", err.Error())
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrCoreStatus{}, fmt.Errorf("status request failed: %s", err.Error())
	} else {
		if res.StatusCode != http.StatusOK {
			return SolrCoreStatus{}, fmt.Errorf("solr returns non-ok response: %s", res.Status)
		}

		var coreStatus SolrCoreList
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&coreStatus); err != nil {
			return SolrCoreStatus{}, fmt.Errorf("failed to unmarshal Solr status response: %s", err.Error())
		}

		status, ok := coreStatus.Status[core.Name]
		if ok {
			return status, nil
		} else {
			return SolrCoreStatus{}, fmt.Errorf("the core %s doesn't exists", core.Name)
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
		return SolrSimpleResponse{}, fmt.Errorf("failed to prepare request: %s", err.Error())
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrSimpleResponse{}, fmt.Errorf("reload request failed: %s", err.Error())
	} else {
		var reloadResponse SolrSimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&reloadResponse); err != nil {
			return SolrSimpleResponse{}, fmt.Errorf("failed to unmarshal Solr status response: %s", err.Error())
		}

		if res.StatusCode != http.StatusOK {
			return reloadResponse, fmt.Errorf("solr returns non-ok response: %+v", reloadResponse)
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
		return SolrSelectResponse[D, F]{}, fmt.Errorf("failed to prepare request: %s", err.Error())
	}

	if res, err := core.client.Do(req); err != nil {
		return SolrSelectResponse[D, F]{}, fmt.Errorf("failed to select request: %s", err.Error())
	} else {
		var selectResponse SolrSelectResponse[D, F]
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&selectResponse); err != nil {
			return SolrSelectResponse[D, F]{}, fmt.Errorf("failed to unmarshal Solr select response: %s", err.Error())
		}

		if res.StatusCode != http.StatusOK {
			return selectResponse, fmt.Errorf("select returns non-ok response: %+v", selectResponse)
		} else {
			return selectResponse, nil
		}
	}
}

func (core *SolrCore[D, F]) Post(body io.Reader, contentType string) (SolrSimpleResponse, error) {
	req, err := http.NewRequest("POST", core.postURL.String(), body)
	if err != nil {
		return SolrSimpleResponse{}, fmt.Errorf("failed to prepare request: %s", err.Error())
	}
	req.Header.Add("Content-Type", contentType)

	if res, err := core.client.Do(req); err != nil {
		return SolrSimpleResponse{}, fmt.Errorf("failed to post request: %s", err.Error())
	} else {
		var response SolrSimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return SolrSimpleResponse{}, fmt.Errorf("failed to unmarshal Solr post response: %s", err.Error())
		}

		if res.StatusCode != http.StatusOK {
			return response, fmt.Errorf("post returns non-ok response: %+v", response)
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
