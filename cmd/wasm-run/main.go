// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
)

func main() {
	log.SetPrefix("wasm-run: ")
	log.SetFlags(0)

	verbose := flag.Bool("v", false, "enable/disable verbose mode")
	verify := flag.Bool("verify-module", false, "run module verification")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	wasm.SetDebugMode(*verbose)

	run(os.Stdout, flag.Arg(0), *verify)
}

func run(w io.Writer, fname string, verify bool) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m, err := wasm.ReadModule(f, importer)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	if verify {
		err = validate.VerifyModule(m)
		if err != nil {
			log.Fatalf("could not verify module: %v", err)
		}
	}

	if m.Export == nil {
		log.Fatalf("module has no export section")
	}

	vm, err := exec.NewVM(m)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	for name, e := range m.Export.Entries {
		i := int64(e.Index)
		fidx := m.Function.Types[int(i)]
		ftype := m.Types.Entries[int(fidx)]
		switch len(ftype.ReturnTypes) {
		case 1:
			fmt.Fprintf(w, "%s() %s => ", name, ftype.ReturnTypes[0])
		case 0:
			fmt.Fprintf(w, "%s() => ", name)
		default:
			log.Printf("running exported functions with more than one return value is not supported")
			continue
		}
		var args []uint64
		if len(ftype.ParamTypes) > 0 {
			log.Printf("enter input parameters for function %s", ftype.String())
			for i := range ftype.ParamTypes {
				arg := uint64(0)
				_, err = fmt.Scan(&arg)
				if err != nil {
					log.Printf("failed to get %dth parameter", i)
				}
				args = append(args, arg)
			}
		}
		o, err := vm.ExecCode(i, args...)
		if err != nil {
			fmt.Fprintf(w, "\n")
			log.Printf("err=%v", err)
			continue
		}
		if len(ftype.ReturnTypes) == 0 {
			fmt.Fprintf(w, "\n")
			continue
		}
		fmt.Fprintf(w, "%[1]v (%[1]T)\n", o)
	}
}

func importer(name string) (*wasm.Module, error) {
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, err := wasm.ReadModule(f, nil)
	if err != nil {
		return nil, err
	}
	err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
