package main

// whisper - indirect syscalls for Linux.
//
// Locates a "syscall; ret" gadget in an external mapped region (vDSO preferred)
// and dispatches syscalls through it via an assembly trampoline.
// The syscall instruction executes in the gadget's memory, not in the Go binary's .text.
//
// Equivalent of Hell's Gate / Halo's Gate for Linux.

import "syscall"

// callGadget is implemented in asm_amd64.s.
// Sets up the Linux syscall ABI (rax=nr, rdi/rsi/rdx/r10/r8/r9=args)
// then jumps to gadgetAddr (syscall; ret) via CALL.
func callGadget(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
func setup_stack()

var gadget uintptr
var gKind uintptr

func initWhisper() error {
	addr, src, kind, err := findGadget()
	if err != nil {
		return err
	}
	gadget = addr
	gKind = uintptr(kind)
	okf("gadget @ %#x  src=%s kind=%s", addr, src, kind)
	return nil
}

func Syscall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := callGadget(gadget, gKind, nr, a1, a2, a3, 0, 0, 0)
	return r1, r2, syscall.Errno(e)
}

func Syscall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := callGadget(gadget, gKind, nr, a1, a2, a3, a4, a5, a6)
	return r1, r2, syscall.Errno(e)
}
