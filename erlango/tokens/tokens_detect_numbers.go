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
	digitsZeroNine := []rune(ABC_Eng_digits)
	digitsZeroNine_underscore := []rune(ABC_Eng_digits+"_")

	digitDot := []rune(".")
	digitHashmark := []rune("#")
	digitDotOrHashmark := []rune(".#")


	digitE := []rune("e")
	digitPlusMinus := []rune("+-")

	const str_digits_abc_ABC = ABC_Eng_digits + ABC_Eng_Lower + ABC_Eng_Upper
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

		0a: the Queen of the numbers. Oh, man. I will die with this:
		   2> 1_6#4ee+4#1.
		   1263
		0b: and this is valid, too: but this is not one number, this is 2.
			1_6#4e+-4#1.
			77

			1_6#4e - 4#1
			77


				1a: 1_6#4fe+3_0. - digitAlphabet, dotOrHashmark, digitAndAlphabet, plusMinus, digitAndAlphabet
		1b: 1_6.4e+3_0. - digitAlphabet, dotOrHashmark, digitAndAlphabet, plusMinus, digitAndAlphabet

		2: 1_6#4f    - digitAlphabet, hashmark, digitAlphabet    (hexa and nondecimal)
		3: 1_2.3_4   - digitOnly, dot, digitOnly                 (float)
		4a: 1_234_5  - digitOnly with underscores                (simple integer)
		4b: 12345    - digitOnly                                 (simple integer)
		5: $A        - $ + oneChar, or 2 char if it is escaped   (char literals)
		*/

		if tokenType == "" {
			detection := "nondecimal_float_scientist 1a, 1b"
			digit______hash_digitAbc      := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine,                            digitDotOrHashmark, digitsZeroNine_abc_ABC                                   , digitE, digitPlusMinus, digitsZeroNine}, "right", "detected_num_digitHashmarkDigitAbc")
			digitUnder_hash_digitAbc      := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDotOrHashmark, digitsZeroNine_abc_ABC                                   , digitE, digitPlusMinus, digitsZeroNine}, "right", "detected_num_digitUnderscoreHashmarkDigitAbc")
			digitUnder_hash_digitAbcUnder := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDotOrHashmark, digitsZeroNine_abc_ABC, digitsZeroNine_abc_ABC_underscore, digitE, digitPlusMinus, digitsZeroNine}, "right", "detected_num_digitUnderscoreHashmarkDigitAbcUnderscore")
			digit______hash_digitAbcUnder := charsGroupsAreMatching(charPos, erlSrcRunes,[]([]rune){digitsZeroNine,                            digitDotOrHashmark, digitsZeroNine_abc_ABC, digitsZeroNine_abc_ABC_underscore, digitE, digitPlusMinus, digitsZeroNine}, "right", "detected_num_digitHashmarkDigitAbcUnderscore")

			detectedMax := max(digitUnder_hash_digitAbc, digit______hash_digitAbc, digitUnder_hash_digitAbcUnder, digit______hash_digitAbcUnder)
			fmt.Println(detection+" detectedMax", detectedMax)

			if detectedMax > 0 {
				tokenType = tokenType_Num_maybeNonDecimal_scientific
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detectedMax - 1
			}
		}




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

	// 16_ is invalid: the last char cannot be an underscore, something is always has to be after the underscore
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

//////////////////////////// FROM HERE, move everything into bignum: //////////////////////////////////////////////////////////////////////////////////






/*
	if decimals are converted to bigNum, that is simple, because the values can be directly loaded
    into the digits, without any calculation

	The general solution is more complex, but important for hexa and other num systems, TODO: maybe next time
*/
func bigNum_from_digits_specialcase_decimalintegergeneral (token Token) (bignum_decimalValue, error) {

	digits := digitList{}
	for pos := 0; pos < len(token.charsInErlSrc); pos++ {
		digit := token.charsInErlSrc[pos]
		fmt.Println("digit[",pos,"] => ", digit, string(digit))

		if digit == '_'	 { // first I removed all _ in the root point. then I realised that I change the characters with that action
			continue	   // so unfortunately, the most native/correct ways is to accept the _ but not to do anything
		}

		// rune -> numberValue conversion, in range 0-9
		digitValueDecimalInteger, errorValueDetection := digitRune_decimalValue(digit)
		if errorValueDetection != nil {
			return bigNum_zero(), errors.New("digit (" + string(digit) + ") value detection error in: " + token.stringRepr())
		}
		digits = append(digits, digitValueDecimalInteger)
	}
	return bignum_decimalValue{digits: digits, exponent: 0}, nil
}







func bigNum_from_digits_general_any_numsystem (token Token) (bignum_decimalValue, error) {
	return bigNum_from_digits_specialcase_decimalintegergeneral(token)
}

////////////////// 1_6#4ee+4 //////////////////////
////////////////// 1_6#4ee+4 //////////////////////
////////////////// 1_6#4ee+4 //////////////////////

// TESTED
// select the scientific part from a number. by default,
// there is no scientific part
func anyNumSystem_charsSelectScientificPart(chars []rune) (bool, []rune) {
	numberCharsScientific := []rune{} // empty, or: e+3_0 | e-2_0_0 typed.
	scientificEsignDetected := false

	splitterPlusDetectedMinimumOnce, _, charsRightPlus:= charsCopySplitAtFirstWithChars(chars, []rune("e+"))

	if splitterPlusDetectedMinimumOnce {
		// everything in the left part is before # sign
		numberCharsScientific = charsRightPlus
		scientificEsignDetected = true

	} else { // if splitter plus is not detected, try splitterMinus
		splitterMinusDetectedMinimumOnce, _, charsRightMinus := charsCopySplitAtFirstWithChars(chars, []rune("e-"))
		if splitterMinusDetectedMinimumOnce {
			numberCharsScientific = charsRightMinus
			scientificEsignDetected = true
		}
	}
	return scientificEsignDetected, numberCharsScientific
}


// select the number system from the token
// basically every number is 10 based, if there is no other sign
func anyNumSystem_detectNumSystem(chars []rune) (bignum_decimalValue, error) {
	numberSystemType := bigNum_ten()

	// _ is removed in 'general_any_numsystem', but if this func
	// is used standalone, the underscores has to be removed, too
	chars = charsCopyRemoveUnwanted(chars, '_')

	fmt.Println(">>> Chars: ", chars)

	isHashmarkDetected, charsBeforeHashMark, _:= charsCopySplitAtFirstWithChars(chars, []rune("#"))
	fmt.Println("is hashmark detected:", isHashmarkDetected)
	fmt.Println("chars before hashmark:", charsBeforeHashMark)
	if isHashmarkDetected {
		if len(charsBeforeHashMark) > 0 {
			/* in the number system section, every digit has to be in 0-9 range */
			digits := digitList{}
			for _, char := range charsBeforeHashMark{
				digitValueDecimalInteger, errorValueDetection := digitRune_decimalValue(char)
				if errorValueDetection != nil { //
					errMsg := "digit->decimalValue conversion error in anyNumsystem, numSystem part: " + string(chars)
					fmt.Println(errMsg, errorValueDetection)
					return numberSystemType, errors.New(errMsg)
				}
				if digitValueDecimalInteger > 9 {
					errMsg := "digit->decimalValue in num system you can use 0-9 digits only - too high value error in anyNumsystem, numSystem part: " + string(chars)
					fmt.Println(errMsg, errorValueDetection)
					return numberSystemType, errors.New(errMsg)
				}
				digits = append(digits, digitValueDecimalInteger)
			}
			fmt.Println("digits, from before hashmark", digits)
			numberSystemType = bigNum_from_digitlist(digits)
			fmt.Println("numberSystemType after bigNum creation", numberSystemType)
		} // len > 0
	} // # is detected
	fmt.Println("numberSystemType at the end", numberSystemType)
	return numberSystemType, nil
}


/* analyse all digits, and calculate a decimal based value from a maybe non-decimal input */
func bigNum_from_digits_general_any_numsystemREFACTORTHIS (token Token) (bignum_decimalValue, error) {

	// if # is in the token, this will be updated and '..#' prefix removed
	// basic situaton: numberSystem is 10 based, and there is NO scientific part
	numberCharsBody := charsCopyRemoveUnwanted(token.charsInErlSrc, '_')
	numberCharsScientific := []rune{} // empty, or: e+3_0 | e-2_0_0 typed.
	///////////////////////////////////////////////////////////////////////////////


	numberSystemType, numberSysErr := anyNumSystem_detectNumSystem(numberCharsBody)
	if numberSysErr != nil { //
		errMsg := "digit->decimalValue conversion error in anyNumsystem, numSystem part: " + token.stringRepr()
		fmt.Println(errMsg, numberSysErr)
		return bigNum_zero(), numberSysErr
	}


	// hashMark is removed, if it was there
	charsCurrentStrRepresentation := string(numberCharsBody)
	scientificPlus := strings.Contains(charsCurrentStrRepresentation, "e+")
	scientificMinus := strings.Contains(charsCurrentStrRepresentation, "e-")

	// fill the scientific part if it is there
	if scientificPlus || scientificMinus {
		// e+ or e- is detected in the string
		charsWithE := charsCopy(numberCharsBody)
		for pos, char := range charsWithE {

			// here I can read the next, because I will stop at e+ or at e-,
			// so if we are at the first char of scientific elem, stop the loop
			charNext := charsWithE[pos+1]
			if char == 'e' && (charNext == '-' || charNext == '+') {
				numberCharsScientific = charsWithE[pos:]
				break
			}
			// collect the real number elems only, without the scientific part
			numberCharsBody = charsWithE[:pos+1]
		}
	}


	////////////////////////////////////////////////////////////
	/* so, here there are 3 things:
	numberCharsBody
	numberSystemType
	numberCharsScientific
	*/

	// with this solution, the token is NOT sensitive if the value is missing
	// after the e- or e+
	scientificMultiply := bigNum_one() // multiply with one doesn't change a number value
	if len(numberCharsScientific) > 2 {
		// in the scientific part, only decimal values can be used
		digitsScientific := digitList{}

		for _, char := range numberCharsScientific {
			valSci, err := digitRune_decimalValue(char)
			if err != nil { //
				errMsg := "digit->decimalValue conversion error in anyNumsystem, Sci part: " + token.stringRepr()
				fmt.Println(errMsg, err)
				return bigNum_zero(), errors.New(errMsg)
			}
			digitsScientific = append(digitsScientific, valSci)
		}
		scientificMultiply = bignum_decimalValue{digits: digitsScientific, exponent: 0, negative: numberCharsScientific[1]=='-'}
	}


	summa := bigNum_zero()

	for posChar, char := range numberCharsBody {

		// Do we need to check the number system's possible valid digit range?
		// I mean in an octal number, 9 is not a valid digit.
		valDigit, err := digitRune_decimalValue(char)
		valDigitBigNum := bigNum_from_digitlist(digitList{valDigit})
		if ! valDigitBigNum.isLessThan(numberSystemType) {
			errMsg := "digit->decimalValue conversion error, the digit has a bigger value than it's number system. " + token.stringRepr()
			fmt.Println(errMsg)
			return bigNum_zero(), errors.New(errMsg)
		}

		if err != nil { //
			errMsg := "digit->decimalValue conversion error in anyNumsystem, number analyse, one digit " + token.stringRepr()
			fmt.Println(errMsg, err)
			return bigNum_zero(), errors.New(errMsg)
		}

		// example digits:      456
		// positions fromBack:  210  the last char is in pos 0 from back
		positionFromBack := bigNum_create_from_int( len(numberCharsBody)-1 - posChar)
		multiplier := bigNum_operator_mul(positionFromBack, numberSystemType)
		digitPositionBasedValue := bigNum_operator_mul(valDigitBigNum, multiplier)
		summa = bigNum_operator_add(summa, digitPositionBasedValue)
	}


	// NOT FINISHED. TOTAL REWRITE IS NECESSARY
	// not OK because the scientific num can be any num system, too, not only decimal
	multiplierScientificVal := bigNum_operator_mul(bigNum_ten(), scientificMultiply)

	summa = bigNum_operator_mul(summa, multiplierScientificVal)
	return summa, nil
}
////////////////// 1_6#4ee+4 //////////////////////
////////////////// 1_6#4ee+4 //////////////////////
////////////////// 1_6#4ee+4 //////////////////////





// if the token is a number, return with a value and and OK
// if it is NOT a number, return with 0 and error
func bigNum_from_token(token Token) (bignum_decimalValue, error)  {

	if token.tokenType == tokenType_Num_int {
		// this is the easiest way
		num, err := bigNum_from_digits_specialcase_decimalintegergeneral(token)

		// and this is the hardest :-)
		// num, err :=  bigNum_from_digits_general_any_numsystem(token)

		return num.normalisedForm_endingZerosIntoExponent(), err
	} // Num_int detected
/*
	if token.tokenType == tokenType_Num_maybeNonDecimal{
		num, err := bigNum_from_digits_general_any_numsystem(token)
		return num.normalisedForm_endingZerosIntoExponent(), err
	} // Num_int detected
*/
	return bigNum_zero(), errors.New("number value detection error ("+token.stringRepr()+")")
}

