#include "textflag.h"

// func callGadget(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
//
// ABI0 (stack-based, FP offsets):
//   in:  gadgetAddr+0  gadgetKind+8  nr+16  a1+24  a2+32  a3+40  a4+48  a5+56  a6+64
//   out: r1+72  r2+80  errno+88
//   frame size: 0 locals, 96 bytes args
TEXT ·callGadget(SB),NOSPLIT,$0-96
    MOVQ gadgetAddr+0(FP), R11
    MOVQ gadgetKind+8(FP), R12
    MOVQ nr+16(FP),         AX
    MOVQ a1+24(FP),         DI
    MOVQ a2+32(FP),         SI
    MOVQ a3+40(FP),         DX
    MOVQ a4+48(FP),        R10
    MOVQ a5+56(FP),         R8
    MOVQ a6+64(FP),         R9

    CMPQ R12, $1
    JEQ  do_frame

do_direct:
    // gadgDirect: syscall; ret
    CALL R11
    JMP  check_err

do_frame:
    // gadgFrame: syscall; pop rbp; ret
    // We use a CALL to a global trampoline to push the return address.
    CALL ·setup_stack(SB)

check_err:
    CMPQ AX, $0xfffffffffffff001
    JLS  ok
    MOVQ $-1, r1+72(FP)
    MOVQ $0,  r2+80(FP)
    NEGQ AX
    MOVQ AX,  errno+88(FP)
    RET
ok:
    MOVQ AX, r1+72(FP)
    MOVQ DX, r2+80(FP)
    MOVQ $0, errno+88(FP)
    RET

// trampoline for gadgFrame
TEXT ·setup_stack(SB),NOSPLIT,$0
    PUSHQ BP
    JMP R11
