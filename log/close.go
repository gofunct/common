package log

var closers []func()

// Close closes cli utilities.
func (s *Service) Close() {
	s.DebugC("closing....")
	for _, f := range closers {
		f()
	}
}

func (s *Service) AddCloseFunc(f func()) {
	closers = append(closers, f)
}
