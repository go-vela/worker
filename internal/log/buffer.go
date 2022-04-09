// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

// This file is based on https://stackoverflow.com/a/44322300 .

import (
	"io"
	"sync"
)

// Buffer is an in-memory log cache that allows multiple reads.
type Buffer interface {
	io.Writer
	// NewReader returns an io.Reader that reads the Buffer from the beginning.
	NewReader() io.Reader
}

// buffer implements the in-memory log cache Buffer.
type buffer struct {
	// data contains the log cache
	data [][]byte
	// RWMutex is required to ensure safe write and read from the same cache.
	sync.RWMutex
}

// Write writes to the shared Buffer cache, implementing io.Writer.
// This manages thread safety using a sync.RWMutex.
func (b *buffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Cannot retain p, so we must copy it:
	p2 := make([]byte, len(p))
	copy(p2, p)
	b.Lock()
	b.data = append(b.data, p2)
	b.Unlock()
	return len(p), nil
}

// bufferReader reads from buffer from the beginning without consuming its contents.
// bufferReader implements the io.Reader returned by buffer.NewReader().
type bufferReader struct {
	buf   *buffer // buffer we read from
	index int     // next slice index
	data  []byte  // current data slice to serve
}

// Read reads from the shared Buffer cache, implementing io.Writer.
// This manages thread safety using a sync.RWMutex.
func (br *bufferReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Do we have data to send?
	if len(br.data) == 0 {
		buf := br.buf
		buf.RLock()
		if br.index < len(buf.data) {
			br.data = buf.data[br.index]
			br.index++
		}
		buf.RUnlock()
	}
	if len(br.data) == 0 {
		return 0, io.EOF
	}

	n = copy(p, br.data)
	br.data = br.data[n:]
	return n, nil
}

// NewReader returns an io.Reader that reads from the shared buffer
// from the beginning without consuming its contents.
func (b *buffer) NewReader() io.Reader {
	return &bufferReader{buf: b}
}

// NewBuffer creates and returns a new empty Buffer ready for Read and Write.
func NewBuffer() Buffer {
	return &buffer{}
}
