package m3

import "io"

type ReaderMiddleware func(io.ReadCloser) io.ReadCloser
type WriterMiddleware func(io.WriteCloser) io.WriteCloser

type Reader interface {
    io.ReadCloser

    // Use provides middleware for the internal `io.ReadCloser`
    Use(...ReaderMiddleware)
}

func NewReader(r io.ReadCloser) Reader {
    return &reader{r}
}

type reader struct {
    rc io.ReadCloser
}

func (r *reader) Read(data []byte) (n int, err error) {
    return r.rc.Read(data)
}

func (r *reader) Close() (err error) {
    return r.rc.Close()
}
func (r *reader) Use(readers ...ReaderMiddleware) {

    // Iterate over all the `ReaderMiddleware` replacing the internal
    // reader with the new one.
    for _, reader := range readers {
        r.rc = reader(r.rc)
    }
}

type Writer interface {
    io.WriteCloser

    // Use provides middleware for the internal `io.WriteCloser`
    Use(...WriterMiddleware)
}

func NewWriter(w io.WriteCloser) Writer {
    return &writer{w}
}

type writer struct {
    wc io.WriteCloser
}

func (w *writer) Write(data []byte) (n int, err error) {
    return w.wc.Write(data)
}

func (w *writer) Close() (err error) {
    return w.wc.Close()
}

func (w *writer) Use(writers ...WriterMiddleware) {

    // Iterate over all the `WriterMiddleware` replacing the internal
    // reader with the new one.
    for _, writer := range writers {
        w.wc = writer(w.wc)
    }
}
