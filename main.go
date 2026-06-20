package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[1;36m[*]\033[0m "+format+"\n", args...)
}

func okf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[1;32m[+]\033[0m "+format+"\n", args...)
}

func errf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\033[1;31m[-]\033[0m "+format+"\n", args...)
}

func main() {
	if err := initWhisper(); err != nil {
		errf("init: %v", err)
		os.Exit(1)
	}

	// getpid
	pid, _, errno := Syscall3(syscall.SYS_GETPID, 0, 0, 0)
	if errno != 0 {
		errf("getpid: %v", errno)
		os.Exit(1)
	}
	okf("SYS_GETPID  -> %d  (os.Getpid()=%d)", pid, os.Getpid())

	// getuid
	uid, _, errno := Syscall3(syscall.SYS_GETUID, 0, 0, 0)
	if errno != 0 {
		errf("getuid: %v", errno)
		os.Exit(1)
	}
	okf("SYS_GETUID  -> %d  (os.Getuid()=%d)", uid, os.Getuid())

	// write to stdout
	msg := []byte("whisper: syscall via vDSO gadget\n")
	n, _, errno := Syscall3(syscall.SYS_WRITE,
		uintptr(syscall.Stdout),
		uintptr(unsafe.Pointer(&msg[0])),
		uintptr(len(msg)))
	if errno != 0 {
		errf("write: %v", errno)
		os.Exit(1)
	}
	okf("SYS_WRITE   -> %d bytes", n)

	// open + read + close /proc/self/comm
	path, _ := syscall.BytePtrFromString("/proc/self/comm")
	fd, _, errno := Syscall3(syscall.SYS_OPEN,
		uintptr(unsafe.Pointer(path)),
		syscall.O_RDONLY, 0)
	if errno != 0 {
		errf("open: %v", errno)
		os.Exit(1)
	}
	buf := make([]byte, 32)
	nr, _, errno := Syscall3(syscall.SYS_READ,
		fd,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)))
	Syscall3(syscall.SYS_CLOSE, fd, 0, 0)
	if errno != 0 {
		errf("read: %v", errno)
		os.Exit(1)
	}
	okf("SYS_READ    /proc/self/comm -> %q", string(buf[:nr]))

	// mmap anonymous page
	addr, _, errno := Syscall6(syscall.SYS_MMAP, 0, 4096,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
		^uintptr(0), 0)
	if errno != 0 {
		errf("mmap: %v", errno)
		os.Exit(1)
	}
	okf("SYS_MMAP    -> %#x", addr)
	Syscall3(syscall.SYS_MUNMAP, addr, 4096, 0)

	okf("whisper demo OK")

	// ── HellsHall mode (JMP-based) ─────────────────────────────────────────
	logf("--- HellsHall (JMP) ---")

	pid2, _, errno := HellsHall3(syscall.SYS_GETPID, 0, 0, 0)
	if errno != 0 {
		errf("HellsHall getpid: %v", errno)
		os.Exit(1)
	}
	okf("SYS_GETPID  -> %d", pid2)

	uid2, _, errno := HellsHall3(syscall.SYS_GETUID, 0, 0, 0)
	if errno != 0 {
		errf("HellsHall getuid: %v", errno)
		os.Exit(1)
	}
	okf("SYS_GETUID  -> %d", uid2)

	addr2, _, errno := HellsHall6(syscall.SYS_MMAP, 0, 4096,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
		^uintptr(0), 0)
	if errno != 0 {
		errf("HellsHall mmap: %v", errno)
		os.Exit(1)
	}
	okf("SYS_MMAP    -> %#x", addr2)
	HellsHall3(syscall.SYS_MUNMAP, addr2, 4096, 0)

	okf("hellshall demo OK")
}
