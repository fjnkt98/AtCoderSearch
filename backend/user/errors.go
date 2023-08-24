package user

import "github.com/morikuni/failure"

const (
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	EncodeError             failure.StringCode = "EncodeError"
	FileOperationError      failure.StringCode = "FileOperationError"
	ReadError               failure.StringCode = "ReadError"
	RequestCreationError    failure.StringCode = "RequestCreationError"
	RequestExecutionError   failure.StringCode = "RequestExecutionError"
	ScrapeError             failure.StringCode = "ScrapeError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
)
