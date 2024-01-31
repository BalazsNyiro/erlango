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

func Test_digitRune_decimalValue(t *testing.T) {
	testName := "Test_digitRune_decimalValue"

	for valueWanted, runeElem := range ABC_Eng_digits + ABC_Eng_Lower {
		fmt.Println("decimal value test,", valueWanted, runeElem)

		valueReceived, err := digitRune_decimalValue(runeElem)
		if err != nil {
			t.Fatalf("\nError in %s, err is not nil, rune value detect (%s) value wanted: %d, valueReceived: %d", testName, string(runeElem), valueWanted, valueReceived)
		}
		if digitElemType(valueWanted) != valueReceived {
			t.Fatalf("\nError in %s, rune (%s) value wanted: %d, valueReceived: %d", testName, string(runeElem), valueWanted, valueReceived)
		}
	}
}
