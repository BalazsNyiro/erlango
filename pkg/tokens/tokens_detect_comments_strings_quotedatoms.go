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

func tokens_detect_comments_strings_quotedatoms(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect_quote_double
	funTokenCloser := token_closer_detect_quote_double
	oneCharacterLongTokenDetection_standaloneCharacterWanted := false // "" a string has minimum 2 chars: an opener and a closer " char.
	charactersInErlSrc, tokensInErlSrc = character_loop(TokenType_id_TextBlockQuotedDouble, oneCharacterLongTokenDetection_standaloneCharacterWanted, charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
	return charactersInErlSrc, tokensInErlSrc
}

// keep it simple. Don't increase the complexity, this is the core of the parser.
func character_loop(
	tokenTypeId_wanted int,
	oneCharacterLongTokenDetection_standaloneCharacterWanted bool,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,

// the opener looks forward, the closer looks backward in the characters.
// the opener/closer elems are part of the token!!!
// so a string has a text, and the boundary too.
// example token content: "string_with_boundary"
// if a long token is detected (so more than one character, the opener can shift the current position.
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (bool, int),
	tokenCloserConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc, int) bool) (CharacterInErlSrcCollector, TokenCollector) {

	backSlashCounterBeforeCurrentChar := 0
	activeTokenDetectionBecauseOpenerConditionTriggered := false

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

		fmt.Printf("charPosition: %d, characterInLoop: %s  isActiveDetection: %t \n", charPositionNowInSrc, charStructNow.stringRepr(), activeTokenDetectionBecauseOpenerConditionTriggered)

		if !activeTokenDetectionBecauseOpenerConditionTriggered {

			openerDetected, positionModifierBecauseLongTokenOpenerCharsAreDetected := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			fmt.Println("opener detected:", openerDetected, positionModifierBecauseLongTokenOpenerCharsAreDetected)

			if openerDetected { ////////////// OPENER DETECT ///////////////
				activeTokenDetectionBecauseOpenerConditionTriggered = true

				charStructNow.tokenDetectedType = tokenTypeId_wanted
				charStructNow.tokenOpenerCharacter = true

				if oneCharacterLongTokenDetection_standaloneCharacterWanted {
					// if we want to detect one character only, which is an opener AND a closer same time,
					// close that immediatelly (and in this case pass a fake/empty closer function.
					charStructNow.tokenCloserCharacter = true
					activeTokenDetectionBecauseOpenerConditionTriggered = false
				}

				// this modifier is used ONLY if the detected token length is longer than 1 char.
				// in that case the modifier value is: (tokenLength-1)
				charPositionNowInSrc += positionModifierBecauseLongTokenOpenerCharsAreDetected

				charactersInErlSrc[charPositionNowInSrc] = charStructNow
				fmt.Printf("charPosition: %d, characterInOpen: %s  opener: %t closer: %t \n", charPositionNowInSrc, charStructNow.stringRepr(), charStructNow.tokenOpenerCharacter, charStructNow.tokenCloserCharacter)
				continue
			} ////////////////////////////////////////////////
		} // not active

		// it is more descriptive instead of a simple else:
		if activeTokenDetectionBecauseOpenerConditionTriggered {

			charStructNow.tokenDetectedType = tokenTypeId_wanted

			closerDetected := tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow, backSlashCounterBeforeCurrentChar)

			if closerDetected {
				charStructNow.tokenCloserCharacter = true
				activeTokenDetectionBecauseOpenerConditionTriggered = false
			} ///////////////////////////////////////////

			charactersInErlSrc[charPositionNowInSrc] = charStructNow

			fmt.Printf("charPosition: %d, characterInClos: %s  opener: %t closer: %t   tokenDetectedType: %d \n",
				charPositionNowInSrc,
				charactersInErlSrc[charPositionNowInSrc].stringRepr(),
				charactersInErlSrc[charPositionNowInSrc].tokenOpenerCharacter,
				charactersInErlSrc[charPositionNowInSrc].tokenCloserCharacter,
				charactersInErlSrc[charPositionNowInSrc].tokenDetectedType)
		}
	}
	return charactersInErlSrc, tokensInErlSrc
}

func token_opener_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) (bool, int) {

	// 0: double quote " opener is 1 char wide, there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '"', 0
}

func token_closer_detect_quote_double(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
	backSlashCounterBeforeCurrentCharInterpretedOnlyInStringLikeClosers int,
) bool {

	// if the remainder is not 0, there is an active escaping before this character
	if backSlashCounterBeforeCurrentCharInterpretedOnlyInStringLikeClosers%2 > 0 {
		return false // so this is
	}
	// 0: double quote " closer is 1 char wide, there is no need to shift the original character loop position
	return charStructNow.runeInErlSrc == '"'
}
