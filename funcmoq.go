package funcmoq

import (
	"github.com/mitchellh/hashstructure"
)

// NewResRegistry .
func NewResRegistry() *ResRegistry {
	return &ResRegistry{
		results: make(map[uint64]*Result),
	}
}

// ResRegistry .
type ResRegistry struct {
	results map[uint64]*Result
}

// Returner .
type Returner interface {
	Returning(args ...interface{})
}

// Retriever .
type Retriever interface {
	Retrieve(args ...interface{}) error
}

// For .
func (h ResRegistry) For(key ...interface{}) Retriever {
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		panic(err) //testcode i think it's ok
	}
	result, exists := h.results[hash]
	if !exists {
		// return &Result{
		// 	err: errors.New("This key wasn't registered"),
		// }
		//todo panic
	}
	return result
}

// With .
func (h ResRegistry) With(key ...interface{}) Returner {
	var br Result
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		panic(err) //testcode i think it's ok
	}
	h.results[hash] = &br
	return &br
}
