#include "textflag.h"

TEXT runtime·mach_task_self_trampoline(SB),NOSPLIT,$0
PUSHQ	BP
MOVQ	SP, BP
CALL	libsystem_mach_task_self(SB)
MOVQ	AX, 0(DI)
POPQ	BP
RET

TEXT runtime·mach_vm_protect_trampoline(SB),NOSPLIT,$0
PUSHQ	BP
MOVQ	SP, BP
MOVQ	target_task+0(FP), DI
MOVQ	address+8(FP), SI
MOVL	size+16(FP), DX
MOVL	setMaximum+24(FP), CX
MOVL	newProt+32(FP), R8
CALL	libsystem_mach_vm_protect(SB)
MOVQ	AX, 0(DI)
POPQ	BP
RET

