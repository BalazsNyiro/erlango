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
	"fmt"
	"testing"
)

func Test_char_block_detect_in_external_erlang_projects(t *testing.T) {
	testName := "Test_char_block_detect_in_external_erlang_projects"

	// erlSrcRunes := []rune(`External  = 1. `)
	fileErl := "../../../erlang_projects_external_sources/rebar3_all.erl"
	charactersInErlSrc, errors := Character_block_detect_in_erl_file("localhost", fileErl, testName+"_rebar3_all")
	// Tokens_detection_print_verbose(charactersInErlSrc, tokensInErlSrc)
	Char_block_detection_print_one_char_per_line(charactersInErlSrc, true)

	fmt.Println("errors in character detection:", errors)
}
