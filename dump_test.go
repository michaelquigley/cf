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
	"fmt"
	"testing"
)

type dumpCfTiny struct {
	Id string
}

type dumpCfNested struct {
	Name string
	Tiny dumpCfTiny
}

type dumpCfNestedPtr struct {
	Name string
	Tiny *dumpCfTiny
}

type dumpCfNestedArray struct {
	Name   string
	Tinies []*dumpCfTiny
}

func TestDumpStruct(t *testing.T) {
	cf := &dumpCfTiny{"testing"}
	fmt.Println(Dump(cf, DefaultOptions()))

	cf2 := &dumpCfNested{Name: "yuu"}
	cf2.Tiny.Id = "yuu_id"
	fmt.Println(Dump(cf2, DefaultOptions()))
}

func TestDumpNilStruct(t *testing.T) {
	cf := &dumpCfNestedPtr{Name: "name"}
	fmt.Println(Dump(cf, DefaultOptions()))
}

func TestDumpNestedArray(t *testing.T) {
	cf := &dumpCfNestedArray{Name: "Testing"}
	cf.Tinies = append(cf.Tinies, &dumpCfTiny{"a"})
	cf.Tinies = append(cf.Tinies, &dumpCfTiny{"b"})
	fmt.Println(Dump(cf, DefaultOptions()))
}
