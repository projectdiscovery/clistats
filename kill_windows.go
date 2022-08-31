//go:build windows
// +build windows

package clistats

import (
	"os"
)

func kill() {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}
	_ = p.Signal(os.Kill)
}
