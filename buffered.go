package m3

import "io"

// ---
// ### **NewBufferedWriter**

// NewBufferedWriter creates a buffered writer with the given buffer
// size. A `WriterMiddleware` is returned which wraps the log's
// internal `io.WriteCloser` in a `bufio.Writer`.
func NewBufferedWriterMiddleware(size int) WriterMiddleware {
    return func(writer io.WriteCloser) io.WriteCloser {
        return bufferedWriteCloser{size, 0, make([]byte, 0, size), writer}
    }
}

// > The buffered middleware is implemented as a `bufferedWriteCloser`
// which requires the newly created `bufio.Writer` as well as the
// log's internal `io.WriteCloser`.
type bufferedWriteCloser struct {
    size   int
    offset int
    buffer []byte
    parent io.WriteCloser
}

// #### Write

// Write writes the data into the buffer.
func (b bufferedWriteCloser) Write(data []byte) (n int, err error) {
    if len(data) > b.size {
        return b.parent.Write(data)
    }

    if len(data)+b.offset > b.size {
        b.parent.Write(b.buffer[:b.offset])

        n = b.offset
        b.offset = 0
        return
    }

    copy(b.buffer[:len(data)], data)
    b.offset += len(data)
    return len(data), nil
}

// #### Close

// Close flushes the buffer and then closes the parent `io.WriteCloser.`
func (b bufferedWriteCloser) Close() error {
    if b.offset > 0 {
        b.parent.Write(b.buffer[:b.offset])
        b.offset = 0
    }
    return b.parent.Close()
}
