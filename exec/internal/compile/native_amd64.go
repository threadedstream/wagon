// Copyright 2019 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"bytes"
	"os"
	"os/exec"
)

func debugPrintAsm(asm []byte) {
	cmd := exec.Command("ndisasm", "-b64", "-")
	cmd.Stdin = bytes.NewReader(asm)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func makeExitIndex(idx int) CompletionStatus {
	return CompletionStatus((int64(idx) << 8) & exitIndexMask)
}

const (
	statusMask    = 15
	exitIndexMask = 0x00000000ffffff00
	unknownIndex  = 0xffffff
)

// JITExitSignal is the value returned from the execution of a native section.
// The bits of this packed 64bit value is encoded as follows:
// [00:03] Completion Status
// [04:07] Reserved
// [08:31] Index of the WASM instruction where the exit occurred.
// [31:63] Status-specific 32bit value.
type JITExitSignal uint64

// CompletionStatus decodes and returns the completion status of the exit.
func (s JITExitSignal) CompletionStatus() CompletionStatus {
	return CompletionStatus(s & statusMask)
}

// Index returns the index to the instruction where the exit happened.
// 0xffffff is returned if the exit was due to normal completion.
func (s JITExitSignal) Index() int {
	return (int(s) & exitIndexMask) >> 8
}
