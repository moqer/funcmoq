package funcmoq

import (
	"testing"

	"github.com/mitchellh/hashstructure"
)

// NewFuncMoq .
func NewFuncMoq(t *testing.T) *FuncMoq {
	return &FuncMoq{
		t:       t,
		results: make(map[uint64]*Store),
	}
}

// FuncMoq .
type FuncMoq struct {
	t       *testing.T
	results map[uint64]*Store
}

// Returner .
type Returner interface {
	Returning(args ...interface{})
}

// Retriever .
type Retriever interface {
	Retrieve(args ...interface{})
}

// For .
func (m FuncMoq) For(key ...interface{}) Retriever {
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		m.t.Fatal("Can't create hash for keys", err)
	}
	result, exists := m.results[hash]
	if !exists {
		m.t.Fatal("This key wasn't registered")
	}
	return result
}

// With .
func (m FuncMoq) With(key ...interface{}) Returner {
	br := NewStore(m.t)
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		m.t.Fatal("Can't create hash for keys", err)
	}
	m.results[hash] = br
	return br
}
