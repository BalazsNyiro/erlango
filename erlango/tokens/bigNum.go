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
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const ABC_Eng_Upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ABC_Eng_Lower string = "abcdefghijklmnopqrstuvwxyz"
const ABC_Eng_digits string = "0123456789"

type digitElemType int8
type digitList []digitElemType




// similar erlang lib: https://github.com/mabrek/erlang-decimal
// general Decimal arithmetic: https://speleotrove.com/decimal/daops.html
// the precisions is not used in my implementation.
type bignum_decimalValue struct {
	// the bignum is ALWAYS a 10 based number, with separated digits representation

	/* example num representation: '12.34'
		digits: [1,2,3,4]
		exponent: -2    because 1234 * (10^-2) = 12.34

	   example 2: 1234000
		digits: [1,2,3,4]
		exponent: 3     because 1234 * (10^3) = 1234000

		with this representation, integers and floats can be represented too, with the same structure
	*/

	// http://www.math.u-szeged.hu/tagok/kurusa/_site/index.php/hu/blogs/maths-blog/98-hun-eng-mat-glos
	// https://www.collinsdictionary.com/word-lists/mathematics-mathematical-terms
	// https://hu.speaklanguages.com/angol/sz%C3%B3szedet/alakzatok-%C3%A9s-matematikai-kifejez%C3%A9sek

	/* index0 is always the first digit, indexLast is the less powerful digit  */
	digits digitList  // decimal digits
	exponent int  // where is the ., or how many 0 are after the digits
	negative bool // false by default: is the number negative?
}

func (bn bignum_decimalValue) isPositive() bool {
	return ! (bn.negative)
}

func (bn bignum_decimalValue) isNegative() bool {
	return bn.negative
}

func (bn bignum_decimalValue) print(msg string) {
	sign := "+"
	if bn.negative {
		sign = "-"
	}
	fmt.Println(msg, sign, bn.digits, bn.exponent)
}

// first digit position is 0. what is the last digit id?
func (bn bignum_decimalValue) digitsIndexLast() int {
	//                 id:  0 1 2
	// if we have 3 digits: 1,2,3
	// then the last id is 2 (the index of the last elem)
	return len(bn.digits)-1
}

// if the position is not in the digit, the value is 0.
// posFromBack: 0 is the last digit, -1 is the second from back
func (bn bignum_decimalValue) digitValueInPosition(posFromBack int) (int, digitElemType) {
	positionLast := bn.digitsIndexLast()
	posAbsolute := positionLast - posFromBack
	var value digitElemType = 0
	if posAbsolute >= 0 && posAbsolute <= positionLast { // so if the digit is in the number...
		value = bn.digits[posAbsolute]
	}
	return posAbsolute, value // give back the real absoule pos of the digit, and the value
							  // (if it is less than 0, then non-real digit was read)
}

func (bn bignum_decimalValue) duplicate() bignum_decimalValue {
	// goal: duplicate the values of the original struct, without any accidental pointer usage
	// with other words: if something is changed in the original, it cannnot be reflected in the duplication

	digitsNew := digitList{}
	for _, digit := range bn.digits {
		digitsNew = append(digitsNew, digit)  // re-build digit list,
	}

	return bignum_decimalValue{digits: digitsNew, exponent: bn.exponent, negative: bn.negative}
}

// collect all useful digits, plus the num of extra zeros at the end.
// with this, it is easy to compare two numbers
func (bn bignum_decimalValue) normalisedForm() bignum_decimalValue {
	allDigits := digitList{}
	zeroCounter := bn.exponent

	lastNonZeroDetected := false
	// check the numbers from the last to the first direction
	for i := len(allDigits) -1; i>=0; i-- {
		digit := allDigits[i]
		if ! lastNonZeroDetected {
			if digit == 0 {
				zeroCounter++
			} else {
				lastNonZeroDetected = true
			}
		}
		// example num: 120300
		// the last 2 chars are 0.
		// when 3 is detected, the last non-zero, from there, collect all valueable digit
		if lastNonZeroDetected {
			allDigits = append(allDigits, digit)
		}
	}

	// I don't trust in pointers. So don't override original values, create a new, normalised form
	return bignum_decimalValue{digits: allDigits, exponent: zeroCounter, negative: bn.negative}
}


func (a bignum_decimalValue) isEqual(b bignum_decimalValue) bool {
	a_normalised := a.normalisedForm()
	b_normalised := b.normalisedForm()

	// deeply equal: https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices

	return 	(a_normalised.exponent == b_normalised.exponent) &&
			(a_normalised.negative == b_normalised.negative) &&
			reflect.DeepEqual(a.digits, b.digits)
}


// TESTED: Test__isLessThan
func (bn bignum_decimalValue) isLessThan(other bignum_decimalValue) bool {
	if bn.isNegative() && other.isPositive() {
		return true
	}
	if bn.isPositive() && other.isNegative(){
		return false
	}
	if bn.isEqual(other) {
		return false
	}

	//
	bignumNew, other_New := bigNum_pair_setSameExponent_decreaseBiggerExponent(bn, other)

	// select the highest value place in the numbers
	positionFromBack__highestPlace := max(bignumNew.digitsIndexLast(), other_New.digitsIndexLast())

	bnAbsoluteValueIsLess := false
	for positionFromBack__highestPlace >=0 {
		_, valueDigitBignum := bignumNew.digitValueInPosition(positionFromBack__highestPlace)
		_, valueDigitOther_ := other_New.digitValueInPosition(positionFromBack__highestPlace)

		if valueDigitBignum < valueDigitOther_{
			bnAbsoluteValueIsLess = true
			break
		}
		if valueDigitBignum > valueDigitOther_{
			bnAbsoluteValueIsLess = false
			break
		}
		positionFromBack__highestPlace--
	}

	if bn.isPositive() && other.isPositive() {
		return bnAbsoluteValueIsLess
	}
	// bn.isNegative(),	other.isNegative() case:
	return ! bnAbsoluteValueIsLess
}



// Tested, from operator_test func
// give back the original int, and the bignum too -
// in automata tests it can be important to know, what was the original integer
// simple implementation: exponent is not modified
func bigNum_from_int(i int) (int, bignum_decimalValue) {
	iOrig := i

	negative := false
	if i < 0 {
		negative = true
		i = -i  // use only the absolute value
	}
	digits := digitList{}

	// copy all digits one by one,
	for _, charElem := range strconv.Itoa(i) {
		digitVal, err := digitRune_decimalValue(charElem)
		if err != nil {
			fmt.Println("Error, character "+string(charElem)+" doesn't have decimal value")
		} else {
			digits = append(digits, digitVal)
		}
	}
	return iOrig, bignum_decimalValue{digits: digits, exponent: 0, negative: negative}
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
func bigNum_from_digitValue__min0_max35(decimalVal digitElemType) bignum_decimalValue {
	// a digit's value is minimum 0, maximum 35. there is no problem with too big integer values
	digit_2 := decimalVal % 10
	digit_1 := (decimalVal - digit_2) / 10
	// simple conversion, NEVER normalize here - normalization can happen at one point only
	if digit_1 > 0 {
		return bignum_decimalValue{digits: digitList{digit_1, digit_2}, exponent: 0}
	} else {
		return bignum_decimalValue{digits: digitList{digit_2}, exponent: 0, negative: false}
	}
}

func bigNum_operator_add(a, b bignum_decimalValue) bignum_decimalValue {
	if a.isPositive() && b.isNegative() {
		b.negative = false
		return bigNum_operator_sub(a, b)
	}
	if a.isNegative() && b.isPositive(){
		a.negative = false
		return bigNum_operator_sub(b, a)
	}
	if a.isNegative() && b.isNegative() {
		a.negative = false
		b.negative = false
		result := internal_used_only__bigNum_add_positive_positive(a, b)
		result.negative = true
		return result
	}

	// basic case: a, b are positive
	return internal_used_only__bigNum_add_positive_positive(a, b)
}


func bigNum_operator_sub(a, b bignum_decimalValue) bignum_decimalValue {
	if a.isPositive() && b.isNegative() {
		// FIXME: WRONG
		b.negative = false
		return bigNum_operator_add(a, b)
	}
	if a.isNegative() && b.isPositive() {
		a.negative = false
		result := bigNum_operator_add(a, b)
		result.negative = true
		return result
	}
	if a.isNegative() && b.isNegative(){
		a.negative = false
		b.negative = false
		result := bigNum_operator_add(a, b)
		result.negative = true
		return result
	}

	fmt.Println("basic sub case, pos/pos", a, b)
	return internal_used_only__bigNum_sub_positive_positive(a, b)
}


// https://www.youtube.com/watch?v=-pPXFvVxlng
func internal_used_only__bigNum_sub_positive_positive(a, b bignum_decimalValue) bignum_decimalValue {

	a, b = bigNum_pair_setSameExponent_decreaseBiggerExponent(a, b)
	digitsReversed := digitList{}

	var overflow = digitElemType(0)
	position := -1

	for {
		position++

		var valueA digitElemType = 0
		var valueB digitElemType = 0

		if position < len(a.digits) {
			valueA = a.digits[position]
		}
		if position < len(b.digits) {
			valueB = b.digits[position]
		}
		valueA = valueA - overflow
		overflow = 0 // because overflow's value was calculated into valueA

		if valueA == 0 && valueB == 0 && overflow == 0 {
			break // exit if there is no more thing to do
		}
		fmt.Println("sub 1 pos/pos valueA:", valueA, "  valueB", valueB, "overflow:", overflow)
		if valueA < valueB {
			valueA += 10
			overflow +=1
		}
		fmt.Println("sub 2 pos/pos valueA:", valueA, "  valueB", valueB, "overflow:", overflow)
		valueDiff := valueA - valueB
		digitsReversed = append(digitsReversed, valueDiff)
		fmt.Println("sub 3 pos/pos valueDiff:", valueDiff)

		// safety exit
		if position > 5{
			break
		}
	}

	summa := bignum_decimalValue{digits: digits_reverse(digitsReversed), exponent: a.exponent}
	return summa
}


//////////// TESTED /////////////////
/* receives 2 numbers. return with 2 numbers, where the exponents are similar*/
func bigNum_pair_setSameExponent_decreaseBiggerExponent(a, b bignum_decimalValue) (bignum_decimalValue, bignum_decimalValue) {
	if a.exponent == b.exponent {
		return a, b
	}
	orderFlipped := false
	numExponentSmaller := a.duplicate()
	numExponentBigger := b.duplicate()

	// addresses tested from  Test_bigNum_pair_set_same_exponent
	// fmt.Printf("bignum decrease, address a %p\n", &a) // address a points to a different mem area
	// fmt.Printf("bignum decrease, address b %p\n", &b)
	// fmt.Printf("bignum decrease, address a.digits[0] %p\n", &(a.digits[0])) // but this is same address with the caller
	// fmt.Printf("bignum decrease, address b.digits[0] %p\n", &(b.digits[0]))


	// fmt.Println("setExponent, common, smaller0:", numExponentSmaller)
	// fmt.Println("setExponent, common, bigger 0:",  numExponentBigger)

	if numExponentBigger.exponent < numExponentSmaller.exponent {
		// fmt.Println("exponent b < a, replace them")
		orderFlipped = true
		numExponentSmaller, numExponentBigger = numExponentBigger, numExponentSmaller
	}
	// fmt.Println("setExponent, common, smaller1:", numExponentSmaller)
	// fmt.Println("setExponent, common, bigger 1:", numExponentBigger)

	for numExponentBigger.exponent > numExponentSmaller.exponent {
		// fmt.Println("for loop, bigger digits:", numExponentBigger.digits, "   exponent:", numExponentBigger.exponent)
		numExponentBigger.exponent--
		numExponentBigger.digits = append(numExponentBigger.digits, 0)
	}
	// fmt.Println("setExponent, common, smaller, end:", numExponentSmaller)
	// fmt.Println("setExponent, common, bigger , end:", numExponentBigger)

	if orderFlipped {
		return numExponentBigger, numExponentSmaller
	}
	return numExponentSmaller, numExponentBigger
}


func internal_used_only__bigNum_add_positive_positive(a, b bignum_decimalValue) bignum_decimalValue {
	// the bignum is ALWAYS decimal number, with separated digits representation

	// add operation can be done ONLY if the exponents are same
	a.print("internal a1")
	b.print("internal b1")
	a, b = bigNum_pair_setSameExponent_decreaseBiggerExponent(a, b)
	a.print("internal a2")
	b.print("internal b2")

	digitsReversed := digitList{}

	var overflow digitElemType = 0

	positionFromBack := -1
	for {
		positionFromBack++ // posA, posB are going from the last (biggest) index to 0 index with -Delta
		posA, valueDigitA := a.digitValueInPosition(positionFromBack)
		posB, valueDigitB := b.digitValueInPosition(positionFromBack)
		fmt.Println("internal add before, pos pos >>> a:", valueDigitA, "  b:", valueDigitB, "overflow:", overflow)

		// the reading started from the highest indexes to index 0, which is the first digit.
		// if both position is before the first digit, process can be stopped
		if posA < 0 && posB < 0 && overflow == 0 {
			break
		}

		valueSum := valueDigitA + valueDigitB + overflow
		digitNew := valueSum % 10
		digitsReversed = append(digitsReversed, digitNew)
		fmt.Println("internal add  after, pos pos >>> valueSum:", valueSum, "   digitNew:", digitNew)

		overflow = (valueSum - digitNew) / 10
	}

	summa := bignum_decimalValue{digits: digits_reverse(digitsReversed), exponent: a.exponent}
	summa.print("summa, after add pos pos")
	return summa
}


// return with a fix value
func bigNum_zero() bignum_decimalValue {
	return bignum_decimalValue{digits: digitList{0}, exponent: 0, negative: false}
}



// TESTED!!!
func bigNum_convert_to_INT_for_testcases(bigNum bignum_decimalValue) int {
	summa := 0
	lenDigits := len(bigNum.digits)
	multiplicator := lenDigits

	for pos := 0; pos < lenDigits; pos++ {
		fmt.Println()

		// positions: 012
		//            123: first multiplicator is 2, second is 1, then 0
		multiplicator -= 1
		// fmt.Println("multiplicator:", multiplicator)

		digitValue := int(bigNum.digits[pos])

		for m:= multiplicator; m>0; m-- {
			digitValue = digitValue * 10
		}

		// fmt.Println("convert to INT: digit[",pos,"] =>", bigNum.digits[pos], "digitValue:", digitValue)
		summa += digitValue
	}
	if bigNum.negative {
		summa = -summa
	}
	return summa
}


// tested
func digits_reverse(digits digitList) digitList {
	digitsReversed := digitList{}
	for pos := len(digits)-1; pos >=0; pos-- {
		digitsReversed = append(digitsReversed, digits[pos])
	}
	return digitsReversed
}

