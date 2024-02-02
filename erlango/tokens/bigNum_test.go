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

//  go test -v -run   Test_digits_reverse
func Test_digits_reverse(t *testing.T) {
	testName := "Test_digits_reverse"

	one := digitElemType(1)
	two := digitElemType(2)
	three := digitElemType(3)
	four := digitElemType(4)

	digits := digitList{four, three, two, one}
	digitsReversed := digits_reverse(digits)

	digitsWantedAfterReverse := digitList{one, two, three, four}

	compare_digits_digits(testName, digitsWantedAfterReverse, digitsReversed, t)
}

//  go test -v -run   Test_bigNum_pair_set_same_exponent
func Test_bigNum_pair_set_same_exponent(t *testing.T) {
	testName := "Test_bigNum_pair_set_same_exponent"

	a := bignum_decimalValue{digits: digitList{4,0}, exponent: 3, negative: false}
	b := bignum_decimalValue{digits: digitList{1,2,3,0,0}, exponent: -1, negative: false}

	fmt.Printf("Test 1, address a %p\n", &a)
	fmt.Printf("Test 1, address a.digits[0] %p\n", &(a.digits[0]))
	aUpdated, bUpdated := bigNum_pair_setSameExponent_decreaseBiggerExponent(a, b)

	compare_ints(testName, aUpdated.exponent, bUpdated.exponent, t)
	aUpdatedWanted := bignum_decimalValue{digits: digitList{4,0,0,0,0,0}, exponent: -1, negative: false}

	compare_digits_digits(testName, aUpdatedWanted.digits, aUpdated.digits, t)
}




//  go test -v -run  Test_bignum_add_positive_positive
func Test_bignum_add_positive_positive(t *testing.T) {
	operator_test("add", 3, 5, t)
	operator_test("add", 0, 9, t)
	operator_test("add", 10, 1, t)
	operator_test("add", 19, 1, t)
	operator_test("add", 99, 1, t)
	operator_test("add", 333, 4444, t)
}


func operator_test(math_operator string, a, b int, t *testing.T) {
	bnA := bigNum_from_int(a)
	bnB := bigNum_from_int(b)

	intResult := 0
	bnResult := bigNum_zero()

	if math_operator == "add" {
		intResult = a + b
		bnResult = bigNum_operator_add(bnA, bnB)
	}

	if math_operator == "sub" {
		intResult = a - b
		bnResult = bigNum_operator_sub(bnA, bnB)
	}

	testName := fmt.Sprintf("math_operator_test__%s__%d_%d", math_operator, a, b)
	compare_bigNum_int(testName, intResult, bnResult, t)

	// and one more test, for  bigNum_convert_to_INT_for_testcases
	compare_int_int(testName, bigNum_convert_to_INT_for_testcases(bnResult), intResult, t)

}



func compare_bigNum_int(testName string, wantedNum int, bn bignum_decimalValue, t *testing.T) {
	received := bigNum_convert_to_INT_for_testcases(bn)
	if received != wantedNum {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_int_int(testName string, wantedNum int, received int, t *testing.T) {
	if received != wantedNum {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_digitElem_digitElem(testName string, wanted digitElemType, received digitElemType, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wanted, received)
	}
}

func compare_digits_digits(testName string, wanted digitList, received digitList, t *testing.T) {
	if len(received) != len(wanted) {
		t.Fatalf("\nError, different LENGTH in digit list comparison %s wanted: %d, received: %d", testName, wanted, received)
	}
	for id, _ := range received {
		compare_digitElem_digitElem(testName, wanted[id], received[id], t)
	}
}
