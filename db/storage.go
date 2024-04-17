package chat_db

import "sync/atomic"

type Storage struct {
	Handler
}

type Handler interface {
	Set(key string, value Serializable, group ...string) error
	Get(key string, value Serializable, group ...string) (bool, error)
}
type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

var defaultStorage atomic.Pointer[Storage]

func Init(s *Storage) {
	defaultStorage.Store(s)
}

func Default() *Storage {
	return defaultStorage.Load()
}
