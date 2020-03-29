package config

import "fmt"

// doesn't actually implement the error interface!
type configError struct {
	cmd        string
	lineNumber int
}

type tooFewArgsError struct {
	configError
}

func tooFewArgs(cmd string, lineNumber int) *tooFewArgsError {
	return &tooFewArgsError{
		configError{
			cmd:        cmd,
			lineNumber: lineNumber,
		},
	}
}

func (err *tooFewArgsError) Error() string {
	return fmt.Sprintf("too few arguments for command '%s' on line %d", err.cmd, err.lineNumber)
}

type emptyFilePathError struct {
	configError
}

func emptyFilePath(cmd string, lineNumber int) *emptyFilePathError {
	return &emptyFilePathError{
		configError{
			cmd:        cmd,
			lineNumber: lineNumber,
		},
	}
}

func (err *emptyFilePathError) Error() string {
	return fmt.Sprintf("empty file path after command '%s' on line %d", err.cmd, err.lineNumber)
}
