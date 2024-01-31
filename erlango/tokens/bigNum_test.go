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


// the digits and char values has to be correct
func Test_digitRune_decimalValue(t *testing.T) {
	testName := "Test_digitRune_decimalValue"

	test_string := func (txt string) {
		for valueWanted, runeElem := range txt {
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

	test_string(ABC_Eng_digits + ABC_Eng_Lower)
	test_string(ABC_Eng_digits + ABC_Eng_Upper)
}


func Test_bigNum_from_digitVal__min0_max35(t *testing.T) {
	testName := "Test_bigNum_from_digitVal__min0_max35"

	for numToConvert := 0; numToConvert <=35; numToConvert++ {
		bigNumReceived := bigNum_from_digitValue__min0_max35(digitElemType(numToConvert))

		valueReceived := -1
		if len(bigNumReceived.digits) == 1 {
			valueReceived = int(bigNumReceived.digits[0])
		}
		if len(bigNumReceived.digits) == 2 {
			valueReceived = int(bigNumReceived.digits[0]*10)
			valueReceived += int(bigNumReceived.digits[1])
		}

		if valueReceived != numToConvert {
			t.Fatalf("\nError in numToConvert conversion to bigNum %s value wanted: %d, valueReceived: %d", testName, numToConvert, valueReceived)
		}
	}

}
