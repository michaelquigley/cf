/*
   Copyright NetFoundry, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package cf

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	basic := &struct {
		StringValue string
	}{}

	var data = map[string]interface{}{
		"string_value": "oh, wow!",
	}

	err := Bind(basic, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "oh, wow!", basic.StringValue)
}

func TestRenaming(t *testing.T) {
	renamed := &struct {
		SomeInt int `cf:"some_int_,+required"`
	}{}

	var data = map[string]interface{}{
		"some_int_": 46,
	}

	err := Bind(renamed, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, 46, renamed.SomeInt)
}

func TestStringArray(t *testing.T) {
	withArray := &struct {
		StringArray []string
	}{}

	var data = map[string]interface{}{
		"string_array": []string{"one", "two", "three"},
	}

	err := Bind(withArray, data, DefaultOptions())
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"one", "two", "three"}, withArray.StringArray)
}

func TestIntArray(t *testing.T) {
	withArray := &struct {
		IntArray []int
	}{}

	var data = map[string]interface{}{
		"int_array": []int{1, 2, 3, 4, 5, 6},
	}

	err := Bind(withArray, data, DefaultOptions())
	assert.Nil(t, err)
}

func TestRequired(t *testing.T) {
	required := &struct {
		Required int `cf:"+required"`
	}{}

	data := make(map[string]interface{})

	err := Bind(required, data, DefaultOptions())
	assert.NotNil(t, err)
}

type nestedType struct {
	Name  string
	Count int
}

func newNestedType() *nestedType {
	return &nestedType{Name: "oh, wow!", Count: 33} // defaults
}

func TestNestedPtr(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions().AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 33, root.Nested.Count)
}

func TestNestedValue(t *testing.T) {
	root := &struct {
		Id     string
		Nested nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions().AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 33, root.Nested.Count)
}

func TestNestedWithTypeWiring(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	opt := DefaultOptions()
	opt.AddInstantiator(reflect.TypeOf(nestedType{}), func() interface{} { return newNestedType() })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
}

func TestNestedWithDefaultInstantiator(t *testing.T) {
	root := &struct {
		Id     string
		Nested *nestedType
	}{}

	var data = map[string]interface{}{
		"id": "TestNested",
		"nested": map[string]interface{}{
			"name": "Different",
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "TestNested", root.Id)
	assert.NotNil(t, root.Nested)
	assert.Equal(t, "Different", root.Nested.Name)
	assert.Equal(t, 0, root.Nested.Count) // type wiring
}

func TestStructTypeArray(t *testing.T) {
	root := &struct {
		Id      string
		Nesteds []*nestedType
	}{}

	var data = map[string]interface{}{
		"id": "StructTypeArray",
		"nesteds": []map[string]interface{}{
			{"name": "a"},
			{"name": "b"},
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "StructTypeArray", root.Id)
	assert.Equal(t, 2, len(root.Nesteds))
	assert.Equal(t, "a", root.Nesteds[0].Name)
	assert.Equal(t, "b", root.Nesteds[1].Name)
}

func TestAnonymousStruct(t *testing.T) {
	root := &struct {
		Id     string
		Nested struct {
			Name string
		}
	}{}

	var data = map[string]interface{}{
		"id": "AnonymousStruct",
		"nested": map[string]interface{}{
			"name": "oh, wow!",
		},
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "AnonymousStruct", root.Id)
	assert.Equal(t, "oh, wow!", root.Nested.Name)
}

func TestUint32(t *testing.T) {
	root := &struct {
		Id    string
		Count uint32
	}{}

	var data = map[string]interface{}{
		"id":    "uint32",
		"count": 33,
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, "uint32", root.Id)
	assert.Equal(t, uint32(33), root.Count)
}

func TestTimeDuration(t *testing.T) {
	root := &struct {
		Duration time.Duration
	}{}

	var data = map[string]interface{}{
		"duration": "30s",
	}

	err := Bind(root, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(30)*time.Second, root.Duration)
}

type flexibleType struct {
	value string
}

func TestFlexibleSetter(t *testing.T) {
	root := &struct {
		Id       string
		Flexible interface{}
	}{}

	var data = map[string]interface{}{
		"id": "nested",
		"flexible": map[string]interface{}{
			"type":  "a",
			"value": "a value",
		},
	}

	opt := DefaultOptions()
	opt = opt.AddFlexibleSetter("a", func(v interface{}, opt *Options) (interface{}, error) { return &flexibleType{"a value"}, nil })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.NotNil(t, root.Flexible)
	assert.Equal(t, reflect.TypeOf(root.Flexible), reflect.TypeOf(&flexibleType{}))
	assert.Equal(t, "a value", root.Flexible.(*flexibleType).value)
}

func TestFlexibleSetterArray(t *testing.T) {
	root := &struct {
		Id            string
		FlexibleArray []interface{}
	}{}

	var data = map[string]interface{}{
		"id": "flexible_setter_array",
		"flexible_array": []map[string]interface{}{
			{
				"type": "a",
			},
			{
				"type": "b",
			},
		},
	}

	opt := DefaultOptions()
	opt = opt.AddFlexibleSetter("a", func(v interface{}, opt *Options) (interface{}, error) { return &flexibleType{"a"}, nil })
	opt = opt.AddFlexibleSetter("b", func(v interface{}, opt *Options) (interface{}, error) { return flexibleType{"b"}, nil })

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.NotNil(t, root.FlexibleArray)
	assert.Equal(t, 2, len(root.FlexibleArray))
	assert.Equal(t, "a", root.FlexibleArray[0].(*flexibleType).value)
	assert.Equal(t, "b", root.FlexibleArray[1].(flexibleType).value)
}

func TestVariableResolver(t *testing.T) {
	root := &struct {
		Id string
	}{}

	data := map[string]interface{}{
		"id": "${id}",
	}

	opt := DefaultOptions()
	opt.AddVariableResolver(func(vname string) (interface{}, bool) {
		if vname == "id" {
			return "oh.wow", true
		}
		return nil, false
	})

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "oh.wow", root.Id)

	data = map[string]interface{}{
		"id": "${missing}",
	}

	err = Bind(root, data, opt)
	assert.NotNil(t, err)

	opt.AddVariableResolver(func(vname string) (interface{}, bool) {
		if vname == "missing" {
			return "found", true
		}
		return nil, false
	})

	err = Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "found", root.Id)

	data = map[string]interface{}{
		"id": "regular",
	}

	err = Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "regular", root.Id)
}

func TestInlineVariables(t *testing.T) {
	root := &struct {
		Id string
	}{}

	data := map[string]interface{}{
		"id": "/hello/world/${a}/and/${b}",
	}

	opt := DefaultOptions()
	opt.AddVariableResolver(func(vname string) (interface{}, bool) {
		switch vname {
		case "a":
			return "oh", true
		case "b":
			return "wow", true
		}
		return nil, false
	})

	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, "/hello/world/oh/and/wow", root.Id)
}

func TestCustomEnum(t *testing.T) {
	type CustomEnum uint8
	const (
		Val1 CustomEnum = iota
		Val2
	)
	type SubStruct struct {
		Val CustomEnum
	}

	root := &struct {
		Subs []*SubStruct
	}{}

	data := map[string]interface{}{
		"subs": []interface{}{
			map[string]interface{}{
				"val": "Val1",
			},
			map[string]interface{}{
				"val": "Val2",
			},
		},
	}

	opt := DefaultOptions()
	opt.AddSetter(reflect.TypeOf(*(new(CustomEnum))), func(v interface{}, f reflect.Value, opt *Options) error {
		if vt, ok := v.(string); ok {
			var vv CustomEnum
			switch vt {
			case "Val1":
				vv = Val1
			case "Val2":
				vv = Val2
			default:
				return errors.Errorf("invalid custom enum '%v'", vt)
			}
			if f.Kind() == reflect.Ptr {
				f.Elem().SetUint(uint64(vv))
			} else {
				f.SetUint(uint64(vv))
			}
			return nil
		}
		return errors.Errorf("got [%s], expected [%s]", reflect.TypeOf(v), f.Type())
	})
	err := Bind(root, data, opt)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(root.Subs))
	assert.Equal(t, Val1, root.Subs[0].Val)
	assert.Equal(t, Val2, root.Subs[1].Val)
}

func TestFloatWithIntData(t *testing.T) {
	basic := &struct {
		FloatValue float64
	}{}

	var data = map[string]interface{}{
		"float_value": 56,
	}

	err := Bind(basic, data, DefaultOptions())
	assert.Nil(t, err)
	assert.Equal(t, float64(56.0), basic.FloatValue)
}
