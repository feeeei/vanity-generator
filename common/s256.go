package common

import (
	"crypto/sha256"
	"hash"
	"sync"
)

var s256pool *sync.Pool

func init() {
	s256pool = &sync.Pool{
		New: func() any {
			return sha256.New()
		},
	}
}

func S256(s []byte) []byte {
	h := s256pool.Get().(hash.Hash)
	h.Write(s)
	bs := h.Sum(nil)
	h.Reset()
	s256pool.Put(h)
	return bs
}
