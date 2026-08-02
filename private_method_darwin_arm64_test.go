//go:build darwin && arm64
// +build darwin,arm64

package gomonkey

import (
	"encoding/binary"
	"testing"
)

func TestBuildArm64BranchDirective(t *testing.T) {
	from := uintptr(0x100000000)
	to := from + 0x100
	code, err := buildArm64BranchDirective(from, to)
	if err != nil {
		t.Fatalf("buildArm64BranchDirective returned an error: %v", err)
	}
	if len(code) != 4 {
		t.Fatalf("branch directive has length %d, want 4", len(code))
	}

	instruction := binary.LittleEndian.Uint32(code)
	if got := instruction & 0xfc000000; got != 0x14000000 {
		t.Fatalf("branch opcode is %#x, want %#x", got, uint32(0x14000000))
	}
	if got := instruction & 0x03ffffff; got != uint32((to-from)/4) {
		t.Fatalf("branch offset is %#x, want %#x", got, uint32((to-from)/4))
	}

	if _, err := buildArm64BranchDirective(from, from+arm64BranchRange); err == nil {
		t.Fatal("out-of-range branch should be rejected")
	}
}
