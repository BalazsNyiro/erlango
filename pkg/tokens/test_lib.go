/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package tokens

import (
	"fmt"
	"testing"
)

//////////////////////////////////////////////////////////////////

func compare_int__int_(testName string, wantedNum int, received int, t *testing.T) {
	if wantedNum != received {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_bool_bool(testName string, wanted bool, received bool, t *testing.T) {
	if wanted != received {
		t.Fatalf("\nError, different bool comparison %s wanted: %t, received: %t", testName, wanted, received)
	}
}

func compare_string_string(callerInfo, strWanted, strReceived string, t *testing.T) {
	if strWanted != strReceived {
		t.Fatalf("\nErr String difference (%s):\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, strWanted, strReceived)
	}
}

func compare_runes_runes(callerInfo string, runesWanted, runesReceived []rune, t *testing.T) {
	errMsg := fmt.Sprintf("\nErr (%s) []rune <>[]rune:\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, string(runesWanted), string(runesReceived))
	if len(runesWanted) != len(runesReceived) {
		t.Fatalf(errMsg)
		return
	}

	for pos, runeWanted := range runesWanted {
		if runeWanted != runesReceived[pos] {
			t.Fatalf(errMsg)
			return
		}
	}
}

func compare_rune_rune(callerInfo string, runeWanted, runeReceived rune, t *testing.T) {
	if runeWanted != runeReceived {
		errMsg := fmt.Sprintf("\nErr (%s) rune <>rune:\n  wanted -->>%s<<-- ??\nreceived -->>%s<<--\n\n", callerInfo, string(runeWanted), string(runeReceived))
		t.Fatalf(errMsg)
	}
}

func is_string_contains_only_0123456789(txt string) bool {
	return is_string_contains_only_allowed_letters(txt, "0123456789")
}

func is_string_contains_only_0123456789Dot(txt string) bool {
	return is_string_contains_only_allowed_letters(txt, "0123456789.")
}

// example: txt = 12.34, allowedChars = "0123456789" - dot is not allowed string elem
func is_string_contains_only_allowed_letters(txt string, allowedCharsInString string) bool {

	for _, textLetter := range txt {
		letterDetectedInAlloweds := false
		for _, allowedLetter := range allowedCharsInString {
			if textLetter == allowedLetter {
				letterDetectedInAlloweds = true
				break
			}
		}

		// if we got the first non-allowed textLetter, return with false
		if !letterDetectedInAlloweds {
			return false
		}
	}

	return true // every letter in txt was allowed
}
