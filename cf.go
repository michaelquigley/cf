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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	"strings"
)

func BindYaml(cf interface{}, path string, opt *Options) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "error reading yaml [%s]", path)
	}
	dataMap := make(map[string]interface{})
	if err := yaml.Unmarshal(data, dataMap); err != nil {
		return errors.Wrapf(err, "error parsing yaml [%s]", path)
	}
	return Bind(cf, dataMap, opt)
}

func Bind(cf interface{}, data map[string]interface{}, opt *Options) error {
	cfV := reflect.ValueOf(cf)
	if cfV.Kind() == reflect.Ptr {
		cfV = cfV.Elem()
	}
	if cfV.Kind() != reflect.Struct {
		return errors.Errorf("provided type [%s] is not a struct", cfV.Type())
	}
	for i := 0; i < cfV.NumField(); i++ {
		if cfV.Field(i).CanInterface() {
			fd := parseFieldData(cfV.Type().Field(i), opt)
			if !fd.skip {
				if v, found := data[fd.name]; found {
					if cfV.Field(i).CanSet() {
						nestedType := cfV.Type().Field(i).Type

						if handler, found := opt.Setters[nestedType]; found {
							// handler-based type
							if err := handler(v, cfV.Field(i)); err != nil {
								return errors.Wrapf(err, "field '%s'", fd.name)
							}
						} else if nestedType.Kind() == reflect.Struct || (nestedType.Kind() == reflect.Ptr && nestedType.Elem().Kind() == reflect.Struct) {
							// nested structure
							nested := instantiateAsPtr(nestedType, opt)
							if subData, ok := v.(map[string]interface{}); ok {
								err := Bind(nested, subData, opt)
								if err != nil {
									return errors.Wrapf(err, "field '%s'", fd.name)
								}
							} else {
								return errors.Errorf("invalid sub map for field '%s'", fd.name)
							}

							if nestedType.Kind() == reflect.Ptr {
								// by pointer
								cfV.Field(i).Set(reflect.ValueOf(nested))
							} else {
								// by value
								cfV.Field(i).Set(reflect.ValueOf(nested).Elem())
							}

						} else if nestedType.Kind() == reflect.Slice {
							sliceType := valueFromPtr(nestedType.Elem())
							if reflect.ValueOf(v).Kind() == reflect.Slice { // the data should also be a slice
								if sliceType.Kind() != reflect.Interface {
									for j := 0; j < reflect.ValueOf(v).Len(); j++ { // iterate over the available data
										elem := instantiateAsPtr(sliceType, opt)
										sliceV := reflect.ValueOf(v).Index(j).Interface()
										if handler, found := opt.Setters[sliceType]; found {
											if err := handler(sliceV, reflect.ValueOf(elem)); err != nil {
												return errors.Wrapf(err, "field '%s'", fd.name)
											}
										} else if sliceType.Kind() == reflect.Struct {
											if subData, ok := sliceV.(map[string]interface{}); ok {
												if err := Bind(elem, subData, opt); err != nil {
													return errors.Wrapf(err, "field '%s'", fd.name)
												}
											} else {
												return errors.Errorf("invalid sub map for field '%s'", fd.name)
											}
										}

										if nestedType.Elem().Kind() == reflect.Ptr {
											cfV.Field(i).Set(reflect.Append(cfV.Field(i), reflect.ValueOf(elem)))
										} else {
											cfV.Field(i).Set(reflect.Append(cfV.Field(i), reflect.ValueOf(elem).Elem()))
										}
									}

								} else {
									for j := 0; j < reflect.ValueOf(v).Len(); j++ { // iterate over the available data
										sliceV := reflect.ValueOf(v).Index(j).Interface()
										if subData, ok := sliceV.(map[string]interface{}); ok {
											if t, ok := subData["type"]; ok {
												if t, ok := t.(string); ok {
													if opt.FlexibleSetters != nil {
														if fs, found := opt.FlexibleSetters[t]; found {
															elem, err := fs(v, opt)
															if err != nil {
																return errors.Wrapf(err, "flexible setter for field '%v'", fd.name)
															}
															cfV.Field(i).Set(reflect.Append(cfV.Field(i), reflect.ValueOf(elem)))

														} else {
															return errors.Errorf("no flexible setter for field '%v'", fd.name)
														}
													} else {
														return errors.Errorf("no flexible setters found for field '%v'", fd.name)
													}
												} else {
													return errors.Errorf("'type' data not string value for field '%v'", fd.name)
												}
											} else {
												return errors.Errorf("'type' data not found string value for field '%v'", fd.name)
											}
										} else {
											return errors.Errorf("not data map for field '%v'", fd.name)
										}
									}
								}

							} else {
								return errors.Errorf("invalid arr map for field '%s' (%v)", fd.name, reflect.TypeOf(v))
							}

						} else if nestedType.Kind() == reflect.Interface {
							if subData, ok := v.(map[string]interface{}); ok {
								if t, ok := subData["type"]; ok {
									if t, ok := t.(string); ok {
										if opt.FlexibleSetters != nil {
											if fs, found := opt.FlexibleSetters[t]; found {
												value, err := fs(v, opt)
												if err != nil {
													return errors.Wrapf(err, "flexible setter error for field '%s'", fd.name)
												}
												cfV.Field(i).Set(reflect.ValueOf(value))
											} else {
												return errors.Errorf("no flexible setter for field '%s'", fd.name)
											}
										} else {
											return errors.Errorf("no flexible setters found for field '%s'", fd.name)
										}
									} else {
										return errors.Errorf("'type' data not string value for field '%s'", fd.name)
									}
								} else {
									return errors.Errorf("no 'type' data found for field '%s'", fd.name)
								}
							} else {
								return errors.Errorf("interface{} field '%s' requires sub data map", fd.name)
							}

						} else {
							return errors.Errorf("no setter for field '%s' of type '%s/%v'", fd.name, nestedType, nestedType.Kind())
						}
					} else {
						return errors.Errorf("non-settable field '%s' of type '%s'", fd.name, cfV.Field(i).Type())
					}
				} else {
					if fd.required {
						return errors.Errorf("no data found for required field '%s'", fd.name)
					}
				}
			}
		}
	}
	// execute wirings for type
	if opt.Wirings != nil {
		if wirings, found := opt.Wirings[valueFromPtr(reflect.TypeOf(cf))]; found {
			for _, wiring := range wirings {
				if err := wiring(cf); err != nil {
					return errors.Wrapf(err, "error wiring [%s]", cfV.Elem().Type().Name())
				}
			}
		}
	}
	return nil
}

type fieldData struct {
	name     string
	skip     bool
	required bool
}

func parseFieldData(v reflect.StructField, opt *Options) fieldData {
	fd := fieldData{name: opt.NameConverter(v), skip: false, required: false}
	data := v.Tag.Get("cf")
	if data != "" {
		tokens := strings.Split(data, ",")
		for _, token := range tokens {
			if token == "+required" {
				fd.required = true
			} else if token == "+skip" {
				fd.skip = true
			} else {
				fd.name = token
			}
		}
	}
	return fd
}

func valueFromPtr(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func instantiateAsPtr(t reflect.Type, opt *Options) interface{} {
	var it = t
	if it.Kind() == reflect.Ptr {
		it = it.Elem()
	}
	if i, found := opt.Instantiators[it]; found {
		return i()
	} else {
		return reflect.New(it).Interface()
	}
}
