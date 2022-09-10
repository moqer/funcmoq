package funcmoq

import (
	"errors"
	"testing"
)

// Returner interface used for setting up function results
type Returner interface {
	Returning(args ...interface{})
}

// Retriever interface used for retrieving function results, all args need to be pointers
type Retriever interface {
	Retrieve(args ...interface{})
}

// Hasher interface used for hashing parameters
type Hasher interface {
	Hash(args ...interface{}) (uint64, error)
}

// New returns a initialized FuncMoq type
func New(t *testing.T) *FuncMoq {
	return &FuncMoq{
		t:       t,
		results: make(map[uint64]*Store),
		Hasher:  mitchellhash{},
	}
}

// FuncMoq boilerplate friendly object for mocking functions.
// The type acts as a dynamic hashmap, to register a key use With method to retrieve a previous set value use For method.
// The api should look like this:
// adding values: funcmoq.With(parm1, param2).Returning(val1, val2, val3)
// retrieving values: funcmoq.For(parm1, param2).Retrieve(&val1, &val2, &val3)
type FuncMoq struct {
	t         *testing.T
	results   map[uint64]*Store
	CallCount int
	Action    func()
	//by default FuncMoq uses github.com/mitchellh/hashstructure for hashing the parameters
	Hasher Hasher
}

var hashErrStr = "Can't create a hash for this set of parameters"

// For  method used to specific the parameters used when retrieving results.
func (m *FuncMoq) For(key ...interface{}) Retriever {
	result, err := m.get(key...)
	if err != nil {
		m.t.Fatal(hashErrStr, err)
	}
	result.retrieveFinished = func() {
		m.CallCount++
		if m.Action != nil {
			m.Action()
		}
	}
	return result
}

// get implementation of For method used to inspect the err
func (m *FuncMoq) get(key ...interface{}) (*Store, error) {
	hash, err := m.Hasher.Hash(key, nil)
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
	result, err := m.set(key...)
	if err != nil {
		m.t.Fatal(hashErrStr, err)
	}
	return result
}

// set implementation of With method used to inspect the err
func (m *FuncMoq) set(key ...interface{}) (Returner, error) {
	br := NewStore(m.t)
	hash, err := m.Hasher.Hash(key, nil)
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
