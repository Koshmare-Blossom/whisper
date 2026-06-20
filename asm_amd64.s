#include "textflag.h"

// ─── whisper: CALL-based (Hell's Gate / Halo's Gate) ────────────────────────
//
// func callGadget(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
//
// ABI0 layout:
//   in:  gadgetAddr+0  gadgetKind+8  nr+16  a1+24  a2+32  a3+40  a4+48  a5+56  a6+64
//   out: r1+72  r2+80  errno+88
TEXT ·callGadget(SB),NOSPLIT,$0-96
    MOVQ gadgetAddr+0(FP), R11
    MOVQ gadgetKind+8(FP), CX
    MOVQ nr+16(FP),         AX
    MOVQ a1+24(FP),         DI
    MOVQ a2+32(FP),         SI
    MOVQ a3+40(FP),         DX
    MOVQ a4+48(FP),        R10
    MOVQ a5+56(FP),         R8
    MOVQ a6+64(FP),         R9
    CMPQ CX, $1
    JEQ  do_frame
do_direct:
    CALL R11
    JMP  check_err
do_frame:
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

TEXT ·setup_stack(SB),NOSPLIT,$0
    PUSHQ BP
    JMP R11

// ─── hellshall: JMP-based ────────────────────────────────────────────────────
//
// func hellsHallCall(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
//
// Instead of CALL gadget (which pushes a return addr into our code), we JMP.
// The gadget's ret returns directly to the Go caller via hhResult<>,
// so whisper never appears in the call stack during the syscall.
//
// Before JMP:
//   R12 = FP of this frame (where return values must be written)
//   R13 = original return addr to Go caller (saved from [rsp])
//   [rsp] is replaced with &hhResult<> so the gadget's ret lands there
TEXT ·hellsHallCall(SB),NOSPLIT,$0-96
    MOVQ gadgetAddr+0(FP), R11   // gadget addr
    MOVQ gadgetKind+8(FP), CX    // kind (consumed before JMP, before syscall clobbers CX)
    MOVQ nr+16(FP),         AX
    MOVQ a1+24(FP),         DI
    MOVQ a2+32(FP),         SI
    MOVQ a3+40(FP),         DX
    MOVQ a4+48(FP),        R10
    MOVQ a5+56(FP),         R8
    MOVQ a6+64(FP),         R9
    MOVQ (SP),         R13       // R13 = original ret addr to caller [survives syscall]
    LEAQ 8(SP),        R12       // R12 = FP = &args [survives syscall]
    LEAQ ·hhResult<>(SB), R15   // R15 = &hhResult handler
    MOVQ R15, (SP)               // replace [rsp] so gadget ret lands at hhResult
    CMPQ CX, $1
    JEQ  hh_frame
hh_direct:
    JMP R11                      // (syscall;ret) -> ret pops &hhResult -> hhResult
hh_frame:
    PUSHQ $0                     // dummy for pop rbp
    JMP R11                      // (syscall;pop rbp;ret) -> pop rbp eats $0, ret -> hhResult

// hhResult<>: shared return handler for hellHallCall.
// Entered by the gadget's ret instruction (not by a regular CALL).
//   AX  = raw syscall return value
//   R12 = FP of hellHallCall (write r1/r2/errno here)
//   R13 = original return addr to Go caller
TEXT ·hhResult<>(SB),NOSPLIT,$0-0
    CMPQ AX, $0xfffffffffffff001
    JLS  hhok
    MOVQ $-1, 72(R12)            // r1
    MOVQ $0,  80(R12)            // r2
    NEGQ AX
    MOVQ AX,  88(R12)            // errno
    JMP  R13                     // -> Go caller (SP already balanced)
hhok:
    MOVQ AX,  72(R12)            // r1
    MOVQ DX,  80(R12)            // r2
    MOVQ $0,  88(R12)            // errno
    JMP  R13
