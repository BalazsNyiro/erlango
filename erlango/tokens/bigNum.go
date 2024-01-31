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

type digitElemType int8
type digitList []digitElemType

// similar erlang lib: https://github.com/mabrek/erlang-decimal
type erlango_bignum_decimalValue struct {
	// the bignum is ALWAYS decimal number, with separated digits representation

	/* example num representation: '12.34'
		digits: 1234
		exponent: -2    because 1234 * (10^-2) = 12.34

	   example 2: 1234000
		digits: 1234
		exponent: 3     because 1234 * (10^3) = 1234000

		with this representation, integers and floats can be represented too, with the same structure
	*/

	digits digitList  // decimal digits
	exponent int  // where is the ., or how many 0 are after the digits
	negative bool // false by default: is the number negative?
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
