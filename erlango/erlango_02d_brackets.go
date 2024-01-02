/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

*/

package erlango

import (
	"fmt"
	"strings"
)


func expression_detect_brackets(tokensOrExpressionsOld TokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated string) TokensOrExpressions {
	if ! (strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "brackets") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "detectAllExpressions")) {
		return tokensOrExpressionsOld
	}


	// locally used functions only
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


	isBracketRoundOpen := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "("
	}

	isBracketRoundClose := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == ")"
	}

	isBracketSquareOpen := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "["
	}

	isBracketSquareClose := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "]"
	}

	isBracketCurlyOpen := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "{"
	}

	isBracketCurlyClose := func (tokenOrExpression TokenOrExpression) bool {
		return tokenOrExpression.token.stringRepresentation() == "}"
	}


	tokensOrExpressionsNew := TokensOrExpressions{}
	lenTokenOrExpressions :=  len(tokensOrExpressionsOld)

	bracketTokenElems := TokensOrExpressions{ }
	if lenTokenOrExpressions > 0 { // so if there is something to check

		idTokenOrExpr := -1
		for {

			// because id is manipulated inside the for loop, we can move forward more with 1,
			// if it is necessary, with more than 1 token processing
			idTokenOrExpr++
			if idTokenOrExpr >= lenTokenOrExpressions { break }


			tokenOrExpressionActual := getTokenOrExpression_fromLot(idTokenOrExpr,   tokensOrExpressionsOld)
			tokenOrExpressionNext1  := getTokenOrExpression_fromLot(idTokenOrExpr+1, tokensOrExpressionsOld)

			isBracket := false


			/////////// 2 chars long brackets //////////////////////////////////////////////////////////
			if ! isBracket {
				// binary boundary <<...>> are a type of brackets, from my persepctive.
				if 	(isLess(tokenOrExpressionActual)    && isLess(tokenOrExpressionNext1)                    ) || // <<
					(isGreater(tokenOrExpressionActual) && isGreater(tokenOrExpressionNext1)                 ) || // >>
					(isHashmark(tokenOrExpressionActual)    && isBracketCurlyOpen(tokenOrExpressionNext1)    ) { // #{

					isBracket = true
					bracketTokenElems = TokensOrExpressions{
						TokenOrExpression{token: tokenOrExpressionActual.token},
						TokenOrExpression{token: tokenOrExpressionNext1.token},
					}
					idTokenOrExpr += 1 // next token is used, it has to be skipped in next turns
				}
				if ! isBracket{
					fmt.Println("NOT 2 chars long bracket")
				}
			}


			/////////// 1 char long brackets //////////////////////////////////////////////////////////
			if ! isBracket {
				if  isBracketCurlyOpen(tokenOrExpressionActual)     || // {
				    isBracketCurlyClose(tokenOrExpressionActual)    || // }
					isBracketSquareOpen(tokenOrExpressionActual)    || // [
					isBracketSquareClose(tokenOrExpressionActual)   || // ]
					isBracketRoundOpen(tokenOrExpressionActual)     || // (
					isBracketRoundClose(tokenOrExpressionActual)    {  // )

					isBracket = true
					bracketTokenElems = TokensOrExpressions{
						TokenOrExpression{token: tokenOrExpressionActual.token},
					}
				}
				if ! isBracket{
					fmt.Println("NOT 1 char long bracket")
				}
			}

			/////////////////////////////////////////////////////////////////////
			if isBracket {

				tokenOrExpressionActual.elemType = tokenOrExpression_thisIsAnExpression
				tokenOrExpressionActual.expression = ErlExpression {
					ExpressionType:			expression_brackets,
					TokensOrExpressions: 	bracketTokenElems,
				}
				// put back tokenOrExpression with modified elemType and expression
				tokensOrExpressionsNew = append(tokensOrExpressionsNew, tokenOrExpressionActual)

			} else {  // not a bracket - put back the tokenOrExpression without any extra change/modification
				tokensOrExpressionsNew = append(tokensOrExpressionsNew, tokenOrExpressionActual)
			}
		} // FOR

	} // len > 0

	return tokensOrExpressionsNew
}

