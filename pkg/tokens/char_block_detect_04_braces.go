/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

func character_block_detect__04_erlang_braces__dotsCommas__operatorBuilders(charactersInErlSrc CharacterInErlSrcCollector) CharacterInErlSrcCollector {

	for _, wantedCharInErl := range []rune("()[]{}") {
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, CharBlock_braces_grouping_elems)
	}

	// TODO:  find double chars << >> -> <- ::  == != <>
	// :: is used in type specification

	for _, wantedCharInErl := range []rune(",;.:") {
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, CharBlock_dots_commas)
	}

	for _, wantedCharInErl := range []rune("=<>+-*/#?|@!") {
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, CharBlock_LanguageElement_operators_specialchars)
	}

	return charactersInErlSrc
}
