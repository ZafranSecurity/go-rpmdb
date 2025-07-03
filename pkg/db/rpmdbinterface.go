package dbi

type Entry struct {
	Value                []byte
	Err                  error
	BdbFirstOverflowPgNo uint32
}

type RpmDBInterface interface {
	Read() <-chan Entry
	Close() error
}
