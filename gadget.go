package main

import (
	"bufio"
	"debug/elf"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type gadgetKind int

const (
	gadgDirect gadgetKind = iota // syscall; ret          (0f 05 c3)
	gadgFrame                    // syscall; pop rbp; ret (0f 05 5d c3)
)

func (k gadgetKind) String() string {
	switch k {
	case gadgDirect:
		return "syscall;ret"
	case gadgFrame:
		return "syscall;pop_rbp;ret"
	default:
		return "unknown"
	}
}

type mapping struct {
	start, end uintptr
	perms      string
	offset     uintptr
	name       string
}

func parseMaps() ([]mapping, error) {
	f, err := os.Open("/proc/self/maps")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var maps []mapping
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 5 {
			continue
		}
		parts := strings.SplitN(fields[0], "-", 2)
		if len(parts) != 2 {
			continue
		}
		start, err := strconv.ParseUint(parts[0], 16, 64)
		if err != nil {
			continue
		}
		end, err := strconv.ParseUint(parts[1], 16, 64)
		if err != nil {
			continue
		}
		offset, err := strconv.ParseUint(fields[2], 16, 64)
		if err != nil {
			continue
		}
		name := ""
		if len(fields) >= 6 {
			name = fields[5]
		}
		maps = append(maps, mapping{
			start:  uintptr(start),
			end:    uintptr(end),
			perms:  fields[1],
			offset: uintptr(offset),
			name:   name,
		})
	}
	return maps, sc.Err()
}

// findGadgetViaELF parses the libc ELF to find known syscall wrappers,
// then scans them for a syscall gadget.
func findGadgetViaELF(maps []mapping) (addr uintptr, src string, kind gadgetKind, found bool) {
	var libcBase uintptr
	var libcPath string

	for _, m := range maps {
		if strings.Contains(m.name, "libc.so") || strings.Contains(m.name, "libc-") {
			if m.offset == 0 && libcBase == 0 {
				libcBase = m.start
				libcPath = m.name
			}
		}
	}

	if libcPath == "" {
		return 0, "", 0, false
	}

	f, err := elf.Open(libcPath)
	if err != nil {
		return 0, "", 0, false
	}
	defer f.Close()

	syms, err := f.DynamicSymbols()
	if err != nil {
		return 0, "", 0, false
	}

	// Known syscall wrappers in libc
	targets := []string{"getpid", "getuid", "read", "write", "open", "close"}
	
	for _, sym := range syms {
		isTarget := false
		for _, t := range targets {
			if sym.Name == t {
				isTarget = true
				break
			}
		}
		if !isTarget || sym.Value == 0 {
			continue
		}

		funcAddr := libcBase + uintptr(sym.Value)
		// We scan up to 128 bytes of the function body
		p := (*[128]byte)(unsafe.Pointer(funcAddr))
		for i := 0; i < 120; i++ {
			if p[i] == 0x0f && p[i+1] == 0x05 && p[i+2] == 0xc3 {
				return funcAddr + uintptr(i), libcPath + "!" + sym.Name, gadgDirect, true
			}
			if p[i] == 0x0f && p[i+1] == 0x05 && p[i+2] == 0x5d && p[i+3] == 0xc3 {
				return funcAddr + uintptr(i), libcPath + "!" + sym.Name, gadgFrame, true
			}
		}
	}
	return 0, "", 0, false
}

// findGadget scans executable regions for known syscall gadget patterns.
// Priority: ELF parsing of libc > [vdso] > foreign .so > anything else r-xp.
func findGadget() (addr uintptr, src string, kind gadgetKind, err error) {
	maps, err := parseMaps()
	if err != nil {
		return 0, "", 0, err
	}

	// 1. Hell's Gate style: Parse libc ELF
	if a, s, k, ok := findGadgetViaELF(maps); ok {
		return a, s, k, nil
	}

	// 2. Blind scan
	var vdso, so, other *mapping
	for i := range maps {
		m := &maps[i]
		if len(m.perms) < 3 || m.perms[0] != 'r' || m.perms[2] != 'x' {
			continue
		}
		switch {
		case m.name == "[vdso]":
			if vdso == nil {
				vdso = m
			}
		case strings.HasSuffix(m.name, ".so") || strings.Contains(m.name, ".so."):
			if so == nil {
				so = m
			}
		default:
			if other == nil && m.name != "" {
				other = m
			}
		}
	}

	patterns := []struct {
		bytes []byte
		kind  gadgetKind
	}{
		{[]byte{0x0f, 0x05, 0xc3}, gadgDirect},          // syscall; ret
		{[]byte{0x0f, 0x05, 0x5d, 0xc3}, gadgFrame},     // syscall; pop rbp; ret
	}

	for _, m := range []*mapping{vdso, so, other} {
		if m == nil {
			continue
		}
		for _, pat := range patterns {
			if a := scanPattern(m, pat.bytes); a != 0 {
				return a, m.name, pat.kind, nil
			}
		}
	}

	// fallback: any r-xp region
	for i := range maps {
		m := &maps[i]
		if len(m.perms) < 3 || m.perms[0] != 'r' || m.perms[2] != 'x' {
			continue
		}
		for _, pat := range patterns {
			if a := scanPattern(m, pat.bytes); a != 0 {
				return a, m.name, pat.kind, nil
			}
		}
	}
	return 0, "", 0, fmt.Errorf("no syscall gadget found in any r-xp region")
}

func scanPattern(m *mapping, pat []byte) uintptr {
	size := m.end - m.start
	patLen := uintptr(len(pat))
	if size < patLen {
		return 0
	}
	p := (*[1 << 30]byte)(unsafe.Pointer(m.start))
	for i := uintptr(0); i <= size-patLen; i++ {
		match := true
		for j, b := range pat {
			if p[i+uintptr(j)] != b {
				match = false
				break
			}
		}
		if match {
			return m.start + i
		}
	}
	return 0
}
