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

type tokenWanted struct {
	tokenName     string
	positionFirst      int
	positionLast       int
	textRepresentation string
}
type tokensWanted []tokenWanted

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

	wanteds := tokensWanted{
		{"tokenTextBlockQuotedDouble", 180, 185, "~p~n" },
		{"tokenComment", 24, 51, "% testfile for basid types, \n"},
		{"tokenComment", 52, 87, "% whitespaces, commas, dots, colons\n"},
		{"tokenComment", 88, 124, "% atom, string, integer, float, hexa\n"},
		{"tokenComment", 138, 168, "% tab used here as indentation\n"},
		{"tokenTextBlockQuotedDouble", 180, 185, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 213, 218, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 247, 252, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 294, 299, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 335, 340, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 371, 376, "\"~p~n\""},
		{"tokenTextBlockQuotedDouble", 408, 413, "\"~p~n\""},
	}

	for _, tokenWanted := range(wanteds) {
		compare_tokenDetected_tokenWanted( "test basic types", sourceTokensExecutables__whitespacesSeparatorsBasicFile.Tokens, tokenWanted, t)
	}
	// is this token, from this start->end pos, with this representation, is in the reply?
}
func compare_tokenDetected_tokenWanted(callerInfo string, tokensDetected ErlTokens, tokenWanted tokenWanted, t *testing.T) {
	tokenDetected, tokenWantedIsInDetected:= tokensDetected[tokenWanted.positionFirst]

	if tokenWantedIsInDetected {
		// theoretically the charPosFirst is always ok here, because the key in map was the same position
		tokenDetected_charPosFirst, tokenDetected_charPosLast := tokenDetected.charPositionFirstLast()
		if tokenDetected_charPosFirst != tokenWanted.positionFirst {
			t.Fatalf("\nErr %s : detected posFirst: %v  is different from wanted posFirst:  %v, error", callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst)
		}
		if tokenDetected_charPosLast != tokenWanted.positionLast {
			t.Fatalf("\nErr %s : detected posFirst: %v  is different from wanted posFirst:  %v, error", callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst)
		}
	} else {
		t.Fatalf("\nErr %s : wanted tokenPosFirst %v is not in detecteds - error", callerInfo, tokenWanted.positionFirst)
	}

}
