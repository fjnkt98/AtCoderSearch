package solr

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/goark/errs"
)

func Ping(core SolrCore) (PingResponse, error) {
	ctx := context.Background()
	return PingWithContext(ctx, core)
}

func PingWithContext(ctx context.Context, core SolrCore) (PingResponse, error) {
	body, err := core.Ping(ctx)
	if err != nil {
		if errs.Is(err, ErrNonOKResponse) {
			var res PingResponse
			defer body.Close()
			if err2 := json.NewDecoder(body).Decode(&res); err2 != nil {
				return PingResponse{}, errs.Join(err, err2)
			} else {
				return res, errs.Wrap(err)
			}
		}

		return PingResponse{}, errs.Wrap(err)
	}

	var res PingResponse
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return PingResponse{}, errs.Wrap(err)
	}

	return res, nil
}

func Status(core SolrCore) (CoreStatus, error) {
	ctx := context.Background()
	return StatusWithContext(ctx, core)
}

func StatusWithContext(ctx context.Context, core SolrCore) (CoreStatus, error) {
	body, err := core.Status(ctx)
	if err != nil {
		if errs.Is(err, ErrNonOKResponse) {
			var res CoreStatuses
			defer body.Close()
			if err2 := json.NewDecoder(body).Decode(&res); err2 != nil {
				return CoreStatus{}, errs.Join(err, err2)
			} else {
				return res.Status[core.Name()], errs.Wrap(err)
			}
		}
	}

	var res CoreStatuses
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return CoreStatus{}, errs.Wrap(err)
	}
	if status, ok := res.Status[core.Name()]; ok {
		return status, nil
	} else {
		return CoreStatus{}, errs.New("core not found", errs.WithContext("name", core.Name()))
	}
}

func Reload(core SolrCore) (SimpleResponse, error) {
	ctx := context.Background()
	return ReloadWithContext(ctx, core)
}

func ReloadWithContext(ctx context.Context, core SolrCore) (SimpleResponse, error) {
	body, err := core.Reload(ctx)
	if err != nil {
		if errs.Is(err, ErrNonOKResponse) {
			var res SimpleResponse
			defer body.Close()
			if err2 := json.NewDecoder(body).Decode(&res); err2 != nil {
				return SimpleResponse{}, errs.Join(err, err2)
			} else {
				return res, errs.Wrap(err)
			}
		}

		return SimpleResponse{}, errs.Wrap(err)
	}

	var res SimpleResponse
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return SimpleResponse{}, errs.Wrap(err)
	}

	return res, nil
}

func Select[D any, F any](core SolrCore, params url.Values) (SelectResponse[D, F], error) {
	ctx := context.Background()
	return SelectWithContext[D, F](ctx, core, params)
}

func SelectWithContext[D any, F any](ctx context.Context, core SolrCore, params url.Values) (SelectResponse[D, F], error) {
	body, err := core.Select(ctx, params)
	if err != nil {
		if errs.Is(err, ErrNonOKResponse) {
			var res SelectResponse[D, F]
			defer body.Close()
			if err2 := json.NewDecoder(body).Decode(&res); err2 != nil {
				return SelectResponse[D, F]{}, errs.Join(err, err2)
			} else {
				return res, errs.Wrap(err)
			}
		}

		return SelectResponse[D, F]{}, errs.Wrap(err)
	}

	var res SelectResponse[D, F]
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return SelectResponse[D, F]{}, errs.Wrap(err)
	}

	return res, nil
}

func Post(core SolrCore, body io.Reader, contentType string) (SimpleResponse, error) {
	ctx := context.Background()
	return PostWithContext(ctx, core, body, contentType)
}

func PostWithContext(ctx context.Context, core SolrCore, body io.Reader, contentType string) (SimpleResponse, error) {
	resBody, err := core.Post(ctx, body, contentType)
	if err != nil {
		if errs.Is(err, ErrNonOKResponse) {
			var res SimpleResponse
			defer resBody.Close()
			if err2 := json.NewDecoder(resBody).Decode(&res); err2 != nil {
				return SimpleResponse{}, errs.Join(err, err2)
			} else {
				return res, errs.Wrap(err)
			}
		}

		return SimpleResponse{}, errs.Wrap(err)
	}

	var res SimpleResponse
	defer resBody.Close()
	if err := json.NewDecoder(resBody).Decode(&res); err != nil {
		return SimpleResponse{}, errs.Wrap(err)
	}

	return res, nil
}

func Commit(core SolrCore) (SimpleResponse, error) {
	ctx := context.Background()
	return CommitWithContext(ctx, core)
}

func CommitWithContext(ctx context.Context, core SolrCore) (SimpleResponse, error) {
	body := strings.NewReader(`{"commit": {}}`)
	return PostWithContext(ctx, core, body, "application/json")
}

func Optimize(core SolrCore) (SimpleResponse, error) {
	ctx := context.Background()
	return OptimizeWithContext(ctx, core)
}

func OptimizeWithContext(ctx context.Context, core SolrCore) (SimpleResponse, error) {
	body := strings.NewReader(`{"optimize": {}}`)
	return PostWithContext(ctx, core, body, "application/json")
}

func Rollback(core SolrCore) (SimpleResponse, error) {
	ctx := context.Background()
	return RollbackWithContext(ctx, core)
}

func RollbackWithContext(ctx context.Context, core SolrCore) (SimpleResponse, error) {
	body := strings.NewReader(`{"rollback": {}}`)
	return PostWithContext(ctx, core, body, "application/json")
}

func Truncate(core SolrCore) (SimpleResponse, error) {
	ctx := context.Background()
	return TruncateWithContext(ctx, core)
}

func TruncateWithContext(ctx context.Context, core SolrCore) (SimpleResponse, error) {
	body := strings.NewReader(`{"delete":{"query": "*:*"}}`)
	return PostWithContext(ctx, core, body, "application/json")
}
