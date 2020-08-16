package funcmoq

import (
	"errors"
	"reflect"
)

// Store .
type Store struct {
	values []interface{}
}

// Returning .
func (r *Store) Returning(args ...interface{}) {
	r.values = args
}

// Retrieve .
func (r *Store) Retrieve(args ...interface{}) error {
	if len(r.values) != len(args) {
		return errors.New("cant convert object")
	}

	for i := range r.values {
		v1 := reflect.ValueOf(r.values[i])
		v2 := reflect.ValueOf(args[i])
		if v2.Kind() != reflect.Ptr {
			return errors.New("has to be a pointer")
		}
		ve := v2.Elem()
		if !ve.CanSet() {
			return errors.New("cant set object")
		}
		if r.values[i] == nil {
			ve.Set(reflect.Zero(ve.Type()))
			continue
		}
		t := reflect.TypeOf(r.values[i])
		if t.AssignableTo(ve.Type()) {
			ve.Set(v1)
		} else if t.ConvertibleTo(ve.Type()) {
			ve.Set(v1.Convert(ve.Type()))
		} else {
			return errors.New("cant assign object")
		}

	}
	return nil
}
