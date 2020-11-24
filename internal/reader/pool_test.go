package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assertions := assert.New(t)
	pool := newPool()

	b := pool.getBuffer()
	assertions.NotNil(b)

	pool.returnBuffer(b)
}
