package act

import "github.com/sirupsen/logrus"

type HasLogLevelError struct {
	LogLevel logrus.Level
	Err      error
}

func (e *HasLogLevelError) Error() string {
	return e.Err.Error()
}

func (e *HasLogLevelError) Unwrap() error {
	return e.Err
}

func debugError(err error) *HasLogLevelError {
	return &HasLogLevelError{
		LogLevel: logrus.DebugLevel,
		Err:      err,
	}
}
