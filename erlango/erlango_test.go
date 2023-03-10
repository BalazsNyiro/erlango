/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in erlango repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package erlango

import "testing"

func Test_erlango_parse(t *testing.T) {
	received := ParseErlangSourceFile()
	wanted := 0
	compare_int_pair(received, wanted, t)
}
func compare_int_pair(received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
	}
}
