package dbi

type Entry struct {
	Value                []byte
	Err                  error
	BdbFirstOverflowPgNo uint32
}

type RpmDBInterface interface {
	Read() <-chan Entry
	Close() error
	GetPgSize() uint32
	GetLastPgNo() uint32
}
