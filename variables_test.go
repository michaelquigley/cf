package cf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVariableReference(t *testing.T) {
	v, found := variableReference("${oh.wow}")
	assert.True(t, found)
	assert.Equal(t, "oh.wow", v)

	v, found = variableReference(int16(1))
	assert.False(t, found)
}
