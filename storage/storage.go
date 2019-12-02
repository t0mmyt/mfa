package storage

type Storage interface {
	Import([]byte) error
	Export() (string, error)
	Get(string) *Key
	Set(string, string) error
	List() ([]string, error)
	Close() error
}
