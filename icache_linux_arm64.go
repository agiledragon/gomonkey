package gomonkey

// flushICache makes freshly written code visible to instruction fetch. arm64
// CPUs without a coherent I-cache (CTR_EL0.DIC=0, e.g. Apple silicon) otherwise
// keep executing the old instructions; the kernel does not flush on mprotect.
func flushICache(addr, n uintptr)
