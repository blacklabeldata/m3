package m3

import (
    "crypto/rand"
    "hash/crc64"
    "io"
    "testing"

    "github.com/stretchr/testify/assert"
)

type CRC64WriterMiddleware struct {
    CRC64  uint64
    Closed bool
    table  *crc64.Table
    parent io.WriteCloser
}

func (c *CRC64WriterMiddleware) Write(data []byte) (int, error) {
    c.CRC64 = crc64.Update(c.CRC64, c.table, data)
    return c.parent.Write(data)
}

func (c *CRC64WriterMiddleware) Close() error {
    c.Closed = true
    return c.parent.Close()
}

type CountingWriteCloser struct {
    Count  int
    Closed bool
}

func (c *CountingWriteCloser) Write(data []byte) (int, error) {
    c.Count += len(data)
    return len(data), nil
}

func (c *CountingWriteCloser) Close() error {
    c.Closed = true
    return nil
}

func TestNewWriter(t *testing.T) {

    // create test data
    buffer := make([]byte, 1024)
    n, err := rand.Read(buffer)
    assert.Nil(t, err)
    assert.Equal(t, n, len(buffer))

    // test utils
    countingWriter := &CountingWriteCloser{}

    // create writer
    writer := NewWriter(countingWriter)

    // test write
    n, err = writer.Write(buffer)
    assert.Nil(t, err)
    assert.Equal(t, n, len(buffer))
    assert.Equal(t, n, countingWriter.Count)

    // test close
    err = writer.Close()
    assert.Nil(t, err)
    assert.True(t, countingWriter.Closed)
}

func TestWriterUse(t *testing.T) {

    // create test data
    buffer := make([]byte, 1024)
    n, err := rand.Read(buffer)
    assert.Nil(t, err)
    assert.Equal(t, n, len(buffer))

    // test utils
    crcTable := crc64.MakeTable(crc64.ISO)
    countingWriter := &CountingWriteCloser{}
    crcWriter := CRC64WriterMiddleware{0, false, crcTable, countingWriter}
    csum := crc64.Checksum(buffer, crcTable)

    // create writer
    writer := NewWriter(countingWriter)

    // add middleware
    writer.Use(func(w io.WriteCloser) io.WriteCloser {
        return &crcWriter
    })

    // test write
    n, err = writer.Write(buffer)
    assert.Nil(t, err)
    assert.Equal(t, n, len(buffer))
    assert.Equal(t, n, countingWriter.Count)

    // ensure crc matches
    assert.Equal(t, crcWriter.CRC64, csum)

    // test close
    err = writer.Close()
    assert.Nil(t, err)
    assert.True(t, crcWriter.Closed)
    assert.True(t, countingWriter.Closed)
}
