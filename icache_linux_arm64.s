#include "textflag.h"

// func flushICache(addr, n uintptr)
// Same sequence as libgcc's __clear_cache: clean D-cache to PoU, invalidate
// I-cache, then synchronize, with line sizes read from CTR_EL0.
TEXT ·flushICache(SB),NOSPLIT,$0-16
	MOVD addr+0(FP), R0
	MOVD n+8(FP), R1
	ADD  R0, R1, R1           // R1 = end
	MRS  CTR_EL0, R2
	MOVD $4, R4
	UBFX $16, R2, $4, R3
	LSL  R3, R4, R3           // R3 = D-cache line size
	AND  $0xf, R2, R5
	LSL  R5, R4, R5           // R5 = I-cache line size

	SUB  $1, R3, R6
	BIC  R6, R0, R7
dloop:
	CMP  R1, R7
	BHS  ddone
	WORD $0xd50b7b27          // DC CVAU, X7
	ADD  R3, R7, R7
	B    dloop
ddone:
	DSB  $11                  // ISH

	SUB  $1, R5, R6
	BIC  R6, R0, R7
iloop:
	CMP  R1, R7
	BHS  idone
	WORD $0xd50b7527          // IC IVAU, X7
	ADD  R5, R7, R7
	B    iloop
idone:
	DSB  $11                  // ISH
	ISB  $15                  // SY
	RET
