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

// exec specific test: go test -run Test_Simple_structures  *.go

func Test_simple_structures(t *testing.T) {
	txt := `(1, 2)`
	chars := ErlSrcChars_from_str(txt)
	ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(chars, false)
	debug_print_ErlSrcChars(chars)

	fmt.Println("Test simple structures")
	compare_int_pair("Test simple structures", 1, 2, t)
}
