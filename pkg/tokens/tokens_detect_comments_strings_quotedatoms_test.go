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

func Test_tokens_detect_in_erl_src(t *testing.T) {
	testName := "Test_tokens_detect_in_erl_src"
	erlSrcRunes, _ := base_toolset.File_read_runes("erl_src/erlang_whitespaces_separators_basic_types.erl", "test_comment_string_quotedatom_1")

	charactersInErlSrc := Runes_to_character_structs(erlSrcRunes)
	tokens := TokenCollector{}

	charactersInErlSrc, tokens = Tokens_detect_in_erl_src(charactersInErlSrc, tokens)

	compare_string_string(testName, "fakeTest", "fakeTest", t)
}
