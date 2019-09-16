package log

type noopLogger struct{}

func (noopLogger) Log(...interface{}) error {
	return nil
}

func NoopLogger() Logger {
	return noopLogger{}
}
