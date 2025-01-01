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
	"erlango.org/erlango/pkg/base_toolset"
	"testing"
)

func Test_file_read(t *testing.T) {
	testName := "first fake test to see if function calling is working"
	erlSrcRunes, _ := base_toolset.File_read_runes("erl_src/erlang_whitespaces_separators_basic_types.erl", "test_comment_string_quotedatom_1")

	charactersInErlSrc := Runes_to_character_structs(erlSrcRunes)
	answer := tokens_detect_comments_strings_quotedatoms(charactersInErlSrc)
	compare_string_string(testName, "org/erlango/pkg/tokens", answer, t)
}
