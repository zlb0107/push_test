package log

import (
	"errors"
)

var ErrMissingValue = errors.New("(MISSING)")

type Logger interface {
	Log(keyvals ...interface{}) error
}

func With(logger Logger, keyvals ...interface{}) Logger {
	if len(keyvals) == 0 {
		return logger
	}
	switch logger.(type) {
	case noopLogger:
		return logger
	}
	l := newContext(logger)
	kvs := append(l.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	return &context{
		logger: l.logger,
		// Limiting the capacity of the stored keyvals ensures that a new
		// backing array is created if the slice must grow in Log or With.
		// Using the extra capacity without copying risks a data race that
		// would violate the Logger interface contract.
		keyvals:   kvs[:len(kvs):len(kvs)],
		hasValuer: l.hasValuer || containsValuer(keyvals),
	}
}

func WithPrefix(logger Logger, keyvals ...interface{}) Logger {
	if len(keyvals) == 0 {
		return logger
	}
	switch logger.(type) {
	case noopLogger:
		return logger
	}
	l := newContext(logger)
	// Limiting the capacity of the stored keyvals ensures that a new
	// backing array is created if the slice must grow in Log or With.
	// Using the extra capacity without copying risks a data race that
	// would violate the Logger interface contract.
	n := len(l.keyvals) + len(keyvals)
	if len(keyvals)%2 != 0 {
		n++
	}
	kvs := make([]interface{}, 0, n)
	kvs = append(kvs, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	kvs = append(kvs, l.keyvals...)
	return &context{
		logger:    l.logger,
		keyvals:   kvs,
		hasValuer: l.hasValuer || containsValuer(keyvals),
	}
}

type context struct {
	logger    Logger
	keyvals   []interface{}
	hasValuer bool
}

func newContext(logger Logger) *context {
	if c, ok := logger.(*context); ok {
		return c
	}
	return &context{logger: logger}
}

// Log replaces all value elements (odd indexes) containing a Valuer in the
// stored context with their generated value, appends keyvals, and passes the
// result to the wrapped Logger.
func (l *context) Log(keyvals ...interface{}) error {
	kvs := append(l.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	if l.hasValuer {
		// If no keyvals were appended above then we must copy l.keyvals so
		// that future log events will reevaluate the stored Valuers.
		if len(keyvals) == 0 {
			kvs = append([]interface{}{}, l.keyvals...)
		}
		bindValues(kvs[:len(l.keyvals)])
	}
	return l.logger.Log(kvs...)
}
