//go:build !wasm && !wasip1

package xruntime

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

const (
	// CacheLineSize is useful for preventing false sharing.
	CacheLineSize = unsafe.Sizeof(cpu.CacheLinePad{})
)
