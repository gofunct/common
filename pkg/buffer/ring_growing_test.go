package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrowth(t *testing.T) {
	t.Parallel()
	x := 10
	g := NewRingGrowing(1)
	for i := 0; i < x; i++ {
		assert.Equal(t, i, g.readable)
		g.WriteOne(i)
	}
	read := 0
	for g.readable > 0 {
		v, ok := g.ReadOne()
		assert.True(t, ok)
		assert.Equal(t, read, v)
		read++
	}
	assert.Equalf(t, x, read, "expected to have read %d items: %d", x, read)
	assert.Zerof(t, g.readable, "expected readable to be zero: %d", g.readable)
	assert.Equalf(t, 16, g.n, "expected N to be 16: %d", g.n)
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	g := NewRingGrowing(1)
	_, ok := g.ReadOne()
	assert.False(t, ok)
}
