//go:build amd64
// +build amd64

package gomonkey

import (
	"encoding/binary"
	"fmt"
	"sync"
)

const (
	amd64TrampolineSize  = 16
	amd64TrampolineCount = 256
	amd64BranchSize      = 5
)

var amd64Trampolines struct {
	sync.Mutex
	used [amd64TrampolineCount]bool
}

// amd64TrampolineBase returns the address of the metadata-safe assembly slots
// declared in trampoline_amd64.s.
func amd64TrampolineBase() uintptr

func buildCorePatchDirective(target, double, trampoline uintptr) ([]byte, uintptr) {
	acquired := false
	if trampoline == 0 {
		trampoline = acquireAmd64Trampoline(target)
		acquired = true
	}
	completed := false
	defer func() {
		if acquired && !completed {
			releaseAmd64Trampoline(trampoline)
		}
	}()

	branch, err := buildAmd64BranchDirective(target, trampoline)
	if err != nil {
		panic(err)
	}

	jump := buildJmpDirective(double)
	code := make([]byte, amd64TrampolineSize)
	for i := range code {
		code[i] = 0x90 // NOP
	}
	copy(code, jump)
	modifyBinary(trampoline, code)
	completed = true
	return branch, trampoline
}

func buildAmd64BranchDirective(from, to uintptr) ([]byte, error) {
	delta := int64(to) - int64(from+amd64BranchSize)
	if int64(int32(delta)) != delta {
		return nil, fmt.Errorf("amd64 trampoline %#x is out of rel32 range from target %#x", to, from)
	}

	code := make([]byte, amd64BranchSize)
	code[0] = 0xe9 // JMP rel32
	binary.LittleEndian.PutUint32(code[1:], uint32(int32(delta)))
	return code, nil
}

func acquireAmd64Trampoline(target uintptr) uintptr {
	amd64Trampolines.Lock()
	defer amd64Trampolines.Unlock()

	base := amd64TrampolineBase()
	for i := range amd64Trampolines.used {
		if amd64Trampolines.used[i] {
			continue
		}
		slot := base + uintptr(i*amd64TrampolineSize)
		if _, err := buildAmd64BranchDirective(target, slot); err != nil {
			continue
		}
		amd64Trampolines.used[i] = true
		return slot
	}
	panic(fmt.Sprintf("no reachable amd64 patch trampoline available for target %#x (%d slots)", target, amd64TrampolineCount))
}

func releaseAmd64Trampoline(trampoline uintptr) {
	base := amd64TrampolineBase()
	if trampoline < base {
		return
	}
	offset := trampoline - base
	if offset%amd64TrampolineSize != 0 {
		return
	}
	index := offset / amd64TrampolineSize
	if index >= amd64TrampolineCount {
		return
	}

	amd64Trampolines.Lock()
	amd64Trampolines.used[index] = false
	amd64Trampolines.Unlock()
}

func releaseCorePatchTrampoline(trampoline uintptr) {
	releaseAmd64Trampoline(trampoline)
}
