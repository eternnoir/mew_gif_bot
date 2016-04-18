package utils

type FileStore interface {
	Put(key, fileid string) error
	Get(key string) ([]string, error)
	GetAll() ([]string, error)
	Hint(key string)
	Reset() error
	GetStatus() (*Status, error)
}

type Status struct {
	TotalGifs   int64
	TotalQuerys int64
	TotalSend   int64
}
