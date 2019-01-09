package utils

var closers []func()

// Close closes cli utilities.
func Close() {
	for _, f := range closers {
		f()
	}
}

func AddCloseFunc(f func()) {
	closers = append(closers, f)
}
