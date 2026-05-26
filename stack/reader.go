// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package stack

import (
	"errors"
	"io"
)

var (
	errBufferFull = errors.New("buffer full")
)

type reader struct {
	buf  [16 * 1024]byte
	rd   io.Reader
	r, w int
	err  error
}

// fill reads a new chunk into the buffer.
func (r *reader) fill() {
	_ = "STUB: not implemented"
	// Slide existing data to beginning.
	return
}

// Read new data: try a limited number of times.

func (r *reader) buffered() []byte { _ = "STUB: not implemented"; return nil }

func (r *reader) readSlice() ([]byte, error) { _ = "STUB: not implemented"; return nil, nil }

// readLine is our own implementation of ReadBytes().
//
// We try to use readSlice() as much as we can but we need to tolerate if an
// input line is longer than the buffer specified at Reader creation. Not using
// the more complicated slice of slices that Reader.ReadBytes() uses since it
// should not happen often here. Instead bootstrap the memory allocation by
// starting with 4x buffer size, which should get most cases with a single
// allocation.
func (r *reader) readLine() ([]byte, error) { _ = "STUB: not implemented"; return nil, nil }
