/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

func tokens_detect_comments_strings_quotedatoms(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect_quote_double
	funTokenCloser := token_closer_detect_quote_double
	oneCharacterLongTokenDetection_standaloneCharacterWanted := false // "" a string has minimum 2 chars: an opener and a closer " char.
	charactersInErlSrc, tokensInErlSrc = character_loop(TokenType_id_TextBlockQuotedDouble, oneCharacterLongTokenDetection_standaloneCharacterWanted, charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
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

// keep it simple. Don't increase the complexity, this is the core of the parser.
func character_loop(
	tokenTypeId_wanted int,
	oneCharacterLongTokenDetection__theCharIsOpenerAndCloserSameTime_closeDetectionImmediately bool,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

	// the opener looks forward, the closer looks backward in the characters.
	// the opener/closer elems are part of the token - so a string has a text, and the boundary too.
	// example token content: "string_with_boundary"
	// if a long token is detected (so more than one character, the opener can shift the current position.
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (bool, int),
	tokenCloserConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool) (CharacterInErlSrcCollector, TokenCollector) {

	activeTokenDetectionBecauseOpenerConditionTriggered := false

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]

		if !activeTokenDetectionBecauseOpenerConditionTriggered {

			openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			// fmt.Println("opener detected:", openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected)

			if openerDetected { ////////////// OPENER DETECT ///////////////
				activeTokenDetectionBecauseOpenerConditionTriggered = true

				charStructNow.tokenDetectedType = tokenTypeId_wanted
				charStructNow.tokenOpenerCharacter = true

				if oneCharacterLongTokenDetection__theCharIsOpenerAndCloserSameTime_closeDetectionImmediately {
					charStructNow.tokenCloserCharacter = true                   // and if you detect 1 char only,
					activeTokenDetectionBecauseOpenerConditionTriggered = false // pass a fake/empty closer function.
				}

				// this modifier>0 ONLY if the detected token length is longer than 1 char.
				// in that case the modifier value is: (tokenLength-1), because the for loop does 1 increasing
				charPositionNowInSrc += positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected
			} ////////////////////////////////////////////////

		} else { // activeTokenDetectionBecauseOpenerConditionTriggered == true:
			charStructNow.tokenDetectedType = tokenTypeId_wanted

			closerDetected := tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			if closerDetected {
				charStructNow.tokenCloserCharacter = true
				activeTokenDetectionBecauseOpenerConditionTriggered = false
			} ///////////////////////////////////////////
		}
		charactersInErlSrc[charPositionNowInSrc] = charStructNow
	} // for charPosition....

	return charactersInErlSrc, tokensInErlSrc
} // func character_loop

func token_opener_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) (bool, int) {

	// 0: double quote " opener is 1 char wide,
	//there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '"', 0
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
