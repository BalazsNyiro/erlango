/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

http://erlang.org/documentation/doc-6.0/doc/reference_manual/data_types.html
*/

package erlango

import (
	"fmt"
	"strings"
)


func expression_detect_operators(tokensOrExpressionsOld TokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated string) TokensOrExpressions {
	if ! (strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "operators") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "detectAllExpressions")) {
		return tokensOrExpressionsOld
	}

	/*
	https://www.erlang.org/doc/reference_manual/expressions.html

	# 3 chars long:
	=:=	Exactly equal to
	=/=	Exactly not equal to

	# 2 chars long:
	==	Equal to
	/=	Not equal to
	=<	Less than or equal to
	>=	Greater than or equal to

	# 2 chars, list operators:
	--
	++
	|| comprehension, [{E,L} || E <- L=[1,2,3]].
	<- comprehension

	# 2 chars, map operator:
	=>   // map val declaration
	:=   // map update existing value

	# 2 chars, fun expression:
	->

	# 2 chars, bitstring generator:
	<=   BitstringPattern <= BitStringExpr



	######### 1 chars long ###########

	:   modul:func
	#   The record expressions Expr#Name.Field and #Name.Field ???   honestly I have never seen # operator

	=	The = character is used to denote two similar but distinct operators: the match operator and the compound pattern operator. Which one is meant is determined by context.

	<	Less than
	>	Greater than

	# arithmetic operators, 1 chars long
	+	Unary +					Number
	-	Unary -					Number
	+	 						number
	-	 						Number
	*	 						Number
	/	Floating point division	Number




	///////////////////////////////////////////////////////////////////////////////////////////
	// from interpreter perspective, these are detected as ATOMS - so they have to be converted
	// TODO: convert these operators, too
	bnot	Unary bitwise NOT			Integer
	div		Integer division			Integer
	rem		Integer remainder of X/Y	Integer
	band	Bitwise AND					Integer
	bor		Bitwise OR					Integer
	bxor	Arithmetic bitwise XOR		Integer
	bsl		Arithmetic bitshift left	Integer
	bsr		Bitshift right				Integer

	TODO: convert these, too:
	not	Unary 	Logical NOT
	and			Logical AND
	or			Logical OR
	xor			Logical XOR


	*/


	/*
	// locally used functions only, for num detection
	isHashmark := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "#"
	}

	isLess := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isLess: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "<"
	}

	isGreater := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isGreater: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == ">"
	}

	isPlus := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isPlus: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "+"
	}

	isMinus := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isMinus: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "-"
	}
	*/
	isEqual := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isEqual: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == "="
	}

	isColon := func (tokenOrExpression TokenOrExpression) bool {
		fmt.Println("isColon: ", tokenOrExpression.token.stringRepresentation())
		return tokenOrExpression.token.stringRepresentation() == ":"
	}


	tokensOrExpressionsNew := TokensOrExpressions{}
	lenTokenOrExpressions :=  len(tokensOrExpressionsOld)

	operatorTokenElems := TokensOrExpressions{ }
	if lenTokenOrExpressions > 0 { // so if there is something to check

		idTokenOrExpr := -1
		for {

			// because id is manipulated inside the for loop, we can move forward more with 1,
			// if it is necessary, with more than 1 token processing
			idTokenOrExpr++
			if idTokenOrExpr >= lenTokenOrExpressions { break }


			tokenOrExpressionActual := getTokenOrExpression_fromLot(idTokenOrExpr,   tokensOrExpressionsOld)
			tokenOrExpressionNext1  := getTokenOrExpression_fromLot(idTokenOrExpr+1, tokensOrExpressionsOld)
			tokenOrExpressionNext2  := getTokenOrExpression_fromLot(idTokenOrExpr+2, tokensOrExpressionsOld)

			fmt.Println("detect numbers - token expression  id token", idTokenOrExpr)

			isOperator := false
			/////////////////////////////////////////////////////////////////////




			if ! isOperator {   //  =:=
				if isEqual(tokenOrExpressionActual) {
					if isColon(tokenOrExpressionNext1){
						if isEqual(tokenOrExpressionNext2) {

							isOperator = true
							operatorTokenElems = TokensOrExpressions{
								TokenOrExpression{token: tokenOrExpressionActual.token},
								TokenOrExpression{token: tokenOrExpressionNext1.token},
								TokenOrExpression{token: tokenOrExpressionNext2.token},
							}
							idTokenOrExpr += 2 // 2 next token were used, they have to be skipped in next turns
						}
					}
				}
				if ! isOperator{
					fmt.Println("NO =:= operator detected")
				}
			}



			if ! isOperator {    //  =
				if isEqual(tokenOrExpressionActual) {
					isOperator = true
					operatorTokenElems = TokensOrExpressions{
						TokenOrExpression{token: tokenOrExpressionActual.token},
					}
				}
				if ! isOperator{
					fmt.Println("NO = operator detected")
				}
			}



			/////////////////////////////////////////////////////////////////////
			if isOperator {

				/* an operator is not a real expression, because a + sign doesn't have a value in itself.
				  the operator and the operandus, together, that is an expression.

				 but operators can be translated to functions: 1 + 2 == add(1, 2)
				 so from this perspective, ADD function can be evaluated, and it can have its own value,
				 if it is calculated.

				 Here I will mark the tokens as expressions, because an expression here means:
				 we know what can we do with those tokens, and value can be calculated.

				 and later when the expressions are interpreted, if it is not a simple value,
				 (not a function or an integer, for example) then the operator is executed, because that is a form of
				 function, too.

				 */

				tokenOrExpressionActual.elemType = tokenOrExpression_thisIsAnExpression
				tokenOrExpressionActual.expression = ErlExpression{
					ExpressionType:			expression_operator,
					TokensOrExpressions: 	operatorTokenElems,
				}
				// put back tokenOrExpression with modified elemType and expression
				tokensOrExpressionsNew = append(tokensOrExpressionsNew, tokenOrExpressionActual)

			} else {  // not an operator - put back the tokenOrExpression without any extra change/modification
				tokensOrExpressionsNew = append(tokensOrExpressionsNew, tokenOrExpressionActual)
			}
		} // FOR

	} // len > 0

	return tokensOrExpressionsNew
}

