// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

import (
	"github.com/maruel/panicparse/v2/stack"
)

// Palette defines the color used.
//
// An empty object Palette{} can be used to disable coloring.
type Palette struct {
	EOLReset string

	// Routine header.
	RoutineFirst string // The first routine printed.
	Routine      string // Following routines.
	CreatedBy    string
	Race         string

	// Call line.
	Package                     string
	SrcFile                     string
	FuncMain                    string
	FuncLocationUnknown         string
	FuncLocationUnknownExported string
	FuncGoMod                   string
	FuncGoModExported           string
	FuncGOPATH                  string
	FuncGOPATHExported          string
	FuncGoPkg                   string
	FuncGoPkgExported           string
	FuncStdLib                  string
	FuncStdLibExported          string
	Arguments                   string
}

// pathFormat determines how much to show.
type pathFormat int

const (
	fullPath pathFormat = iota
	relPath
	basePath
)

func (pf pathFormat) formatCall(c *stack.Call) string { _ = "STUB: not implemented"; return "" }

func (pf pathFormat) createdByString(s *stack.Signature) string {
	_ = "STUB: not implemented"
	return ""
}

// calcBucketsLengths returns the maximum length of the source lines and
// package names.
func calcBucketsLengths(a *stack.Aggregated, pf pathFormat) (int, int) {
	_ = "STUB: not implemented"
	return 0, 0
}

// calcGoroutinesLengths returns the maximum length of the source lines and
// package names.
func calcGoroutinesLengths(s *stack.Snapshot, pf pathFormat) (int, int) {
	_ = "STUB: not implemented"
	return 0, 0
}

// functionColor returns the color to be used for the function name based on
// the type of package the function is in.
func (p *Palette) functionColor(c *stack.Call) string { _ = "STUB: not implemented"; return "" }

func (p *Palette) funcColor(l stack.Location, main, exported bool) string {
	_ = "STUB: not implemented"
	return ""
}

// routineColor returns the color for the header of the goroutines bucket.
func (p *Palette) routineColor(first, multipleBuckets bool) string {
	_ = "STUB: not implemented"
	return ""
}

// BucketHeader prints the header of a goroutine signature.
func (p *Palette) BucketHeader(b *stack.Bucket, pf pathFormat, multipleBuckets bool) string {
	_ = "STUB: not implemented"
	return ""
}

// GoroutineHeader prints the header of a goroutine.
func (p *Palette) GoroutineHeader(g *stack.Goroutine, pf pathFormat, multipleGoroutines bool) string {
	_ = "STUB: not implemented"
	return ""
}

// callLine prints one stack line.
func (p *Palette) callLine(line *stack.Call, srcLen, pkgLen int, pf pathFormat) string {
	_ = "STUB: not implemented"
	return ""
}

// StackLines prints one complete stack trace, without the header.
func (p *Palette) StackLines(signature *stack.Signature, srcLen, pkgLen int, pf pathFormat) string {
	_ = "STUB: not implemented"
	return ""
}
