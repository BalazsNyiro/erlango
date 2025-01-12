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
	funTokenOpener := token_opener_and_closer_look_forward__detect__char_literal
	printVerboseOpenerDetectMsg := true
	charactersInErlSrc, tokensInErlSrc = character_loop__pattern_detection__one_or_more_char__like_a_regexp(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

func token_opener_and_closer_look_forward__detect__char_literal(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc) (int, bool, int) { // ret: tokenType, openerDetected, positionModifier

	/////////////////////////////////////////////////////////
	tokenTypeId := TokenType_id_unknown
	openerDetected := false
	positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := 0
	/////////////////////////////////////////////////////////

	charStructNext1, charNext1CanBeDetected_notOverindexed := charactersInErlSrc.char_get_by_index(charPositionNowInSrc + 1)
	_, charNext2CanBeDetected_notOverindexed := charactersInErlSrc.char_get_by_index(charPositionNowInSrc + 2)

	if charStructNow.runeInErlSrc == '$' && charNext1CanBeDetected_notOverindexed {

		if charStructNext1.runeInErlSrc != '\\' {
			tokenTypeId = TokenType_id_Num_charLiterals
			openerDetected = true
			positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected = 1
		} else {
			// the next1 char is backslash.
			if charNext2CanBeDetected_notOverindexed {
				tokenTypeId = TokenType_id_Num_charLiterals
				openerDetected = true
				positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected = 2
			}
		}
	}
	return tokenTypeId, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected
}

func general_pattern__is_dollar_rune(r rune) bool {
	return r == '$'
}

func general_pattern__is__one_nonescaped__or__backslash_and_one_escaped__char(r rune) bool {
	return r == '$'
}
