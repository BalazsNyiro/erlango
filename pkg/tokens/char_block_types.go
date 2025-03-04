/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.



*/

package tokens

import "strconv"

const CharBlock_unknown = -1

// value 0 is not used, because 0 is the golang integer default value
// if a new char is created,
// and if the value is not updated, 0 is set by default,
// I would like to see if somebody changes the value or not.

const CharBlock_Comment = 1
const CharBlockQuotedSingle = 2
const CharBlockQuotedDouble = 3
const CharBlockQuotedTriple = 4

const CharBlock_braces_grouping_elems = 5
const CharBlock_dots_commas = 6
const CharBlock_LanguageElement_operators_specialchars = 8

// a number can be represented in a lot of forms.
// integer, float with dots, hexadeciamls - in first step,
// I simply try to see if the runes are part of any number, or not.
const CharBlock_AlphaNumeric = 10

// $A $B - this is a really special char representation,
// and can be detected directly, in a simple way
const CharBlock_Num_charLiterals = 20

const CharBlock_WhitespaceInLine_ErlSrc = 30
const CharBlock_WhitespaceNewLine_ErlSrc = 31

func CharBlockReprShort(wantedTokenTypeNum int) rune {
	the_key_is_not_defined_in_repr_map := -99
	var map_types_repr = map[int]rune{
		the_key_is_not_defined_in_repr_map:               'K',
		CharBlock_unknown:                                '?',
		CharBlock_Comment:                                '%',
		CharBlockQuotedSingle:                            '\'',
		CharBlockQuotedDouble:                            '"',
		CharBlockQuotedTriple:                            '3',
		CharBlock_AlphaNumeric:                           'A',
		CharBlock_Num_charLiterals:                       'L',
		CharBlock_WhitespaceInLine_ErlSrc:                'w',
		CharBlock_braces_grouping_elems:                  '(',
		CharBlock_dots_commas:                            '.',
		CharBlock_LanguageElement_operators_specialchars: 's',
	}
	repr, ok := map_types_repr[wantedTokenTypeNum]

	if ok {
		return repr
	} else { // if the wantedTokenType not in map, than that is unknown
		return map_types_repr[the_key_is_not_defined_in_repr_map]
	}
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

// in char_get_by_index() if the char position is overindexed/incorrect, an error is given back.
// here if the char cannot be indexed, a default value is returned.
// sometime the important question is only what is the next char (if it exists, and it is not important if it doesn't exist)
// and with this a lot of unnecessary condition/validation step can be avoided
func (collector CharacterInErlSrcCollector) char_get_by_index___give_fake_empty_space_char_if_no_real_char_in_position(index int) CharacterInErlSrc {

	if index < len(collector) && index >= 0 {
		return collector[index]
	}
	return CharacterInErlSrc{runeInErlSrc: ' '}
}

type CharacterInErlSrc struct {
	charBlockDetectedType    int // the id of the parent token in the file, IF the token is detected.
	runeInErlSrc             rune
	charBlockOpenerCharacter bool
	charBlockCloserCharacter bool
	positionInErlSrc         int // 0 based char position in the source

	// POSSIBLE OPTIONS:
	// file:<sourceComputerCoord>:/struct/in/computers/file/system
	// storageSystem:<kafkaOrDatabaseOrOtherStorageCoord>:/info/about/localization/of/the/data/in/the/system
	// insertedManually:<whyWhenHowWhereInserted>
	srcPathInformation_fromWhereIsItComing string
}

func (chr CharacterInErlSrc) stringRepr() string {
	return string(chr.runeInErlSrc)
}

func (chr CharacterInErlSrc) stringReprDetailed() string {
	return "characterOrigin:" + chr.srcPathInformation_fromWhereIsItComing + "positionInErlSrc:" + strconv.Itoa(chr.positionInErlSrc) + "  stringRepr: " + string(chr.runeInErlSrc)
}

func (chr CharacterInErlSrc) charBlockIsNotDetected() bool {
	return chr.charBlockDetectedType == CharBlock_unknown
}

func (chr CharacterInErlSrc) charBlockIsDetected() bool {
	return chr.charBlockDetectedType != CharBlock_unknown
}

func Runes_to_character_structs(runesAll []rune, srcPathInformation string) []CharacterInErlSrc {
	CharactersAll := []CharacterInErlSrc{}
	for posInErlSrc, runeInErlSrc := range runesAll {
		CharactersAll = append(CharactersAll, CharacterInErlSrc{
			charBlockDetectedType:                  CharBlock_unknown,
			runeInErlSrc:                           runeInErlSrc,
			charBlockOpenerCharacter:               false,
			charBlockCloserCharacter:               false,
			positionInErlSrc:                       posInErlSrc,
			srcPathInformation_fromWhereIsItComing: srcPathInformation,
		})
	}
	return CharactersAll
}

// unocode stop table, https://www.amp-what.com/unicode/search/stop
// octagonal sign, &#128721; //U+1f6d1
var unicode_stop_table = '🛑'

type errorMessages []string
