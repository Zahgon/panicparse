// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate go run regen.go

package stack

import (
	"html/template"
	"io"
	"regexp"
)

// ToHTML formats the aggregated buckets as HTML to the writer.
//
// Use footer to add custom HTML at the bottom of the page.
func (a *Aggregated) ToHTML(w io.Writer, footer template.HTML) error {
	_ = "STUB: not implemented"
	return nil
}

// ToHTML formats the snapshot as HTML to the writer.
//
// Use footer to add custom HTML at the bottom of the page.
func (s *Snapshot) ToHTML(w io.Writer, footer template.HTML) error {
	_ = "STUB: not implemented"
	return nil
}

// Private stuff.

func toHTML(w io.Writer, data map[string]interface{}) error { _ = "STUB: not implemented"; return nil }

var reMethodSymbol = regexp.MustCompile(`^\(\*?([^)]+)\)(\..+)$`)

func funcClass(c *Call) template.HTML { _ = "STUB: not implemented"; return *new(template.HTML) }

/* #nosec G203 */

func minus(i, j int) int {
	_ = "STUB: not implemented"

	// pkgURL returns a link to the godoc for the call.
	return 0
}

func pkgURL(c *Call) template.URL {
	_ = "STUB: not implemented"

	// Check for vendored code first.
	return *new(template.URL)
}

// This always links to the latest release, past releases are not online.
// That's somewhat unfortunate.

// TODO(maruel): Leverage Location.
// Use pkg.go.dev when there's a version (go module) and godoc.org when
// there's none (implies branch master).

// srcURL returns an URL to the sources.
//
// TODO(maruel): Support custom local godoc server as it serves files too.
func srcURL(c *Call) template.URL { _ = "STUB: not implemented"; return *new(template.URL) }

func escape(s string) template.URL {
	_ = "STUB: not implemented"
	// That's the only way I found to get the kind of escaping I wanted, where
	// '/' is not escaped.
	return *new(template.URL)
}

/* #nosec G203 */

// getSrcBranchURL returns a link to the source on the web and the tag name for
// the package version, if possible.
func getSrcBranchURL(c *Call) (template.URL, template.URL) {
	_ = "STUB: not implemented"
	return *new(template.URL), *new(template.URL)
}

// TODO(maruel): This is not strictly speaking correct. The remote could be
// running a different Go version from the current executable.

/* #nosec G203 */

// TODO(maruel): Leverage Location.

// Check for vendored code first.

// Specialized support for github and golang.org. This will cover a fair
// share of the URLs, but it'd be nice to support others too. Please submit
// a PR (including a unit test that I was too lazy to add yet).

/* #nosec G203 */

// https://github.com/golang/build/blob/HEAD/repos/repos.go lists all
// the golang.org/x/<foo> packages.

// parts is: "x", <project@version>, <path inside the repo>.

// The source of truth is are actually go.googlesource.com, but
// github.com has nicer syntax highlighting.

/* #nosec G203 */

// For example gopkg.in. In this case there's no known way to find the
// link to the source files, but we can still try to extract the version
// if fetched from a go module.
// Do a best effort to find a version by searching for a '@'.

/* #nosec G203 */

/* #nosec G203 */

func splitHost(s string) (string, string) { _ = "STUB: not implemented"; return "", "" }

// "v0.0.0-20200223170610-d5e6a3e2c0ae"
var reVersion = regexp.MustCompile(`v\d+\.\d+\.\d+\-\d+\-([a-f0-9]+)`)

func splitTag(s string) (string, string, template.URL) {
	_ = "STUB: not implemented"
	// Default to branch master for non-versioned dependencies. It's not
	// optimal but it's better than nothing?
	// TODO(maruel): Replace with HEAD.
	return "", "", *new(template.URL)
}

// No tag was found.

// We got a versioned go module.

/* #nosec G203 */

// symbol is the hashtag to use to refer to the symbol when looking at
// documentation.
//
// All of godoc/gddo, pkg.go.dev and golang.org/godoc use the same symbol
// reference format.
func symbol(f *Func) template.URL { _ = "STUB: not implemented"; return *new(template.URL) }

// Transform the method form.

/* #nosec G203 */
