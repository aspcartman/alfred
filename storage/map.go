package storage

import "strings"

type inmem struct {
	keypath []string
	data    map[string]interface{}
}

func InMem() Store {
	return inmem{nil, map[string]interface{}{}}
}

func (m inmem) Sub(key string) Store {
	m.keypath = append(m.keypath, key)
	return m
}

func (m inmem) Get(key string) []byte {
	data, ok := m.data[m.traverse(key)].([]byte)
	if !ok {
		return nil
	}
	return data
}

func (m inmem) Set(key string, value []byte) {
	m.data[m.traverse(key)] = value
}

func (m inmem) traverse(path string) string {
	return strings.Join(append(m.keypath, path), "/")
}
