// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package internal implements the handlers for panicweb so they are in a
// separate package than "main".
package internal

import (
	"net/http"
)

// Unblock unblocks one http server handler.
var Unblock = make(chan struct{})

// GetAsync does an HTTP GET to the URL but leaves the actual fetching to a
// goroutine.
func GetAsync(url string) {
	_ = "STUB: not implemented"
	/* #nosec G107 */ return
}

// URL1Handler is a http.HandlerFunc that hangs.
func URL1Handler(w http.ResponseWriter, req *http.Request) {
	_ = "STUB: not implemented"
	// Respond the HTTP header to unblock the http.Get() function.
	return
}

// URL2Handler is a http.HandlerFunc that hangs.
func URL2Handler(w http.ResponseWriter, req *http.Request) {
	_ = "STUB: not implemented"
	// Respond the HTTP header to unblock the http.Get() function.
	return
}
