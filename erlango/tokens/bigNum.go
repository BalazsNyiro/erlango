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
	return ! bn.negative
}

func (bn bignum_decimalValue) isNegative() bool {
	return bn.negative
}

func (bn bignum_decimalValue) print(msg string) {
	sign := "+"
	if bn.isNegative() {
		sign = "-"
	}
	fmt.Println(msg, sign, bn.digits, bn.exponent)
}

// first digit position is 0. what is the last digit id?
func (bn bignum_decimalValue) digitsIndexLast() int {
	//                 id:  0 1 2
	// if we have 3 digits: 1,2,3
	// then the last id is 2 (the index of the last elem)
	// if digits is empty, retval is -1
	return len(bn.digits)-1
}

func (bn bignum_decimalValue) digitsExport() digitList{
	// export a bigNum's digits into a new data structure
	// I want a separated list of digits, not connected to the original num
	digits := digitList{}
	for _, digit := range bn.digits{
		digits = append(digits, digit)
	}
	return digits
}


// if the position is not in the digit, the value is 0.
// posFromBack: 0 is the last digit, 1 is the second from back, 2 is the third from back
// so here there is no negative position usage
// this is indexing, from the END, not from the start, as we use normal indexing
// the perspective is changed only: indexing from right->left direction
func (bn bignum_decimalValue) digitValueInPositionFromBack(posFromBack int) (int, digitElemType) {
	positionLast := bn.digitsIndexLast()
	posAbsolute := positionLast - posFromBack
	var value digitElemType = 0
	if posAbsolute >= 0 && posAbsolute <= positionLast { // so if the digit is in the number...
		value = bn.digits[posAbsolute]
	}
	return posAbsolute, value // give back the real absolute pos of the digit, and the value
							  // (if it is less than 0, then non-real digit was read)
}

func (bn bignum_decimalValue) duplicate() bignum_decimalValue {
	// goal: duplicate the values of the original struct, without any accidental pointer usage
	// with other words: if something is changed in the original, it cannot be reflected in the duplication

	digitsNew := digitList{}
	for _, digit := range bn.digits {
		digitsNew = append(digitsNew, digit)  // re-build digit list,
	}

	return bignum_decimalValue{digits: digitsNew, exponent: bn.exponent, negative: bn.negative}
}


//set the exponent to zero, and load the ending zeros back into the digits
// the returned value is a Copy, not the original number
func (bn bignum_decimalValue) normalisedForm_exponentZerosIntoDigits() bignum_decimalValue {
	// if the exponent > 0, load back the 0 digits into the num

	digits := bn.digitsExport()
	exponent := bn.exponent

	if exponent <= 0 {
		return bignum_decimalValue{digits: digits, exponent: exponent, negative: bn.negative}
	}

	for exponent > 0 {
		digits = append(digits, 0)
		exponent--
	}

	return bignum_decimalValue{digits: digits, exponent: 0, negative: bn.negative}
}


// collect all useful digits, plus the num of extra zeros at the end.
// with this, it is easy to compare two numbers
// normalisedForm means: the number digits and exponents are re-organised, by a method.
// TESTED
func (bn bignum_decimalValue) normalisedForm_endingZerosIntoExponent() bignum_decimalValue {
	allDigitsReversed := digitList{}
	zeroCounter := bn.exponent
	// fmt.Println("	zeroCounter:", zeroCounter)

	lastNonZeroDetected := false
	// check the numbers from the last to the first direction
	for i := bn.digitsIndexLast(); i>=0; i-- {
		digit := bn.digits[i]
		if ! lastNonZeroDetected {
			if digit == 0 {
				zeroCounter++
			} else {
				lastNonZeroDetected = true
			}
		}
		// fmt.Println("normalise: pos", i, "  digit:", digit, "   lastNonZeroDetected:", lastNonZeroDetected, "  zeroCounter:", zeroCounter, "   allDigits:", allDigits)
		// example num: 120300
		// the last 2 chars are 0.
		// when 3 is detected, the last non-zero, from there, collect all valueable digit
		if lastNonZeroDetected {
			allDigitsReversed = append(allDigitsReversed, digit)
		}
	}

	allDigits := digitsReverse(allDigitsReversed) // reverse back the order to normal
	allDigitsNoLeadingZeros := digitsCleaning_leadingZerosRemoval(allDigits)

	// I don't trust in pointers. So don't override original values, create a new, normalised form
	return bignum_decimalValue{digits: allDigitsNoLeadingZeros, exponent: zeroCounter, negative: bn.negative}
}


// TESTED
func (bn bignum_decimalValue) isEqual(other bignum_decimalValue) bool {
	bigNum_normalised := bn.normalisedForm_endingZerosIntoExponent()
	other__normalised := other.normalisedForm_endingZerosIntoExponent()

	// deeply equal: https://stackoverflow.com/questions/15311969/checking-the-equality-of-two-slices

	return 	(bigNum_normalised.exponent == other__normalised.exponent) &&
			(bigNum_normalised.negative == other__normalised.negative) &&
			reflect.DeepEqual(bigNum_normalised.digits, other__normalised.digits)
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
	bignumNew, other_New := bigNum_exponentsSetSame_decreaseBiggerExponent(bn, other)

	// select the highest value place in the numbers: the first digit, with relative positions from BACK!
	// with other words: from back, which is the longest index num to go forward?
	//                   the last index num is totally the num that you need to use
	//                   to reach the first elem from the back, so this is correct:
	positionFromBack__firstDigitPositionFromBackViewPoint := max(bignumNew.digitsIndexLast(), other_New.digitsIndexLast())
	/*
                     01234   <- digit indexes
	   example, bn = 45678   last digit's index is 4, so: bn.digits[4] == 8
	     other =       123

	here we select the greater index, because from the back, that is a pointer to the first char:


	reversed indexes:43210 <----   from the back position, the first digit index is 4, too.
					 45678
	                 00123
	*/

	bnAbsoluteValueIsLess := false
	for positionFromBack__firstDigitPositionFromBackViewPoint >=0 {

		// Select the first character, with relative index from the last char:
		_, valueDigitBignum := bignumNew.digitValueInPositionFromBack(positionFromBack__firstDigitPositionFromBackViewPoint)
		_, valueDigitOther_ := other_New.digitValueInPositionFromBack(positionFromBack__firstDigitPositionFromBackViewPoint)

		// if the highest place's digit is different, it is easy to decide, which num is less
		if valueDigitBignum < valueDigitOther_{
			bnAbsoluteValueIsLess = true
			break
		}
		if valueDigitBignum > valueDigitOther_{
			bnAbsoluteValueIsLess = false
			break
		}
		positionFromBack__firstDigitPositionFromBackViewPoint--
	}

	if bn.isPositive() && other.isPositive() {
		return bnAbsoluteValueIsLess
	}
	// bn.isNegative(),	other.isNegative() case:
	return ! bnAbsoluteValueIsLess
}


// TODO: check this later: is is used in hexa nums, or not? if not, delete this
// this is important when digit values are checked, in non-decimal char processing
func bigNum_create_from_digitValue__min0_max35(decimalVal digitElemType) bignum_decimalValue {
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


/// used in multiply
func bigNum_create_from_digits__positiveZeroExponent(digits digitList) bignum_decimalValue {
	return bigNum_create_from_digits(digits, false, 0)
}

func bigNum_create_from_digits(digits digitList, negative bool, exponent int) bignum_decimalValue {
	return bignum_decimalValue{digits: digits, negative: negative, exponent: exponent}
}


// Tested, from operator_test func
// simple implementation: exponent is not modified
// integer conversion to bigNum
func bigNum_create_from_int(i int) bignum_decimalValue {
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
	return bignum_decimalValue{digits: digits, exponent: 0, negative: negative}
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

	// add extra 0 numbers, if exponent > 0:
	for bigNum.exponent > 0 {
		summa = summa*10
		bigNum.exponent -= 1
	}

	// float->int conversion, be careful, maybe important digits are lost from the end, this is a to INT conversion!!!
	for bigNum.exponent < 0 {
		summa = summa/10
		bigNum.exponent += 1
	}

	return summa
}




func bigNum_operator_mul(bigNum, mul bignum_decimalValue) bignum_decimalValue {
	fmt.Println("bigNum operator mul, bigNum:", bigNum, "  mul:", mul)
	if bigNum.isPositive() && mul.isNegative() {
		mul.negative = false
		result := internal_used_only__bigNum_mul_positive_positive(bigNum, mul)
		result.negative = true
		return result
	}
	if bigNum.isNegative() && mul.isPositive(){
		bigNum.negative = false
		result := internal_used_only__bigNum_mul_positive_positive(bigNum, mul)
		result.negative = true
		return result
	}
	if bigNum.isNegative() && mul.isNegative() {
		bigNum.negative = false
		mul.negative = false
		result := internal_used_only__bigNum_mul_positive_positive(bigNum, mul)
		return result
	}

	// basic case: bigNum, mul are positive
	return internal_used_only__bigNum_mul_positive_positive(bigNum, mul)
}


//////////// TESTED /////////////////
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


//////////// TESTED /////////////////
func bigNum_operator_sub(a, b bignum_decimalValue) bignum_decimalValue {
	if a.isPositive() && b.isNegative() {
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
		// negative - negative = negative + positive = positive + negative, handled in ADD
		b.negative = false
		return bigNum_operator_add(a, b)
	}

	// fmt.Println("basic sub case, positive/positive", a, b)
	return internal_used_only__bigNum_sub_positive_positive(a, b)
}

//////////// TESTED /////////////////
func bigNum_operator_div(a, b bignum_decimalValue) (bignum_decimalValue, bignum_decimalValue, error) {
	if a.isPositive() && b.isNegative() {
		/*
		9> 9 / -7.
		-1.2857142857142858
		10> 9 div -7.
		-1
		11> 9 rem -7.
		2

		*/

		b.negative = false
		quotient, remainder, err := internal_used_only_bigNum_div_positivePositive_FULL_ALGORITHM(a, b)
		quotient.negative = true
		remainder.negative = false
		return quotient, remainder, err

	}
	if a.isNegative() && b.isPositive() {
		/*
		6> -9 / 7.
		-1.2857142857142858
		7> -9 div 7.
		-1
		8> -9 rem 7.
		-2

		*/
		a.negative = false
		quotient, remainder, err := internal_used_only_bigNum_div_positivePositive_FULL_ALGORITHM(a, b)
		quotient.negative = true
		remainder.negative = true
		return quotient, remainder, err
	}

	if a.isNegative() && b.isNegative(){
		/*
		3> -9 / -7.
		1.2857142857142858
		4> -9 div -7.
		1
		5> -9 rem -7.
		-2

		*/
		a.negative = false
		b.negative = false
		quotient, remainder, err := internal_used_only_bigNum_div_positivePositive_FULL_ALGORITHM(a, b)
		remainder.negative = true
		return quotient, remainder, err
	}

	return internal_used_only_bigNum_div_positivePositive_FULL_ALGORITHM(a, b)
}




//////////// TESTED /////////////////
/* receives 2 numbers. return with 2 numbers, where the exponents are similar*/
func bigNum_exponentsSetSame_decreaseBiggerExponent(a, b bignum_decimalValue) (bignum_decimalValue, bignum_decimalValue) {
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



// return with a fix value
func bigNum_zero() bignum_decimalValue {
	return bignum_decimalValue{digits: digitList{0}, exponent: 0, negative: false}
}

func bigNum_one() bignum_decimalValue {
	return bignum_decimalValue{digits: digitList{1}, exponent: 0, negative: false}
}




// tested
func digitsReverse(digits digitList) digitList {
	digitsReversed := digitList{}
	for pos := len(digits)-1; pos >=0; pos-- {
		digitsReversed = append(digitsReversed, digits[pos])
	}
	return digitsReversed
}

// TESTED
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

// example: generate five zeros, this: [0, 0, 0, 0, 0]
//          repeat: 5, digitNum: 0
func digits_series_simple_generate(repeat int, digitNum digitElemType) digitList {
	// repeat >= 0
	result := digitList{}
	for i:=0; i<repeat; i++ {
		result = append(result, digitNum)
	}
	return result
}



// Tested from  Test_normaliseExponent_endingZerosRemove
func digitsCleaning_leadingZerosRemoval(digits digitList) digitList {
	// don't store prefix possible emtpy zeros, that doesn't have value
	// for example: digitList{0,1,2,0,3,0,0}
	allDigitsNoLeadingZeros := digitList{}
	copyEverything := false

	for _, digit := range digits {
		if digit != digitElemType(0) { // from the first non-zero char, copy everything
			copyEverything = true
		}
		if copyEverything {
			allDigitsNoLeadingZeros = append(allDigitsNoLeadingZeros, digit)
		}
	}

	if len(allDigitsNoLeadingZeros) == 0 {        // so if there are totally no digits here:
		allDigitsNoLeadingZeros = digitList{0}   // minimum one zero has to be there, without digits there is no number
	}
	return allDigitsNoLeadingZeros
}


/////////////////// MULTIPLY //////////////
// TESTED
func internal_used_only__bigNum_mul_positive_positive(bigNum, multiply bignum_decimalValue) bignum_decimalValue {
	bigNum = bigNum.normalisedForm_endingZerosIntoExponent()
	multiply = multiply.normalisedForm_endingZerosIntoExponent()
	// fmt.Println("  bigNum:", bigNum)
	// fmt.Println("multiply:", multiply)
	result := bigNum_zero()
	/*
	Algorithm demo:
	345 * 12
	---- calc: ---
	345     Level1 values
	 690    Level2 values
	------sum:--
	4140


	Algorithm:
	 - calculate one digit:
		345*1
		345 -> 3450, because 1 is in position 10
	 - calculate next digit:
	    345 * 2
	    690 -> 690*1, because 2 is in position 1

	*/

	extraZeroCounter_becauseOfDigitPosition := -1

	// Dear reader: here is a gift for you:
	// https://open.spotify.com/track/34BArHFcgaf7nXm487HJqf?si=f60d0393137c4fb7
	// a really good music, if you want to manage math operations...
	// Mari Samuelsen: Sequence(Four) for Solo Violin
	for idxMul := multiply.digitsIndexLast(); idxMul>=0; idxMul-- {

		digitMul := multiply.digits[idxMul]
		/* if you do 12*345
		- first  mul digit is 5, decimalValueRange:   1
		- first  mul digit is 4, decimalValueRange:  10
		- second mul digit is 3, decimalValueRange: 100
		*/
		extraZeroCounter_becauseOfDigitPosition++ // 0,1,2,3



		////////////////////////////////////////////
		///// CONCENTRATE ON ONE OPERATION HERE  ///
		digitsReversed := digitList{}
		var overflow digitElemType = 0
		bnPositionFromBack := -1

		// use the actual digit of MultiplyNum, do the calculation:
		for {
			bnPositionFromBack++ // positions are checked from the last (higher) index to 0 index with -Delta
			bnPosAbs, bnValueDigit := bigNum.digitValueInPositionFromBack(bnPositionFromBack)
			// fmt.Print("\n\n\n")
			// fmt.Println("         multiplyDigit:", digitMul)
			// fmt.Println("         bnPosFromBack:", bnPositionFromBack)
			// fmt.Println("         bnPosAbs     :", bnPosAbs)
			// fmt.Println("         bnValueDigit :", bnValueDigit)

			if bnPosAbs < 0 && overflow == 0 {
				break
			}

			valueAfterMultiply := bnValueDigit * digitMul + overflow
			digitNew := valueAfterMultiply % 10
			digitsReversed = append(digitsReversed, digitNew)

			overflow = (valueAfterMultiply - digitNew) / 10

			// fmt.Println("bnValue After multiply:", valueAfterMultiply)
			// fmt.Println("              digitNew:", digitNew)
			// fmt.Println("              overflow:", overflow)
		} ///// CONCENTRATE ON ONE OPERATION HERE //
		////////////////////////////////////////////


		digitsActualResult := digitsReverse(digitsReversed)
		// fmt.Println("END1 digitsActualResult:", digitsActualResult)
		// fmt.Println("END1 extra zero counter:", extraZeroCounter_becauseOfDigitPosition)

		extraZerosFromPlace := digits_series_simple_generate(extraZeroCounter_becauseOfDigitPosition, 0)
		digitsActualResult = append(digitsActualResult, extraZerosFromPlace...)
		// fmt.Println("END2 digitsActualResult:", digitsActualResult)

		// as the bigNum*digitMul are calculated, the actual step's result has to be accumulated
		resultOfActualStep := bigNum_create_from_digits__positiveZeroExponent(digitsActualResult)
		result = bigNum_operator_add(result, resultOfActualStep)
		// fmt.Println("result updated:", result)
	} // for, idxMul

	// the 10^exponent values has to be added to the number:
	result.exponent += bigNum.exponent + multiply.exponent
	return result
} // end of multiply


///////////////// DIVISION ///////////////
// floating challenge: https://stackoverflow.com/questions/7662850/precision-in-erlang
// Alpha = math:acos((4*4 + 5*5 - 3*3) / (2*4*5))
// Area = 1/2 * 4 * 5 * math:sin(Alpha)


///////////// ADD ///////////// TESTED
func internal_used_only__bigNum_add_positive_positive(a, b bignum_decimalValue) bignum_decimalValue {
	// the bignum is ALWAYS decimal number, with separated digits representation

	// add operation can be done ONLY if the exponents are same
	a, b = bigNum_exponentsSetSame_decreaseBiggerExponent(a, b)
	if a.negative || b.negative {
		// this case is not possible, if this fun is not called directly
		fmt.Println("Error, only positive numbers are accepted:", a, b)
	}

	digitsReversed := digitList{}

	var overflow digitElemType = 0
	positionFromBack := -1

	for {
		positionFromBack++ // posA, posB are going from the last (biggest) index to 0 index with -Delta
		posA, valueDigitA := a.digitValueInPositionFromBack(positionFromBack)
		posB, valueDigitB := b.digitValueInPositionFromBack(positionFromBack)
		// fmt.Println("internal add before, pos pos >>> a:", valueDigitA, "  b:", valueDigitB, "overflow:", overflow)

		// the reading started from the highest indexes to index 0, which is the first digit.
		// if both position is before the first digit, process can be stopped
		if posA < 0 && posB < 0 && overflow == 0 {
			break
		}

		valueSum := valueDigitA + valueDigitB + overflow
		digitNew := valueSum % 10
		digitsReversed = append(digitsReversed, digitNew)
		// fmt.Println("internal add  after, pos pos >>> valueSum:", valueSum, "   digitNew:", digitNew)

		overflow = (valueSum - digitNew) / 10
	}

	summa := bignum_decimalValue{digits: digitsReverse(digitsReversed), exponent: a.exponent}
	// summa.print("summa, after add pos pos")
	return summa
}


// https://www.youtube.com/watch?v=-pPXFvVxlng
/////////// SUB////////// TESTED /////////////////
func internal_used_only__bigNum_sub_positive_positive(a, b bignum_decimalValue) bignum_decimalValue {

	a, b = bigNum_exponentsSetSame_decreaseBiggerExponent(a, b)
	// fmt.Println("A, bignum sub positive positive: ", a, b)
	if a.negative || b.negative {
		// this case is not possible, if this fun is not called directly
		fmt.Println("Error, only positive numbers are accepted:", a, b)
	}

	// the algorithm works if a >= b. so swap them, if not.
	// the difference is same between the two numbers, only the sign is different
	negativeResult := false
	if ! b.isLessThan(a) {
		a, b = b, a
		negativeResult = true
	}
	// fmt.Println("b, bignum sub positive positive: ", a, b)

	digitsReversed := digitList{}

	var overflow = digitElemType(0)
	positionFromBack := -1

	for {
		positionFromBack++
		posA, valueA := a.digitValueInPositionFromBack(positionFromBack)
		posB, valueB := b.digitValueInPositionFromBack(positionFromBack)

		valueA = valueA - overflow
		overflow = 0 // because overflow's value was calculated into valueA
		// fmt.Println("sub 1  pos/pos valueA:", valueA, "  valueB", valueB, "overflow:", overflow)

		if posA < 0 && posB < 0 {
			break // exit if there is no more thing to do
		}

		// fmt.Println("sub 2a pos/pos valueA:", valueA, "  valueB", valueB, "overflow:", overflow)
		if valueA < valueB {
			// fmt.Println("sub 2b valueA", valueA, " < ", valueB, "valueB")
			valueA += 10
			overflow +=1
		}
		// fmt.Println("sub 3  pos/pos valueA:", valueA, "  valueB", valueB, "overflow:", overflow)
		valueDiff := valueA - valueB
		digitsReversed = append(digitsReversed, valueDiff)
		// fmt.Println("sub 4  pos/pos valueDiff:", valueDiff, "APPENDED")
	}

	/* during a subtraction, it is normal if leading zeros are inserted:
		71
	   -65
	    06 is the result, with the algorithm, because 7-7 at the first pos is 0
	but later it has to be removed

	if the num would be 81:
	    81
	   -65
	    16  // so the 0 in the previous example is the natural output.
	*/

	digitsNormal := digitsCleaning_leadingZerosRemoval(digitsReverse(digitsReversed))
	summa := bignum_decimalValue{digits: digitsNormal, exponent: a.exponent, negative: negativeResult}
	// fmt.Println("sub 5 positive positive summa: ", summa)
	return summa
}


/////////////// DIVISION ///////////////


/* the normal division operation needs a naive func that can tell 10:2, 86:23 -
so a division operator, that can work with relatively small numbers, or
when the ratio between the numbers is not too big, for example: 1536:93
with other words:
	- when with a few, naive addition, the result can be calculated.
	- when the result can be calculated in max 10-20 steps, that would be the ideal

in this case: 1234567891122333: 36, in this case this algorithm is not effective, because it is slow.
So use it only if you can calculate the result in a few step
*/
func internal_used_only_bigNum_div_positivePositive__for_relative_small_numbers(bigNum, divisor bignum_decimalValue) (quotient, remainder bignum_decimalValue, err error)  {
	// the function is planned to work with positive bigNum and divisor ONLY.
	// and with relative small ratio between bigNum and divisor, because this is a naive (but stable :-) solution

	quotient = bigNum_zero()
	remainder = bigNum
	err = nil

	if divisor.isEqual(bigNum_zero()) {
		err = errors.New("zero division")
		return bigNum_zero(), remainder, err
	}

	for {
		// fmt.Println("quotient:", quotient, " remainder:", remainder,  "divisor:", divisor)
		if remainder.isLessThan(divisor) {
			break
		}
		quotient = bigNum_operator_add(quotient, bigNum_one())
		remainder = bigNum_operator_sub(remainder, divisor)
	}

	return quotient, remainder, err
}


/*
 https://www.youtube.com/watch?v=yLknFrMrdAM&pp=ygUGb3N6dGFz
 algorithm: https://www.youtube.com/watch?v=p8KSnecgfHs

this is not a naive implementation - but in a few step, with small numbers, I need to use the
small_nubmer focused func

// https://golangbyexample.com/remainder-modulus-go-golang/
*/
func internal_used_only_bigNum_div_positivePositive_FULL_ALGORITHM(bigNum, divisor bignum_decimalValue) (quotient, remainder bignum_decimalValue, err error) {

	bigNum = bigNum.normalisedForm_exponentZerosIntoDigits()
	divisor = divisor.normalisedForm_exponentZerosIntoDigits()
	fmt.Println("bigNum:", bigNum, "divisor:", divisor)
	quotientDigits:= digitList{}
	remainder = bigNum
	err = nil

	if divisor.isEqual(bigNum_zero()) {
		err = errors.New("zero division")
		return bigNum_zero(), remainder, err
	} ///////////////////////////////////////////

	digitsToAnalyseOneByOne := bigNum.digitsExport()

	/* algorithm example: take the first X chars, and if it is bigger than divisor, do the math
	            vv
	            9815:65 =
          step1:33   ---> 1
	      step2:331  --->  5
	      step3:  65 --->   1

		result: 9815:65 = 151, remainder=0
	*/

	digitsTmp := digitList{}
	// index:
	for len(digitsToAnalyseOneByOne) > 0 {
		// relocate one digit from the orig digit list into the temporary value (first 98) where we do the first division
		digitsTmp = append(digitsTmp, digitsToAnalyseOneByOne[0])
		digitsToAnalyseOneByOne = digitsToAnalyseOneByOne[1:]
		// fmt.Println("A, tmp:", digitsTmp, " <- ", digitsToAnalyseOneByOne, "  quotientDigits", quotientDigits)

		numTmp := bignum_decimalValue{digits: digitsTmp, exponent: 0, negative: false}

		if numTmp.isLessThan(divisor) {
			continue
		} else {
			// fmt.Println("Temporary num ", numTmp.digits, ">= than divisor", divisor.digits)
			// 98:65 step -> calculate a SubQuotient and subRemainder
			subQuotient, subRemainder, subErr := internal_used_only_bigNum_div_positivePositive__for_relative_small_numbers(numTmp, divisor)
			// fmt.Println("subQuotient:", subQuotient, "  subRemainder:", subRemainder)

			if subErr != nil {
				return bigNum_zero(), subRemainder, subErr
			}

			quotientDigits = append(quotientDigits, subQuotient.digitsExport()...)
			digitsTmp = subRemainder.digitsExport()
		}
	}

	// the value in digitsTmp is less than divisor, so after the for loop, this is the remainder
	remainder = bignum_decimalValue{digits: digitsCleaning_leadingZerosRemoval(digitsTmp), exponent: 0, negative: false}

	/////////////////////////////////////////////
	quotient = bignum_decimalValue{digits: digitsCleaning_leadingZerosRemoval(quotientDigits), exponent: 0, negative: false}

	quotient = quotient.normalisedForm_exponentZerosIntoDigits()
	remainder = remainder.normalisedForm_endingZerosIntoExponent()

	return quotient, remainder.normalisedForm_endingZerosIntoExponent(), err
}