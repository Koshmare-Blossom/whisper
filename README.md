# whisper

Indirect syscalls for Linux using `syscall; ret` gadgets discovered at runtime via `/proc/self/maps` and ELF parsing of libc, written in Go and Assembly.

Equivalent of **Hell's Gate** / **Halo's Gate** for Linux.

## How it works

1.  **Map Discovery**: Parses `/proc/self/maps` to find the memory regions for `[vdso]` and the system `libc`.
2.  **ELF Parsing**: Opens the `libc` binary on disk and parses its dynamic symbol table (using `debug/elf`) to find known syscall wrappers (e.g., `getpid`, `read`, `write`).
3.  **Gadget Extraction**: Scans the machine code of these wrappers to find `syscall; ret` or `syscall; pop rbp; ret` sequences.
4.  **Indirect Invocation**: Executes syscalls by jumping to these gadgets via an assembly trampoline. This ensures the `syscall` instruction is executed from within a legitimate library's memory space, bypassing simple static analysis or instruction-pointer-based monitoring.

## Features

- **Dynamic Gadget Hunting**: No hardcoded offsets; works across different libc versions and distributions.
- **Stack-Safe Trampoline**: Correctly handles gadgets that modify the stack (e.g., `pop rbp`).
- **Minimal Footprint**: Core implementation in ~3 Go files.

## References

*   Inspired by [Hell's Gate](https://github.com/am0nsec/HellsGate), a technique to bypass EDRs by dynamically retrieving syscall numbers.
*   Inspired by [Halo's Gate](https://blog.rebelit.net/?p=163), an evolution of Hell's Gate that bypasses EDR hooks by scanning for gadgets.
