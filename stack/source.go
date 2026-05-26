// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// This file contains the code to process sources, to be able to deduct the
// original types.

package stack

import (
	"go/ast"
)

// Private stuff.

// cacheAST is a cache of parsed Go sources.
type cacheAST struct {
	files  map[string][]byte
	parsed map[string]*parsedFile
}

// augmentGoroutine processes source files to improve call to be more
// descriptive.
//
// It modifies the routine.
func (c *cacheAST) augmentGoroutine(g *Goroutine) error { _ = "STUB: not implemented"; return nil }

// Only load the AST if there's an argument to process.

//log.Printf("%s", err)

// loadFile loads a Go source file and parses the AST tree.
func (c *cacheAST) loadFile(fileName string) error { _ = "STUB: not implemented"; return nil }

// Do not attempt to parse the same file twice.

// Ignore C and assembly.

/* #nosec G304 */

// lineToByteOffsets extract the line number into raw file offset.
//
// Inserts a dummy 0 at offset 0 so line offsets can be 1 based.
func lineToByteOffsets(src []byte) []int { _ = "STUB: not implemented"; return nil }

// parsedFile is a processed Go source file.
type parsedFile struct {
	lineToByteOffset []int
	parsed           *ast.File
}

// getFuncAST gets the callee site function AST representation for the code
// inside the function f at line l.
func (p *parsedFile) getFuncAST(f string, l int) (d *ast.FuncDecl, err error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// The line number in the stack trace line does not exist in the file. That
// can only mean that the sources on disk do not match the sources used to
// build the binary.

// Walk the AST to find the lineToByteOffset that fits the line number.

// Inspect() goes depth first. This means for example that a function like:
// func a() {
//   b := func() {}
//   c()
// }
//
// Were we are looking at the c() call can return confused values. It is
// important to look at the actual ast.Node hierarchy.

// We are expecting a ast.CallExpr node. It can be harder to figure out
// when there are multiple calls on a single line, as the stack trace
// doesn't have file byte offset information, only line based.
// gofmt will always format to one function call per line but there can
// be edge cases, like:
//   a = A{Foo(), Bar()}

//p.processNode(call, n)

func name(n ast.Node) string { _ = "STUB: not implemented"; return "" }

// fieldToType returns the type name and whether if it's an ellipsis.
func fieldToType(f *ast.Field) (string, bool) { _ = "STUB: not implemented"; return "", false }

// Array.

// Slice.

// Do not print the function signature to not overload the trace.

// TODO(maruel): Implement anything missing.

// extractArgumentsType returns the name of the type of each input argument.
func extractArgumentsType(f *ast.FuncDecl) ([]string, bool) {
	_ = "STUB: not implemented"
	return nil, false
}

// If it is an object receiver (vs a pointer receiver), its address is not
// printed in the stack trace so it needs to be ignored.

// Assert that ellipsis is only set on the last item of fields?

// augmentCall walks the function and populate call accordingly.
func augmentCall(call *Call, f *ast.FuncDecl) { _ = "STUB: not implemented"; return }

// These are unexpected value! Print them as hex.

// If top-level argument is an aggregate-type, include each
// of its sub-arguments.

// Assumes it's an interface. For now, discard the object
// value, which is probably not a good idea.
