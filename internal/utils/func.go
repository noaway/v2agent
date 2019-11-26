package utils

import (
	"reflect"
)

func NewFunction(fn interface{}) *Function {
	return &Function{
		fnType:  reflect.TypeOf(fn),
		fnValue: reflect.ValueOf(fn),
	}
}

type Function struct {
	fnType  reflect.Type
	fnValue reflect.Value
}

// An empty value is assigned when the missing parameter calls a function
func (f *Function) args(args ...interface{}) []reflect.Value {
	injmap := make(map[reflect.Type]reflect.Value)
	for i := range args {
		injmap[reflect.TypeOf(args[i])] = reflect.ValueOf(args[i])
	}
	count := f.fnType.NumIn()
	inValues := make([]reflect.Value, count)
loop:
	for i := 0; i < count; i++ {
		at := f.fnType.In(i)
		val, ok := injmap[at]
		if ok {
			inValues[i] = val
			continue loop
		}

		if at.Kind() == reflect.Interface {
			for k, v := range injmap {
				if k.Implements(at) {
					inValues[i] = v
					continue loop
				}
			}
		}

		inValues[i] = reflect.Zero(f.fnType.In(i))
	}
	return inValues
}

func (f *Function) IsFunc() bool {
	return f.fnType.Kind() == reflect.Func
}

func (f *Function) Invoke(fnArgs ...interface{}) []reflect.Value {
	if !f.IsFunc() {
		return []reflect.Value{}
	}
	return f.fnValue.Call(f.args(fnArgs...))
}

func (f *Function) GetType() reflect.Type {
	return f.fnType
}

func (f *Function) GetValue() reflect.Value {
	return f.fnValue
}
