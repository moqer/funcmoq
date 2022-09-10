package funcmoq

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type tester interface {
	Helper()
	Errorf(str string, v ...interface{})
}

// NewStore .
func NewStore(t tester) *Store {
	return &Store{
		t: t,
	}
}

// Store .
type Store struct {
	t                tester
	values           []interface{}
	setLocation      string
	retrieveFinished func()
}

// Returns .
func (r *Store) Returns(args ...interface{}) {
	r.values = args
	_, file, line, ok := runtime.Caller(1)
	if ok {
		r.setLocation = fmt.Sprintf("%v:%v", file, line)
	}
}

// Retrieve .
func (r *Store) Retrieve(args ...interface{}) {
	r.t.Helper()
	if err := r.retrieve(args...); err != nil {
		r.t.Errorf("\n%s", err.Error())
		return
	}

	if r.retrieveFinished != nil {
		r.retrieveFinished()
	}
}

func getTypes(args ...interface{}) (string, error) {
	var sb strings.Builder
	if _, err := sb.WriteRune('('); err != nil {
		return "", err
	}

	for i, val := range args {
		if _, err := sb.WriteString(fmt.Sprintf("%T", val)); err != nil {
			return "", err
		}
		if i != len(args)-1 {
			if _, err := sb.WriteRune(','); err != nil {
				return "", err
			}
		}
	}
	if _, err := sb.WriteRune(')'); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func formatDiff(registered, retrieved []interface{}, setLocation, getLocation string) (string, error) {
	reg, err := getTypes(registered...)
	if err != nil {
		return "", err
	}
	ret, err := getTypes(retrieved...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Registered%s: %v\nRetrieving%s: %v", reg, setLocation, ret, getLocation), nil
}

// retrieve implementation of Retrieve method, used to expose the error for inspection
func (r *Store) retrieve(args ...interface{}) error {
	r.t.Helper()
	if len(r.values) != len(args) {
		_, file, line, _ := runtime.Caller(2)
		diff, err := formatDiff(r.values, args, r.setLocation, fmt.Sprintf("%v:%v", file, line))
		if err != nil {
			return err
		}

		return fmt.Errorf("Provided %d values while trying to retrieve %d\n"+diff, len(r.values), len(args))
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
