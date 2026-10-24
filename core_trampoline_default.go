//go:build !amd64 && (!darwin || !arm64)
// +build !amd64
// +build !darwin !arm64

package gomonkey

func buildCorePatchDirective(_, double, _ uintptr) ([]byte, uintptr) {
	return buildJmpDirective(double), 0
}

func releaseCorePatchTrampoline(trampoline uintptr) {
	releasePrivateMethodTrampoline(trampoline)
}
