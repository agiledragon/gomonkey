//go:build linux && !arm64
// +build linux,!arm64

package gomonkey

func flushICache(addr, n uintptr) {}
