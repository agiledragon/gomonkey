package test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

// r10ReadResult flattens to 9 register slots:
//
//	Data map[string]interface{} -> 1 (pointer)
//	Body []byte                 -> 3 (ptr + len + cap)
//	Name string                 -> 2 (ptr + len)
//	Code int                    -> 1
//	OK   bool                   -> 1
//	Ver  int64                  -> 1
//
// 9 slots in total, so the data word of r10Parse's "out any" argument lands
// in R10 (the 11th argument slot).
type r10ReadResult struct {
	Data map[string]interface{}
	Body []byte
	Name string
	Code int
	OK   bool
	Ver  int64
}

type r10License struct {
	Status int
}

//go:noinline
func r10Parse(in r10ReadResult, out interface{}) error {
	return nil
}

// TestApplyFuncArm64R10 is a regression test for the arm64 R10 clobber.
// Before the fix, on darwin/arm64 and linux/arm64 the mock writes to its
// own text segment and crashes with SIGBUS; with the fix it passes.
// amd64 uses jmp_amd64.go and is unaffected.
func TestApplyFuncArm64R10(t *testing.T) {
	patches := gomonkey.ApplyFunc(r10Parse, func(in r10ReadResult, out interface{}) error {
		lic := out.(*r10License)
		lic.Status = 42
		return nil
	})
	defer patches.Reset()

	lic := &r10License{}
	if err := r10Parse(r10ReadResult{}, lic); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lic.Status != 42 {
		t.Fatalf("lic.Status = %d, want 42", lic.Status)
	}
}
