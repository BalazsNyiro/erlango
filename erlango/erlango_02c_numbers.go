/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite


http://erlang.org/documentation/doc-6.0/doc/reference_manual/data_types.html

1> 42.
42
2> $A.
65
3> $\n.
10
4> 2#101.
5
5> 16#1f.
31
6> 2.3.
2.3
7> 2.3e3.
2.3e3
8> 2.3e-3.
0.0023

*/

package erlango

import (
	"fmt"
	"strings"
)


func expression_detect_numbers(tokensOrExpressionsOld TokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated string) TokensOrExpressions {
	if ! (strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "numbers") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "detectAllExpressions")) {
		// if number detection is not a wanted operation, then don't do that
		return tokensOrExpressionsOld
	}

	/*
		https://www.erlang.org/doc/reference_manual/data_types.html

		Number variations - the + or - unary operators are NOT detected here
		12
		12_34


		$A      1 char after $
		$\n     2 char after $

		2#101   base#value, integer result
		2#1_01  base#value_with_underscore
		16#1f   base#value, characters can be interpreted as num elems (f)
		16#1F	base#value, CAPITAL chars are interpreted, too
		1_6#1f  base_with_underscore#value
		1_6#1_f  base_with_underscore#value

		12.34
		12_34.56
		12_34.56_78

		2.3e3
		2.3e+3
		2.3e-3
		2_3.4e+3
		2_3.4e+3_0   result: 2.34e31

		2.0E3   # capital E

		F = 0_13.
		13

		G = 16#11111111111111111111111111111.
		5538449982437149470432529417834769

		scientific notation can be mixed with non-decimal numbers:
		16#1e-4.
		26



		_ cannot be the first elem of a number:
			E = _2_3.
			* 1:5: variable '_2_3' is unbound


		An interesting article: https://erlang.org/pipermail/erlang-questions/2019-March/097474.html




		=== General number detection algorithm: ===

		definitions, that are used in number detection:

		numberBlockDecimal: token, which
			- contains only 0123456789_
			- doesn't start with _


		numberBlockDecimal_abc_ABC:
			- contains only 0123456789_abcdefghijklmnopqrstuvwxyz
			- doesn't start with _

		numberBlockSeparators:
			e E . # e+ e- E+ E-


		tokenPrev1, tokenPrev2, tokenPrevN:
			- previous token, previous-previous token...

		tokenNext1, tokenNext2, tokenNextN:
			- next, next-next...




		=== What is a number? ===

		NUM_SEPARATOR_USAGE == 0
		integer - simple (123):
			- tokenPrev1 is not numberBlockSeparator
			- numberBlockDecimal,
			- tokenNext1 is not numberBlockSeparator

		NUM_SEPARATOR_USAGE == 1  (there is one separator)
		float - simple (16.1) or hexa simple (16#1):
			- tokenPrev1 is not numberBlockSeparator
			- numberBlockDecimal,
			- tokenNext1 is numberBlockSeparator
			- tokenNext2 is numberBlockDecimal
			- tokenNext3 is not numberBlockSeparator



		NUM_SEPARATOR_USAGE == 2  (there are two separators)
		example: 2.3e+3
		the first separator is '.'   second separator is 'e+'  in this case

		A possible number representation:
			- tokenPrev1 is not numberBlockSeparator
			- numberBlockDecimal,
			- tokenNext1 is numberBlockSeparator
			- tokenNext2 is numberBlockDecimal_abc_ABC
			- tokenNext1 is numberBlockSeparator
			- tokenNext2 is numberBlockDecimal_abc_ABC

		something is MAYBE a number if it is matching with these rules.
		So, these will be parsed, and the numbers detected - or not detected.

	*/



	/* Number detection is tricky, because more tokens have to be analysed, same time.

	number detection happens in more turn.

	What a fun: the numbers are very close to atoms and strings, as I see it now.
	and I think number detection is the most complex part of the game.
	*/

	// NUMBER_BLOCK_SEPARATORS := []string{".", "#", "e", "E", "e+", "e-", "E+", "E-"}
	const DIGITS_UNDERSCORE = digitsDecimal + "_"
	const DIGITS_UNDERSCORE_abc_ABC = DIGITS_UNDERSCORE + abcEngLower + abcEngUpper



	// locally used functions only, for num detection
	isHashmark := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "#"
	}

	isDot:= func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isDot: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "."
	}

	isDigitAndAlphabetBlock := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isDigitAndAlphabetBlock: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.charAllInPassedCharacterSet(DIGITS_UNDERSCORE_abc_ABC)
	}

	isDigitOnlyBlock := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isDigitOnlyBlock: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.charAllInPassedCharacterSet(DIGITS_UNDERSCORE)
	}

	isPlusMinus := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isPlusMinus: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "+" || tokenOrExpression.token.stringRepresentation() == "-"
	}






	///////////////////// PREPARE NUM DETECTION WITH BLOCKS ////////////////////////////////////////
	// for number detections, atoms are important, because in some case an atom can be the part of a number.
	// for example: 16#ff  ff can be an atom, 16 can be a num, # can be an operator, from a special perspective.

	tokensOrExpressionsNew_blockDetection := TokensOrExpressions{}
	for _, tokenOrExpression := range(tokensOrExpressionsOld) {

		if tokenOrExpression.isExpression() {  // if it is a previously detected expression, there is nothing to do
			tokensOrExpressionsNew_blockDetection= append(tokensOrExpressionsNew_blockDetection, tokenOrExpression)
			continue
		}

		fmt.Println("detect DIGITS_UNDERSCORE blocks:", tokenOrExpression)

		// in numbers _ can be used BUT _ cannot be the first character
		if tokenOrExpression.token.charFirstRuneVal() != '_' {

			// atoms and strings were detected previously - so if there are digits+mixed characters, it can be a hexadecimal number block
			if tokenOrExpression.token.charAllInPassedCharacterSet(DIGITS_UNDERSCORE_abc_ABC) {
				tokenOrExpression.token.TokenType = tokenType_DigitAndAlphabet_Block__numberDetectionStepInitial
			}
		}

		tokensOrExpressionsNew_blockDetection = append(tokensOrExpressionsNew_blockDetection, tokenOrExpression)

	} // FOR, block detection



	/////////////////////////////////////////////////////////////

	tokensOrExpressionsNew_numsDetected := TokensOrExpressions{}
	lenTokenOrExpressions :=  len(tokensOrExpressionsNew_blockDetection)

	numberTokenElems := TokensOrExpressions{ }
	if lenTokenOrExpressions > 0 { // so if there is something to check

		idTokenOrExpr := -1
		for {

			// because id is manipulated inside the for loop, we can move forward more with 1,
			// if it is necessary, with more than 1 token processing
			idTokenOrExpr++
			if idTokenOrExpr >= lenTokenOrExpressions { break }


			tokenOrExpressionActual := getTokenOrExpression_fromLot(idTokenOrExpr,   tokensOrExpressionsNew_blockDetection)
			tokenOrExpressionNext1  := getTokenOrExpression_fromLot(idTokenOrExpr+1, tokensOrExpressionsNew_blockDetection)
			tokenOrExpressionNext2  := getTokenOrExpression_fromLot(idTokenOrExpr+2, tokensOrExpressionsNew_blockDetection)
			tokenOrExpressionNext3  := getTokenOrExpression_fromLot(idTokenOrExpr+3, tokensOrExpressionsNew_blockDetection)
			tokenOrExpressionNext4  := getTokenOrExpression_fromLot(idTokenOrExpr+4, tokensOrExpressionsNew_blockDetection)

			fmt.Println("detect numbers - token expression  id token", idTokenOrExpr)


			isNum := false
			/////////////////////////////////////////////////////////////////////








			/*
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
			*/

			// detect only if num was NOT detected in this turn
			if (! isNum) {
				if isDigitAndAlphabetBlock(tokenOrExpressionActual) {
					if isDot(tokenOrExpressionNext1) || isHashmark(tokenOrExpressionNext1) {
						if isDigitAndAlphabetBlock(tokenOrExpressionNext2) {
							if isPlusMinus(tokenOrExpressionNext3) {
								if isDigitAndAlphabetBlock(tokenOrExpressionNext4) {

									isNum = true
									numberTokenElems = TokensOrExpressions{
										// Actual + next 4 elem forms a number!
										TokenOrExpression{token: tokenOrExpressionActual.token},
										TokenOrExpression{token: tokenOrExpressionNext1.token},
										TokenOrExpression{token: tokenOrExpressionNext2.token},
										TokenOrExpression{token: tokenOrExpressionNext3.token},
										TokenOrExpression{token: tokenOrExpressionNext4.token},
									}

									// move the id from Actual => Actual + 4, so for loop will skip the next 4 elems
									idTokenOrExpr += 4
								}
							}
						}
					}
				}
				if ! isNum {
					fmt.Println("NO - num detection, 4 elem")
				}
			}






			// 1_6#4e
			// detect only if num was NOT detected in this turn
			if (! isNum) {
				fmt.Println("hexa try?")
				if isDigitAndAlphabetBlock(tokenOrExpressionActual) {
					if isHashmark(tokenOrExpressionNext1) {
						if isDigitAndAlphabetBlock(tokenOrExpressionNext2) {
							if isDigitAndAlphabetBlock(tokenOrExpressionNext3) {

								isNum = true
								numberTokenElems = TokensOrExpressions{
									// Actual + next 4 elem forms a number!
									TokenOrExpression{token: tokenOrExpressionActual.token},
									TokenOrExpression{token: tokenOrExpressionNext1.token},
									TokenOrExpression{token: tokenOrExpressionNext2.token},
									TokenOrExpression{token: tokenOrExpressionNext3.token},
								}
								idTokenOrExpr += 3 // move the id from Actual => Actual + 3
							}
						}
					}
				}
				if ! isNum {
					fmt.Println("hexa try: NO - num detection, like hexa ")
				}
			}





			// 1_6.4
			// detect only if num was NOT detected in this turn
			if (! isNum) {
				fmt.Println("float try: ")
				if isDigitOnlyBlock(tokenOrExpressionActual) {
					if isDot(tokenOrExpressionNext1){
						if isDigitOnlyBlock(tokenOrExpressionNext2) {
							isNum = true
							numberTokenElems = TokensOrExpressions{
								TokenOrExpression{token: tokenOrExpressionActual.token},
								TokenOrExpression{token: tokenOrExpressionNext1.token},
								TokenOrExpression{token: tokenOrExpressionNext2.token},
							}
							idTokenOrExpr += 2  // 3-1 = 2, after the Actual,
							// 2 next token were used, they have to be skipped in next turns
						}
					}
				}
			}
			if ! isNum {
				fmt.Println("float try: NO - num detection ")
			}


			// integers
			if (! isNum) {
				fmt.Println("int try: ")
				if isDigitOnlyBlock(tokenOrExpressionActual) {
					isNum = true
					numberTokenElems = TokensOrExpressions{
						TokenOrExpression{token: tokenOrExpressionActual.token},
					}
				}
			}
			if ! isNum {
				fmt.Println("int try: NO - num detection ")
			}



			/////////////////////////////////////////////////////////////////////
			if isNum {
				tokenOrExpressionActual.elemType = tokenOrExpression_thisIsAnExpression
				tokenOrExpressionActual.expression = ErlExpression{
					ExpressionType:			expression_num,
					TokensOrExpressions: numberTokenElems,
				}
				// put back tokenOrExpression with modified elemType and expression
				tokensOrExpressionsNew_numsDetected = append(tokensOrExpressionsNew_numsDetected, tokenOrExpressionActual)

			} else {  // not an atom - put back the tokenOrExpression without any extra change/modification
				tokensOrExpressionsNew_numsDetected = append(tokensOrExpressionsNew_numsDetected, tokenOrExpressionActual)
			}
		} // FOR

	} // len > 0

	return tokensOrExpressionsNew_numsDetected
}

