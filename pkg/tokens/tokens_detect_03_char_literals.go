/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

// minimum 2 char wide,
//   - non-escaped: 'A=$B.'
//   - escaped: 'C=$\n.'
//     Char can be any unicode char: 'U=$ち.'

func tokens_detect_03a_erlang_char_literals_nonescaped__dollar_plus_character(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {

	charactersInErlSrc = character_loop__set_one_char_tokentype('$', charactersInErlSrc, TokenType_id_Num_charLiterals)

	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		charStructNext1, charNext1CanBeDetected_notOverindexed := charactersInErlSrc.char_get_by_index(charPositionNowInSrc + 1)
		charStructNext2, charNext2CanBeDetected_notOverindexed := charactersInErlSrc.char_get_by_index(charPositionNowInSrc + 2)
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		if charStructNow.tokenDetectedType == TokenType_id_Num_charLiterals && charStructNow.runeInErlSrc == '$' && charNext1CanBeDetected_notOverindexed {

			if charStructNext1.tokenNotDetected() {
				charStructNext1.tokenDetectedType = TokenType_id_Num_charLiterals

				if charStructNext1.runeInErlSrc == '\\' && charNext2CanBeDetected_notOverindexed {
					if charStructNext2.tokenNotDetected() {
						charStructNext2.tokenDetectedType = TokenType_id_Num_charLiterals
					}
				} // if, next 2

			} // if, next1 token is not detected
		} // current char is $

		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		charactersInErlSrc[charPositionNowInSrc] = charStructNow
		if charNext1CanBeDetected_notOverindexed {
			charactersInErlSrc[charPositionNowInSrc+1] = charStructNext1
		}
		if charNext2CanBeDetected_notOverindexed {
			charactersInErlSrc[charPositionNowInSrc+2] = charStructNext2
		}
	}

	character_loop__opener_closer_sections_set__if_more_separated_isolated_elems_are_next_to_each_other(charactersInErlSrc, TokenType_id_Num_charLiterals)

	return charactersInErlSrc, tokensInErlSrc
}
