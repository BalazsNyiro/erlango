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

// value 0 is not used, because 0 is the golang integer default value
// if a new char is created,
// and if the value is not updated, 0 is set by default,
// I would like to see if somebody changes the value or not.

const TokenType_id_Comment = 1
const TokenType_id_TextBlockQuotedSingle = 2
const TokenType_id_TextBlockQuotedDouble = 3
const TokenType_id_TextBlockQuotedTriple = 4

const TokenType_id_braces_grouping_elems = 5
const TokenType_id_dots_commas = 6
const TokenType_id_LanguageElement_operators_specialchars = 8

// a number can be represented in a lot of forms.
// integer, float with dots, hexadeciamls - in first step,
// I simply try to see if the runes are part of any number, or not.
const TokenType_id_AlphaNumeric = 10

// $A $B - this is a really special char representation,
// and can be detected directly, in a simple way
const TokenType_id_Num_charLiterals = 20

const TokenType_id_WhitespaceInLine_ErlSrc = 30
const TokenType_id_WhitespaceNewLine_ErlSrc = 31

func TokenTypeReprShort(wantedTokenTypeNum int) rune {
	the_key_is_not_defined_in_repr_map := -99
	var map_types_repr = map[int]rune{
		the_key_is_not_defined_in_repr_map:                  'K',
		TokenType_id_unknown:                                '?',
		TokenType_id_Comment:                                '%',
		TokenType_id_TextBlockQuotedSingle:                  '\'',
		TokenType_id_TextBlockQuotedDouble:                  '"',
		TokenType_id_TextBlockQuotedTriple:                  '3',
		TokenType_id_AlphaNumeric:                           'A',
		TokenType_id_Num_charLiterals:                       'L',
		TokenType_id_WhitespaceInLine_ErlSrc:                'w',
		TokenType_id_braces_grouping_elems:                  '(',
		TokenType_id_dots_commas:                            '.',
		TokenType_id_LanguageElement_operators_specialchars: 's',
	}
	repr, ok := map_types_repr[wantedTokenTypeNum]

	if ok {
		return repr
	} else { // if the wantedTokenType not in map, than that is unknown
		return map_types_repr[the_key_is_not_defined_in_repr_map]
	}
}

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

// return with wanted char from the collection, or with a default value if the char doesn't exist - and
func (collector CharacterInErlSrcCollector) char_get_by_index(index int) (CharacterInErlSrc, bool) {

	var charInErl CharacterInErlSrc
	indexIsExistInValidRange__charWasDetectedCorrectly__not_over_or_under_indexed := false

	if index < len(collector) && index >= 0 {
		indexIsExistInValidRange__charWasDetectedCorrectly__not_over_or_under_indexed = true
		charInErl = collector[index]
	}

	return charInErl, indexIsExistInValidRange__charWasDetectedCorrectly__not_over_or_under_indexed
}

func (collector CharacterInErlSrcCollector) char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(index int) CharacterInErlSrc {

	if index < len(collector) && index >= 0 {
		return collector[index]
	}
	return CharacterInErlSrc{runeInErlSrc: ' '}
}

type CharacterInErlSrc struct {
	tokenDetectedType    int // the id of the parent token in the file, IF the token is detected.
	runeInErlSrc         rune
	tokenOpenerCharacter bool
	tokenCloserCharacter bool
	positionInErlSrc     int // 0 based char position in the source
}

func (chr CharacterInErlSrc) stringRepr() string {
	return string(chr.runeInErlSrc)
}

func (chr CharacterInErlSrc) tokenNotDetected() bool {
	return chr.tokenDetectedType == TokenType_id_unknown
}

func (chr CharacterInErlSrc) tokenIsDetected() bool {
	return chr.tokenDetectedType != TokenType_id_unknown
}

func Runes_to_character_structs(runesAll []rune) []CharacterInErlSrc {
	CharactersAll := []CharacterInErlSrc{}
	for posInErlSrc, runeInErlSrc := range runesAll {
		CharactersAll = append(CharactersAll, CharacterInErlSrc{
			tokenDetectedType:    TokenType_id_unknown,
			runeInErlSrc:         runeInErlSrc,
			tokenOpenerCharacter: false,
			tokenCloserCharacter: false,
			positionInErlSrc:     posInErlSrc,
		})
	}
	return CharactersAll
}

// unocode stop table, https://www.amp-what.com/unicode/search/stop
// octagonal sign, &#128721; //U+1f6d1
var unicode_stop_table = '🛑'
