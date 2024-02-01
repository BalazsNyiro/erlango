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
			valueReceived, err := digitRune_decimalValue(runeElem)
			fmt.Println("decimal value test,", string(runeElem), valueWanted, valueReceived)

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

//  go test -v -run  Test_bignum_add_positive_positive
func Test_bignum_add_positive_positive(t *testing.T) {
	testName := "Test_bignum_add_positive_positive"

	bnA := bigNum_from_int(3)
	bnB := bigNum_from_int(5)
	bnResult := bigNum_operator_add(bnA, bnB)

	compare_bigNum_int(testName+"_3_5", 8, bnResult, t)
}

func compare_bigNum_int(testName string, wantedNum int, bn erlango_bignum_decimalValue, t *testing.T) {
	received := bigNum_convert_to_INT_for_testcases(bn)
	if received != wantedNum {
		t.Fatalf("\nError in %s wanted: %d, result: %d", testName, wantedNum, received)
	}

}
