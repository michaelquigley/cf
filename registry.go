/*
 * Cf
 * Extensible, idiomatic, opinionated wiring for golang applications.
 * Copyright (c) Michael Quigley
 */

package cf

func Register(handle string, v interface{}) {
	if registry == nil {
		registry = make(map[string]interface{})
	}
	registry[handle] = v
}

func Wire(handle string) interface{} {
	if v, found := registry[handle]; found {
		return v
	}
	return nil
}

var registry map[string]interface{}
