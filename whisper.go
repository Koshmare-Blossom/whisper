package main

// whisper - indirect syscalls for Linux.
//
// Locates a "syscall; ret" gadget in an external mapped region (vDSO preferred)
// and dispatches syscalls through it via an assembly trampoline.
// The syscall instruction executes in the gadget's memory, not in the Go binary's .text.
//
// Equivalent of Hell's Gate / Halo's Gate for Linux.

import "syscall"

// callGadget: CALL-based trampoline (Hell's Gate / Halo's Gate).
// whisper's code appears in the call stack during the syscall.
func callGadget(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
func setup_stack()

// hellsHallCall: JMP-based trampoline (HellsHall).
// The gadget's ret returns directly to the Go caller via hhResult<>,
// so whisper is absent from the call stack during the syscall.
func hellsHallCall(gadgetAddr, gadgetKind, nr, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)

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

// Syscall3/Syscall6: Hell's Gate mode (CALL-based).
func Syscall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := callGadget(gadget, gKind, nr, a1, a2, a3, 0, 0, 0)
	return r1, r2, syscall.Errno(e)
}

func Syscall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := callGadget(gadget, gKind, nr, a1, a2, a3, a4, a5, a6)
	return r1, r2, syscall.Errno(e)
}

// HellsHall3/HellsHall6: HellsHall mode (JMP-based, whisper absent from call stack).
func HellsHall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := hellsHallCall(gadget, gKind, nr, a1, a2, a3, 0, 0, 0)
	return r1, r2, syscall.Errno(e)
}

func HellsHall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno) {
	r1, r2, e := hellsHallCall(gadget, gKind, nr, a1, a2, a3, a4, a5, a6)
	return r1, r2, syscall.Errno(e)
}
