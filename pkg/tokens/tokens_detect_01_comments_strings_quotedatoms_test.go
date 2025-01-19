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

func Test_tokens_detect_in_erl_src(t *testing.T) {
	testName := "Test_tokens_detect_strings_1"
	// erlSrcRunes, _ := base_toolset.File_read_runes("erl_src/erlang_whitespaces_separators_basic_types.erl", "Test_tokens_detect_in_erl_src")

	erlSrcRunes := []rune(`A  = "B\"". % this is a string plus a comment
                           % comment in newline
                           Long_String1 = """This is a "quote"  """.
                           Long_String2 = """ \n\n very\t"Long" string   \"""".
                           Num = (1 + 2.3 / 4 * 5).
                           List = [6,7].
                            `)
	charactersInErlSrc := Runes_to_character_structs(erlSrcRunes)
	tokensInErlSrc := TokenCollector{}

	charactersInErlSrc, tokensInErlSrc = tokens_detect_01_erlang_strings__quoted_atoms__comments(charactersInErlSrc, tokensInErlSrc)
	charactersInErlSrc, tokensInErlSrc = tokens_detect_02_erlang_whitespaces(charactersInErlSrc, tokensInErlSrc)
	charactersInErlSrc, tokensInErlSrc = tokens_detect_03_alphanumerics(charactersInErlSrc, tokensInErlSrc)
	charactersInErlSrc, tokensInErlSrc = tokens_detect_04__braces__dotsCommas__operatorBuilders(charactersInErlSrc, tokensInErlSrc)

	Tokens_detection_print_verbose(charactersInErlSrc, tokensInErlSrc)

	// line   0 >>> ============================
	// line   0 >>> 012345678901234567890
	// line   0 >>>  oc 2o   c 2o
	// line   0 >>> ?ww?w"""""?w%%%%%%%%%
	// line   0 >>> A  = "B\"". % comment

	charNow := charactersInErlSrc[4]
	compare_bool_bool(testName, true, charNow.tokenOpenerCharacter, t)
	compare_bool_bool(testName, true, charNow.tokenCloserCharacter, t)
	// compare_string_string(testName, charactersInErlSrc[4].tokenOpenerCharacter, true, t)

	charNow = charactersInErlSrc[5]
	compare_bool_bool(testName, true, charNow.tokenOpenerCharacter, t)
	compare_bool_bool(testName, false, charNow.tokenCloserCharacter, t)
	compare_int__int_(testName, TokenType_id_TextBlockQuotedDouble, charNow.tokenDetectedType, t)

	charNow = charactersInErlSrc[6]
	compare_bool_bool(testName, false, charNow.tokenOpenerCharacter, t)
	compare_bool_bool(testName, false, charNow.tokenCloserCharacter, t)
	compare_int__int_(testName, TokenType_id_TextBlockQuotedDouble, charNow.tokenDetectedType, t)

	charNow = charactersInErlSrc[8]
	compare_bool_bool(testName, false, charNow.tokenOpenerCharacter, t)
	compare_bool_bool(testName, false, charNow.tokenCloserCharacter, t)
	compare_int__int_(testName, TokenType_id_TextBlockQuotedDouble, charNow.tokenDetectedType, t)

	charNow = charactersInErlSrc[9]
	compare_bool_bool(testName, false, charNow.tokenOpenerCharacter, t)
	compare_bool_bool(testName, true, charNow.tokenCloserCharacter, t)
	compare_int__int_(testName, TokenType_id_TextBlockQuotedDouble, charNow.tokenDetectedType, t)
}
