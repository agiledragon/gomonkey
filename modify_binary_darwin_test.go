package gomonkey

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func TestValidateWritePageSize(t *testing.T) {
	if err := validateWritePageSize(writeIsolationSize); err != nil {
		t.Fatalf("supported page size was rejected: %v", err)
	}

	unsupported := writeIsolationSize * 2
	err := validateWritePageSize(unsupported)
	if err == nil {
		t.Fatalf("page size %d should be rejected", unsupported)
	}
	if !strings.Contains(err.Error(), "unsupported system page size") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// A 24-byte jmp patch that begins close to the end of a page must make both
// the current and the next page writable, otherwise writing it triggers SIGBUS
// (see issue #186).
func TestProtectSizeAcrossPageBoundary(t *testing.T) {
	ps := syscall.Getpagesize()
	base := uintptr(ps * 4) // any page-aligned address

	cases := []struct {
		name   string
		target uintptr
		n      int
		want   int
	}{
		{"patch fully inside one page", base + 16, 24, ps},
		{"patch ends on the last byte of the page", base + uintptr(ps) - 24, 24, ps},
		{"24-byte patch starts 16 bytes before boundary", base + uintptr(ps) - 16, 24, 2 * ps},
		{"patch starts on the very last byte", base + uintptr(ps) - 1, 24, 2 * ps},
	}
	for _, c := range cases {
		if got := protectSize(c.target, c.n); got != c.want {
			t.Errorf("%s: protectSize=%d, want %d", c.name, got, c.want)
		}
	}
}

// End-to-end reproduction of issue #186: writing a 24-byte patch whose tail
// spills into an adjacent, non-writable page must not crash. Runs in a child
// process because a SIGBUS from a bad write cannot be recovered.
func TestModifyBinaryAcrossPageBoundary(t *testing.T) {
	if os.Getenv("GOMONKEY_REPRO_186") == "1" {
		reproCrossPageWrite()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestModifyBinaryAcrossPageBoundary$")
	cmd.Env = append(os.Environ(), "GOMONKEY_REPRO_186=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("modifyBinary crashed on a cross-page write (issue #186): %v\n%s", err, out)
	}
}

func reproCrossPageWrite() {
	ps := syscall.Getpagesize()
	// Two consecutive pages, both writable to start with...
	mem, err := syscall.Mmap(-1, 0, 2*ps,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		os.Exit(2)
	}
	base := PtrOf(mem)
	// ...then revoke write on the SECOND page, so a cross-page write faults.
	if _, _, e := syscall.Syscall(syscall.SYS_MPROTECT,
		base+uintptr(ps), uintptr(ps), uintptr(syscall.PROT_READ)); e != 0 {
		os.Exit(2)
	}
	// target sits 16 bytes before the boundary -> a 24-byte patch spills 8
	// bytes into the read-only second page.
	target := base + uintptr(ps) - 16
	modifyBinary(target, make([]byte, 24)) // SIGBUS here before the fix
	os.Exit(0)
}
