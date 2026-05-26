// Copyright 2020 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:build windows

package main

func sysHang() {
	_ = "STUB: not implemented"
	// 49.7 days is enough for everyone.
	return
}
