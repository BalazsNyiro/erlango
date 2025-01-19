/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import "fmt"

/*
in a comment, there can be a string:  % example "string" in comment
in a string, there can be a % sign:  "taxes are increased with 10% in this year"

and in quoted atoms, there can be other signs:
A='atom_with_double_quote"'.
'atom_with_double_quote"'

So these 3 has to be handled in one func.
*/
func tokens_detect_01_erlang_strings__quoted_atoms__comments(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect__quoteTriple_quoteDouble_quoteSingle_comment
	printVerboseOpenerDetectMsg := false
	charactersInErlSrc, tokensInErlSrc = character_loop_openers_closers__detect_minimum_2_chars_with_welldefined_opener_closer_section(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

// TODO: test this
func is_escaped_char(charPositionInSrc int, charactersInErlSrc CharacterInErlSrcCollector) bool {

	charPositionPrev := charPositionInSrc - 1
	backSlashCounterBeforeCurrentChar := 0

	for charPositionPrev >= 0 {
		charStructPrev := charactersInErlSrc[charPositionPrev]

		if charStructPrev.runeInErlSrc == '\\' {
			backSlashCounterBeforeCurrentChar++
			charPositionPrev--
		} else { // if prev is not backslash
			break
		}
	} // >= 0

	isEscaped := false // it is more readable if named variable is created here
	if backSlashCounterBeforeCurrentChar%2 != 0 {
		isEscaped = true
	}
	return isEscaped
}

// keep it simple. Don't increase the complexity
func character_loop_openers_closers__detect_minimum_2_chars_with_welldefined_opener_closer_section(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

	// the opener looks forward, the closer looks backward in the characters.
	// the opener/closer elems are part of the token - so a string has a text, and the boundary too.
	// example token content: "string_with_boundary"
	// if a long token is detected (so more than one character, the opener can shift the current position.
	// the closer func is returned from the opener func, because sometime an opener can detect
	// more than one type (string|quotedAtom|comment) and this info is created only in the opener state
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (int, bool, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, int),
	printVerboseOpenerDetectMsg bool) (CharacterInErlSrcCollector, TokenCollector) {

	tokenCloserConditionFun := token_closer_always_false___never_close

	activeTokenDetectionBecauseOpenerConditionTriggered := false
	tokenTypeId_now := TokenType_id_unknown

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); {

		if !activeTokenDetectionBecauseOpenerConditionTriggered {
			charStructNow := charactersInErlSrc[charPositionNowInSrc]

			if charStructNow.tokenIsDetected() { // READING ONLY of current char OPERATION
				continue // if the char was detected and has a TokenType_id, there is no more to do.
			} //don't start new detection if the current char was detected once

			// the opener func can look forward, and sometime more than one char is processed and detected as an opener part
			tokenTypeId_fromOpener, openerDetected, tokenCloserConditionFunFromOpener, positionShifter__usedCharNumDuringDetection := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			tokenCloserConditionFun = tokenCloserConditionFunFromOpener
			tokenTypeId_now = tokenTypeId_fromOpener

			if printVerboseOpenerDetectMsg {
				fmt.Println("opener detected:", openerDetected, charStructNow.stringRepr()) // READING
			}

			for shift := 0; shift < positionShifter__usedCharNumDuringDetection; shift++ {

				if openerDetected { ////////////// OPENER DETECT ///////////////

					charStructNow = charactersInErlSrc[charPositionNowInSrc]
					charStructNow.tokenDetectedType = tokenTypeId_now
					if shift == 0 { // the first char is marked as an opener
						charStructNow.tokenOpenerCharacter = true
					}
					charactersInErlSrc[charPositionNowInSrc] = charStructNow

					activeTokenDetectionBecauseOpenerConditionTriggered = true
				}
				charPositionNowInSrc++
			}

		} else { // opener was detected previously - the loop is in activeTokenDetectionBecauseOpenerConditionTriggered == true:

			charStructNow := charactersInErlSrc[charPositionNowInSrc]
			charStructNow.tokenDetectedType = tokenTypeId_now

			closerDetected := tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			if closerDetected {
				charStructNow.tokenCloserCharacter = true

				// set_noActiveTokenDetection__tokenTypeUnknown
				activeTokenDetectionBecauseOpenerConditionTriggered = false
				tokenTypeId_now = TokenType_id_unknown

			} ///////////////////////////////////////////
			charactersInErlSrc[charPositionNowInSrc] = charStructNow
			charPositionNowInSrc++
		}

	} // for charPosition....

	return charactersInErlSrc, tokensInErlSrc
} // func character_loop_openers_closers__detect_minimum_2_chars_with_welldefined_opener_closer_section

func token_opener_detect__quoteTriple_quoteDouble_quoteSingle_comment(
	charPositionNowInSrc int, //                      this opener uses ONLY the actual character,
	charactersInErlSrc CharacterInErlSrcCollector, // there is no need to look forward/back in src
	charStructNow CharacterInErlSrc) (int, bool, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, int) {

	positionShifter__usedCharNumDuringDetection := 1

	tokenTypeId := TokenType_id_unknown
	funCloser := token_closer_always_false___never_close
	openerDetected := false

	// 0: double quote " opener is 1 char wide,
	//there is no need to shift the original character loop position

	if charStructNow.runeInErlSrc == '"' {
		openerDetected = true
		funCloser = token_closer_detect_quote_double
		tokenTypeId = TokenType_id_TextBlockQuotedDouble

		charStructNext1 := charactersInErlSrc.char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(charPositionNowInSrc + 1)
		charStructNext2 := charactersInErlSrc.char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(charPositionNowInSrc + 2)

		if charStructNext1.runeInErlSrc == '"' && charStructNext2.runeInErlSrc == '"' {
			funCloser = token_closer_detect_quote_triple
			tokenTypeId = TokenType_id_TextBlockQuotedTriple
			positionShifter__usedCharNumDuringDetection = 3
		}

	}

	if charStructNow.runeInErlSrc == '\'' {
		funCloser = token_closer_detect_quote_single
		openerDetected = true
		tokenTypeId = TokenType_id_TextBlockQuotedSingle
	}

	if charStructNow.runeInErlSrc == '%' {
		funCloser = token_closer_detect_comment_end
		openerDetected = true
		tokenTypeId = TokenType_id_Comment
	}

	if charStructNow.runeInErlSrc == '$' {
		// minimum 2 char wide,
		//   - non-escaped: 'A=$B.'
		//   - escaped: 'C=$\n.'
		//     Char can be any unicode char: 'U=$ち.'

		// possible problems: in a char literal, there can be ", ', % signs.
		// for example: $", $', $% - and these will be detedted by quote/comma detection.
		// So the char literal has to be in the same logical step with the string detection

		funCloser = token_closer_always_close_at_next_char
		openerDetected = true
		tokenTypeId = TokenType_id_Num_charLiterals

		// check if the next char is a backslash or not
		charStructNext1 := charactersInErlSrc.char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(charPositionNowInSrc + 1)

		if charStructNext1.runeInErlSrc == '\\' {
			positionShifter__usedCharNumDuringDetection = 2 // '$\' was detected, and the next char will be closed immediatelly (one char can be after $ or $\
		} // if there is backslash/escape sign...

	}

	return tokenTypeId, openerDetected, funCloser, positionShifter__usedCharNumDuringDetection
}

func token_closer_detect_comment_end(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return charStructNow.runeInErlSrc == '\n' // the end of a comment is a newline char
}

func token_closer_detect_quote_triple(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {

	if charPositionNowInSrc >= 3 {
		// tripple quote needs minimum 3 previous opener chars, in pos 0,1,2.
		// so the first closer can be minimum in position 3
		if !is_escaped_char(charPositionNowInSrc-2, charactersInErlSrc) {
			charStructPrev1 := charactersInErlSrc.char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(charPositionNowInSrc - 1)
			charStructPrev2 := charactersInErlSrc.char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(charPositionNowInSrc - 2)

			if charStructNow.runeInErlSrc == '"' && charStructPrev1.runeInErlSrc == '"' && charStructPrev2.runeInErlSrc == '"' {
				return true
			}

		}
	}
	return false
}
func token_closer_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {

	if is_escaped_char(charPositionNowInSrc, charactersInErlSrc) {
		return false // so this cannot be a " closer, because escaped
	}
	// 0: double quote " closer is 1 char wide, there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '"'
}

func token_closer_detect_quote_single(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {

	if is_escaped_char(charPositionNowInSrc, charactersInErlSrc) {
		return false // so this cannot be a ' closer, because escaped
	}
	// 0: single quote ' closer is 1 char wide, there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '\''
}

func token_closer_always_close_at_next_char(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return true
}

func token_closer_always_false___never_close(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return false
}
