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
	"math/rand"
	"testing"
)


//  go test -v -run  Test_isEqual
func Test_isEqual(t *testing.T) {
	testName := "Test_isEqual"

	bn1 := bignum_decimalValue{digits: digitList{1,2,0,3,0,0}, exponent: 1, negative: false}
	bn2 := bignum_decimalValue{digits: digitList{1,2,0,3,0}, exponent: 2, negative: false}
	bn3 := bignum_decimalValue{digits: digitList{1,2,0,3}, exponent: 3, negative: false}

	compare_bool_bool(testName, true, bn1.isEqual(bn2), t)
	compare_bool_bool(testName, true, bn1.isEqual(bn3), t)
}

//  go test -v -run   Test_normaliseExponent_endingZerosRemove
func Test_normaliseExponent_endingZerosRemove(t *testing.T) {
	testName := "Test_normaliseExponent_endingZerosRemove"

	// leading unsignificant 0 will be removed, too
	bn := bignum_decimalValue{digits: digitList{0,1,2,0,3,0,0}, exponent: 1, negative: false}
	bnNormalised := bn.normalisedForm_endingZerosIntoExponent()
	digitsWantedAfterReverse := digitList{1,2,0,3}

	compare_digits_digits(testName, digitsWantedAfterReverse, bnNormalised.digits, t)
	compare_int_int(testName, 3, bnNormalised.exponent, t)

}

func Test_digitsCleaning_leadingZerosRemoval(t *testing.T) {
	testName := "Test_digitsCleaning_leadingZerosRemoval"
	// leading unsignificant 0 will be removed, too
	digits := digitList{0,0,0,1,2,0,3,0,0}
	digitsLeadingZerosRemoved := digitsCleaning_leadingZerosRemoval(digits)
	digitsWanted := digitList{1,2,0,3,0,0}
	compare_digits_digits(testName, digitsLeadingZerosRemoved, digitsWanted, t)
}


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

//  go test -v -run   Test_isLessThan
func Test_isLessThan(t *testing.T) {

	lessTest := func(a, b int) {
		bignumA := bigNum_from_int(a)
		bignumB := bigNum_from_int(b)
		aLessThanOther := bignumA.isLessThan(bignumB)
		testName := fmt.Sprintf("Test_isLessThan %d %d", a, b)
		compare_bool_bool(testName, a<b, aLessThanOther, t)
	}

	lessTest(8,13)
	lessTest(13,8)

	lessTest(0,0)

	lessTest(0,1)
	lessTest(1,0)
	lessTest(0,-1)
	lessTest(-1,0)


	lessTest(-5,-8)
	lessTest(-8,-5)

	lessTest(-5,+5)
	lessTest(+5,-5)
}






//  go test -v -run   Test_digits_reverse
func Test_digits_reverse(t *testing.T) {
	testName := "Test_digits_reverse"

	one := digitElemType(1)
	two := digitElemType(2)
	three := digitElemType(3)
	four := digitElemType(4)

	digits := digitList{four, three, two, one}
	digitsReversed := digitsReverse(digits)

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

	compare_int_int(testName, aUpdated.exponent, bUpdated.exponent, t)
	aUpdatedWanted := bignum_decimalValue{digits: digitList{4,0,0,0,0,0}, exponent: -1, negative: false}

	compare_digits_digits(testName, aUpdatedWanted.digits, aUpdated.digits, t)
}



//  go test -v -run  Test_bignum_operators
func Test_bignum_operators(t *testing.T) {
	operations := []string{"sub", "add"}

	for _, op := range operations {
		// zero, negative
		operator_test(op,  0, -2, t)
		operator_test(op, -2 , 0, t)

		// zero, positive
		operator_test(op,  0,  2, t)
		operator_test(op,  2 , 0, t)

		// negative, positive
		operator_test(op, -2, +2, t)
		operator_test(op, +2, -2, t)

		// negative, negative
		operator_test(op, -2, -3, t)
		operator_test(op, +2, +3, t)

		// positive, positive
		operator_test(op, 3, 5, t)
		operator_test(op, 0, 9, t)
		operator_test(op, 10, 1, t)
		operator_test(op, 19, 1, t)
		operator_test(op, 99, 1, t)
		operator_test(op, 333, 4444, t)

		// I saw problems with these in random tests:
		operator_test(op,-908, 105, t)
	}

	// RANDOM MATH TESTS against all operations /////////////////////////////////////////
	sign := +1
	limit := 2000
	for a := -limit; a < limit; a++ {

		if rand.Intn(limit) % 2 == 1 {
			sign = -1
		} else {
			sign = +1
		}

		b := sign * rand.Intn(limit)

		for _, op := range operations {
			operator_test(op, a, b, t)
		}
	}

}

// if something is wrong, debug it here:
//  go test -v -run  Test_bignum_debug
func Test_bignum_debug(t *testing.T) {
	// operator_test("sub", 10, 1, t)
	operator_test("mul", 10, 1, t)
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

	if math_operator == "mul" {
		intResult = a * b
		bnResult = bigNum_operator_mul(bnA, bnB)
	}


	testName := fmt.Sprintf("math_operator_test__%s__%d_%d", math_operator, a, b)
	compare_bigNum_int(testName, intResult, bnResult, t)

	// and one more test, for  bigNum_convert_to_INT_for_testcases
	compare_int_int(testName, bigNum_convert_to_INT_for_testcases(bnResult), intResult, t)

}



func compare_bigNum_int(testName string, wantedNum int, bn bignum_decimalValue, t *testing.T) {
	received := bigNum_convert_to_INT_for_testcases(bn)
	compare_int_int(testName, wantedNum, received, t)
}

func compare_int_int(testName string, wantedNum int, received int, t *testing.T) {
	if wantedNum != received {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wantedNum, received)
	}
}

func compare_digitElem_digitElem(testName string, wanted digitElemType, received digitElemType, t *testing.T) {
	if wanted != received {
		t.Fatalf("\nError in %s wanted: %d, received: %d", testName, wanted, received)
	}
}

func compare_digits_digits(testName string, wanted digitList, received digitList, t *testing.T) {
	if len(wanted) != len(received) {
		t.Fatalf("\nError, different LENGTH in digit list comparison %s wanted: %d, received: %d", testName, wanted, received)
	}
	for id, _ := range received {
		compare_digitElem_digitElem(testName, wanted[id], received[id], t)
	}
}
func compare_bool_bool(testName string, wanted bool, received bool, t *testing.T) {
	if wanted != received {
		t.Fatalf("\nError, different bool comparison %s wanted: %t, received: %t", testName, wanted, received)
	}
}
