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
	CoreNotFound failure.StringCode = "CoreNotFound"
	NotOK        failure.StringCode = "NotOK"
	InvalidHost  failure.StringCode = "InvalidHost"
	RequestError failure.StringCode = "RequestError"
	DecodeError  failure.StringCode = "DecodeError"
)

type Core struct {
	name    string
	baseURL url.URL
	client  *http.Client
}

func NewSolrCore(coreName string, baseURL string) (*Core, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, failure.Translate(err, InvalidHost, failure.Context{"coreName": coreName, "baseURL": baseURL}, failure.Message("failed to create solr core"))
	}
	solrBaseURL := url.URL{Scheme: parsedBaseURL.Scheme, Host: parsedBaseURL.Host}

	return &Core{
		name:    coreName,
		baseURL: solrBaseURL,
		client:  &http.Client{},
	}, nil
}

func Ping(core *Core) (PingResponse, error) {
	u := core.baseURL.JoinPath("solr", core.name, "ping").String()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return PingResponse{}, failure.Translate(err, RequestError, failure.Context{"url": u}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return PingResponse{}, failure.Translate(err, RequestError, failure.Context{"url": u}, failure.Message("ping failed"))
	} else {
		var body PingResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return PingResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": u}, failure.Message("failed to decode solr ping response"))
		}

		if res.StatusCode != http.StatusOK {
			return body, failure.New(NotOK, failure.Context{"url": u}, failure.Message("ping failed"))
		}

		return body, nil
	}
}

func Status(core *Core) (CoreStatus, error) {
	v := url.Values{}
	v.Set("action", "STATUS")
	v.Set("core", core.name)
	u := core.baseURL.JoinPath("solr", "admin", "cores")
	u.RawQuery = v.Encode()
	url := u.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CoreStatus{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return CoreStatus{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("status request failed"))
	} else {
		var body CoreList
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return CoreStatus{}, failure.Translate(err, DecodeError, failure.Context{"url": url}, failure.Messagef("failed to decode solr status response"))
		}

		if res.StatusCode != http.StatusOK {
			return CoreStatus{}, failure.New(NotOK, failure.Context{"url": url}, failure.Message("couldn't get core status"))
		}

		status, ok := body.Status[core.name]
		if ok {
			return status, nil
		} else {
			return CoreStatus{}, failure.New(CoreNotFound, failure.Context{"url": url}, failure.Messagef("the core `%s` doesn't exists", core.name))
		}
	}
}

func Reload(core *Core) (SimpleResponse, error) {
	v := url.Values{}
	v.Set("action", "RELOAD")
	v.Set("core", core.name)
	u := core.baseURL.JoinPath("solr", "admin", "cores")
	u.RawQuery = v.Encode()
	url := u.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SimpleResponse{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return SimpleResponse{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("reload request failed"))
	} else {
		var body SimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return SimpleResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": url}, failure.Message("failed to decode solr status response"))
		}

		if res.StatusCode != http.StatusOK {
			return body, failure.New(NotOK, failure.Context{"url": url}, failure.Messagef("failed to reload core: %s", body.Error.Msg))
		} else {
			return body, nil
		}
	}
}

func Select[D any, F any](core *Core, params url.Values) (SelectResponse[D, F], error) {
	u := core.baseURL.JoinPath("solr", core.name, "select")
	u.RawQuery = params.Encode()
	url := u.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SelectResponse[D, F]{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("failed to prepare request"))
	}

	if res, err := core.client.Do(req); err != nil {
		return SelectResponse[D, F]{}, failure.Translate(err, RequestError, failure.Context{"url": url}, failure.Message("failed to select request"))
	} else {
		var body SelectResponse[D, F]
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return SelectResponse[D, F]{}, failure.Translate(err, DecodeError, failure.Context{"url": url}, failure.Message("failed to decode Solr select response"))
		}

		if res.StatusCode != http.StatusOK {
			return body, failure.New(NotOK, failure.Context{"url": url}, failure.Messagef("select request failed: %s", body.Error.Msg))
		} else {
			return body, nil
		}
	}
}

func Post(core *Core, body io.Reader, contentType string) (SimpleResponse, error) {
	u := core.baseURL.JoinPath("solr", core.name, "update")
	url := u.String()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return SimpleResponse{}, failure.Translate(err, RequestError, failure.Context{"url": url, "Content-Type": contentType}, failure.Message("failed to prepare request"))
	}
	req.Header.Add("Content-Type", contentType)

	if res, err := core.client.Do(req); err != nil {
		return SimpleResponse{}, failure.Translate(err, RequestError, failure.Context{"url": url, "Content-Type": contentType}, failure.Message("failed to post request"))
	} else {
		var body SimpleResponse
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return SimpleResponse{}, failure.Translate(err, DecodeError, failure.Context{"url": url, "Content-Type": contentType}, failure.Message("failed to decode Solr post response"))
		}

		if res.StatusCode != http.StatusOK {
			return body, failure.New(NotOK, failure.Context{"url": url, "Content-Type": contentType}, failure.Messagef("post failed: %s", body.Error.Msg))
		} else {
			return body, nil
		}
	}
}

func Commit(core *Core) (SimpleResponse, error) {
	body := strings.NewReader(`{"commit": {}}`)
	return Post(core, body, "application/json")
}

func Optimize(core *Core) (SimpleResponse, error) {
	body := strings.NewReader(`{"optimize": {}}`)
	return Post(core, body, "application/json")
}

func Rollback(core *Core) (SimpleResponse, error) {
	body := strings.NewReader(`{"rollback": {}}`)
	return Post(core, body, "application/json")
}

func Truncate(core *Core) (SimpleResponse, error) {
	body := strings.NewReader(`{"delete":{"query": "*:*"}}`)
	return Post(core, body, "application/json")
}
