// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate go install golang.org/x/tools/cmd/stringer@v0.24.0
//go:generate stringer -type state
//go:generate stringer -type Location

package stack

import (
	"io"
	"path/filepath"
	"regexp"
)

// Opts represents options to process the snapshot.
type Opts struct {
	// LocalGOROOT is GOROOT with "/" as path separator. No trailing "/". Can be
	// unset.
	LocalGOROOT string
	// LocalGOPATHs is GOPATH with "/" as path separator. No trailing "/". Can be
	// unset.
	LocalGOPATHs []string

	// NameArguments tells panicparse to find the recurring pointer values and
	// give them pseudo 'names'.
	//
	// Since the algorithm is O(n²), this can be worth disabling on live servers.
	NameArguments bool

	// GuessPaths tells panicparse to guess local RemoteGOROOT and GOPATH for
	// what was found in the snapshot.
	//
	// Initializes in Snapshot the following members: RemoteGOROOT,
	// RemoteGOPATHs, LocalGomoduleRoot and GomodImportPath.
	//
	// This is done by scanning the local disk, so be warned of performance
	// impact.
	GuessPaths bool

	// AnalyzeSources tells panicparse to processes source files to improve calls
	// to be more descriptive.
	//
	// Requires GuessPaths to be true.
	AnalyzeSources bool

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// DefaultOpts returns default options to process the snapshot.
func DefaultOpts() *Opts { _ = "STUB: not implemented"; return nil }

func (o *Opts) isValid() bool { _ = "STUB: not implemented"; return false }

// Snapshot is a parsed runtime.Stack() or race detector dump.
type Snapshot struct {
	// Goroutines is the Goroutines found.
	//
	// They are in the order that they were printed.
	Goroutines []*Goroutine

	// LocalGOROOT is copied from Opts.
	LocalGOROOT string
	// LocalGOPATHs is copied from Opts.
	LocalGOPATHs []string

	// The following members are initialized when Opts.GuessPaths is true.

	// RemoteGOROOT is the GOROOT as detected in the traceback, not the on the
	// host.
	//
	// It can be empty if no root was determined, for example the traceback
	// contains only non-stdlib source references.
	RemoteGOROOT string
	// RemoteGOPATHs is the GOPATH as detected in the traceback, with the value
	// being the corresponding path mapped to the host if found.
	//
	// It can be empty if only stdlib code is in the traceback or if no local
	// sources were matched up. In the general case there is only one entry in
	// the map.
	RemoteGOPATHs map[string]string

	// LocalGomods are the root directories containing go.mod or that directly
	// contained source code as detected in the traceback, with the value being
	// the corresponding import path found in the go.mod file.
	//
	// Uses "/" as path separator. No trailing "/".
	//
	// Because of the "replace" statement in go.mod, there can be multiple root
	// directories. A file run by "go run" is also considered a go module to (a
	// certain extent).
	//
	// It is initialized by findRoots().
	//
	// Unlike GOROOT and GOPATH, it only works with stack traces created in the
	// local file system, hence "Local" prefix.
	LocalGomods map[string]string

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// ScanSnapshot scans the Reader for the output from runtime.Stack() in br.
//
// Returns nil *Snapshot if no stack trace was detected.
//
// If a Snapshot is returned, you can call the function again to find another
// trace, or do io.Copy(br, out) to flush the rest of the stream.
//
// ParseSnapshot processes the output from runtime.Stack() or the race detector.
//
// Returns a nil *Snapshot if no stack trace was detected and SearchSnapshot()
// was a false positive.
//
// Returns io.EOF if all of reader was read.
//
// The suffix of the stack trace is returned as []byte.
//
// It pipes anything not detected as a panic stack trace from r into out. It
// assumes there is junk before the actual stack trace. The junk is streamed to
// out.
func ScanSnapshot(in io.Reader, prefix io.Writer, opts *Opts) (*Snapshot, []byte, error) {
	_ = "STUB: not implemented"
	return nil, nil, nil
}

// TODO(maruel): Validate opts.

// IsRace returns true if a race detector stack trace was found.
//
// Otherwise, it is a normal goroutines snapshot.
//
// When a race condition was detected, it is preferable to not call Aggregate().
func (s *Snapshot) IsRace() bool { _ = "STUB: not implemented"; return false }

func (s *Snapshot) guessPaths() bool { _ = "STUB: not implemented"; return false }

// Note that this is important to call it even if
// s.RemoteGOROOT == s.LocalGOROOT.

// augment processes source files to improve calls to be more descriptive.
//
// It modifies goroutines in place. It requires calling guessPaths() to work
// properly.
//
// Returns the last error that occurred while processing files.
func (s *Snapshot) augment() error { _ = "STUB: not implemented"; return nil }

// Private stuff.

const pathSeparator = string(filepath.Separator)

var (
	lockedToThread = []byte("locked to thread")
	// gotRaceHeader1, done
	raceHeaderFooter = []byte("==================")
	// gotRaceHeader2
	raceHeader             = []byte("WARNING: DATA RACE")
	crlf                   = []byte("\r\n")
	lf                     = []byte("\n")
	commaSpace             = []byte(", ")
	writeCap               = []byte("Write")
	writeLow               = []byte("write")
	threeDots              = []byte("...")
	underscore             = []byte("_")
	inaccurateQuestionMark = []byte("?")
)

// These are effectively constants.
var (
	// gotRoutineHeader
	reRoutineHeader = regexp.MustCompile("^([ \t]*)goroutine (\\d+)(?: gp=[^ ]+ m=[^ ]+(?: mp=[^ ]+)?)? \\[([^\\]]+)\\]\\:$")
	reMinutes       = regexp.MustCompile(`^(\d+) minutes$`)

	// gotUnavail
	reUnavail = regexp.MustCompile("^(?:\t| +)goroutine running on other thread; stack unavailable")

	// gotFileFunc, gotRaceOperationFile, gotRaceGoroutineFile
	// See gentraceback() in src/runtime/traceback.go for more information.
	// - Sometimes the source file comes up as "<autogenerated>". It is the
	//   compiler than generated these, not the runtime.
	// - The tab may be replaced with spaces when a user copy-paste it, handle
	//   this transparently.
	// - "runtime.gopanic" is explicitly replaced with "panic" by gentraceback().
	// - The +0x123 byte offset is printed when frame.pc > _func.entry. _func is
	//   generated by the linker.
	// - The +0x123 byte offset is not included with generated code, e.g. unnamed
	//   functions "func·006()" which is generally go func() { ... }()
	//   statements. Since the _func is generated at runtime, it's probably why
	//   _func.entry is not set.
	// - C calls may have fp=0x123 sp=0x123 appended. I think it normally happens
	//   when a signal is not correctly handled. It is printed with m.throwing>0.
	//   These are discarded.
	// - For cgo, the source file may be "??".
	reFile = regexp.MustCompile("^(?:\t| +)(\\?\\?|\\<autogenerated\\>|.+\\.(?:c|go|s))\\:(\\d+)(?:| \\+0x[0-9a-f]+)(?:| fp=0x[0-9a-f]+ sp=0x[0-9a-f]+(?:| pc=0x[0-9a-f]+))$")

	// gotCreated
	// Sadly, it doesn't note the goroutine number so we could cascade them per
	// parenthood.
	reCreated = regexp.MustCompile("^created by (.+)$")

	// gotFunc, gotRaceOperationFunc, gotRaceGoroutineFunc
	reFunc = regexp.MustCompile(`^(.+)\((.*)\)$`)

	// Race:
	// See https://github.com/llvm/llvm-project/blob/HEAD/compiler-rt/lib/tsan/rtl/tsan_report.cpp
	// for the code generating these messages. Please note only the block in
	//   #else  // #if !SANITIZER_GO
	// is used.
	// TODO(maruel): "    [failed to restore the stack]\n\n"
	// TODO(maruel): "Global var %s of size %zu at %p declared at %s:%zu\n"

	// gotRaceOperationHeader
	reRaceOperationHeader = regexp.MustCompile(`^(Read|Write) at (0x[0-9a-f]+) by goroutine (\d+):$`)

	// gotRaceOperationHeader
	reRacePreviousOperationHeader = regexp.MustCompile(`^Previous (read|write) at (0x[0-9a-f]+) by goroutine (\d+):$`)

	// gotRaceGoroutineHeader
	reRaceGoroutine = regexp.MustCompile(`^Goroutine (\d+) \((running|finished)\) created at:$`)

	// TODO(maruel): Use it.
	//reRacePreviousOperationMainHeader = regexp.MustCompile("^Previous (read|write) at (0x[0-9a-f]+) by main goroutine:$")
)

// state is the state of the scan to detect and process a stack trace.
type state int

// Initial state is looking. Other states are when a stack trace is detected.
const (
	// Haven't found a stack trace yet.
	// to: gotRoutineHeader, raceHeader1
	looking state = iota

	// Done processing a stack trace.
	done

	// Panic stack trace:

	// Signature: ""
	// An empty line between goroutines.
	// from: gotFileCreated, gotFileFunc
	// to: gotRoutineHeader, done
	betweenRoutine
	// Regexp: reRoutineHeader
	// Signature: "goroutine 1 [running]:"
	// Goroutine header was found.
	// from: looking
	// to: gotUnavail, gotFunc
	gotRoutineHeader
	// Regexp: reFunc
	// Signature: "main.main()"
	// Function call line was found.
	// from: gotRoutineHeader
	// to: gotFileFunc
	gotFunc
	// Regexp: reCreated
	// Signature: "created by main.init.func4"
	// Goroutine creation line was found.
	// from: gotFileFunc
	// to: gotFileCreated
	gotCreated
	// Regexp: reFile
	// Signature: "\t/foo/bar/baz.go:116 +0x35"
	// File header was found.
	// from: gotFunc
	// to: gotFunc, gotCreated, betweenRoutine, done
	gotFileFunc
	// Regexp: reFile
	// Signature: "\t/foo/bar/baz.go:116 +0x35"
	// File header was found.
	// from: gotCreated
	// to: betweenRoutine, done
	gotFileCreated
	// Regexp: reUnavail
	// Signature: "goroutine running on other thread; stack unavailable"
	// State when the goroutine stack is instead is reUnavail.
	// from: gotRoutineHeader
	// to: betweenRoutine, gotCreated
	gotUnavail

	// Race detector:

	// Constant: raceHeaderFooter
	// Signature: "=================="
	// from: looking
	// to: done, gotRaceHeader2
	gotRaceHeader1
	// Constant: raceHeader
	// Signature: "WARNING: DATA RACE"
	// from: gotRaceHeader1
	// to: done, gotRaceOperationHeader
	gotRaceHeader2
	// Regexp: reRaceOperationHeader, reRacePreviousOperationHeader
	// Signature: "Read at 0x00c0000e4030 by goroutine 7:"
	// A race operation was found.
	// from: gotRaceHeader2
	// to: done, gotRaceOperationFunc
	gotRaceOperationHeader
	// Regexp: reFunc
	// Signature: "  main.panicRace.func1()"
	// Function that caused the race.
	// from: gotRaceOperationHeader
	// to: done, gotRaceOperationFile
	gotRaceOperationFunc
	// Regexp: reFile
	// Signature: "\t/foo/bar/baz.go:116 +0x35"
	// File header that caused the race.
	// from: gotRaceOperationFunc
	// to: done, betweenRaceOperations, gotRaceOperationFunc
	gotRaceOperationFile
	// Signature: ""
	// Empty line between race operations or just after.
	// from: gotRaceOperationFile
	// to: done, gotRaceOperationHeader, gotRaceGoroutineHeader
	betweenRaceOperations

	// Regexp: reRaceGoroutine
	// Signature: "Goroutine 7 (running) created at:"
	// Goroutine header.
	// from: betweenRaceOperations, betweenRaceGoroutines
	// to: done, gotRaceOperationHeader
	gotRaceGoroutineHeader
	// Regexp: reFunc
	// Signature: "  main.panicRace.func1()"
	// Function that caused the race.
	// from: gotRaceGoroutineHeader
	// to: done, gotRaceGoroutineFile
	gotRaceGoroutineFunc
	// Regexp: reFile
	// Signature: "\t/foo/bar/baz.go:116 +0x35"
	// File header that caused the race.
	// from: gotRaceGoroutineFunc
	// to: done, betweenRaceGoroutines
	gotRaceGoroutineFile
	// Signature: ""
	// Empty line between race stack traces.
	// from: gotRaceGoroutineFile
	// to: done, gotRaceGoroutineHeader
	betweenRaceGoroutines
)

// scanningState is the state of the scan to detect and process a stack trace
// and stores the traces found.
type scanningState struct {
	*Snapshot
	state          state
	prefix         []byte
	goroutineIndex int
}

func isFramesElidedLine(line []byte) bool {
	_ = "STUB: not implemented"
	// before go1.21:
	// ...additional frames elided...
	//
	// go1.21 and newer:
	// print("...", elide, " frames elided...\n")
	return false
}

// scan scans one line, updates goroutines and move to the next state.
//
// Returns true if the line was processed and thus should not be printed out.
//
// TODO(maruel): Handle corrupted stack cases:
// - missed stack barrier
// - found next stack barrier at 0x123; expected
// - runtime: unexpected return pc for FUNC_NAME called from 0x123
func (s *scanningState) scan(line []byte) (bool, error) {
	_ = "STUB: not implemented"
	/* This is very useful to debug issues in the state machine.
	defer func() {
		log.Printf("scan(%q) -> %s", line, s.state)
	}()
	//*/return false, nil
}

// It's the end of the stream and it's not terminating with EOL character.

// Let it flow. It's possible the last line was trimmed and we still want
// to parse it.

// This can only be the case if s.state != looking | done or the line is
// empty.

// We could look for '^panic:' but this is more risky, there can be a lot
// of junk between this and the stack dump.

// Look for a goroutine header.

// See runtime/traceback.go.
// "<state>, \d+ minutes, locked to thread"

// Look for duration, if any.

// Increase performance by always allocating 4 goroutines minimally.

// Switch to race detection mode.

// TODO(maruel): We should buffer it in case the next line is not a
// WARNING so we can output it back.

// Generate a fake stack entry.

// Next line is expected to be an empty line.

// Increase performance by always allocating 4 calls minimally.

// cur.Stack.Calls is guaranteed to have at least one item.

// This initializes ImportPath.

// TODO(maruel): New state.

// Increase performance by always allocating 4 calls minimally.

// Race detector.

// TODO(maruel): We should buffer it in case the next line is not a
// WARNING so we can output it back.

// TODO(maruel): While this shouldn't error out, it should still force the
// output of raceHeaderFooter.

// Increase performance by always allocating 4 calls minimally.

// Look for other previous race data operations.

// Race stack traces

// TODO(maruel): Set s.Goroutines[].CreatedBy.

// parseFunc only return an error if it also returns true.
//
// Uses reFunc.
func parseFunc(c *Call, line []byte) (bool, error) { _ = "STUB: not implemented"; return false, nil }

// It is also done in c.init() but do it here in case of a corrupted trace
// for the file section.

// parseArgs parses a collection of comma-separated arguments into an Args
// struct.
func parseArgs(line []byte) (Args, error) {
	_ = "STUB: not implemented"
	// 5 from traceback.go, +1 for top level
	return *new(Args), nil
}

// Assume the stack was generated with the same bitness (32 vs 64) as
// the code processing it.

// parseFile only return an error if also processing a Call.
//
// Uses reFile.
func parseFile(c *Call, line []byte) (bool, error) { _ = "STUB: not implemented"; return false, nil }

// hasPrefix returns true if any of s is the prefix of p.
func hasPrefix(p string, s map[string]string) bool { _ = "STUB: not implemented"; return false }

// hasSrcPrefix returns true if any of s is the prefix of p with /src/ or
// /pkg/mod/.
func hasSrcPrefix(p string, s map[string]string) bool { _ = "STUB: not implemented"; return false }

// getFiles returns all the source files deduped and ordered.
func getFiles(goroutines []*Goroutine) []string { _ = "STUB: not implemented"; return nil }

// splitPath splits a path using "/" as separator into its components.
//
// The first item has its initial path separator kept.
func splitPath(p string) []string { _ = "STUB: not implemented"; return nil }

// isFile returns true if the path is a valid file.
func isFile(p string) bool {
	_ = "STUB: not implemented"
	// TODO(maruel): Is it faster to open the file or to stat it? Worth a perf
	// test on Windows.
	return false
}

// isRootedIn returns a root if the file split in parts exists under root.
//
// Uses "/" as path separator.
func isRootedIn(root string, parts []string) string { _ = "STUB: not implemented"; return "" }

// reModule find the module line in a go.mod file. It works even on CRLF file.
var reModule = regexp.MustCompile(`(?m)^module\s+([^\n\r]+)\r?$`)

type gomodCache map[string]struct{}

// isGoModule returns the string to the directory containing a go.mod file, and
// the go import path it represents, if found.
func (g *gomodCache) isGoModule(parts []string) (string, string) {
	_ = "STUB: not implemented"
	return "", ""
}

// Was already looked up.

/* #nosec G304 */

// findRoots sets member RemoteGOROOT, RemoteGOPATHs and LocalGomods.
//
// This causes disk I/O as it checks for file presence.
//
// Returns the number of missing files.
func (s *Snapshot) findRoots() int {
	_ = "STUB: not implemented"
	// TODO(maruel): Reduce memory allocations in this function.
	return 0
}

// TODO(maruel): Could a stack dump have mixed cases? I think it's
// possible, need to confirm and handle.
//log.Printf("  Analyzing %s", f)

// First checks skip file I/O.

// stdlib.

// $GOPATH/src or go.mod dependency in $GOPATH/pkg/mod.

// At this point, disk will be looked up.

// Initializes RemoteGOROOT.

//log.Printf("Found RemoteGOROOT=%s", s.RemoteGOROOT)

// Initializes RemoteGOPATHs.

//log.Printf("Found RemoteGOPATHs[%s] = %s", r[:len(r)-len(src)], l)

//log.Printf("Found RemoteGOPATHs[%s] = %s", r[:len(r)-len(pkgmod)], l)

// Initializes localGomods.

// Search upward looking for a go.mod.

// Assumes "go run" was used, thus is package main. Still consider it a
// "go module" but in the weakest sense.

// If the source is not found, just too bad.
//log.Printf("Failed to find locally: %s", f)

// getGOPATHs returns parsed GOPATH or its default, using "/" as path separator.
func getGOPATHs() []string { _ = "STUB: not implemented"; return nil }

// Disallow non-absolute paths?

// Trim trailing "/".

// atou is a fast Atoi() function.
//
// It is a very simplified version of strconv.Atoi() that it never go into the
// slow path and it operates on []byte instead of string so it doesn't do
// memory allocation. It will fail on edge cases like prefix of zeros and other
// things that the panic stack trace generator never outputs.
//
// It doesn't handle negative values.
func atou(s []byte) (int, bool) { _ = "STUB: not implemented"; return 0, false }

// trimLeftSpace is the faster equivalent of bytes.TrimLeft(s, "\t ").
func trimLeftSpace(s []byte) []byte { _ = "STUB: not implemented"; return nil }

// trimCurlyBrackets is the faster equivalent of
// bytes.TrimRight(bytes.TrimLeft(s, "{"), "}"). The function
// also returns the number of curly brackets trimmed from the
// left and the right.
func trimCurlyBrackets(s []byte) (int, []byte, int) { _ = "STUB: not implemented"; return 0, nil, 0 }

// unsafeString performs an unsafe conversion from a []byte to a string. The
// returned string will share the underlying memory with the []byte which thus
// allows the string to be mutable through the []byte. We're careful to use
// this method only in situations in which the []byte will not be modified.
//
// A workaround for the absence of https://github.com/golang/go/issues/2632.
func unsafeString(b []byte) string {
	_ = "STUB: not implemented"
	/* #nosec G103 */ return ""
}
