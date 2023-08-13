package cmd

import "github.com/morikuni/failure"

const (
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	FileOperationError      failure.StringCode = "FileOperationError"
	ReadError               failure.StringCode = "ReadError"
	RequestCreationError    failure.StringCode = "RequestCreationError"
	RequestExecutionError   failure.StringCode = "RequestExecutionError"
	ScrapeError             failure.StringCode = "ScrapeError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
)
