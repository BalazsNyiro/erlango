/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package erlango

import (
	"fmt"
	"testing"
)

// exec specific test: go test -run Test_simple_structures  *.go

func Test_simple_structures(t *testing.T) {
	txt := `M = #{9 => "nine"}, ID = (1+2)*3, maps:find(ID, M).`

	prg := Prg{callStackDisplay: true}

	chars := ErlSrcChars_from_str(txt)
	chars, _ = ParseErlangSourceCode(prg, chars, "__all__")
	debug_print_ErlSrcChars(chars)

	fmt.Println("Test simple structures - printed only in an error :-)")
	// compare_int_pair("Test simple structures", 2, 1, t)
}
