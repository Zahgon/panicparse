// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package webstack provides a http.HandlerFunc that serves a snapshot similar
// to net/http/pprof.Index().
//
// Contrary to net/http/pprof, the handler is not automatically registered.
package webstack

import (
	"net/http"

	"github.com/maruel/panicparse/v2/stack"
)

// SnapshotHandler implements http.HandlerFunc to returns a panicparse HTML
// format for a snapshot of the current goroutines.
//
// For best results, compile the executable with optimization (-N) and inlining
// (-l) disabled with -gcflags '-N -l'.
//
// Arguments are passed as form values. If you want to change the default,
// override the form values in a wrapper as shown in the example.
//
// The implementation is designed to be reasonably fast, it currently does a
// small amount of disk I/O only for file presence.
//
// It is a direct replacement for "/debug/pprof/goroutine?debug=2" handler in
// net/http/pprof.
//
// augment: (default: 1) When set to 1, panicparse tries to find the sources on
// disk to improve the display of arguments based on type information. This is
// slower and should be avoided on high utilization server.
//
// maxmem: (default: 67108864) maximum amount of temporary memory to use to
// generate a snapshot. In practice at least the double of this is used.
// Minimum is 1048576.
//
// similarity: (default: "anypointer") Can be one of stack.Similarity value in
// lowercase: "exactflags", "exactlines", "anypointer" or "anyvalue".
func SnapshotHandler(w http.ResponseWriter, req *http.Request) { _ = "STUB: not implemented"; return }

// snapshot returns a Context based on the snapshot of the stacks of the
// current process.
func snapshot(maxmem int, opts *stack.Opts) (*stack.Snapshot, error) {
	_ = "STUB: not implemented"
	// We don't know how big the buffer needs to be to collect all the
	// goroutines. Start with 1 MB and try a few times, doubling each time. Give
	// up and use a truncated trace if maxmem is not enough.
	return nil, nil
}

// That's expected.
