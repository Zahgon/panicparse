// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internaltest

import (
	"errors"
	"sync"
)

// PanicwebOutput returns the output of panicweb with inlining disabled.
//
// The function panics if any internal error occurs.
func PanicwebOutput() []byte { _ = "STUB: not implemented"; return nil }

// PanicOutputs returns a map of the output of every subcommands.
//
// panic is built with inlining disabled.
//
// The subcommand "race" is built with the race detector. Others are built
// without. In particular "asleep" doesn't work with the race detector.
//
// The function panics if any internal error occurs.
func PanicOutputs() map[string][]byte { _ = "STUB: not implemented"; return nil }

// Extracts the subcommands, then run each of them individually.

// The odd of this failing is close to nil.

// Race detector is not supported on this platform.

// Collect the subcommands.

// Collect the output of each subcommand.

// Race detector is not supported.

// StaticPanicwebOutput returns a constant version of panicweb output for use
// in benchmarks.
func StaticPanicwebOutput() []byte { _ = "STUB: not implemented"; return nil }

// StaticPanicRaceOutput returns a constant version of 'panic race' output.
func StaticPanicRaceOutput() []byte { _ = "STUB: not implemented"; return nil }

// IsUsingModules returns if go modules are enabled.
//
// It reads the current value of GO111MODULES.
func IsUsingModules() bool { _ = "STUB: not implemented"; return false }

//

var (
	panicwebOnce     sync.Once
	panicwebOutput   []byte
	panicOutputsOnce sync.Once
	panicOutputs     map[string][]byte
)

// GetGoMinorVersion returns the Go1 minor version.
//
// Returns 0 for a developer build, panics if can't parse the version.
//
// Ignores the revision (go1.<minor>.<revision>).
func GetGoMinorVersion() int { _ = "STUB: not implemented"; return 0 }

// This will break on go2. Please submit a PR to fix this once Go2 is
// released.

// build builds to a temporary file and returns the path to it.
func build(tool string, race bool) string { _ = "STUB: not implemented"; return "" }

var errNoRace = errors.New("platform does not support -race")

// Compile compiles sources into an executable.
func Compile(in, exe, cwd string, disableInlining, race bool) error {
	_ = "STUB: not implemented"
	// Disable optimization (-N) and inlining (-l) otherwise the inlining varies
	// between local execution and remote execution. This can be observed as
	// Elided being true without any argument.
	return nil
}

/* #nosec G204 */

// execRun runs a command and returns the combined output.
//
// It ignores the exit code, since it's meant to run panic, which crashes by
// design.
func execRun(cmd ...string) []byte {
	_ = "STUB: not implemented"
	/* #nosec G204 */ return nil
}
