/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

func tokens_detect_02_erlang_whitespaces(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {

	for _, wantedCharInErl := range []rune{'\n', '\r', '\t', ' '} {
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, TokenType_id_WhitespaceInLine_ErlSrc)
	}
	character_loop__opener_closer_sections_set__if_more_separated_isolated_elems_are_next_to_each_other(charactersInErlSrc, TokenType_id_WhitespaceInLine_ErlSrc)

	charactersInErlSrc = character_loop__set_one_char_tokentype('\n', charactersInErlSrc, TokenType_id_WhitespaceNewLine_ErlSrc)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop__set_one_char_tokentype(
	wantedCharInErlSrc rune,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokenTypeId_ifDetected int,
) CharacterInErlSrcCollector {

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		if charStructNow.tokenIsDetected() {
			continue // if the char was detected and has a TokenType_id, there is no more to do.
		} //don't start new detection if the current char was detected once

		if charStructNow.runeInErlSrc == wantedCharInErlSrc {
			charStructNow.tokenDetectedType = tokenTypeId_ifDetected
			charStructNow.tokenOpenerCharacter = true // the current char is opener and closer first time,
			charStructNow.tokenCloserCharacter = true // later the same TokenTypes will be merged
			charactersInErlSrc[charPositionNowInSrc] = charStructNow
		}
	}

	return charactersInErlSrc
} // func character_loop patterns
///////////////////////////////////////////////////

// if there are more whitespaces next to each other, or more numbers for example, merge them into one group
func character_loop__opener_closer_sections_set__if_more_separated_isolated_elems_are_next_to_each_other(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokenTypeId_toMergeIntoOneOpenerCloserGroup int,
) CharacterInErlSrcCollector {

	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructPrev, charStructPrev_detectedCorrectly := charactersInErlSrc.char_get_by_index(charPositionNowInSrc - 1)
		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		charStructNext, charStructNext_detectedCorrectly := charactersInErlSrc.char_get_by_index(charPositionNowInSrc + 1)

		if charStructNow.tokenDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {

			if charStructPrev_detectedCorrectly {
				if charStructPrev.tokenDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {
					charStructNow.tokenOpenerCharacter = false
				}
			}

			if charStructNext_detectedCorrectly {
				if charStructNext.tokenDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {
					charStructNow.tokenCloserCharacter = false
				}
			}

			charactersInErlSrc[charPositionNowInSrc] = charStructNow
		} // current token type is the wanted
	}
	return charactersInErlSrc
} // func character_loop patterns
///////////////////////////////////////////////////
