package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type patch struct {
	offset   int64
	original []byte
	patched  []byte
}

var patches = []patch{
	{
		0x0002A7,
		[]byte{0xE0},
		[]byte{0xC0},
	},
	{
		0x14A842,
		[]byte{0x3B, 0x3D, 0xB0, 0xA5, 0xBE, 0x00},
		[]byte{0x83, 0xFF, 0x50, 0x90, 0x90, 0x90},
	},
	{
		0x14A85C,
		[]byte{0x83, 0xF8, 0x02, 0x7D, 0x1F},
		[]byte{0x83, 0xF8, 0x50, 0x7F, 0x5F},
	},
	{
		0x364D40,
		[]byte{0x74},
		[]byte{0xEB},
	},
	{
		0x46B8E0,
		[]byte{0x55, 0x8B, 0xEC},
		[]byte{0x33, 0xC0, 0xC3},
	},
	{
		0x0614A0,
		[]byte{0x55, 0x8B, 0xEC},
		[]byte{0x33, 0xC0, 0xC3},
	},
	{
		0x46AC43,
		[]byte{0x76, 0x0C},
		[]byte{0x90, 0x90},
	},
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: patcher /path/to/Wow.exe")
	}

	path := os.Args[1]

	wow, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read Wow.exe: %v\n", err)
	}

	changed := false
	for i, p := range patches {
		applied, err := applyPatch(wow, p)
		if err != nil {
			log.Fatalf("patch %d: %v\n", i, err)
		}
		if applied {
			changed = true
			fmt.Printf("patch %d applied\n", i)
		} else {
			fmt.Printf("patch %d already patched\n", i)
		}
	}

	if !changed {
		fmt.Println("Wow.exe already fully patched, nothing to do")
		return
	}

	if err := backup(path); err != nil {
		log.Fatalf("Failed to back up Wow.exe: %v\n", err)
	}

	if err := os.WriteFile(path, wow, 0644); err != nil {
		log.Fatalf("Failed to write Wow.exe: %v\n", err)
	}

	fmt.Println("Wow.exe was patched")
}

func applyPatch(wow []byte, p patch) (bool, error) {
	end := p.offset + int64(len(p.original))
	if end > int64(len(wow)) {
		return false, fmt.Errorf("offset 0x%X out of range (file is %d bytes)", p.offset, len(wow))
	}

	cur := wow[p.offset:end]
	switch {
	case bytes.Equal(cur, p.original):
		copy(cur, p.patched)
		return true, nil
	case bytes.Equal(cur, p.patched):
		return false, nil
	default:
		return false, fmt.Errorf("unexpected bytes at 0x%X: got % X, want original % X or patched % X",
			p.offset, cur, p.original, p.patched)
	}
}

func backup(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	newPath := strings.Replace(path, "Wow.exe", "Wow.exe.backup", 1)

	dst, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
