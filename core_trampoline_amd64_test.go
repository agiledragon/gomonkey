//go:build amd64
// +build amd64

package gomonkey

import (
	"encoding/binary"
	"io"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

//go:noinline
func amd64MetadataTarget(value int) int { return value + 1 }

//go:noinline
func amd64MetadataReplacement(value int) int { return value + 2 }

func TestBuildAmd64BranchDirective(t *testing.T) {
	from := uintptr(0x100000000)
	to := from + 0x100
	code, err := buildAmd64BranchDirective(from, to)
	if err != nil {
		t.Fatalf("buildAmd64BranchDirective returned an error: %v", err)
	}
	if len(code) != amd64BranchSize || code[0] != 0xe9 {
		t.Fatalf("branch directive = %x, want a 5-byte JMP rel32", code)
	}
	if got, want := int32(binary.LittleEndian.Uint32(code[1:])), int32(to-from-amd64BranchSize); got != want {
		t.Fatalf("branch displacement = %d, want %d", got, want)
	}

	if _, err := buildAmd64BranchDirective(from, from+uintptr(1<<31)+uintptr(amd64BranchSize)); err == nil {
		t.Fatal("out-of-range branch should be rejected")
	}
}

func TestApplyFuncKeepsRuntimeMetadataValidDuringProfiling(t *testing.T) {
	patches := ApplyFunc(amd64MetadataTarget, amd64MetadataReplacement)
	defer patches.Reset()

	if got := amd64MetadataTarget(10); got != 12 {
		t.Fatalf("patched call returned %d, want 12", got)
	}

	if err := pprof.StartCPUProfile(io.Discard); err != nil {
		t.Fatalf("start CPU profile: %v", err)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				_ = amd64MetadataTarget(10)
			}
		}
	}()

	time.Sleep(1500 * time.Millisecond)
	close(stop)
	<-done
	pprof.StopCPUProfile()
	runtime.GC()

	patches.Reset()
	if got := amd64MetadataTarget(10); got != 11 {
		t.Fatalf("restored call returned %d, want 11", got)
	}
}

func TestAmd64TrampolineHasGoRuntimeFunctionMetadata(t *testing.T) {
	base := amd64TrampolineBase()
	fn := runtime.FuncForPC(base)
	if fn == nil {
		t.Fatal("trampoline slot is not represented in Go runtime function metadata")
	}
	second := runtime.FuncForPC(base + 2*amd64TrampolineSize)
	if second == nil || second.Entry() != fn.Entry() {
		t.Fatal("trampoline slots do not share the metadata-bearing assembly function")
	}
}
