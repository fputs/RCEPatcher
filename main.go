package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	patchOffset int64 = 0x2A7
	patchValue  byte  = 0xC0
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) != 2 {
		log.Fatal("Usage: patcher /path/to/Wow.exe")
	}

	path := os.Args[1]

	_, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to find Wow.exe: %v\n", err)
	}

	patched, err := isPatched(path)
	if err != nil {
		log.Fatalf("Failed to check if patched: %v\n", err)
	}
	if patched {
		fmt.Println("Already patched!")
		os.Exit(0)
	}

	err = patch(path)
	if err != nil {
		log.Fatalf("Failed to patch: %v\n", err)
	}

	fmt.Println("Wow.exe was patched")
}

func isPatched(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Seek(patchOffset, 0)
	if err != nil {
		return false, err
	}

	b := make([]byte, 1)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}

	if b[0] == patchValue {
		return true, nil
	}
	return false, nil
}

func patch(path string) error {
	err := backup(path)
	if err != nil {
		return fmt.Errorf("backup error: %v", err)
	}
	wow, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	wow[patchOffset] = patchValue

	err = os.WriteFile(path, wow, 0644)
	if err != nil {
		return err
	}
	return nil
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
