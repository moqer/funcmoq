package funcmoq

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

type testingT interface {
	Helper()
	Errorf(str string, v ...interface{})
	Fatalf(str string, v ...interface{})
}

// hasher interface used for hashing parameters
type hasher interface {
	Hash(args ...interface{}) (uint64, error)
}

// New returns a initialized FuncMoq type
func New(t *testing.T) *FuncMoq {
	return NewWithCallback(t, nil)
}

// New returns a initialized FuncMoq type
func NewWithCallback(t testingT, callback func(key ...interface{})) *FuncMoq {
	return &FuncMoq{
		t:            t,
		results:      make(map[uint64]*Store),
		hasher:       mitchellhash{},
		Action:       callback,
		setLocations: make([]string, 0),
	}
}

// FuncMoq boilerplate friendly object for mocking functions.
// The type acts as a dynamic hashmap, to register a key use With method to retrieve a previous set value use For method.
// The api should look like this:
// adding values: funcmoq.With(parm1, param2).Returning(val1, val2, val3)
// retrieving values: funcmoq.For(parm1, param2).Retrieve(&val1, &val2, &val3)
type FuncMoq struct {
	t            testingT
	results      map[uint64]*Store
	CallCount    int
	Action       func(key ...interface{})
	setLocations []string
	hasher       hasher
}

// For  method used to specific the parameters used when retrieving results.
func (m *FuncMoq) For(key ...interface{}) Retriever {
	m.t.Helper()
	result, err := m.get(key...)
	if err != nil {
		m.t.Fatalf("No return is set for the parameters: %v\ninner error: %v\nSetup locations:\n%s",
			key, err, strings.Join(m.setLocations, " \n"))
	}
	result.retrieveFinished = func() {
		m.CallCount++
		if m.Action != nil {
			m.Action(key)
		}
	}
	return result
}

// get implementation of For method used to inspect the err
func (m *FuncMoq) get(key ...interface{}) (*Store, error) {
	hash, err := m.hasher.Hash(key, nil)
	if err != nil {
		return nil, err
	}
	result, exists := m.results[hash]
	if !exists {
		return nil, errors.New("This set of parameters weren't registered")
	}
	return result, nil
}

// With method used to specific the parameters used when storing results.
func (m *FuncMoq) With(key ...interface{}) Returner {
	_, file, line, _ := runtime.Caller(1)
	m.setLocations = append(m.setLocations, fmt.Sprintf("%s:%d\n\tparams:%v", file, line, key))
	result, err := m.set(key...)
	if err != nil {
		m.t.Fatalf("\nCan't create a hash for this set of parameters\n inner error: %v", err)
	}
	return result
}

// set implementation of With method used to inspect the err
func (m *FuncMoq) set(key ...interface{}) (Returner, error) {
	br := NewStore(m.t)
	hash, err := m.hasher.Hash(key, nil)
	if err != nil {
		return nil, err
	}
	m.results[hash] = br
	return br, nil
}

// Called true if the results were retrieved at least once
func (m *FuncMoq) Called() bool {
	return m.CallCount > 0
}
