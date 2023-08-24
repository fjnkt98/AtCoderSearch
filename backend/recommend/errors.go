package recommend

import "github.com/morikuni/failure"

const (
	DBError                 failure.StringCode = "DBError"
	DecodeError             failure.StringCode = "DecodeError"
	EncodeError             failure.StringCode = "EncodeError"
	FileOperationError      failure.StringCode = "FileOperationError"
	RequestCreationError    failure.StringCode = "RequestCreationError"
	RequestExecutionError   failure.StringCode = "RequestExecutionError"
	SearcherInitializeError failure.StringCode = "SearcherInitializeError"
)
