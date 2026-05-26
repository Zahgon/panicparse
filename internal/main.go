// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package internal implements panicparse
//
// It is mostly useful on servers will large number of identical goroutines,
// making the crash dump harder to read than strictly necessary.
//
// Colors:
//   - Magenta: first goroutine to be listed.
//   - Yellow: main package.
//   - Green: standard library.
//   - Red: other packages.
//
// Bright colors are used for exported symbols.
package internal

import (
	"html/template"
	"io"
	"regexp"

	"github.com/maruel/panicparse/v2/stack"
	"github.com/mgutz/ansi"
)

// resetFG is similar to ansi.Reset except that it doesn't reset the
// background color, only the foreground color and the style.
//
// That much for the "ansi" abstraction layer...
const resetFG = ansi.DefaultFG + "\033[m"

// defaultPalette is the default recommended palette.
var defaultPalette = Palette{
	EOLReset:                    resetFG,
	RoutineFirst:                ansi.ColorCode("magenta+b"),
	CreatedBy:                   ansi.LightBlack,
	Race:                        ansi.LightRed,
	Package:                     ansi.ColorCode("default+b"),
	SrcFile:                     resetFG,
	FuncMain:                    ansi.ColorCode("yellow+b"),
	FuncLocationUnknown:         ansi.White,
	FuncLocationUnknownExported: ansi.ColorCode("white+b"),
	FuncGoMod:                   ansi.Red,
	FuncGoModExported:           ansi.ColorCode("red+b"),
	FuncGOPATH:                  ansi.Cyan,
	FuncGOPATHExported:          ansi.ColorCode("cyan+b"),
	FuncGoPkg:                   ansi.Blue,
	FuncGoPkgExported:           ansi.ColorCode("blue+b"),
	FuncStdLib:                  ansi.Green,
	FuncStdLibExported:          ansi.ColorCode("green+b"),
	Arguments:                   resetFG,
}

func writeBucketsToConsole(out io.Writer, p *Palette, a *stack.Aggregated, pf pathFormat, needsEnv bool, filter, match *regexp.Regexp) error {
	_ = "STUB: not implemented"
	return nil
}

func writeGoroutinesToConsole(out io.Writer, p *Palette, s *stack.Snapshot, pf pathFormat, needsEnv bool, filter, match *regexp.Regexp) error {
	_ = "STUB: not implemented"
	return nil
}

type toHTMLer interface {
	ToHTML(io.Writer, template.HTML) error
}

func toHTML(h toHTMLer, p string, needsEnv bool) error {
	_ = "STUB: not implemented"
	/* #nosec G304 */ return nil
}

func processInner(out io.Writer, p *Palette, s stack.Similarity, pf pathFormat, html string, filter, match *regexp.Regexp, c *stack.Snapshot, first bool) error {
	_ = "STUB: not implemented"
	return nil
}

// Bucketing should only be done if no data race was detected.

// It's a data race.

// process copies stdin to stdout and processes any "panic: " line found.
//
// If html is used, a stack trace is written to this file instead.
func process(in io.Reader, out io.Writer, p *Palette, s stack.Similarity, pf pathFormat, parse, rebase bool, html string, filter, match *regexp.Regexp) error {
	_ = "STUB: not implemented"
	return nil
}

// Process it even if an error occurred.

// This means the whole buffer was not read, loop again.

// Parts of the input will be lost.

func showBanner() bool { _ = "STUB: not implemented"; return false }

// Main is implemented here so both 'pp' and 'panicparse' executables can be
// compiled. This is to work around the Perl Package manager 'pp' that is
// preinstalled on some OSes.
func Main() error { _ = "STUB: not implemented"; return nil }

// Console only.

// HTML only.

// Explicitly silence SIGQUIT, as it is useful to gather the stack dump
// from the piped command.

// Do not handle SIGQUIT when passed a file to process.

/* #nosec G304 */

/* #nosec G307 */
