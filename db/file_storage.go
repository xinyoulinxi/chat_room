package chat_db

import (
	"os"
	"path"
	"sync"
	"web_server/utils"
)

type fileStorage struct {
	dir    string
	keyMap sync.Map
}

func NewFileStorage(dir string) (*Storage, error) {
	handler := &fileStorage{
		dir:    dir,
		keyMap: sync.Map{},
	}
	if err := utils.EnsureDir(dir); err != nil {
		return nil, err
	}
	return &Storage{handler}, nil
}

func (f *fileStorage) getKey(key string, group ...string) string {
	elem := make([]string, 1, len(group)+2)
	elem[0] = f.dir
	for _, g := range group {
		elem = append(elem, g)
	}
	elem = append(elem, key+".json")
	return path.Join(elem...)
}

func (f *fileStorage) Set(key string, value Serializable, group ...string) error {
	key = f.getKey(key, group...)
	_, ok := f.keyMap.Load(key)
	if !ok {
		if err := utils.EnsureFileExist(key); err != nil {
			return err
		}
		f.keyMap.Store(key, struct{}{})
	}
	bytes, err := value.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(key, bytes, 0644)
}

func (f *fileStorage) Get(key string, value Serializable, group ...string) (bool, error) {
	key = f.getKey(key, group...)
	if _, ok := f.keyMap.Load(key); !ok {
		if err := utils.EnsureFileExist(key); err != nil {
			return false, err
		}
		f.keyMap.Store(key, struct{}{})
	}
	bytes, err := os.ReadFile(key)
	if err != nil {
		return false, err
	}
	if len(bytes) == 0 {
		return false, err
	}
	return true, value.Deserialize(bytes)
}
