# whisper

Indirect syscalls for Linux using `syscall; ret` gadgets discovered at runtime, written in Go and Plan9 assembly.

## Methods

### Hell's Gate
Parses the system `libc` ELF on disk (via `debug/elf`) and scans known syscall wrapper symbols (`getpid`, `read`, `write`, ...) to locate a `syscall; ret` or `syscall; pop rbp; ret` gadget. Preferred strategy.

### Halo's Gate
Fallback when ELF parsing fails. Blind-scans the `[vdso]` and all other `r-xp` mapped regions from `/proc/self/maps` for the same gadget patterns.

### HellsHall
JMP-based variant. Instead of `CALL gadget` (which leaves whisper on the call stack), the trampoline patches `[rsp]` with a result handler then `JMP`s to the gadget. The gadget's `ret` lands directly in the caller - whisper is invisible in the call stack during the syscall.

## API

```go
// Hell's Gate / Halo's Gate - CALL-based
Syscall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
Syscall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno)

// HellsHall - JMP-based, whisper absent from call stack
HellsHall3(nr, a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
HellsHall6(nr, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, syscall.Errno)
```

## References

- Inspired by [Hell's Gate](https://github.com/am0nsec/HellsGate)
- Inspired by [Halo's Gate](https://blog.rebelit.net/?p=163)
- Inspired by [HellsHall](https://github.com/Maldev-Academy/HellHall)
