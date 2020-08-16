package funcmoq

import (
	"errors"
	"log"
	"reflect"
)

// Result .
type Result struct {
	err    error
	values []interface{}
}

// Returning .
func (r *Result) Returning(err error, args ...interface{}) {
	if err != nil {
		r.err = err
	}
	r.values = args
}

// Retrieve .
func (r *Result) Retrieve(err error, args ...interface{}) error {
	err = r.err
	if len(r.values) != len(args) {
		return errors.New("cant convert object")
	}

	for i := range r.values {
		v := reflect.ValueOf(r.values[i])
		isNil := r.values[i] == nil || v.IsNil()
		t := reflect.TypeOf(r.values[i])
		x := v.Kind()
		v1 := reflect.ValueOf(args[i]).Elem()
		t1 := reflect.TypeOf(args[i])
		x1 := v1.Kind()
		log.Println(v1.CanSet(), x, t1, x1)
		if !v1.CanSet() {
			return errors.New("cant set object")
		}
		if isNil {
			v1.Set(reflect.Zero(v1.Type()))
		} else if t.AssignableTo(v1.Type()) {
			v1.Set(v)
		} else if t.ConvertibleTo(v1.Type()) {
			v1.Set(v.Convert(v1.Type()))
		} else {
			return errors.New("cant assign object")
		}

	}
	return nil
}
