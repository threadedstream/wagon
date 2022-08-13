//go:build !appengine
// +build !appengine

package compile

import "unsafe"

func jitcall(_ unsafe.Pointer, _, _, _ *[]uint64, _ *[]byte) uint64 {
	panic("not implemented")
}
