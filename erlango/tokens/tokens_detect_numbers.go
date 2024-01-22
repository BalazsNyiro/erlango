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

import "fmt"

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
	//digitHashmark := []rune("#")
	//digitE := []rune("e")
	//digitPlusMinus := []rune("+-")

	//# const abcEngLower = "abcdefghijklmnopqrstuvwxyz"
	//const abcEngUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//const DIGITS_UNDERSCORE_abc_ABC = DIGITS_UNDERSCORE + abcEngLower + abcEngUpper
	///////////////////////////////////////////////////////////////








	erlSrcRunes := []rune(erlSrc)
	// be careful: charPos can be increased inside the loop!!!!
	for charPos := 0; charPos < len(erlSrcRunes); charPos++ {

		tokenType := ""
		detectionCharPosFirst := -1
		detectionCharPosLast := -1
		//............................................................................................





		/* num representations, from more complicated to simple direction:

		1: 2_3#4e+3_0 - digitAlphabet, dotOrHashmark, digitAndAlphabet, plusMinus, digitAndAlphabet

			   vvv digitAndAlphabet
			      v isDot OR isHashmark
			       vv digitOrAlphabet (e is incluced in the alphabet)
			         v plusminus
			          vvv digitOrAlphabet
			   2_3.4e+3_0,

			      v isDot OR isHashmark
			   2_3#4e+3_0,

			 erl shell> 2_3#4e+3_0.
			            136

		2: 1_6#4e    - digitAlphabet, hashmark, digitAlphabet    (hexa and nondecimal)
		3: 1_2.3_4   - digitOnly, dot, digitOnly                 (float)
		4a: 1_234_5  - digitOnly with underscores                (simple integer)
		4b: 12345    - digitOnly                                 (simple integer)
		5: $A        - $ + oneChar, or 2 char if it is escaped   (char literals)
		*/

		if tokenType == "" { // hexa and non-decimal detection (2)
			fmt.Println("hexa/nondecimal detections >>" + string(erlSrcRunes[charPos:]) + "<<")
			/* check possible variations:
			A = 16#4e.
			B = 1_6#4e.
			C = 1_6#4_e.
			D = 16#4_e.
			*/

			detected_num_digit_dot_digitUnderscore           := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitDot, digitsZeroNine, digitsZeroNine_underscore}, "right", "detected_num_digit_dot_digitUnderscore")
			detected_num_digit_dot_digit                     := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitDot, digitsZeroNine}, "right", "detected_num_digit_dot_digit")
			detected_num_digitUnderscore_dot_digit           := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDot, digitsZeroNine}, "right", "detected_num_digitUnderscore_dot_digit")
			detected_num_digitUnderscore_dot_digitUnderscore := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore, digitDot, digitsZeroNine, digitsZeroNine_underscore}, "right", "detected_num_digitUnderscore_dot_digitUnderscore")

			detectedMax := max(detected_num_digit_dot_digitUnderscore, detected_num_digit_dot_digit, detected_num_digitUnderscore_dot_digit, detected_num_digitUnderscore_dot_digitUnderscore)

			if detectedMax> 0 {
				tokenType = tokenType_Num_maybeNonDecimal
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detectedMax - 1
			}
		} // float



		if tokenType == "" { // float detection (3)
			fmt.Println("float detections >>" + string(erlSrcRunes[charPos:]) + "<<")
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
		} // float


		if tokenType == "" { // simple INT detection (4a): digit+|digit_underscore* (one or more digits|one ore more digitsAndUnderscore
			detected_num_of_digitsZeroNine_underscore := charsGroupsAreMatching( charPos, erlSrcRunes,[]([]rune){digitsZeroNine, digitsZeroNine_underscore}, "right", "simple INT 4a detection")
			if detected_num_of_digitsZeroNine_underscore > 0 {
				tokenType = tokenType_Num_int
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detected_num_of_digitsZeroNine_underscore - 1
			}
		}


		if tokenType == "" { // simple INT detection (4b) digit+ (one or more digits)
			detected_num_of_digitsZeroNine := charsHowManyAreInTheGroup(charPos, erlSrcRunes,digitsZeroNine, "right")
			if detected_num_of_digitsZeroNine > 0 {
				tokenType = tokenType_Num_int
				detectionCharPosFirst = charPos
				detectionCharPosLast = charPos + detected_num_of_digitsZeroNine - 1
			}
		}


		if tokenType == "" { // char literals (5)
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
		}


		/////////////// GENERAL TOKEN SAVE SECTION //////////////////////
		if tokenType != "" { // Token was detected, the type is NOT empty

			tokenNow := Token{ 	positionCharFirst: detectionCharPosFirst,
								positionCharLast: detectionCharPosLast,
								tokenType: tokenType}

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
