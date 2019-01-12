package render

import "github.com/google/wire"

var Set = wire.NewSet(
	NewConfig,
	NewRenderer,
)
