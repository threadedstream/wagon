package compile

import (
	"bytes"
	"os"
	"os/exec"
)

// TODO(gildarov): to be honest, I have absolutely no clue what these masks refer to, so I decided to
// copy code from native_amd64, since it saves me time figuring out how to compile this project

func debugPrintAsm(asm []byte) {
	cmd := exec.Command("ndisasm", "-b64", "-")
	cmd.Stdin = bytes.NewReader(asm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
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

func (s JITExitSignal) CompletionStatus() CompletionStatus {
	return CompletionStatus(s & statusMask)
}

func (s JITExitSignal) Index() int {
	return (int(s) & exitIndexMask) >> 8
}
