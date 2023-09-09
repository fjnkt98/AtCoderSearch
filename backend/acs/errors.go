package acs

import "github.com/morikuni/failure"

const (
	ConvertError            failure.StringCode = "ConvertError"
	CookieInitializeError   failure.StringCode = "CookieInitializeError"
	CrawlError              failure.StringCode = "CrawlError"
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	EncodeError             failure.StringCode = "EncodeError"
	ExtractError            failure.StringCode = "ExtractError"
	FileOperationError      failure.StringCode = "FileOperationError"
	GenerateError           failure.StringCode = "GenerateError"
	Interrupt               failure.StringCode = "Interrupt"
	InvalidURL              failure.StringCode = "InvalidURL"
	LoginError              failure.StringCode = "LoginError"
	MinifyHTMLError         failure.StringCode = "MinifyHTMLError"
	PostError               failure.StringCode = "PostError"
	ReadError               failure.StringCode = "ReadError"
	RequestError            failure.StringCode = "RequestError"
	ScrapeError             failure.StringCode = "ScrapeError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
	UpdateIndexError        failure.StringCode = "UpdateIndexError"
)
