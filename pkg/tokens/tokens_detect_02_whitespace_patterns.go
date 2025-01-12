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

func tokens_detect_02_erlang_whitespaces(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_and_closer_look_forward__detect__whitespaces
	printVerboseOpenerDetectMsg := true
	charactersInErlSrc, tokensInErlSrc = character_loop__pattern_detection__one_or_more_char__like_a_regexp(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop__pattern_detection__one_or_more_char__like_a_regexp(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,
	tokenOpenerConditionFun func(int, CharacterInErlSrcCollector, CharacterInErlSrc) (int, bool, int),
	printVerboseOpenerDetectMsg bool) (CharacterInErlSrcCollector, TokenCollector) {

	tokenTypeId_now := TokenType_id_unknown

	cleaningAfterTokenClose_set_back_default_values := func() {
		tokenTypeId_now = TokenType_id_unknown
	}

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		if charStructNow.tokenDetectedType != TokenType_id_unknown {
			continue // if the char was detected and has a TokenType_id, there is no more to do.
		} //don't start new detection if the current char was detected once

		tokenTypeId_fromOpener, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
		tokenTypeId_now = tokenTypeId_fromOpener

		if printVerboseOpenerDetectMsg {
			fmt.Println("opener detected:", openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, charStructNow.stringRepr())
		}

		if openerDetected { ////////////// OPENER DETECT ///////////////

			charStructNow.tokenDetectedType = tokenTypeId_now
			charStructNow.tokenOpenerCharacter = true
			charactersInErlSrc[charPositionNowInSrc] = charStructNow

			// this modifier>0 ONLY if the detected token length is longer than 1 char.
			// in that case the modifier value is: (tokenLength-1), because the for loop does 1 increasing
			charPositionNewWanted := charPositionNowInSrc + positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected
			for charPositionNowInSrc < charPositionNewWanted {
				charPositionNowInSrc++

				charStructNow = charactersInErlSrc[charPositionNowInSrc]
				charStructNow.tokenDetectedType = tokenTypeId_now
				charactersInErlSrc[charPositionNowInSrc] = charStructNow
			}
			// close Immediately
			charStructNow.tokenCloserCharacter = true                // close the last charStructNow elem,
			charactersInErlSrc[charPositionNowInSrc] = charStructNow // if the previous loop updated more chars.
			cleaningAfterTokenClose_set_back_default_values()        // maybe that is not the starter one,
		} // if openerDetected

	} // for charPosition....

	return charactersInErlSrc, tokensInErlSrc
} // func character_loop patterns
///////////////////////////////////////////////////

func token_opener_and_closer_look_forward__detect__whitespaces(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc) (int, bool, int) { // ret: tokenType, openerDetected, positionModifier

	generalCharOpenerDetector := general_pattern__is_whitespace_rune_inside_line
	generalCharNextAcceptableDetector := general_pattern__is_whitespace_rune_inside_line
	tokenTypeIfActiveDetection := TokenType_id_WhitespaceInLine_ErlSrc

	if charStructNow.runeInErlSrc == '\n' { // newline handled separately, I want to close that at detection
		generalCharOpenerDetector = general_pattern__is_whitespace_rune_newline
		generalCharNextAcceptableDetector = general_pattern__false_always // every newline chars are separated, no next char check
		tokenTypeIfActiveDetection = TokenType_id_WhitespaceNewLine_ErlSrc
	}

	return general_look_forward_accepted_chars_detector(
		charPositionNowInSrc,
		charactersInErlSrc,
		charStructNow,
		generalCharOpenerDetector,
		generalCharNextAcceptableDetector,
		tokenTypeIfActiveDetection,
	)
}

// /////////////////////////////////////////////////
// sometime a char has to be closed immediately, there is no need to analyse the next char.
func general_pattern__false_always(_ rune) bool {
	return false
}

func general_pattern__is_whitespace_rune_inside_line(r rune) bool {
	return r == ' ' || r == '\r' || r == '\t'
}

func general_pattern__is_whitespace_rune_newline(r rune) bool {
	return r == '\n'
}

// this is a generic 'look forward' detector
func general_look_forward_accepted_chars_detector(
	charPositionNowInSrc int,                      //                      this opener uses ONLY the actual character,
	charactersInErlSrc CharacterInErlSrcCollector, // there is no need to look forward/back in src
	charStructNow CharacterInErlSrc,

// the first char rules are sometime different from the next char rules
	generalCharNowAcceptableDetector func(rune) bool,
	generalCharNextAcceptableDetector func(rune) bool,

	tokenTypeIfActiveDetection int) (int, bool, int) {

	// 0 means: there is no need to shift the original character loop position
	positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := 0

	tokenTypeId := TokenType_id_unknown
	openerDetected := false

	if generalCharNowAcceptableDetector(charStructNow.runeInErlSrc) {

		openerDetected = true
		tokenTypeId = tokenTypeIfActiveDetection

		charPositionNextInSrc := charPositionNowInSrc
		for true {
			charPositionNextInSrc++

			charStructNext, charWasDetectedCorrectly := charactersInErlSrc.char_get_by_index(charPositionNextInSrc)

			if charWasDetectedCorrectly {
				if generalCharNextAcceptableDetector(charStructNext.runeInErlSrc) {
					positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected++
				} else {
					break
				}

			} else {
				break
			}

		}
	}
	return tokenTypeId, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected
}
