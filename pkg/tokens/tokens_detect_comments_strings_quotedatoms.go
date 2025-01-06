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
func tokens_detect_erlang_strings__quoted_atoms__comments(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect__quoteDouble__quoteSinge_comment
	printVerboseOpenerDetectMsg := false
	charactersInErlSrc, tokensInErlSrc = character_loop_openers_closers(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

// TODO: test this
func is_escaped_the_current_char(charPositionInSrc int, charactersInErlSrc CharacterInErlSrcCollector) bool {

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
func character_loop_openers_closers(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

	// the opener looks forward, the closer looks backward in the characters.
	// the opener/closer elems are part of the token - so a string has a text, and the boundary too.
	// example token content: "string_with_boundary"
	// if a long token is detected (so more than one character, the opener can shift the current position.
	// the closer func is returned from the opener func, because sometime an opener can detect
	// more than one type (string|quotedAtom|comment) and this info is created only in the opener state
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (int, bool, int, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, bool),
	printVerboseOpenerDetectMsg bool) (CharacterInErlSrcCollector, TokenCollector) {

	tokenCloserConditionFun := token_closer_fake_placeholder_fun

	activeTokenDetectionBecauseOpenerConditionTriggered := false
	tokenTypeId_now := TokenType_id_unknown


	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		if charStructNow.tokenDetectedType != TokenType_id_unknown {
			continue // if the char was detected and has a TokenType_id, there is no more to do.
		}

		if !activeTokenDetectionBecauseOpenerConditionTriggered {

			tokenTypeId_fromOpener, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, tokenCloserConditionFunFromOpener, _ := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			tokenCloserConditionFun = tokenCloserConditionFunFromOpener
			tokenTypeId_now = tokenTypeId_fromOpener

			if printVerboseOpenerDetectMsg {
				fmt.Println("opener detected:", openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, charStructNow.stringRepr())
			}

			if openerDetected { ////////////// OPENER DETECT ///////////////
				activeTokenDetectionBecauseOpenerConditionTriggered = true
				charStructNow.tokenDetectedType = tokenTypeId_now
				charStructNow.tokenOpenerCharacter = true
				charactersInErlSrc[charPositionNowInSrc] = charStructNow
			} ///////////////////////////////// OPENER DETECT ///////////////

		} else { // opener was detected previously - the loop is in activeTokenDetectionBecauseOpenerConditionTriggered == true:
			charStructNow.tokenDetectedType = tokenTypeId_now

			closerDetected := tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			if closerDetected {
				charStructNow.tokenCloserCharacter = true

				// set_noActiveTokenDetection__tokenTypeUnknown
				activeTokenDetectionBecauseOpenerConditionTriggered = false
				tokenTypeId_now = TokenType_id_unknown

			} ///////////////////////////////////////////

			charactersInErlSrc[charPositionNowInSrc] = charStructNow
		}
	} // for charPosition....

	return charactersInErlSrc, tokensInErlSrc
} // func character_loop_openers_closers

func token_opener_detect__quoteDouble__quoteSinge_comment(
	charPositionNowInSrc int, //                      this opener uses ONLY the actual character,
	charactersInErlSrc CharacterInErlSrcCollector, // there is no need to look forward/back in src
	charStructNow CharacterInErlSrc) (int, bool, int, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, bool) {

	oneCharacterLongTokenDetection__theCharIsOpenerAndCloserSameTime_closeDetectionImmediately := false
	positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := 0

	tokenTypeId := TokenType_id_unknown
	funCloser := token_closer_fake_placeholder_fun
	openerDetected := false

	// 0: double quote " opener is 1 char wide,
	//there is no need to shift the original character loop position

	if charStructNow.runeInErlSrc == '"' {
		openerDetected = true
		funCloser = token_closer_detect_quote_double
		tokenTypeId = TokenType_id_TextBlockQuotedDouble
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

	return tokenTypeId, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, funCloser, oneCharacterLongTokenDetection__theCharIsOpenerAndCloserSameTime_closeDetectionImmediately
}

func token_closer_detect_comment_end(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return charStructNow.runeInErlSrc == '\n' // the end of a comment is a newline char
}

func token_closer_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {

	if is_escaped_the_current_char(charPositionNowInSrc, charactersInErlSrc) {
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

	if is_escaped_the_current_char(charPositionNowInSrc, charactersInErlSrc) {
		return false // so this cannot be a ' closer, because escaped
	}
	// 0: single quote ' closer is 1 char wide, there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '\''
}

func token_closer_fake_placeholder_fun(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return false
}
