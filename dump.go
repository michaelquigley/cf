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
	"reflect"
)

func Dump(cf interface{}, opt *Options) string {
	return dump(reflect.ValueOf(cf), -1, opt)
}

func dump(v reflect.Value, indent int, opt *Options) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		return dumpStruct(v, indent+1, opt)
	case reflect.Slice:
		return dumpSlice(v, indent+1, opt)
	default:
		return dumpValue(v)
	}
}

func dumpStruct(v reflect.Value, indent int, opt *Options) string {
	format := fmt.Sprintf("%%-%ds", maxFieldNameLength(v, opt))
	out := "{\n"
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			fd := parseFieldData(v.Type().Field(i), opt)
			out += nTabs(indent+1) + fmt.Sprintf(format, fd.name) + " = " + dump(v.Field(i), indent, opt) + "\n"
		}
	}
	out += nTabs(indent) + "}"
	return out
}

func dumpSlice(v reflect.Value, indent int, opt *Options) string {
	out := "[\n"
	for i := 0; i < v.Len(); i++ {
		out += nTabs(indent+1) + dump(v.Index(i), indent, opt) + "\n"
	}
	out += nTabs(indent) + "]"
	return out
}

func dumpValue(v reflect.Value) string {
	if v.Kind() == reflect.String {
		return fmt.Sprintf("\"%v\"", v.Interface())
	} else {
		if v.IsValid() {
			return fmt.Sprintf("%v", v.Interface())
		} else {
			return "<nil>"
		}
	}
}

func maxFieldNameLength(cfV reflect.Value, opt *Options) int {
	maxFieldNameLength := 0
	for i := 0; i < cfV.NumField(); i++ {
		fd := parseFieldData(cfV.Type().Field(i), opt)
		length := len(fd.name)
		if length > maxFieldNameLength {
			maxFieldNameLength = length
		}
	}
	return maxFieldNameLength
}

func nTabs(n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += "\t"
	}
	return out
}
