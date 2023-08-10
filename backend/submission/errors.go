package submission

import "github.com/morikuni/failure"

const (
	CrawlError              failure.StringCode = "CrawlError"
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	FileOperationError      failure.StringCode = "FileOperationError"
	RequestCreationError    failure.StringCode = "RequestCreationError"
	RequestExecutionError   failure.StringCode = "RequestExecutionError"
	ScrapeError             failure.StringCode = "ScrapeError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
)
