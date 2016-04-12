package utils

type FileStore interface {
	Put(key, fileid string) error
	Get(key string) ([]string, error)
	GetAll() ([]string, error)
	Hint(key string)
}
