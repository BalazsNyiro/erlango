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
	charactersInErlSrc, tokensInErlSrc = character_loop(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
	return charactersInErlSrc, tokensInErlSrc
}

func tokens_detect_erlang_whitespaces(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect__whitespaces_inside_lines
	printVerboseOpenerDetectMsg := true
	charactersInErlSrc, tokensInErlSrc = character_loop(charactersInErlSrc, tokensInErlSrc, funTokenOpener, printVerboseOpenerDetectMsg)
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
				// do an opening OR a closing, because the parsing's first step was string/comment/quotedAtom detection,
				// and in those cases there are well-defined and differently positioned opening/closing elems.

				// so, if you have a token which needs to be detected in one step (operators for example)
				// then you need to use the opening+closing method, not the separated opening/closing funs

				// I have the feeling that there is a way to convert a mixed solution to be separated,
				// but maybe different approach can be used when it is guaranteed that the token is minimum 2 char long, or not.

				// if the token is one char wide, this special section is added to do an immediate closing:

				if openerAndCloserSameTime_closeDetectionImmediately {
					charStructNow.tokenCloserCharacter = true                               // close the last charStructNow elem,
					charactersInErlSrc[charPositionNowInSrc] = charStructNow                // if the previous loop updated more chars.
					set_noActiveTokenDetection__tokenTypeUnknown__cleaningAfterTokenClose() // maybe that is not the starter one,
				}
			} ////////////////////////////////////////////////

		} else { // activeTokenDetectionBecauseOpenerConditionTriggered == true:
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

func token_opener_detect__quoteDouble__quoteSinge_comment(
	charPositionNowInSrc int,                      //                      this opener uses ONLY the actual character,
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

func token_opener_detect__whitespaces_inside_lines(
	charPositionNowInSrc int,                      //                      this opener uses ONLY the actual character,
	charactersInErlSrc CharacterInErlSrcCollector, // there is no need to look forward/back in src
	charStructNow CharacterInErlSrc) (int, bool, int, func(int, CharacterInErlSrcCollector, CharacterInErlSrc) bool, bool) {

	oneCharacterLongTokenDetection__openerAndCloserSameTime_closeDetectionImmediately := false
	positionModifierBecauseLongerThanOneTokenOpenerCharsAreDetected := 0

	tokenTypeId := TokenType_id_unknown
	funCloser := token_closer_fake_placeholder_fun
	openerDetected := false

	// 0: double quote " opener is 1 char wide,
	//there is no need to shift the original character loop position

	is_whitespace_rune := func(r rune) bool {
		// I don't want to mix newlines here - that is a line separator, these can be in the middle of the src.
		return (r == ' ' || r == '\r' || r == '\t')
	}

	if is_whitespace_rune(charStructNow.runeInErlSrc) {

		// all whitespaces are detected in this opener step, the detection process can be closed immediatelly
		oneCharacterLongTokenDetection__openerAndCloserSameTime_closeDetectionImmediately = true

		openerDetected = true
		funCloser = token_closer_fake_placeholder_fun
		tokenTypeId = TokenType_id_WhitespaceInSrc

		charPositionNextInSrc := charPositionNowInSrc
		for true {
			charPositionNextInSrc++

			if charPositionNextInSrc < len(charactersInErlSrc) {
				charStructNext := charactersInErlSrc[charPositionNextInSrc]

				if is_whitespace_rune(charStructNext.runeInErlSrc) {
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
