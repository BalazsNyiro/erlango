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

func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {

	prg := new_program_state()
	prg = cli_argument_detect(prg)  // all arguments are parsed, placed in prg

	fileNameBasic :=  "test/parse/erlang_whitespaces_separators_basic_types.erl"
	fileNameGarden := "test/parse/erlang_language_garden.erl"
	prg = prg_cli_argument_append_from_list(prg, []string{"--files", fileNameBasic, fileNameGarden},  []string{})
	// Erlang_program_exec(prg)


	fileNamesOfErlangSources := filenames_erlang_sources_collect_from_cli_params(prg)
	sourcesTokensExecutables_all := SourcesTokensExecutables_map{}
	sourcesTokensExecutables_all = step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_all, fileNamesOfErlangSources)

	fmt.Println("TEST, tokens check answer length:", len(sourcesTokensExecutables_all))
	for _, sourceTokensExecutables_answer := range(sourcesTokensExecutables_all) {
		fmt.Println("=== TEST, tokens check in: ===", sourceTokensExecutables_answer.PathErlFile)
	}

	sourceTokensExecutables__whitespacesSeparatorsBasicFile := sourcesTokensExecutables_all["test/parse/erlang_whitespaces_separators_basic_types.erl"]

	fmt.Println("TEST, tokens commits, textblocks:")
	sourceTokensExecutables__whitespacesSeparatorsBasicFile.tokens_print()
}
