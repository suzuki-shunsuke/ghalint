package act

import "log/slog"

type HasLogLevelError struct {
	LogLevel slog.Level
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
		LogLevel: slog.LevelDebug,
		Err:      err,
	}
}
