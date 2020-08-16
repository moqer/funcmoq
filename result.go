package funcmoq

import (
	"reflect"
	"testing"
)

// NewStore .
func NewStore(t *testing.T) *Store {
	return &Store{
		t: t,
	}
}

// Store .
type Store struct {
	t      *testing.T
	values []interface{}
}

// Returning .
func (r *Store) Returning(args ...interface{}) {
	r.values = args
}

// Retrieve .
func (r *Store) Retrieve(args ...interface{}) {
	if len(r.values) != len(args) {
		r.t.Fatal("Can't convert object")
	}

	for i := range r.values {
		v1 := reflect.ValueOf(r.values[i])
		v2 := reflect.ValueOf(args[i])
		if v2.Kind() != reflect.Ptr {
			r.t.Fatal("has to be a pointer")
		}
		ve := v2.Elem()
		if !ve.CanSet() {
			r.t.Fatal("cant set object")
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
			r.t.Fatal("cant assign object")
		}
	}
}
