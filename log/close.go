package log

var closers []func()

// Close closes cli utilities.
func (l *Logger) Close() {
	l.Print("closing....")
	for _, f := range closers {
		f()
	}
}

func (l *Logger) AddCloseFunc(f func()) {
	closers = append(closers, f)
}
