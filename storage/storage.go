package storage

import "sync"

type Store interface {
	Sub(key string) Store
	Get(key string) []byte
	Set(key string, value []byte)
}

