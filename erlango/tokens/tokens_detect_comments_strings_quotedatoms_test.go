/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.3, third total rewrite
*/

package tokens

import (
	"fmt"
	"testing"
)

func print_string_runes_diff(txt1, txt2 string) {
	// the length has to be same
	if len(txt1) != len(txt2) {
		fmt.Println("different string lengths:", len(txt1), len(txt2))
	}
	for pos1, runeVal1 := range txt1{
		runeVal2 := []rune(txt2)[pos1]

		fmt.Println(pos1, runeVal1, string(runeVal1), "     ", runeVal2, string(runeVal2))
	}
}

func print_tokens(msg string, tokens Tokens) {
	fmt.Println(msg)
	for _, token := range tokens {
		fmt.Println("token print:", token.positionCharFirst, token.tokenType, string(token.charsInErlSrc))
	}
}


func Test_charsHowManyAreInTheGroup(t *testing.T) {
	funName := "Test_charsHowManyAreInTheGroup"

	setNumbers := []rune(ABC_Eng_digits)
	setAlphabetLow := []rune(ABC_Eng_Lower)
	setAlphabetUp := []rune(ABC_Eng_Upper)
	setCommas := []rune(",;")

	//    position:   0123456789ABCDEFG
	erlSrc := []rune(`ABabcdE123456 ,;:`)

	charsDetectedInGroup_counter := charsHowManyAreInTheGroup(0, erlSrc, setAlphabetLow, "right" )
	compare_int_int(funName + "_test_alphabet_0", charsDetectedInGroup_counter, 0, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(1, erlSrc, setAlphabetLow, "right" )
	compare_int_int(funName + "_test_alphabet_1", charsDetectedInGroup_counter, 0, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(2, erlSrc, setAlphabetLow, "right" )
	compare_int_int(funName + "_test_alphabet_2", charsDetectedInGroup_counter, 4, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(3, erlSrc, setAlphabetLow, "right" )
	compare_int_int(funName + "_test_alphabet_3", charsDetectedInGroup_counter, 3, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(6, erlSrc, setAlphabetUp, "right" )
	compare_int_int(funName + "_test_alphabet_4", charsDetectedInGroup_counter, 1, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(12, erlSrc, setNumbers, "left" )
	compare_int_int(funName + "_test_nums_1", charsDetectedInGroup_counter, 6, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(13, erlSrc, setNumbers, "left" )
	compare_int_int(funName + "_test_nums_2", charsDetectedInGroup_counter, 0, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(16, erlSrc, setCommas, "left" )
	compare_int_int(funName + "_test_nums_3", charsDetectedInGroup_counter, 0, t)

	charsDetectedInGroup_counter = charsHowManyAreInTheGroup(15, erlSrc, setCommas, "left" )
	compare_int_int(funName + "_test_nums_4", charsDetectedInGroup_counter, 2, t)
}

func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {
	testName := "01 simple atomQuoted, string, comment detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `[AtomVal, IntVal1, StringVal, IntVal2] = ['atomQuoted', 2, "txt", 4]. % comment` + "\n%comment2"
	erlSrcWantedAfterTokenDetect :=  `[AtomVal, IntVal1, StringVal, IntVal2] = [            , 2,      , 4].          ` + "\n         "

	erlSrc_received_after_tokenDetect, tokensTable_02_textBlocksDetected := Tokens_0_detect_text_blocks(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table" + testName, tokensTable_02_textBlocksDetected)


	fmt.Println("erlSrc, without strings, quoted atoms", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_textBlocksDetected)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_string_string(testName, "atomQuoted", tokensTable_02_textBlocksDetected[42].stringRepr(), t)
	compare_string_string(testName, "txt", tokensTable_02_textBlocksDetected[59].stringRepr(), t)
	compare_string_string(testName, " comment", tokensTable_02_textBlocksDetected[70].stringRepr(), t)
}

func Test_parse_comments_textDoubleQuoted_textSingleQuoted_escaping(t *testing.T) {
	testName := "02 escaping"

	tokensTable := Tokens{}
	erlSrcOrig :=                    "Atom='atom\nQuoted\\ escapeAgan: \\\t,  end:\\\\',"
	//                                           \n is one char in representation
	erlSrcWantedAfterTokenDetect :=  "Atom=                                      ,"

	erlSrc_received_after_tokenDetect, _ := Tokens_0_detect_text_blocks(erlSrcOrig, tokensTable)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
}


////////// GENERAL TEST FUNCTIONS //////////


