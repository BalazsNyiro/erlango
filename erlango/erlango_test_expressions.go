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



// go test -run Test_expression_detection
// expressions are in focus here:
func Test_expression_detection(t *testing.T) {
	prg := new_program_state(true)
	prg = cli_argument_detect(prg)  // all arguments are parsed, placed in prg

	fileNameBasic :=  "test/parse/erlang_whitespaces_separators_basic_types.erl"
	prg = prg_cli_argument_append_from_list(prg, []string{"--files", fileNameBasic},  []string{})

	fileNamesOfErlangSources := filenames_erlang_sources_collect_from_cli_params(prg)

	// sourcesTokensExecutables_all can be empty (like here), or it can have existing elements - in a running system new src can be loaded, next to the existing ones
	sourcesTokensExecutables_all := SourcesTokensExecutables_map{}
	sourcesTokensExecutables_all = step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_all, fileNamesOfErlangSources, true)

	step_02_expressions_from_tokens(sourcesTokensExecutables_all, fileNamesOfErlangSources)
}

