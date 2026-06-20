# whisper

Indirect syscalls for Linux. Locates `syscall; ret` gadgets in `[vdso]` or `libc` at runtime and dispatches syscalls through them - the `syscall` instruction always executes from a legitimate mapped region, never from the Go binary's `.text`.

Written in Go and Plan9 assembly, zero CGO.

## Gadget discovery

### Hell's Gate strategy
Parses the system `libc` ELF on disk (via `debug/elf`) and scans known syscall wrapper symbols (`getpid`, `read`, `write`, ...) for a `syscall; ret` or `syscall; pop rbp; ret` gadget. Preferred.

### Halo's Gate strategy
Fallback when ELF parsing fails. Blind-scans `[vdso]` then all `r-xp` regions from `/proc/self/maps` for the same patterns.

## Invocation methods

### Whisper (CALL-based)
Standard indirect syscall. `CALL gadget` - whisper's trampoline appears in the call stack during the syscall.

### HellsHall (JMP-based)
Patches `[rsp]` with a result handler then `JMP`s to the gadget. The gadget's `ret` returns directly to the Go caller - whisper is invisible in the call stack during the syscall.

## API

```go
// Whisper - CALL-based indirect syscall
Syscall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
Syscall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno)

// HellsHall - JMP-based, trampoline absent from call stack
HellsHall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
HellsHall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno)
```

## References

- Inspired by [Hell's Gate](https://github.com/am0nsec/HellsGate)
- Inspired by [Halo's Gate](https://github.com/boku7/AsmHalosGate)
- Inspired by [HellsHall](https://github.com/Maldev-Academy/HellHall)
