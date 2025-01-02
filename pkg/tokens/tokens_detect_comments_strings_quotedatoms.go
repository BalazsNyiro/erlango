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
	"fmt"
)

func tokens_detect_in_erl_src(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	charactersInErlSrc2, tokensInErlSrc2 := tokens_detect_comments_strings_quotedatoms(charactersInErlSrc, tokensInErlSrc)
	return charactersInErlSrc2, tokensInErlSrc2
}

func tokens_detect_comments_strings_quotedatoms(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect_quote_double
	funTokenCloser := token_closer_detect_quote_double
	oneCharacterLongTokenDetection_standaloneCharacterWanted := false // "" a string has minimum 2 chars: an opener and a closer " char.
	charactersInErlSrc, tokensInErlSrc = character_loop(TokenType_id_TextBlockQuotedDouble, oneCharacterLongTokenDetection_standaloneCharacterWanted, charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop(
	tokenTypeId_wanted int,
	oneCharacterLongTokenDetection_standaloneCharacterWanted bool,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

// the opener looks forward, the closer looks backward in the characters.
// the opener/closer elems are part of the token!!!
// so a string has a text, and the boundary too.
// example token content: "string_with_boundary"
// if a long token is detected (so more than one characters, the opener can shift the current position.
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (bool, int),
	tokenCloserConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool) (CharacterInErlSrcCollector, TokenCollector) {

	backSlashCounterBeforeCurrentChar := 0
	inActiveTokenDetectionBecauseOpenerConditionTriggered := false

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	// so I think it is safer to not use a range here (containter is updated inside the loop)
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]

		if charPositionNowInSrc > 0 { // backslash
			charStructPrev := charactersInErlSrc[charPositionNowInSrc-1]
			if charStructPrev.runeInErlSrc == '\\' {
				backSlashCounterBeforeCurrentChar++
			} else { // if prev is not backslash reset the counter
				backSlashCounterBeforeCurrentChar = 0
			}
		} // > 0

		fmt.Printf("charPosition: %d, characterLoop: %s\n", charPositionNowInSrc, charStructNow.stringRepr())

		////////////// OPENER DETECT ///////////////
		if !inActiveTokenDetectionBecauseOpenerConditionTriggered {

			openerDetected, positionModifierBecauseLongTokenOpenerCharsAreDetected := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)

			if openerDetected {
				inActiveTokenDetectionBecauseOpenerConditionTriggered = true

				charactersInErlSrc[charPositionNowInSrc].tokenDetectedType = tokenTypeId_wanted
				charactersInErlSrc[charPositionNowInSrc].tokenOpenerCharacter = true

				if oneCharacterLongTokenDetection_standaloneCharacterWanted {
					// if we want to detect one character only, which is an opener AND a closer same time,
					// close that immediatelly (and in this case pass a fake/empty closer function.
					charactersInErlSrc[charPositionNowInSrc].tokenCloserCharacter = true
					inActiveTokenDetectionBecauseOpenerConditionTriggered = false
				}

				// this modifier is used ONLY if the detected token length is longer than 1 char.
				// in that case the modifier value is: (tokenLength-1)
				charPositionNowInSrc += positionModifierBecauseLongTokenOpenerCharsAreDetected
				continue
			}
		} ////////////////////////////////////////////////

		// it is more descriptive instead of a simple else:
		if inActiveTokenDetectionBecauseOpenerConditionTriggered {

			charactersInErlSrc[charPositionNowInSrc].tokenDetectedType = tokenTypeId_wanted

			////////////// CLOSER DETECT ///////////////
			closerDetected := tokenCloserConditionFun(
				charPositionNowInSrc, charactersInErlSrc, charStructNow)

			if closerDetected {
				charactersInErlSrc[charPositionNowInSrc].tokenCloserCharacter = true
				inActiveTokenDetectionBecauseOpenerConditionTriggered = false
				continue
			} ///////////////////////////////////////////
		}

	}

	return charactersInErlSrc, tokensInErlSrc
}

func token_opener_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc) (bool, int) {

	// 0: double quote " opener is 1 char wide, there is no need to shift the original character loop position
	return true, 0
}

func token_closer_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc) bool {

	return true
}
