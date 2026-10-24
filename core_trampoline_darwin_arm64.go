//go:build darwin && arm64
// +build darwin,arm64

package gomonkey

func buildCorePatchDirective(target, double, trampoline uintptr) ([]byte, uintptr) {
	return buildPrivateMethodDirective(target, double, trampoline)
}

func releaseCorePatchTrampoline(trampoline uintptr) {
	releasePrivateMethodTrampoline(trampoline)
}
