package errs

import "errors"

var (
	ErrorUnmatachedHeaders       = errors.New("Expected headers don't match file's headers.")
	ErrorHeaderNotFound          = errors.New("Header '%s' not found in source document.")
	ErrorNoHeadersFound          = errors.New("No headers found in source document.")
	ErrorUnknownSourceFile       = errors.New("Unknown source file provided.")
	ErrorNoSourceFileName        = errors.New("No source filename provided for reporting.")
	ErrorEmptyRowFound           = errors.New("Empty row found in file.")
	ErrorFailedToCreateDirectory = errors.New("Failed to create directory. %s")
	ErrorFieldIsEmpty            = errors.New("'%s' Field cannot be empty.")
	ErrorFieldNotValid           = errors.New("'%s' Field is not valid.")
	ErrorCreditLimitInvalid      = errors.New("'%s' Field is invalid.")
	ErrorArgsDirSpecified        = errors.New("A directory cannot be passed as an argument.")
)
