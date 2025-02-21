//go:build wasm || wasip1

package xruntime

// For WebAssembly, we assume a default cache line size.
// This value is a fallback and may not reflect any real hardware cache.
const CacheLineSize = 64
