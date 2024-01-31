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

import "errors"

const ABC_Eng_Upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ABC_Eng_Lower string = "abcdefghijklmnopqrstuvwxyz"
const ABC_Eng_digits string = "0123456789"

type digitElemType int8
type digitList []digitElemType

// similar erlang lib: https://github.com/mabrek/erlang-decimal
// general Decimal arithmetic: https://speleotrove.com/decimal/daops.html
// the precisions is not used in my implementation.
type erlango_bignum_decimalValue struct {
	// the bignum is ALWAYS a 10 based number, with separated digits representation

	/* example num representation: '12.34'
		digits: [1,2,3,4]
		exponent: -2    because 1234 * (10^-2) = 12.34

	   example 2: 1234000
		digits: [1,2,3,4]
		exponent: 3     because 1234 * (10^3) = 1234000

		with this representation, integers and floats can be represented too, with the same structure
	*/

	digits digitList  // decimal digits
	exponent int  // where is the ., or how many 0 are after the digits
	negative bool // false by default: is the number negative?
}

func (bn erlango_bignum_decimalValue) isPositive() bool {
	return ! (bn.negative)
}

func (bn erlango_bignum_decimalValue) isNegative() bool {
	return bn.negative
}


func digitRune_decimalValue(digit rune) (digitElemType, error) {
	valueMap := map[rune]digitElemType{
		'0'	: 0,	'a' : 10,   'A' : 10,
		'1' : 1,	'b' : 11,   'B' : 11,
		'2' : 2,	'c' : 12,   'C' : 12,
		'3' : 3,	'd' : 13,   'D' : 13,
		'4' : 4,	'e' : 14,   'E' : 14,
		'5' : 5,	'f' : 15,   'F' : 15,
		'6' : 6,	'g' : 16,   'G' : 16,
		'7' : 7,	'h' : 17,   'H' : 17,
		'8' : 8,	'i' : 18,   'I' : 18,
		'9' : 9,	'j' : 19,   'J' : 19,
		'k' : 20,   'K' : 20,
		'l' : 21,   'L' : 21,
		'm' : 22,   'M' : 22,
		'n' : 23,   'N' : 23,
		'o' : 24,   'O' : 24,
		'p' : 25,   'P' : 25,
		'q' : 26,   'Q' : 26,
		'r' : 27,   'R' : 27,
		's' : 28,   'S' : 28,
		't' : 29,   'T' : 29,
		'u' : 30,   'U' : 30,
		'v' : 31,   'V' : 31,
		'w' : 32,   'W' : 32,
		'x' : 33,   'X' : 33,
		'y' : 34,   'Y' : 34,
		'z' : 35,   'Z' : 35,
	}
	val, digitInMap := valueMap[digit]
	if ! digitInMap {
		return 0, errors.New("digit value is not detected (" + string(digit) + ")")
	}
	return val, nil
}


// this is important when digit values are checked, in non-decimal char processing
func bigNum_from_digitValue__min0_max35(decimalVal digitElemType) erlango_bignum_decimalValue {
	// a digit's value is minimum 0, maximum 35. there is no problem with too big integer values
	digit_2 := decimalVal % 10
	digit_1 := (decimalVal - digit_2) / 10
	// simple conversion, NEVER normalize here - normalization can happen at one point only
	if digit_1 > 0 {
		return erlango_bignum_decimalValue{digits: digitList{digit_1, digit_2}, exponent: 0}
	} else {
		return erlango_bignum_decimalValue{digits: digitList{digit_2}, exponent: 0}
	}
}

func bigNum_add(a, b erlango_bignum_decimalValue) erlango_bignum_decimalValue {
	if a.isPositive() && b.isNegative() {
		b.negative = false
		return bigNum_sub(a, b)
	}
	if a.isNegative() && b.isPositive(){
		a.negative = false
		return bigNum_sub(b, a)
	}
	if a.isNegative() && b.isNegative() {
		a.negative = false
		b.negative = false
		result := bigNum_add_positive_positive(a, b)
		result.negative = true
		return result
	}

	// basic case: a, b are positive
	return bigNum_add_positive_positive(a, b)
}


func bigNum_sub(a, b erlango_bignum_decimalValue) erlango_bignum_decimalValue {
	if !a.negative && b.negative {
		// FIXME: WRONG
		b.negative = false
		return bigNum_add(a, b)
	}
	if a.negative && !b.negative {
		a.negative = false
		result := bigNum_add(a, b)
		result.negative = true
		return result
	}
	if a.negative && b.negative {
		a.negative = false
		b.negative = false
		result := bigNum_add(a, b)
		result.negative = true
		return result
	}

	// basic case: a, b are positive
	return bigNum_sub_positive_positive(a, b)

}


func bigNum_add_positive_positive(a, b erlango_bignum_decimalValue) erlango_bignum_decimalValue {
	// the bignum is ALWAYS decimal number, with separated digits representation

	// add operation can be done ONLY if the exponents are same
	a, b = bigNum_pair_set_same_exponent(a, b)

	digitsReversed := digitList{}

	var overflow digitElemType = 0
	position := -1

	for {
		position++

		var valueA digitElemType = 0  // a decimal digit value is between 0-9, so a byte can store that
		var valueB digitElemType = 0

		if position < len(a.digits) {
			valueA = a.digits[position]
		}
		if position < len(b.digits) {
			valueB = b.digits[position]
		}

		if valueA == 0 && valueB == 0 && overflow == 0 {
			break // exit if there is no more thing to do
		}

		valueSum := valueA + valueB + overflow
		digitNew := valueSum % 10

		overflow = (valueSum - digitNew) / 10
	}

	summa := erlango_bignum_decimalValue{digits: digits_reverse(digitsReversed), exponent: a.exponent}
	return summa
}


// FIXME: this is not OK, completely rewrite this:
func bigNum_sub_positive_positive(a, b erlango_bignum_decimalValue) erlango_bignum_decimalValue {
	// the bignum is ALWAYS decimal number, with separated digits representation

	// add operation can be done ONLY if the exponents are same
	a, b = bigNum_pair_set_same_exponent(a, b)

	digitsReversed := digitList{}

	var overflow digitElemType = 0
	position := -1

	for {
		position++

		var valueA digitElemType = 0  // a decimal digit value is between 0-9, so a byte can store that
		var valueB digitElemType = 0

		if position < len(a.digits) {
			valueA = a.digits[position]
		}
		if position < len(b.digits) {
			valueB = b.digits[position]
		}

		if valueA == 0 && valueB == 0 && overflow == 0 {
			break // exit if there is no more thing to do
		}

		valueSum := valueA + valueB + overflow
		digitNew := valueSum % 10

		overflow = (valueSum - digitNew) / 10
	}

	summa := erlango_bignum_decimalValue{digits: digits_reverse(digitsReversed), exponent: a.exponent}
	return summa
}
