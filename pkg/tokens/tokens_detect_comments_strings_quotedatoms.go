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
	charactersInErlSrc, tokensInErlSrc = character_loop(charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc, bool) bool,
	tokenCloserConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc, bool) bool) (CharacterInErlSrcCollector, TokenCollector) {

	backSlashCounterBeforeCurrentChar := 0
	isActiveTokenDetectionBecauseOpenerConditionTriggered := false

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	// so I think it is safer to not use a range here (containter is updated inside the loop)
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]

		if charPositionNowInSrc > 0 {
			charStructPrev := charactersInErlSrc[charPositionNowInSrc-1]
			if charStructPrev.runeInErlSrc == '\\' {
				backSlashCounterBeforeCurrentChar++
			} else { // if prev is not backslash reset the counter
				backSlashCounterBeforeCurrentChar = 0
			}
		} // > 0

		fmt.Printf("charPosition: %d, characterLoop: %s\n", charPositionNowInSrc, charStructNow.stringRepr())

		// TODO: opener/closer func usage
		tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow,
			isActiveTokenDetectionBecauseOpenerConditionTriggered)

		tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow,
			isActiveTokenDetectionBecauseOpenerConditionTriggered)

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
