//go:build linux || darwin
// +build linux darwin

package clistats

import "syscall"

func kill() {
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
