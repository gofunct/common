package runtime

import (
	"net"
)

// Serve provides an interface for starting and stopping the server.
type Serve interface {
	Serve(l net.Listener) error
	Shutdown()
}
