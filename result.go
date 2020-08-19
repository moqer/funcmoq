package funcmoq

import (
	"errors"
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
	t                *testing.T
	values           []interface{}
	retrieveFinished func()
}

// Returning .
func (r *Store) Returning(args ...interface{}) {
	r.values = args
}

// Retrieve .
func (r *Store) Retrieve(args ...interface{}) {
	if err := r.retrieve(args...); err != nil {
		r.t.Fatal(err)
	}
	r.retrieveFinished()
}

// retrieve implementation of Retrieve method, used to expose the error for inspection
func (r *Store) retrieve(args ...interface{}) error {
	if len(r.values) != len(args) {
		return errors.New("Can't convert object")
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
