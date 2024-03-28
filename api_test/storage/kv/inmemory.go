package kv

import (
	"errors"
	"sync"

	"github.com/spf13/cast"
)

type InMemory struct {
	inst sync.Map
}

func NewInMemoryInst() *InMemory {
	return &InMemory{}
}

func (iM *InMemory) Set(key string, value string, seconds int) error {
	iM.inst.Store(key, value)
	return nil
}

func (iM *InMemory) Get(key string) (string, error) {
	value, ok := iM.inst.Load(key)
	if !ok {
		return "", errors.New("not found")
	}

	return cast.ToString(value), nil
}

func (iM *InMemory) Delete(key string) error {
	iM.inst.Delete(key)
	return nil
}

func (iM *InMemory) List() (map[string]string, error) {
	pairs := make(map[string]string)

	iM.inst.Range(func(key, value interface{}) bool {
		pairs[cast.ToString(key)] = cast.ToString(value)
		return true
	})

	return pairs, nil
}
