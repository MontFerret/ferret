package encoding_test

type (
	nilWriterFunc  func([]byte) (int, error)
	nilReaderFunc  func([]byte) (int, error)
	nilWriterSlice []byte
	nilWriterMap   map[string]byte
)

func (w nilWriterFunc) Write(p []byte) (int, error) { return w(p) }
func (r nilReaderFunc) Read(p []byte) (int, error)  { return r(p) }
func (nilWriterSlice) Write([]byte) (int, error)    { panic("nil slice writer used") }
func (nilWriterMap) Write([]byte) (int, error)      { panic("nil map writer used") }
