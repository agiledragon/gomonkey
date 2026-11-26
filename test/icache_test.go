package test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

//go:noinline
func icacheTarget() int { return 1 }

// Calling the target first pulls it into the icache; a patch applied afterwards
// must still be seen. On CPUs without coherent icaches (CTR_EL0.DIC=0) this
// needs explicit cache maintenance after the code is written.
func TestPatchAfterTargetRanIsVisible(t *testing.T) {
	stalePatch, staleReset := 0, 0
	for i := 0; i < 2000; i++ {
		for j := 0; j < 100; j++ {
			icacheTarget()
		}
		p := gomonkey.ApplyFunc(icacheTarget, func() int { return 2 })
		if icacheTarget() != 2 {
			stalePatch++
		}
		p.Reset()
		if icacheTarget() != 1 {
			staleReset++
		}
	}
	if stalePatch+staleReset > 0 {
		t.Fatalf("stale code executed: %d/2000 after patch, %d/2000 after reset", stalePatch, staleReset)
	}
}
