#include "textflag.h"

TEXT ·mach_task_self_trampoline(SB),NOSPLIT,$0
	PUSHQ	BP            // make a frame; keep stack aligned
	MOVQ	SP, BP
	CALL	libsystem_mach_task_self(SB)
	MOVQ	AX, 0(DI)     // return value arg1 ret
	POPQ	BP
	RET

TEXT ·mach_vm_protect_trampoline(SB),NOSPLIT,$0
	PUSHQ	BP            // make a frame; keep stack aligned
	MOVQ	SP, BP
	MOVQ	DI, BX        // BX is caller-save
	MOVQ	0(BX), DI     // arg 1 targetTask
    MOVQ	8(BX), SI     // arg 2 address
    MOVL	16(BX), DX    // arg 3 size
    MOVL	24(BX), CX    // arg 4 setMaximum
    MOVL	32(BX), R8    // arg 5 newProtection
	CALL	libsystem_mach_vm_protect(SB)
	MOVQ	AX, 40(BX)    // return value arg6 ret
	POPQ	BP
	RET
