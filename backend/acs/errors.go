package acs

import "github.com/morikuni/failure"

const (
	FileOperationError failure.StringCode = "FileOperationError"
	ConvertError       failure.StringCode = "ConvertError"
	WriteError         failure.StringCode = "WriteError"
	GenerateError      failure.StringCode = "GenerateError"
	PostError          failure.StringCode = "PostError"
	Interrupt          failure.StringCode = "Interrupt"
)
