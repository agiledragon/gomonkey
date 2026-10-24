//go:build amd64
// +build amd64

#include "textflag.h"

#define NOP8 BYTE $0x90; BYTE $0x90; BYTE $0x90; BYTE $0x90; BYTE $0x90; BYTE $0x90; BYTE $0x90; BYTE $0x90;
#define NOP64 NOP8; NOP8; NOP8; NOP8; NOP8; NOP8; NOP8; NOP8;
#define NOP512 NOP64; NOP64; NOP64; NOP64; NOP64; NOP64; NOP64; NOP64;
#define NOP4096 NOP512; NOP512; NOP512; NOP512; NOP512; NOP512; NOP512; NOP512;

// amd64TrampolineBase exposes the start of the fixed-size patch trampoline
// slots. The slots are part of a NOSPLIT|NOFRAME function with zero stack delta,
// so Go's runtime metadata remains valid while a patched function tail-jumps
// through one of them.
TEXT ·amd64TrampolineBase(SB),NOSPLIT,$0-8
	LEAQ amd64TrampolineSpace<>(SB), AX
	MOVQ AX, ret+0(FP)
	RET

// 256 slots of 16 bytes each: 12 bytes for MOVABS/JMP plus 4 bytes of padding.
TEXT amd64TrampolineSpace<>(SB),NOSPLIT|NOFRAME,$0-0
	NOP4096
	RET
