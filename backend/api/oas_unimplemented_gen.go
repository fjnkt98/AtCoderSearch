// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// APICategoryGet implements GET /api/category operation.
//
// GET /api/category
func (UnimplementedHandler) APICategoryGet(ctx context.Context) (r *APICategoryGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APIContestGet implements GET /api/contest operation.
//
// GET /api/contest
func (UnimplementedHandler) APIContestGet(ctx context.Context, params APIContestGetParams) (r *APIContestGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APIHealthGet implements GET /api/health operation.
//
// GET /api/health
func (UnimplementedHandler) APIHealthGet(ctx context.Context) (r *APIHealthGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APILanguageGet implements GET /api/language operation.
//
// GET /api/language
func (UnimplementedHandler) APILanguageGet(ctx context.Context, params APILanguageGetParams) (r *APILanguageGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APIProblemGet implements GET /api/problem operation.
//
// GET /api/problem
func (UnimplementedHandler) APIProblemGet(ctx context.Context, params APIProblemGetParams) (r *APIProblemGetOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APIProblemPost implements POST /api/problem operation.
//
// POST /api/problem
func (UnimplementedHandler) APIProblemPost(ctx context.Context, req *APIProblemPostReq) (r *APIProblemPostOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APISubmissionPost implements POST /api/submission operation.
//
// POST /api/submission
func (UnimplementedHandler) APISubmissionPost(ctx context.Context, req *APISubmissionPostReq) (r *APISubmissionPostOK, _ error) {
	return r, ht.ErrNotImplemented
}

// APIUserPost implements POST /api/user operation.
//
// POST /api/user
func (UnimplementedHandler) APIUserPost(ctx context.Context, req *APIUserPostReq) (r *APIUserPostOK, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *ErrorResponseStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorResponseStatusCode) {
	r = new(ErrorResponseStatusCode)
	return r
}
