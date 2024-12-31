/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import (
	"testing"
)

func Test_file_read(t *testing.T) {
	testName := "first fake test to see if function calling is working"
	answer := tokens_detect_comments_strings_quotedatoms("something")
	compare_string_string(testName, "org/erlango/pkg/tokens", answer, t)
}
