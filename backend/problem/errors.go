package problem

import "github.com/morikuni/failure"

const (
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	ExtractError            failure.StringCode = "ExtractError"
	FileOperationError      failure.StringCode = "FileOperationError"
	MinifyHTMLError         failure.StringCode = "MinifyHTMLError"
	RequestCreationError    failure.StringCode = "RequestCreationError"
	RequestExecutionError   failure.StringCode = "RequestExecutionError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
)
