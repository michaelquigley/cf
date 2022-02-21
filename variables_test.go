package cf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VariableReference(t *testing.T) {
	v, found := variableReference("${oh.wow}")
	assert.True(t, found)
	assert.Equal(t, "oh.wow", v)

	v, found = variableReference(int16(1))
	assert.False(t, found)
}

func Test_Variables_Multiple(t *testing.T) {
	type simple struct {
		Name string
	}
	obj := &simple{}
	data := map[string]interface{}{
		"name": "name:${a}/${b}/${c}",
	}
	opt := DefaultOptions()
	opt.AddVariableResolver(func(name string) (interface{}, bool) {
		switch name {
		case "a":
			return "AAA", true
		case "b":
			return "BBB", true
		case "c":
			return "CCC", true
		default:
			return "", false
		}
	})
	err := Bind(obj, data, opt)
	assert.NoError(t, err)
	assert.Equal(t, "name:AAA/BBB/CCC", obj.Name)
}
