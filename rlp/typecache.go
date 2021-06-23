package rlp

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var cache = newTypeCache()

type decoder func(*Decoder, reflect.Value) error

type typeinfo struct {
	decoder   decoder
	decodeErr error
}

func (i *typeinfo) generate(t reflect.Type) {
	i.decoder, i.decodeErr = createDecoder(t)
}

type typekey struct {
	reflect.Type
}

type typecache struct {
	cur  atomic.Value
	mu   sync.Mutex
	next map[typekey]*typeinfo
}

func (t *typecache) info(rt reflect.Type) *typeinfo {
	key := typekey{Type: rt}
	cur := t.cur.Load().(map[typekey]*typeinfo)
	if info := cur[key]; info != nil {
		return info
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// copy
	t.next = make(map[typekey]*typeinfo, len(cur)+1)
	for k, v := range cur {
		t.next[k] = v
	}

	info := new(typeinfo)
	t.next[key] = info
	info.generate(rt)

	// t.next -> t.cur
	t.cur.Store(t.next)
	t.next = nil

	return info
}

func newTypeCache() *typecache {
	c := new(typecache)
	c.cur.Store(make(map[typekey]*typeinfo))

	return c
}

func cachedDecoder(t reflect.Type) (decoder, error) {
	info := cache.info(t)
	return info.decoder, info.decodeErr
}
