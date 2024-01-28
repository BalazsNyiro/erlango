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
	"slices"
	"strconv"
	"strings"
)

const DebugNumDetection_printDetails bool = true

/*
Receives Erlang source code - return with non-detected source code and detected Tokens.

First I planned to use the same approach that was used with strings, too - to find openers, and closers.
But numbers can be represented with a lot of forms.

So this section will be a little different :-)

The main idea:
 - take the next character to analyse (actual char is selected)
 - look forward, find matching character ranges.
 - if you find something that is matching with a segment of a number-representation, look forward again

So with other words, the actual char is analysed one-by one, and the func always looks forward.
If a number representation form is detected, the whole block is removed from the src, and added to the tokens.

If the look-forward is not successful, then take the next char, and start again the detection

*/
func Tokens_detect_numbers(erlSrc string, tokensTable Tokens) (string, Tokens) {

	tokensTableUpdated := tokensTable.deepCopy()
	var erlSrcTokenDetectionsRemoved []rune
	///////////////////////////////////////////////////////////////
	digitsZeroNine := []rune("0123456789")
	digitsZeroNine_underscore := []rune("0123456789_")
	digitDot := []rune(".")
	digitHashmark := []rune("#")


	//digitE := []rune("e")
	//digitPlusMinus := []rune("+-")

	const str_abcEngLower = "abcdefghijklmnopqrstuvwxyz"
	const str_abcEngUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const str_digits_abc_ABC = "0123456789" + str_abcEngLower + str_abcEngUpper
	const str_digits_abc_ABC_underscore = str_digits_abc_ABC + "_"

	digitsZeroNine_abc_ABC := []rune(str_digits_abc_ABC)
	digitsZeroNine_abc_ABC_underscore := []rune(str_digits_abc_ABC_underscore)
	///////////////////////////////////////////////////////////////








	erlSrcRunes := []rune(erlSrc)
	// be careful: charPos can be increased inside the loop!!!!
	for charPos := 0; charPos < len(erlSrcRunes); charPos++ {

		tokenType := ""
		detectionCharPosFirst := -1
		detectionCharPosLast := -1
		//............................................................................................





		/* num representations, from more complicated to simple direction:

		1a: 1_6#4fe+3_0. - digitAlphabet, dotOrHashmark, digitAndAlphabet, plusMinus, digitAndAlphabet
		1b: 1_6.4e+3_0. - digitAlphabet, dotOrHashmark, digitAndAlphabet, plusMinus, digitAndAlphabet

		2: 1_6#4f    - digitAlphabet, hashmark, digitAlphabet    (hexa and nondecimal)
		3: 1_2.3_4   - digitOnly, dot, digitOnly                 (float)
		4a: 1_234_5  - digitOnly with underscores                (simple integer)
		4b: 12345    - digitOnly                                 (simple integer)
		5: $A        - $ + oneChar, or 2 char if it is escaped   (char literals)
		*/

		if tokenType == "" { // hexa and non-decimal detection (2)
			detection := "hexa/nondecimal detections"
			fmt.Println(detection+" >>" + string(erlSrcRunes[charPos:]) + "<<")
			/* check possible variations:
			A = 16#4f.
			B = 1_6#4f.
			C = 1_6#4_f.
			D = 16#4_f.
			*/

			digit______hash_digitAbc      := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine,                            digitHashmark, digitsZeroNine_abc_ABC                                   }, "right", "detected_num_digitHashmarkDigitAbc")
			digitUnder_hash_digitAbc      := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitHashmark, digitsZeroNine_abc_ABC                                   }, "right", "detected_num_digitUnderscoreHashmarkDigitAbc")
			digitUnder_hash_digitAbcUnder := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitHashmark, digitsZeroNine_abc_ABC, digitsZeroNine_abc_ABC_underscore}, "right", "detected_num_digitUnderscoreHashmarkDigitAbcUnderscore")
			digit______hash_digitAbcUnder := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine,                            digitHashmark, digitsZeroNine_abc_ABC, digitsZeroNine_abc_ABC_underscore}, "right", "detected_num_digitHashmarkDigitAbcUnderscore")

			printIfDebug(detection+" 1: digit______hash_digitAbc      ", digit______hash_digitAbc)
			printIfDebug(detection+" 2: digitUnder_hash_digitAbc      ", digitUnder_hash_digitAbc)
			printIfDebug(detection+" 3: digitUnder_hash_digitAbcUnder ", digitUnder_hash_digitAbcUnder)
			printIfDebug(detection+" 4: digit______hash_digitAbcUnder ", digit______hash_digitAbcUnder)

			detectedMax := max(digitUnder_hash_digitAbc, digit______hash_digitAbc, digitUnder_hash_digitAbcUnder, digit______hash_digitAbcUnder)
			fmt.Println(detection+" detectedMax", detectedMax)

			if detectedMax > 0 {
				tokenType = tokenType_Num_maybeNonDecimal
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detectedMax - 1
			}
			printIfDebug(detection+" tokenType:", tokenType)
		} // hexa/nondecimal



		if tokenType == "" { // float detection (3)
			detection := "float detection (3)"
			printIfDebug(detection + " >>" + string(erlSrcRunes[charPos:]) + "<<")

			/* check possible 4 float variations:
				12.3_4
				12.34
				1_2.34
				1_2.3_4

			*/
			detected_num_digit_dot_digitUnderscore           := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitDot, digitsZeroNine, digitsZeroNine_underscore}, "right", "detected_num_digit_dot_digitUnderscore")
			detected_num_digit_dot_digit                     := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitDot, digitsZeroNine}, "right", "detected_num_digit_dot_digit")
			detected_num_digitUnderscore_dot_digit           := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDot, digitsZeroNine}, "right", "detected_num_digitUnderscore_dot_digit")
			detected_num_digitUnderscore_dot_digitUnderscore := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDot, digitsZeroNine, digitsZeroNine_underscore}, "right", "detected_num_digitUnderscore_dot_digitUnderscore")

			// fmt.Println("detected_num_digit_dot_digitUnderscore           ", detected_num_digit_dot_digitUnderscore           )
			// fmt.Println("detected_num_digit_dot_digit                     ", detected_num_digit_dot_digit                     )
			// fmt.Println("detected_num_digitUnderscore_dot_digit           ", detected_num_digitUnderscore_dot_digit           )
			// fmt.Println("detected_num_digitUnderscore_dot_digitUnderscore ", detected_num_digitUnderscore_dot_digitUnderscore )

			detectedMax := max(detected_num_digit_dot_digitUnderscore, detected_num_digit_dot_digit, detected_num_digitUnderscore_dot_digit, detected_num_digitUnderscore_dot_digitUnderscore)

			if detectedMax> 0 {
				tokenType = tokenType_Num_float
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detectedMax - 1
			}
			printIfDebug(detection+" tokenType:", tokenType, "\n")
		} // float


		if tokenType == "" { // simple INT detection (4a): digit+|digit_underscore* (one or more digits|one ore more digitsAndUnderscore
			detection := "simple INT detection (4a)"
			printIfDebug(detection + " >>" + string(erlSrcRunes[charPos:]) + "<<")

			detected_num_of_digitsZeroNine_underscore := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore}, "right", "simple INT 4a detection")
			if detected_num_of_digitsZeroNine_underscore > 0 {
				tokenType = tokenType_Num_int
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detected_num_of_digitsZeroNine_underscore - 1
			}

			printIfDebug(detection+" tokenType:", tokenType, "\n")
		}


		if tokenType == "" { // simple INT detection (4b) digit+ (one or more digits)
			detection := "simple INT detection (4b)"
			printIfDebug(detection + " >>" + string(erlSrcRunes[charPos:]) + "<<")
			detected_num_of_digitsZeroNine := charsHowManyAreInTheGroup(charPos, erlSrcRunes,digitsZeroNine, "right")

			if detected_num_of_digitsZeroNine > 0 {
				tokenType = tokenType_Num_int
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detected_num_of_digitsZeroNine - 1
			}
			printIfDebug(detection+" tokenType:", tokenType, "\n")
		}


		if tokenType == "" { // char literals (5)
			detection := "char literals (5)"
			printIfDebug(detection + " >>" + string(erlSrcRunes[charPos:]) + "<<")

			charNow := erlSrcRunes[charPos]
			charNext1, next1InSrc := charRuneNext(charPos, +1, erlSrcRunes)
			_,         next2InSrc := charRuneNext(charPos, +2, erlSrcRunes)

			if charNow == '$' {
				tokenType = tokenType_Num_charLiterals
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos

				if next1InSrc {
					detectionCharPosLast = charPos + 1
				}

				if charNext1 == '\\' && next2InSrc { // if the value is escaped, and charNext2 is in the source code, read that, too
					detectionCharPosLast = charPos + 2
				}
			} // $ detected
			printIfDebug(detection+" tokenType:", tokenType, "\n")
		}


		/////////////// GENERAL TOKEN SAVE SECTION //////////////////////

		if tokenType != "" { // Token was detected, the type is NOT empty


			selectedSrc := string(erlSrcRunes[detectionCharPosFirst:detectionCharPosLast+1])
			tokenTypeUpdated, msgFromParser := number_token_validation(tokenType, selectedSrc)

			tokenNow := Token{ 	positionCharFirst: detectionCharPosFirst,
								positionCharLast: detectionCharPosLast,
								tokenType: tokenTypeUpdated, msgFromParser: msgFromParser}

			// tokenType is set IF there is minimum one char, so this loop is always executed, minimum ONCE
			for posTokenChar := detectionCharPosFirst; posTokenChar <= detectionCharPosLast; posTokenChar++ { // 0..9 digit block is detected
				tokenNow.charsInErlSrc = append(tokenNow.charsInErlSrc, erlSrcRunes[posTokenChar])
				erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, ' ')
				charPos ++
			}
			charPos -- // in the foor loop, one unnecessary Pos increasing happens,
			// because if one char is detected only, charPos doesn't need to be changed

			tokensTableUpdated[tokenNow.positionCharFirst] = tokenNow

		} else {
			// there is no token detection - save the rune back into the original src
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, erlSrcRunes[charPos])
		}
		/////////////// GENERAL TOKEN SAVE SECTION //////////////////////


	} // for charPos

	///////////////////////////////////////////////////////////////
	return string(erlSrcTokenDetectionsRemoved), tokensTableUpdated
}



/*
	receives a tokenType and a string representation of the token.

	return with a tokenType and a message:
		- if the num is valid, then with the original tokenType,
		- otherwise with a syntax error
 */
func number_token_validation(tokenType, erlSrc string) (string, string) {

	errMsgForUser := ""
	//////////////////////////////////////
	// 16#4__a is invalid
	if strings.Contains(erlSrc, "__"){
		printIfDebug("validation: __")
		errMsgForUser += "more than 1 underscore in the number: "+erlSrc+" ;"
	}

	// 16_#4a  is invalid
	if strings.Contains(erlSrc, "_#"){
		printIfDebug("validation: _#")
		errMsgForUser += "invalid non-decimal number (underscore/hashmark): "+erlSrc+" ;"
	}
	// 16#_4a  is invalid
	if strings.Contains(erlSrc, "#_"){
		printIfDebug("validation: #_")
		errMsgForUser += "invalid non-decimal number (hashmark/underscore): "+erlSrc+" ;"
	}


	// _16 is invalid: this is an unbound variable!!! and cannot be detected as a number
	// if strings.HasPrefix(erlSrc, "_"){
	// 	errMsgForUser += "a number cannot start with _ sign: "+erlSrc+" ;"
	// }

	// 16_ is invalid
	if strings.HasSuffix(erlSrc, "_"){
		printIfDebug("validation: suffix _")
		errMsgForUser += "a number cannot end with _ sign: "+erlSrc+" ;"
	}


	if len(erlSrc) < 1 {
		errMsg := "a number token needs minimum 1 char, cannot be empty!"
		printIfDebug(errMsg)
		errMsgForUser += errMsg
	}


	/////////////////// # in number:
	if strings.Contains(erlSrc, "#") {

		printIfDebug("validation: # in number...")

		// const ABC_Eng_Upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		const ABC_Eng_Lower string = "abcdefghijklmnopqrstuvwxyz"
		const ABC_Eng_digits string = "0123456789"

		const digitsAll string = ABC_Eng_digits + ABC_Eng_Lower

		elems := strings.Split(erlSrc, "#")

		baseStr := elems[0]
		numStr := strings.ToLower(elems[1]) // use only lowercase representation
		// remove _ from base/num - example case: 1_6#4_e.
		baseStr = numDetect_removeUnderscoreFromString(baseStr)
		numStr  = numDetect_removeUnderscoreFromString(numStr)

		base , errIntConversation := strconv.Atoi(baseStr)
		if errIntConversation != nil {
			printIfDebug("validation: incorrect base before #: " + erlSrc)
			errMsgForUser += "incorrect base before # separator: "+erlSrc+" ;"

		} else { // base can be converted to int
			digitsAllAccepted := []rune(digitsAll[0:base]) // 10 based last accepted: 9, 16 based last accepted: f

			for _, runeDigit := range numStr {
					if slices.Contains(digitsAllAccepted, runeDigit) {
						// digit is valid, nothing to do
					} else {
						printIfDebug("validation: incorrect digit: " + string(runeDigit))
						errMsgForUser += "incorrect digit ("+string(runeDigit)+") in num representation, after # separator: "+erlSrc+" ;"
					} // digit is not accepted

			} // for, runeDigit

		} // else, base converted to int

	} // if


	if len(errMsgForUser) > 0 {
		tokenType = tokenType_SyntaxError
	}
	//////////////////////////////////////
	return tokenType, errMsgForUser

}

func printIfDebug(a ...any) {
	if DebugNumDetection_printDetails {
		fmt.Println(a...)
	}
}
func numDetect_removeUnderscoreFromString(txt string) string {
	return strings.Replace(txt, "_", "", -1)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
// similar erlang lib: https://github.com/mabrek/erlang-decimal
type erlango_bignum struct {
	/* example num representation: '12.34'
		digits: 1234
		decimalPoint: -2    because 1234 * (10^-2) = 12.34

	   example 2: 1234000
		digits: 1234
		decimalPoint: 3     because 1234 * (10^3) = 1234000

		with this representation, integers and floats can be represented too, with the same structure
	 */

	digits []int
	exponent int  // where is the ., or how many 0 are after the digits
}


func digit_value(digit rune) (int, error) {
	valueMap := map[rune]int{
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

// if the token is a number, return with a value and and OK
// if it is NOT a number, return with 0 and error
func number_value_integerDecimal_pair_detect_from_token(token Token) (erlango_bignum, error)  {
	valueDetectedFromToken := false

	zero := erlango_bignum{digits: []int{0}, exponent: 0}
	partInteger := []int{}

	if token.tokenType == tokenType_Num_int {
		valueDetectedFromToken = true

		// naive number parsing,
		//  position 210  <- the lastValue = 7*10^0, middle: 1*10^1, firstVal: 2*10^2
		//  example: 567
		// the first digit is 5, first pos is 2:

		// the token needs to have minimum 1 char - there is a validation against that
		for _, digit := range token.charsInErlSrc {

			// rune -> numberValue conversion
			// values are represented in decimal Num system!!!
			// so f = 15, a = 10
			digitVal_inDecimalNumSystem, errorValueDetection := digit_value(digit)

			if errorValueDetection != nil {
				return zero, errors.New("number value detection error ("+token.stringRepr()+")")
			} else {
				partInteger = append(partInteger, digitVal_inDecimalNumSystem)
			}

		}

	}

	if ! valueDetectedFromToken {
		return zero, errors.New("number value detection error ("+token.stringRepr()+")")
	}
	exponent := 0 // here we set the integer part only now
	// TODO: later normalize the number here! so (1000, 0) -> (1, 3)!!
	return erlango_bignum{digits: partInteger, exponent: exponent}, nil
}