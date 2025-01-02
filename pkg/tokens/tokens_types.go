/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.



*/

package tokens

const TokenType_id_unknown = -1

const TokenType_id_Comment = 0
const TokenType_id_TextBlockQuotedSingle = 1
const TokenType_id_TextBlockQuotedDouble = 2

// numbers are in 1X block
const TokenType_id_Num_int = 10
const TokenType_id_Num_float = 11

// $A $B
const TokenType_id_Num_charLiterals = 20

type TokenCollector []TokenInErlSrc

type TokenInErlSrc struct {
	tokenTypeId int

	// "quoted" string's first " is the tokens' first position!
	// 'atoms'  fist ' char is the tokens first pos!
	// so ALL character is included from the src into the token position range
	// 0 based position: the first char pos in the whole src is 0.
	// the first/last positions means included positions.
	// so first is the included first char, last is the included last char.
	// if the token length == 1, the positionCharFirst == positionCharLast
	positionCharFirst int
	positionCharLast  int

	sourceOfCode string // can be a file, a user-typed-terminal-input,
	// or maybe a dynamically generated, then evaluated string
	charsInErlSrc []rune
}

type CharacterInErlSrcCollector []CharacterInErlSrc

type CharacterInErlSrc struct {
	tokenDetectedForThisCharacterInErlSrc bool
	tokenIdNumInSrc                       int // the id of the parent token in the file, IF the token is detected.
	runeInErlSrc                          rune
}

func (chr CharacterInErlSrc) stringRepr() string {
	return string(chr.runeInErlSrc)
}

func Runes_to_character_structs(runesAll []rune) []CharacterInErlSrc {
	CharactersAll := []CharacterInErlSrc{}
	for _, runeInErlSrc := range runesAll {
		CharactersAll = append(CharactersAll, CharacterInErlSrc{
			tokenDetectedForThisCharacterInErlSrc: false,
			tokenIdNumInSrc:                       -1, // no valid real token for this character
			runeInErlSrc:                          runeInErlSrc,
		})
	}
	return CharactersAll
}
