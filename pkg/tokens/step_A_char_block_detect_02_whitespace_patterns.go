/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

func character_block_detect__02_erlang_whitespaces(charactersInErlSrc CharacterInErlSrcCollector) CharacterInErlSrcCollector {

	for _, wantedCharInErl := range []rune{'\r', '\t', ' '} {
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, CharBlock_WhitespaceInLine_ErlSrc)
	}
	character_loop__opener_closer_sections_set__if_more_separated_isolated_elems_are_next_to_each_other(charactersInErlSrc, CharBlock_WhitespaceInLine_ErlSrc)

	charactersInErlSrc = character_loop__set_one_char_tokentype('\n', charactersInErlSrc, CharBlock_WhitespaceNewLine_ErlSrc)
	return charactersInErlSrc
}

func character_loop__set_one_char_tokentype(
	wantedCharInErlSrc rune,
	charactersInErlSrc CharacterInErlSrcCollector,
	tokenTypeId_ifDetected int,
) CharacterInErlSrcCollector {

	// use the slice position only, because in the for loop, charactersInErlSrc will be updated/modified,
	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]
		if charStructNow.charBlockIsDetected() {
			continue // if the char was detected and has a TokenType_id, there is no more to do.
		} //don't start new detection if the current char was detected once

		if charStructNow.runeInErlSrc == wantedCharInErlSrc {
			charStructNow.charBlockDetectedType = tokenTypeId_ifDetected
			charStructNow.charBlockOpenerCharacter = true // the current char is opener and closer first time,
			charStructNow.charBlockCloserCharacter = true // later the same TokenTypes will be merged
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

		if charStructNow.charBlockDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {

			if charStructPrev_detectedCorrectly {
				if charStructPrev.charBlockDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {
					charStructNow.charBlockOpenerCharacter = false
				}
			}

			if charStructNext_detectedCorrectly {
				if charStructNext.charBlockDetectedType == tokenTypeId_toMergeIntoOneOpenerCloserGroup {
					charStructNow.charBlockCloserCharacter = false
				}
			}

			charactersInErlSrc[charPositionNowInSrc] = charStructNow
		} // current token type is the wanted
	}
	return charactersInErlSrc
} // func character_loop patterns
///////////////////////////////////////////////////
