package repo

type RedisStorageI interface {
	Set(string, string) error
	Get(string) (interface{}, error)
	SetWithTTL(string, string, int) error
}
