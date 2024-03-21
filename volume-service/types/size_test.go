package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSize(t *testing.T) {
	t.Run("test Size", testSize)
}

func testSize(t *testing.T) {
	size := SizeStringToNum("1024")
	assert.Equal(t, 1024, size)

	size = SizeStringToNum("1024K")
	assert.Equal(t, 1024*1024, size)

	size = SizeStringToNum("1024kb")
	assert.Equal(t, 1024*1024, size)

	size = SizeStringToNum("1024Mb")
	assert.Equal(t, 1024*1024*1024, size)

	size = SizeStringToNum("1024gb")
	assert.Equal(t, 1024*1024*1024*1024, size)
}
