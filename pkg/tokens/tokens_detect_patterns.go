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

func tokens_detect_erlang_whitespaces(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_and_closer_look_forward__detect__whitespaces
	printVerboseOpenerDetectMsg := true
	charactersInErlSrc, tokensInErlSrc = character_loop_patterns(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop_patterns(
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

	tokenCloserConditionFun := token_closer_fake_placeholder_fun_pattern

	activeTokenDetectionBecauseOpenerConditionTriggered := false
	tokenTypeId_now := TokenType_id_unknown
	openerAndCloserSameTime_closeDetectionImmediately := false

	set_noActiveTokenDetection__tokenTypeUnknown__cleaningAfterTokenClose := func() {
		activeTokenDetectionBecauseOpenerConditionTriggered = false
		tokenTypeId_now = TokenType_id_unknown
		openerAndCloserSameTime_closeDetectionImmediately = false
	}

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		if charStructNow.tokenDetectedType != TokenType_id_unknown {
			continue // if the char was detected and has a TokenType_id, there is no more to do.
		}

		if !activeTokenDetectionBecauseOpenerConditionTriggered {

			tokenTypeId_fromOpener, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, tokenCloserConditionFunFromOpener, closeImmediately := tokenOpenerConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			tokenCloserConditionFun = tokenCloserConditionFunFromOpener
			tokenTypeId_now = tokenTypeId_fromOpener
			openerAndCloserSameTime_closeDetectionImmediately = closeImmediately

			if printVerboseOpenerDetectMsg {
				fmt.Println("opener detected:", openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, charStructNow.stringRepr())
			}

			if openerDetected { ////////////// OPENER DETECT ///////////////
				activeTokenDetectionBecauseOpenerConditionTriggered = true

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

				if openerAndCloserSameTime_closeDetectionImmediately {
					// one char long operators (+,-,*,/), commas and other elems are only ONE char wide elems, they need to be closed when they opened
					// or more than one char were processed in the opener, and positionModifier was used
					// honestly the separated tokenCloser is typically used for strings and comments,
					// other elems are easier to handled in one step, when opener/closer are processed in one step,

					// BUT: if the opener/closer are handled in one func, that is more complicated,
					// in one word: try to use which method is more nature in a given situation (separated opener/closer or mixed solution).

					// if you can, use separated opener/closer functions.
					// this can be a problem in a situation when the (active-1) so the previous character
					// is the closer. The for loop goes forward, so it is harder to look back from the closer fun,
					// and modify a previously processed character again.

					// new suggestion: separated opener/closer can be used easily, if the knowing of actual character
					// is enough, and you don't need to modify back a closing property.

					// ======================================================
					// I try to explain it in a different way (and maybe this is the best):
					// a section with whitespaces can be one char long, or more char long.

					// the separated opener/closer approach cannot be used when the token is one char wide,
					// because when the opening is detected, a token closing is necessary, too.

					// the current character_loop() solution made a choice in first level:
					// do an opening OR a closing (is activeTokenDetection or not), because the parsing's first step was string/comment/quotedAtom detection,
					// and in those cases there are well-defined and differently positioned opening/closing elems.

					// so, if you have a token which needs to be detected in one step (operators for example)
					// then you need to use the opening+closing method, not the separated opening/closing funs,
					// because with the opening, a closing is necessary too for the actual character

					// I have the feeling that there is a way to convert a mixed solution to be separated,
					// but this solution seems to give the flexibility: to use the easier method which is the more natural.

					// if the token has external boundaries, the opener/closer approach is simple and work (string/quotedAtoms/comments)
					// if the token hasn't got boundaries, but it has general firs char + other chars rules, the mixed opener/closer is useful.
					// these are SOFT RULES ONLY because the mixed-opener/closer solution can be written from separated opener/closer solution,

					// important: the string/quotedAtom/Comment question is problematic,
					// because they need to be detected together (there can be a comment in a string, or a string in a comment)
					// and based on the real situation, the closer function is changing.

					// !! the separated opener/closer has an option to look back OR look forward, same time.
					// !! the mixed option can look forward only, so in some situation that is not enough.

					// if the token is one char wide, this special section is added to do an immediate closing:

					charStructNow.tokenCloserCharacter = true                               // close the last charStructNow elem,
					charactersInErlSrc[charPositionNowInSrc] = charStructNow                // if the previous loop updated more chars.
					set_noActiveTokenDetection__tokenTypeUnknown__cleaningAfterTokenClose() // maybe that is not the starter one,
				}
			} ////////////////////////////////////////////////

		} else { // opener was detected previously - the loop is in activeTokenDetectionBecauseOpenerConditionTriggered == true:
			charStructNow.tokenDetectedType = tokenTypeId_now

			closerDetected := tokenCloserConditionFun(charPositionNowInSrc, charactersInErlSrc, charStructNow)
			if closerDetected {
				charStructNow.tokenCloserCharacter = true
				set_noActiveTokenDetection__tokenTypeUnknown__cleaningAfterTokenClose()
			} ///////////////////////////////////////////

			charactersInErlSrc[charPositionNowInSrc] = charStructNow
		}
	} // for charPosition....

	return charactersInErlSrc, tokensInErlSrc
} // func character_loop
///////////////////////////////////////////////////

func token_opener_and_closer_look_forward__detect__whitespaces(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc) (int, bool, int, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, bool) {

	generalCharOpenerDetector := general_pattern__is_whitespace_rune
	generalCharNextAcceptableDetector := general_pattern__is_whitespace_rune
	tokenTypeIfActiveDetection := TokenType_id_WhitespaceInSrc

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
func general_pattern__is_whitespace_rune(r rune) bool {
	return (r == ' ' || r == '\r' || r == '\t' || r == '\n')
}

// this is a generic 'look forward' detector
func general_look_forward_accepted_chars_detector(
	charPositionNowInSrc int,                      //                      this opener uses ONLY the actual character,
	charactersInErlSrc CharacterInErlSrcCollector, // there is no need to look forward/back in src
	charStructNow CharacterInErlSrc,

// the first char rules are sometime different from the next char rules
	generalCharNowAcceptableDetector func(rune) bool,
	generalCharNextAcceptableDetector func(rune) bool,

	tokenTypeIfActiveDetection int) (int, bool, int, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, bool) {

	oneCharacterLongTokenDetection__openerAndCloserSameTime_closeDetectionImmediately := false
	// 0 means: there is no need to shift the original character loop position
	positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := 0

	tokenTypeId := TokenType_id_unknown
	funCloser := token_closer_fake_placeholder_fun_pattern
	openerDetected := false

	if generalCharNowAcceptableDetector(charStructNow.runeInErlSrc) {

		// all whitespaces are detected in this opener step, the detection process can be closed immediatelly
		oneCharacterLongTokenDetection__openerAndCloserSameTime_closeDetectionImmediately = true

		openerDetected = true
		funCloser = token_closer_fake_placeholder_fun_pattern // this is a general opener/closer fun,
		tokenTypeId = tokenTypeIfActiveDetection              // ^^^^ no need to user closer later

		charPositionNextInSrc := charPositionNowInSrc
		for true {
			charPositionNextInSrc++

			if charPositionNextInSrc < len(charactersInErlSrc) {
				charStructNext := charactersInErlSrc[charPositionNextInSrc]

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
	return tokenTypeId, openerDetected, positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected, funCloser, oneCharacterLongTokenDetection__openerAndCloserSameTime_closeDetectionImmediately
}

// /////////////////////////////////////////////////
func token_closer_fake_placeholder_fun_pattern(
	charPositionNowInSrc int,
	charactersInErlSrc CharacterInErlSrcCollector,
	charStructNow CharacterInErlSrc,
) bool {
	return false
}
