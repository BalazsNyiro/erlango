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

func Test_tokens_detect_in_external_erlang_projects(t *testing.T) {
	testName := "Test_tokens_detect_in_external_erlang_projects"

	// erlSrcRunes := []rune(`External  = 1. `)
	erlSrcRunes, _ := base_toolset.File_read_runes("../../../erlang_projects_external_sources/rebar3_all.erl", testName+"_rebar3_all")
	// erlSrcRunes, _ := base_toolset.File_read_runes("../../../erlang_projects_1external_sources/rebar3_one_problem.erl", testName+"_rebar3_all")

	charactersInErlSrc := Runes_to_character_structs(erlSrcRunes)
	tokensInErlSrc := TokenCollector{}

	charactersInErlSrc, tokensInErlSrc = Tokens_detect_in_erl_src(charactersInErlSrc, tokensInErlSrc)

	// Tokens_detection_print_verbose(charactersInErlSrc, tokensInErlSrc)
	Tokens_detection_print_one_char_per_line(charactersInErlSrc, tokensInErlSrc, true)

}
