package funcmoq

import "github.com/mitchellh/hashstructure"

type mitchellhash struct {
}

func (m mitchellhash) Hash(obj ...interface{}) (uint64, error) {
	return hashstructure.Hash(obj, nil)
}
