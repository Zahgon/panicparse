// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package stack analyzes stack dump of Go processes and simplifies it.
//
// It is mostly useful on servers will large number of identical goroutines,
// making the crash dump harder to read than strictly necessary.
package stack

import (
	"os"
	"strings"
)

// Func is a function call in a goroutine stack trace.
type Func struct {
	// Complete is the complete reference. It can be ambiguous in case where a
	// path contains dots.
	Complete string
	// ImportPath is the directory name for this function reference, or "main" if
	// it was in package main. The package name may not match.
	ImportPath string
	// DirName is the directory name containing the package in which the function
	// is. Normally this matches the package name, but sometimes there's smartass
	// folks that use a different directory name than the package name.
	DirName string
	// Name is the function name or fully quality method name.
	Name string
	// IsExported is true if the function is exported.
	IsExported bool
	// IsPkgMain is true if it is in the main package.
	IsPkgMain bool

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// Init parses the raw function call line from a goroutine stack trace.
//
// Go stack traces print a mangled function call, this wrapper unmangle the
// string before printing and adds other filtering methods.
//
// The main caveat is that for calls in package main, the package import URL is
// left out.
func (f *Func) Init(raw string) error {
	_ = "STUB: not implemented"
	// Format can be:
	//   - gopkg.in/yaml%2ev2.(*Struct).Method  (handling dots is tricky)
	//   - main.func·001  (go statements)
	//   - foo  (C code)
	//
	// The function is optimized to reduce its memory usage.
	return nil
}

// Cut the path elements.

// It's fine if there's no dot, it happens in C code in go1.4 and lower.

// Only the path part is escaped.

// Update the index in the unescaped string.

// TODO(go1.20): switch to strings.CutSuffix

// Consider main.main to be exported.

// String returns Complete.
func (f *Func) String() string {
	_ = "STUB: not implemented"

	// Arg is an argument on a Call.
	return ""
}

type Arg struct {
	// IsAggregate is true if the argument is an aggregate type. If true, the
	// argument does not contain a value itself, but contains a set of nested
	// argument fields. If false, the argument contains a single scalar value.
	IsAggregate bool

	// The following are set if IsAggregate == false.

	// Name is a pseudo name given to the argument.
	Name string
	// Value is the raw value as found in the stack trace
	Value uint64
	// IsPtr is true if we guess it's a pointer. It's only a guess, it can be
	// easily confused by a bitmask.
	IsPtr bool
	// IsOffsetTooLarge is true if the argument's frame offset was too large,
	// preventing the argument from being printed in the stack trace.
	IsOffsetTooLarge bool

	// IsInaccurate determines if Value is inaccurate. Stacks could have inaccurate values
	// for arguments passed in registers. Go 1.18 prints a ? for these values.
	IsInaccurate bool

	// The following are set if IsAggregate == true.

	// Fields are the fields/elements of aggregate-typed arguments.
	Fields Args

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

const zeroToNine = "0123456789"

// String prints the argument as the name if present, otherwise as the value.
func (a *Arg) String() string { _ = "STUB: not implemented"; return "" }

const (
	// The pointer floor can be below 1MiB (!) on Windows.
	// Assumes all values are above 512KiB and positive are pointers; assuming
	// that above half the memory is kernel memory.
	//
	// This is not always true but this should be good enough to help
	// implementing AnyPointer.
	pointerFloor = 512 * 1024
	// Assume the stack was generated with the same bitness (32 vs 64) than the
	// code processing it.
	pointerCeiling = uint64((^uint(0)) >> 1)
)

// equal returns true only if both arguments are exactly equal.
func (a *Arg) equal(r *Arg) bool { _ = "STUB: not implemented"; return false }

// similar returns true if the two Arg are equal or almost but not quite equal.
func (a *Arg) similar(r *Arg, similar Similarity) bool { _ = "STUB: not implemented"; return false }

// Args is a series of function call arguments.
type Args struct {
	// Values is the arguments as shown on the stack trace. They are mangled via
	// simplification.
	Values []Arg
	// Processed is the arguments generated from processing the source files. It
	// can have a length lower than Values.
	Processed []string
	// Elided when set means there was a trailing ", ...".
	Elided bool

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

func (a *Args) String() string { _ = "STUB: not implemented"; return "" }

// equal returns true only if both arguments are exactly equal.
func (a *Args) equal(r *Args) bool { _ = "STUB: not implemented"; return false }

// similar returns true if the two Args are equal or almost but not quite
// equal.
func (a *Args) similar(r *Args, similar Similarity) bool { _ = "STUB: not implemented"; return false }

// merge merges two similar Args, zapping out differences.
func (a *Args) merge(r *Args) Args { _ = "STUB: not implemented"; return *new(Args) }

// walk traverses all non-aggregate arguments in the Args struct, calling the
// provided visitor function with each Arg.
func (a *Args) walk(visitor func(arg *Arg)) { _ = "STUB: not implemented"; return }

// Location is the source location, if determined.
type Location int

const (
	// LocationUnknown is the default value when Opts.GuessPaths was false.
	LocationUnknown Location = iota
	// GoMod is a go module, it is outside $GOPATH and is inside a directory
	// containing a go.mod file. This is considered a local copy.
	GoMod
	// GOPATH is in $GOPATH/src. This is either a dependency fetched via
	// GO111MODULE=off or intentionally fetched this way. There is no guaranteed
	// that the local copy is pristine.
	GOPATH
	// GoPkg is in $GOPATH/pkg/mod. This is a dependency fetched via go module.
	// It is considered to be an unmodified external dependency.
	GoPkg
	// Stdlib is when it is a Go standard library function. This includes the 'go
	// test' generated main executable.
	Stdlib

	lastLocation
)

// Call is an item in the stack trace.
//
// All paths in this struct are in POSIX format, using "/" as path separator.
type Call struct {
	// The following are initialized on the first line of the call stack.

	// Func is the fully qualified function name (encoded).
	Func Func
	// Args is the call arguments.
	Args Args

	// The following are initialized on the second line of the call stack.

	// RemoteSrcPath is the full path name of the source file as seen in the
	// trace.
	RemoteSrcPath string
	// Line is the line number.
	Line int
	// SrcName is the base file name of the source file.
	SrcName string
	// DirSrc is one directory plus the file name of the source file. It is a
	// subset of RemoteSrcPath.
	DirSrc string

	// The following are only set if Opts.GuessPaths was set.

	// LocalSrcPath is the full path name of the source file as seen in the host,
	// if found.
	LocalSrcPath string
	// RelSrcPath is the relative path to GOROOT, GOPATH or LocalGoMods.
	RelSrcPath string
	// ImportPath is the fully qualified import path as found on disk (when
	// Opts.GuessPaths was set). Defaults to Func.ImportPath otherwise.
	//
	// In the case of package "main", it returns the underlying path to the main
	// package instead of "main" if Opts.GuessPaths was set.
	ImportPath string
	// Location is the source location, if determined.
	Location Location

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// Init initializes RemoteSrcPath, SrcName, DirName and Line.
//
// For test main, it initializes Location only with Stdlib.
//
// It does its best educated guess for ImportPath.
func (c *Call) init(srcPath string, line int) {
	c.Line = line
	if srcPath != "" {
		c.RemoteSrcPath = srcPath
		if i := strings.LastIndexByte(c.RemoteSrcPath, '/'); i != -1 {
			c.SrcName = c.RemoteSrcPath[i+1:]
			if i = strings.LastIndexByte(c.RemoteSrcPath[:i], '/'); i != -1 {
				c.DirSrc = c.RemoteSrcPath[i+1:]
			}
		}
		if c.DirSrc == testMainSrc {
			// Consider _test/_testmain.go as stdlib since it's injected by "go test".
			c.Location = Stdlib
		}
	}
	c.ImportPath = c.Func.ImportPath
}

const testMainSrc = "_test" + string(os.PathSeparator) + "_testmain.go"

// updateLocations initializes LocalSrcPath, RelSrcPath, Location and ImportPath.
//
// goroot, localgoroot, localgomod, gomodImportPath and gopaths are expected to
// be in "/" format even on Windows. They must not have a trailing "/".
//
// Returns true if a match was found.
func (c *Call) updateLocations(goroot, localgoroot string, localgomods, gopaths map[string]string) bool {
	_ = "STUB: not implemented"
	// TODO(maruel): Reduce memory allocations.
	return false
}

// Check GOROOT first.

// Replace remote GOROOT with local GOROOT.

// Check GOPATH.
// TODO(maruel): Sort for deterministic behavior?

// For modules, the path has to be altered, as it contains the version.

// Check Go modules.
// Go module path detection only works with stack traces created on the local
// file system.

// Maybe the path is just absolute and exists?

// equal returns true only if both calls are exactly equal.
func (c *Call) equal(r *Call) bool { _ = "STUB: not implemented"; return false }

// similar returns true if the two Call are equal or almost but not quite
// equal.
func (c *Call) similar(r *Call, similar Similarity) bool { _ = "STUB: not implemented"; return false }

// merge merges two similar Call, zapping out differences.
func (c *Call) merge(r *Call) Call { _ = "STUB: not implemented"; return *new(Call) }

// Stack is a call stack.
type Stack struct {
	// Calls is the call stack. First is original function, last is leaf
	// function.
	Calls []Call
	// Elided is set when there's >100 items in Stack, currently hardcoded in
	// package runtime.
	Elided bool

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// equal returns true on if both call stacks are exactly equal.
func (s *Stack) equal(r *Stack) bool { _ = "STUB: not implemented"; return false }

// similar returns true if the two Stack are equal or almost but not quite
// equal.
func (s *Stack) similar(r *Stack, similar Similarity) bool { _ = "STUB: not implemented"; return false }

// merge merges two similar Stack, zapping out differences.
func (s *Stack) merge(r *Stack) *Stack {
	_ = "STUB: not implemented"
	// Assumes similar stacks have the same length.
	return nil
}

// less compares two Stack, where the ones that are less are more
// important, so they come up front.
//
// A Stack with more private functions is 'less' so it is at the top.
// Inversely, a Stack with only public functions is 'more' so it is at the
// bottom.
func (s *Stack) less(r *Stack) bool { _ = "STUB: not implemented"; return false }

// Check unknown code type last.

// Stack lengths are the same and they are mostly of the same kind of location.

// Stacks are the same.

// updateLocations calls updateLocations on each call frame and returns true if
// they were all resolved.
func (s *Stack) updateLocations(goroot, localgoroot string, localgomods, gopaths map[string]string) bool {
	_ = "STUB: not implemented"
	// If there were none, it was "resolved".
	return false
}

// Signature represents the signature of one or multiple goroutines.
//
// It is effectively the stack trace plus the goroutine internal bits, like
// it's state, if it is thread locked, which call site created this goroutine,
// etc.
type Signature struct {
	// State is the goroutine state at the time of the snapshot.
	//
	// Use git grep 'gopark(|unlock)\(' to find them all plus everything listed
	// in runtime/traceback.go. Valid values includes:
	//     - chan send, chan receive, select
	//     - finalizer wait, mark wait (idle),
	//     - Concurrent GC wait, GC sweep wait, force gc (idle)
	//     - IO wait, panicwait
	//     - semacquire, semarelease
	//     - sleep, timer goroutine (idle)
	//     - trace reader (blocked)
	// Stuck cases:
	//     - chan send (nil chan), chan receive (nil chan), select (no cases)
	// Runnable states:
	//    - idle, runnable, running, syscall, waiting, dead, enqueue, copystack,
	// Scan states:
	//    - scan, scanrunnable, scanrunning, scansyscall, scanwaiting, scandead,
	//      scanenqueue
	//
	// When running under the race detector, the values are 'running' or
	// 'finished'.
	State string
	// CreatedBy is the call stack that created this goroutine, if applicable.
	//
	// Normally, the stack is a single Call.
	//
	// When the race detector is enabled, a full stack snapshot is available.
	CreatedBy Stack
	// SleepMin is the wait time in minutes, if applicable.
	//
	// Not set when running under the race detector.
	SleepMin int
	// SleepMax is the wait time in minutes, if applicable.
	//
	// Not set when running under the race detector.
	SleepMax int
	// Stack is the call stack.
	Stack Stack
	// Locked is set if the goroutine was locked to an OS thread.
	//
	// Not set when running under the race detector.
	Locked bool

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// equal returns true only if both signatures are exactly equal.
func (s *Signature) equal(r *Signature) bool { _ = "STUB: not implemented"; return false }

// similar returns true if the two Signature are equal or almost but not quite
// equal.
func (s *Signature) similar(r *Signature, similar Similarity) bool {
	_ = "STUB: not implemented"
	return false
}

// merge merges two similar Signature, zapping out differences.
func (s *Signature) merge(r *Signature) *Signature { _ = "STUB: not implemented"; return nil }

// Drop right side.
// Drop right side.

// TODO(maruel): This is weirdo.

// less compares two Signature, where the ones that are less are more
// important, so they come up front. A Signature with more private functions is
// 'less' so it is at the top. Inversely, a Signature with only public
// functions is 'more' so it is at the bottom.
func (s *Signature) less(r *Signature) bool { _ = "STUB: not implemented"; return false }

// SleepString returns a string "N-M minutes" if the goroutine(s) slept for a
// long time.
//
// Returns an empty string otherwise.
func (s *Signature) SleepString() string { _ = "STUB: not implemented"; return "" }

// updateLocations calls updateLocations on both CreatedBy and Stack and
// returns true if they were both resolved.
func (s *Signature) updateLocations(goroot, localgoroot string, localgomods, gopaths map[string]string) bool {
	_ = "STUB: not implemented"
	return false
}

// Goroutine represents the state of one goroutine, including the stack trace.
type Goroutine struct {
	// Signature is the stack trace, internal bits, state, which call site
	// created it, etc.
	Signature
	// ID is the goroutine id.
	ID int
	// First is the goroutine first printed, normally the one that crashed.
	First bool

	// RaceWrite is true if a race condition was detected, and this goroutine was
	// race on a write operation, otherwise it was a read.
	RaceWrite bool
	// RaceAddr is set to the address when a data race condition was detected.
	// Otherwise it is 0.
	RaceAddr uint64

	// Disallow initialization with unnamed parameters.
	_ struct{}
}

// Private stuff.

// nameArguments is a post-processing step where Args are 'named' with numbers.
func nameArguments(goroutines []*Goroutine) {
	_ = "STUB: not implemented"
	// Set a name for any pointer occurring more than once.
	return
}

// Enumerate all the arguments.

// CreatedBy.Args is never set.

// Now do the rest. This is done so the output is deterministic.

// Process the remaining pointers, they were not referenced by primary
// thread so will have higher IDs.

func pathJoin(s ...string) string { _ = "STUB: not implemented"; return "" }

type uint64Slice []uint64

func (a uint64Slice) Len() int           { _ = "STUB: not implemented"; return 0 }
func (a uint64Slice) Swap(i, j int)      { _ = "STUB: not implemented"; return }
func (a uint64Slice) Less(i, j int) bool { _ = "STUB: not implemented"; return false }
