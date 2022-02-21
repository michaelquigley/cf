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
	"reflect"
	"time"
)

type Instantiator func() interface{}
type Setter func(v interface{}, f reflect.Value, opt *Options) error
type FlexibleSetter func(v interface{}, opt *Options) (interface{}, error)
type NameConverter func(f reflect.StructField) string
type VariableResolver func(name string) (interface{}, bool)

type Options struct {
	Instantiators         map[reflect.Type]Instantiator
	Setters               map[reflect.Type]Setter
	FlexibleSetters       map[string]FlexibleSetter
	NameConverter         NameConverter
	VariableResolverChain []VariableResolver
}

func DefaultOptions() *Options {
	var td time.Duration
	opt := &Options{
		Setters: map[reflect.Type]Setter{
			reflect.TypeOf(0):          intSetter,
			reflect.TypeOf(int8(0)):    int8Setter,
			reflect.TypeOf(uint8(0)):   uint8Setter,
			reflect.TypeOf(int16(0)):   int16Setter,
			reflect.TypeOf(uint16(0)):  uint16Setter,
			reflect.TypeOf(int32(0)):   int32Setter,
			reflect.TypeOf(uint32(0)):  uint32Setter,
			reflect.TypeOf(int64(0)):   int64Setter,
			reflect.TypeOf(uint64(0)):  uint64Setter,
			reflect.TypeOf(float64(0)): float64Setter,
			reflect.TypeOf(true):       boolSetter,
			reflect.TypeOf(""):         stringSetter,
			reflect.TypeOf(td):         timeDurationSetter,
		},
		NameConverter: SnakeCaseNameConverter,
	}
	return opt
}

func (opt *Options) AddInstantiator(t reflect.Type, i Instantiator) *Options {
	if opt.Instantiators == nil {
		opt.Instantiators = make(map[reflect.Type]Instantiator)
	}
	opt.Instantiators[t] = i
	return opt
}

func (opt *Options) AddSetter(t reflect.Type, s Setter) *Options {
	if opt.Setters == nil {
		opt.Setters = make(map[reflect.Type]Setter)
	}
	opt.Setters[t] = s
	return opt
}

func (opt *Options) AddFlexibleSetter(typeName string, fs FlexibleSetter) *Options {
	if opt.FlexibleSetters == nil {
		opt.FlexibleSetters = make(map[string]FlexibleSetter)
	}
	opt.FlexibleSetters[typeName] = fs
	return opt
}

func (opt *Options) SetNameConverter(nc NameConverter) *Options {
	opt.NameConverter = nc
	return opt
}

func (opt *Options) AddVariableResolver(vr VariableResolver) *Options {
	opt.VariableResolverChain = append(opt.VariableResolverChain, vr)
	return opt
}

func (opt *Options) resolveVariable(vname string) (interface{}, bool) {
	for _, vr := range opt.VariableResolverChain {
		if vvalue, found := vr(vname); found {
			return vvalue, true
		}
	}
	return nil, false
}
