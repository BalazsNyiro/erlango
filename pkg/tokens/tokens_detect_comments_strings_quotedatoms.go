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
	charactersInErlSrc, tokensInErlSrc = character_loop(TokenType_id_TextBlockQuotedDouble, charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop(
	tokenTypeId_wanted int,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

// the opener looks forward, the closer looks backward in the characters.
// the opener/closer elems are part of the token!!!
// so a string has a text, and the boundary too.
// example token content: "string_with_boundary"
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc, bool) bool,
	tokenCloserConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc, bool) bool) (CharacterInErlSrcCollector, TokenCollector) {

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

			openerDetected := tokenOpenerConditionFun(
				charPositionNowInSrc, charactersInErlSrc, charStructNow,
				inActiveTokenDetectionBecauseOpenerConditionTriggered)

			if openerDetected {
				inActiveTokenDetectionBecauseOpenerConditionTriggered = true

				charactersInErlSrc[charPositionNowInSrc].tokenDetectedType = tokenTypeId_wanted
				charactersInErlSrc[charPositionNowInSrc].tokenOpenerCharacter = true
				continue
			}

		} ////////////////////////////////////////////

		// I feel to declare the ELSE case with a verbose condition instead of an else,
		// it is more descriptive
		if inActiveTokenDetectionBecauseOpenerConditionTriggered {

			charactersInErlSrc[charPositionNowInSrc].tokenDetectedType = tokenTypeId_wanted

			////////////// CLOSER DETECT ///////////////
			closerDetected := tokenCloserConditionFun(
				charPositionNowInSrc, charactersInErlSrc, charStructNow,
				inActiveTokenDetectionBecauseOpenerConditionTriggered)

			if closerDetected {
				charactersInErlSrc[charPositionNowInSrc].tokenCloserCharacter = true
				inActiveTokenDetectionBecauseOpenerConditionTriggered = false
				continue
			}
		}

	}

	return charactersInErlSrc, tokensInErlSrc
}

func token_opener_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
	isActiveTokenDetectionBecauseOpenerConditionTriggered bool) bool {

	return true
}

func token_closer_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
	isActiveTokenDetectionBecauseOpenerConditionTriggered bool) bool {

	return true
}
