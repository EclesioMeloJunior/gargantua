package rlp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

var cache = newTypeCache()

type decoder func(*Decoder, reflect.Value) error
type writer func(reflect.Value, *Encoder) (int, error)

type typeinfo struct {
	decoder   decoder
	decodeErr error

	writer    writer
	writerErr error
}

func (i *typeinfo) generate(t reflect.Type) {
	i.decoder, i.decodeErr = createDecoder(t)
	i.writer, i.writerErr = createWriter(t)
}

type typekey struct {
	reflect.Type
}

type typecache struct {
	cur  atomic.Value
	mu   sync.Mutex
	next map[typekey]*typeinfo
}

// info will get from cache or generate a new typeinfo and then
// update the cache with the new one
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

func cachedWriter(t reflect.Type) (writer, error) {
	info := cache.info(t)
	return info.writer, info.writerErr
}

type tags struct {
	ignored bool
	tail    bool
}

type field struct {
	index int
	info  *typeinfo
	tags  tags
}

func getStructFields(t reflect.Type) ([]field, error) {

	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i); f.PkgPath == "" {
			tags, err := parseStructTag(t, f)
			if err != nil {
				return nil, err
			}

			if tags.ignored {
				continue
			}

		}
	}
}

func parseStructTag(typ reflect.Type, f reflect.StructField) (tags, error) {
	var tg tags
	for _, t := range strings.Split(f.Tag.Get("rlp"), ",") {
		switch t = strings.TrimSpace(t); t {
		case "":
		case "-":
			tg.ignored = true
		case "tail":
			tg.tail = true
			if !isLastPublicField(typ, f) {
				return tags{}, errors.New("rlp: field must be the last field to be tail")
			}
			if f.Type.Kind() != reflect.Slice {
				return tags{}, errors.New("rlp: tail field type is not a slice")
			}
		default:
			return tags{}, fmt.Errorf("rlp: unkonw struct tag %q on %v.%s", t, f.Type, f.Name)
		}
	}

	return tg, nil
}

func isLastPublicField(t reflect.Type, f reflect.StructField) bool {
	var field reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath == "" {
			field = t.Field(i)
		}
	}

	return reflect.DeepEqual(field, f)
}
