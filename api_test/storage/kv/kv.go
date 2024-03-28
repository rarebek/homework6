package kv

type KV interface {
	Set(key string, value string, seconds int) error
	Get(key string) (string, error)
	Delete(key string) error
	List() (map[string]string, error)
}

var inst KV

func Init(store KV) {
	inst = store
}

func Set(key string, value string, seconds int) error {
	return inst.Set(key, value, seconds)
}

func Get(key string) (string, error) {
	return inst.Get(key)
}

func Delete(key string) error {
	return inst.Delete(key)
}

func List() (map[string]string, error) {
	return inst.List()
}
