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


func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {
	testName := "01 simple atomQuoted, string, comment detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `[AtomVal, IntVal1, StringVal, IntVal2] = ['atomQuoted', 2, "txt", 4]. % comment` + "\n%comment2"
	erlSrcWantedAfterTokenDetect :=  `[AtomVal, IntVal1, StringVal, IntVal2] = [            , 2,      , 4].          ` + "\n         "

	erlSrc_received_after_tokenDetect, tokensTable_02_textBlocksDetected := Tokens_detect_text_blocks(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table, after comment, atomQuoted, string detect", tokensTable_02_textBlocksDetected)


	fmt.Println("erlSrc, without strings, quoted atoms", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_textBlocksDetected)

	compare_strings(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_strings(testName, "atomQuoted", tokensTable_02_textBlocksDetected[42].stringRepr(), t)
	compare_strings(testName, "txt", tokensTable_02_textBlocksDetected[59].stringRepr(), t)
	compare_strings(testName, " comment", tokensTable_02_textBlocksDetected[70].stringRepr(), t)
}

func compare_strings(callerInfo, strWanted, strReceived string, t *testing.T) {
	if strWanted != strReceived {
		t.Fatalf("\nErr String difference (%s):\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, strWanted, strReceived)
	}
}

/*
func compare_tokenDetected_tokenWanted(callerInfo string, tokensDetected ErlTokens, tokenWanted tokenWanted, t *testing.T) {
	tokenDetected, tokenWantedIsInDetected:= tokensDetected[tokenWanted.positionFirst]

	if tokenWantedIsInDetected {
		// theoretically the charPosFirst is always ok here, because the key in map was the same position
		tokenDetected_charPosFirst, tokenDetected_charPosLast := tokenDetected.charPositionFirstLast()
		if tokenDetected_charPosFirst != tokenWanted.positionFirst {
			t.Fatalf("\nErr First: %s : detected posFirst: %v  is different from wanted posFirst:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected_charPosLast != tokenWanted.positionLast {
			t.Fatalf("\nErr Last: %s : detected posLast: %v  is different from wanted posLast:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosLast, tokenWanted.positionLast, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected.stringRepresentation() != tokenWanted.textRepresentation {
			t.Fatalf("\nErr repr %s : startPos:%v  detected string representation: %v  is different from wanted representation:  %v, error",
				callerInfo, tokenDetected_charPosFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
	} else {
		t.Fatalf("\nErr %s : wanted tokenPosFirst %v is not in detecteds - error", callerInfo, tokenWanted.positionFirst)
	}

}

 */
