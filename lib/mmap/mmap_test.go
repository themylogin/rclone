package mmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllocFree(t *testing.T) {
	const Size = 4096

	b := Alloc(Size)
	assert.Equal(t, Size, len(b))

	// check we can write to all the memory
	for i := range b {
		b[i] = byte(i)
	}

	// Now free the memory
	Free(b)
}
