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


// tokens are in focus here:
func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {

	prg := new_program_state(true)
	prg = cli_argument_detect(prg)  // all arguments are parsed, placed in prg

	fileNameBasic :=  "test/parse/erlang_whitespaces_separators_basic_types.erl"
	fileNameGarden := "test/parse/erlang_language_garden.erl"
	prg = prg_cli_argument_append_from_list(prg, []string{"--files", fileNameBasic, fileNameGarden},  []string{})

	fileNamesOfErlangSources := filenames_erlang_sources_collect_from_cli_params(prg)
	sourcesTokensExecutables_all := SourcesTokensExecutables_map{}
	sourcesTokensExecutables_all = step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_all, fileNamesOfErlangSources, true)

	fmt.Println("TEST, tokens check answer length:", len(sourcesTokensExecutables_all))
	for _, sourceTokensExecutables_answer := range(sourcesTokensExecutables_all) {
		fmt.Println("=== TEST, tokens check in: ===", sourceTokensExecutables_answer.WhereTheCodeIsStored)
	}

	sourceTokensExecutables__whitespacesSeparatorsBasicFile := sourcesTokensExecutables_all["test/parse/erlang_whitespaces_separators_basic_types.erl"]

	fmt.Println("TEST, tokens commits, textblocks:")
	sourceTokensExecutables__whitespacesSeparatorsBasicFile.tokens_print()

	wanteds := tokensWanted{


		{"tokenOtherPunctuation", 0, 0, "#"},
		{"tokenOtherPunctuation", 1, 1, "!"},
		{"tokenOtherPunctuation", 2, 2, "/"},
		{"tokenAbcFullWith_Underscore_At_numbers", 3, 5, "usr"},
		{"tokenOtherPunctuation", 6, 6, "/"},
		{"tokenAbcFullWith_Underscore_At_numbers", 7, 9, "bin"},
		{"tokenOtherPunctuation", 10, 10, "/"},
		{"tokenAbcFullWith_Underscore_At_numbers", 11, 13, "env"},
		{"tokenAbcFullWith_Underscore_At_numbers", 15, 21, "escript"},
		{"tokenAbcFullWith_Underscore_At_numbers", 127, 130, "main"},
		{"tokenOtherPunctuation", 131, 131, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 132, 132, "_"},
		{"tokenOtherPunctuation", 133, 133, ")"},
		{"tokenOtherPunctuation", 135, 135, "-"},
		{"tokenOtherPunctuation", 136, 136, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 171, 172, "io"},
		{"tokenOtherPunctuation", 173, 173, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 174, 179, "fwrite"},
		{"tokenOtherPunctuation", 180, 180, "("},
		{"tokenTextBlockQuotedDouble", 181, 186, "\"~p~n\""},
		{"tokenOtherPunctuation", 187, 187, ","},
		{"tokenOtherPunctuation", 189, 189, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 190, 196, "example"},
		{"tokenOtherPunctuation", 197, 197, "("},
		{"tokenOtherPunctuation", 198, 198, ")"},
		{"tokenOtherPunctuation", 199, 199, "]"},
		{"tokenOtherPunctuation", 200, 200, ")"},
		{"tokenOtherPunctuation", 201, 201, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 204, 205, "io"},
		{"tokenOtherPunctuation", 206, 206, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 207, 212, "fwrite"},
		{"tokenOtherPunctuation", 213, 213, "("},
		{"tokenTextBlockQuotedDouble", 214, 219, "\"~p~n\""},
		{"tokenOtherPunctuation", 220, 220, ","},
		{"tokenOtherPunctuation", 222, 222, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 223, 229, "example"},
		{"tokenOtherPunctuation", 230, 230, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 231, 231, "1"},
		{"tokenOtherPunctuation", 232, 232, ")"},
		{"tokenOtherPunctuation", 233, 233, "]"},
		{"tokenOtherPunctuation", 234, 234, ")"},
		{"tokenOtherPunctuation", 235, 235, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 238, 239, "io"},
		{"tokenOtherPunctuation", 240, 240, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 241, 246, "fwrite"},
		{"tokenOtherPunctuation", 247, 247, "("},
		{"tokenTextBlockQuotedDouble", 248, 253, "\"~p~n\""},
		{"tokenOtherPunctuation", 254, 254, ","},
		{"tokenOtherPunctuation", 256, 256, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 257, 263, "example"},
		{"tokenOtherPunctuation", 264, 264, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 265, 274, "1234567890"},
		{"tokenOtherPunctuation", 275, 275, ")"},
		{"tokenOtherPunctuation", 276, 276, "]"},
		{"tokenOtherPunctuation", 277, 277, ")"},
		{"tokenOtherPunctuation", 278, 278, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 285, 286, "io"},
		{"tokenOtherPunctuation", 287, 287, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 288, 293, "fwrite"},
		{"tokenOtherPunctuation", 294, 294, "("},
		{"tokenTextBlockQuotedDouble", 295, 300, "\"~p~n\""},
		{"tokenOtherPunctuation", 301, 301, ","},
		{"tokenOtherPunctuation", 303, 303, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 304, 310, "example"},
		{"tokenOtherPunctuation", 311, 311, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 312, 313, "12"},
		{"tokenOtherPunctuation", 314, 314, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 315, 316, "34"},
		{"tokenOtherPunctuation", 317, 317, ")"},
		{"tokenOtherPunctuation", 318, 318, "]"},
		{"tokenOtherPunctuation", 319, 319, ")"},
		{"tokenOtherPunctuation", 320, 320, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 326, 327, "io"},
		{"tokenOtherPunctuation", 328, 328, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 329, 334, "fwrite"},
		{"tokenOtherPunctuation", 335, 335, "("},
		{"tokenTextBlockQuotedDouble", 336, 341, "\"~p~n\""},
		{"tokenOtherPunctuation", 342, 342, ","},
		{"tokenOtherPunctuation", 344, 344, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 345, 347, "add"},
		{"tokenOtherPunctuation", 348, 348, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 349, 349, "2"},
		{"tokenOtherPunctuation", 350, 350, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 352, 352, "4"},
		{"tokenOtherPunctuation", 353, 353, ")"},
		{"tokenOtherPunctuation", 354, 354, "]"},
		{"tokenOtherPunctuation", 355, 355, ")"},
		{"tokenOtherPunctuation", 356, 356, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 362, 363, "io"},
		{"tokenOtherPunctuation", 364, 364, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 365, 370, "fwrite"},
		{"tokenOtherPunctuation", 371, 371, "("},
		{"tokenTextBlockQuotedDouble", 372, 377, "\"~p~n\""},
		{"tokenOtherPunctuation", 378, 378, ","},
		{"tokenOtherPunctuation", 380, 380, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 381, 386, "double"},
		{"tokenOtherPunctuation", 387, 387, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 388, 388, "9"},
		{"tokenOtherPunctuation", 389, 389, ")"},
		{"tokenOtherPunctuation", 390, 390, "]"},
		{"tokenOtherPunctuation", 391, 391, ")"},
		{"tokenOtherPunctuation", 392, 392, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 399, 400, "io"},
		{"tokenOtherPunctuation", 401, 401, ":"},
		{"tokenAbcFullWith_Underscore_At_numbers", 402, 407, "fwrite"},
		{"tokenOtherPunctuation", 408, 408, "("},
		{"tokenTextBlockQuotedDouble", 409, 414, "\"~p~n\""},
		{"tokenOtherPunctuation", 415, 415, ","},
		{"tokenOtherPunctuation", 417, 417, "["},
		{"tokenAbcFullWith_Underscore_At_numbers", 418, 421, "half"},
		{"tokenOtherPunctuation", 422, 422, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 423, 424, "10"},
		{"tokenOtherPunctuation", 425, 425, ")"},
		{"tokenOtherPunctuation", 426, 426, "]"},
		{"tokenOtherPunctuation", 427, 427, ")"},
		{"tokenOtherPunctuation", 428, 428, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 435, 436, "ok"},
		{"tokenOtherPunctuation", 437, 437, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 440, 446, "example"},
		{"tokenOtherPunctuation", 447, 447, "("},
		{"tokenOtherPunctuation", 448, 448, ")"},
		{"tokenOtherPunctuation", 450, 450, "-"},
		{"tokenOtherPunctuation", 451, 451, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 457, 463, "example"},
		{"tokenOtherPunctuation", 464, 464, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 526, 532, "example"},
		{"tokenOtherPunctuation", 533, 533, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 534, 534, "1"},
		{"tokenOtherPunctuation", 535, 535, ")"},
		{"tokenOtherPunctuation", 537, 537, "-"},
		{"tokenOtherPunctuation", 538, 538, ">"},
		{"tokenTextBlockQuotedDouble", 540, 586, "\"case 1 \\\\\\\" complex string \\\" with \\n newline\""},
		{"tokenOtherPunctuation", 587, 587, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 641, 647, "example"},
		{"tokenOtherPunctuation", 648, 648, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 649, 658, "1234567890"},
		{"tokenOtherPunctuation", 659, 659, ")"},
		{"tokenOtherPunctuation", 661, 661, "-"},
		{"tokenOtherPunctuation", 662, 662, ">"},
		{"tokenTextBlockQuotedDouble", 664, 680, "\"case 1234567890\""},
		{"tokenOtherPunctuation", 681, 681, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 683, 689, "example"},
		{"tokenOtherPunctuation", 690, 690, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 691, 692, "16"},
		{"tokenOtherPunctuation", 693, 693, "#"},
		{"tokenAbcFullWith_Underscore_At_numbers", 694, 701, "af6bfa23"},
		{"tokenOtherPunctuation", 702, 702, ")"},
		{"tokenOtherPunctuation", 704, 704, "-"},
		{"tokenOtherPunctuation", 705, 705, ">"},
		{"tokenTextBlockQuotedDouble", 707, 716, "\"hexa num\""},
		{"tokenOtherPunctuation", 717, 717, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 749, 755, "example"},
		{"tokenOtherPunctuation", 756, 756, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 757, 758, "12"},
		{"tokenOtherPunctuation", 759, 759, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 760, 761, "34"},
		{"tokenOtherPunctuation", 762, 762, ")"},
		{"tokenOtherPunctuation", 764, 764, "-"},
		{"tokenOtherPunctuation", 765, 765, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 767, 768, "12"},
		{"tokenOtherPunctuation", 769, 769, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 770, 771, "34"},
		{"tokenOtherPunctuation", 772, 772, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 804, 810, "example"},
		{"tokenOtherPunctuation", 811, 811, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 812, 822, "atom_direct"},
		{"tokenOtherPunctuation", 823, 823, ")"},
		{"tokenOtherPunctuation", 825, 825, "-"},
		{"tokenOtherPunctuation", 826, 826, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 828, 838, "atom_direct"},
		{"tokenOtherPunctuation", 839, 839, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 841, 847, "example"},
		{"tokenOtherPunctuation", 848, 848, "("},
		{"tokenTextBlockQuotedSingle", 849, 874, "'text_block_single_quoted'"},
		{"tokenOtherPunctuation", 875, 875, ")"},
		{"tokenOtherPunctuation", 877, 877, "-"},
		{"tokenOtherPunctuation", 878, 878, ">"},
		{"tokenTextBlockQuotedSingle", 880, 941, "'text_block_single_quoted_reply with \\' escape and \\n newline'"},
		{"tokenOtherPunctuation", 942, 942, ";"},
		{"tokenAbcFullWith_Underscore_At_numbers", 944, 950, "example"},
		{"tokenOtherPunctuation", 951, 951, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 952, 952, "_"},
		{"tokenOtherPunctuation", 953, 953, ")"},
		{"tokenOtherPunctuation", 955, 955, "-"},
		{"tokenOtherPunctuation", 956, 956, ">"},
		{"tokenTextBlockQuotedDouble", 958, 970, "\"case others\""},
		{"tokenOtherPunctuation", 971, 971, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 974, 976, "add"},
		{"tokenOtherPunctuation", 977, 977, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 978, 978, "A"},
		{"tokenOtherPunctuation", 979, 979, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 981, 981, "B"},
		{"tokenOtherPunctuation", 982, 982, ")"},
		{"tokenOtherPunctuation", 984, 984, "-"},
		{"tokenOtherPunctuation", 985, 985, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 991, 996, "Result"},
		{"tokenOtherPunctuation", 998, 998, "="},
		{"tokenOtherPunctuation", 1000, 1000, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 1001, 1001, "A"},
		{"tokenOtherPunctuation", 1003, 1003, "+"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1005, 1005, "B"},
		{"tokenOtherPunctuation", 1006, 1006, ")"},
		{"tokenOtherPunctuation", 1007, 1007, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 1013, 1018, "Result"},
		{"tokenOtherPunctuation", 1019, 1019, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 1022, 1025, "diff"},
		{"tokenOtherPunctuation", 1026, 1026, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 1027, 1027, "A"},
		{"tokenOtherPunctuation", 1028, 1028, ","},
		{"tokenAbcFullWith_Underscore_At_numbers", 1030, 1030, "B"},
		{"tokenOtherPunctuation", 1031, 1031, ")"},
		{"tokenOtherPunctuation", 1033, 1033, "-"},
		{"tokenOtherPunctuation", 1034, 1034, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1036, 1036, "A"},
		{"tokenOtherPunctuation", 1038, 1038, "-"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1040, 1040, "B"},
		{"tokenOtherPunctuation", 1041, 1041, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 1044, 1049, "double"},
		{"tokenOtherPunctuation", 1050, 1050, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 1051, 1051, "A"},
		{"tokenOtherPunctuation", 1052, 1052, ")"},
		{"tokenOtherPunctuation", 1054, 1054, "-"},
		{"tokenOtherPunctuation", 1055, 1055, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1057, 1057, "A"},
		{"tokenOtherPunctuation", 1059, 1059, "*"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1061, 1061, "2"},
		{"tokenOtherPunctuation", 1062, 1062, "."},
		{"tokenAbcFullWith_Underscore_At_numbers", 1065, 1068, "half"},
		{"tokenOtherPunctuation", 1069, 1069, "("},
		{"tokenAbcFullWith_Underscore_At_numbers", 1070, 1070, "B"},
		{"tokenOtherPunctuation", 1071, 1071, ")"},
		{"tokenOtherPunctuation", 1073, 1073, "-"},
		{"tokenOtherPunctuation", 1074, 1074, ">"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1076, 1076, "B"},
		{"tokenOtherPunctuation", 1078, 1078, "/"},
		{"tokenAbcFullWith_Underscore_At_numbers", 1080, 1080, "2"},
		{"tokenOtherPunctuation", 1081, 1081, "."},

	}

	for _, tokenWanted := range(wanteds) {
		compare_tokenDetected_tokenWanted( "test basic types", sourceTokensExecutables__whitespacesSeparatorsBasicFile.Tokens, tokenWanted, t)
	}
	sourceTokensExecutables__whitespacesSeparatorsBasicFile.CharsFromErlFile.print_with_tokens(sourceTokensExecutables__whitespacesSeparatorsBasicFile.Tokens)
	// is this token, from this start->end pos, with this representation, is in the reply?
}
func compare_tokenDetected_tokenWanted(callerInfo string, tokensDetected ErlTokens, tokenWanted tokenWanted, t *testing.T) {
	tokenDetected, tokenWantedIsInDetected:= tokensDetected[tokenWanted.positionFirst]

	if tokenWantedIsInDetected {
		// theoretically the charPosFirst is always ok here, because the key in map was the same position
		tokenDetected_charPosFirst, tokenDetected_charPosLast := tokenDetected.charPositionFirstLast()
		if tokenDetected_charPosFirst != tokenWanted.positionFirst {
			t.Fatalf("\nErr First: %s : detected posFirst: %v  is different from wanted posFirst:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected_charPosLast != tokenWanted.positionLast {
			t.Fatalf("\nErr Last: %s : detected posLast: %v  is different from wanted posLast:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosLast, tokenWanted.positionLast, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected.stringRepresentation() != tokenWanted.textRepresentation {
			t.Fatalf("\nErr repr %s : startPos:%v  detected string representation: %v  is different from wanted representation:  %v, error",
				callerInfo, tokenDetected_charPosFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
	} else {
		t.Fatalf("\nErr %s : wanted tokenPosFirst %v is not in detecteds - error", callerInfo, tokenWanted.positionFirst)
	}

}
