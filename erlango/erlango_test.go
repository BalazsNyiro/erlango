/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package erlango

import (
	"testing"
)

func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {

	prg := new_program_state()
	prg = cli_argument_detect(prg)  // all arguments are parsed, placed in prg

	fileName :=  "test/parse/erlang_whitespaces_separators_basic_types.erl"
	prg = prg_cli_argument_append_from_list(prg, []string{"--files", fileName},  []string{})
	// Erlang_program_exec(prg)


	fileNamesOfErlangSources := filenames_erlang_sources_collect_from_cli_params(prg)
	sourcesTokensExecutables_list := SourcesTokensExecutables_list{}
	sourcesTokensExecutables_list = step_01_tokens_from_source_code(sourcesTokensExecutables_list, fileNamesOfErlangSources)
}
