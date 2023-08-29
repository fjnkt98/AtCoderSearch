package acs

import "github.com/morikuni/failure"

const (
	ConvertError       failure.StringCode = "ConvertError"
	DBError            failure.StringCode = "DBError"
	FileOperationError failure.StringCode = "FileOperationError"
	GenerateError      failure.StringCode = "GenerateError"
	Interrupt          failure.StringCode = "Interrupt"
	PostError          failure.StringCode = "PostError"
)
