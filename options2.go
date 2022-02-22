/*
 * Cf
 * Extensible, idiomatic, opinionated wiring for golang applications.
 *
 * Copyright (c) Michael Quigley
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the “Software”), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cf

import "reflect"

// Instantiator2 allows the framework client to insert bespoke functionality when instantiating new instances.
//
type Instantiator2 func() interface{}

// Setter2 is used to assign values to struct fields when binding.
//
type Setter2 func(v interface{}, f reflect.Value, opt *Options2) error

// Options2 is used to specify the configuration for the framework client.
//
type Options2 struct {
	Instantiators map[reflect.Type]Instantiator2
}
