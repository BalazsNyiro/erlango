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

	tokensOrExpressionsNew_01_numsDetected := TokensOrExpressions{}


	for _, tokenOrExpression := range(tokensOrExpressionsOld) {
		fmt.Println("detect numbers - token expression", tokenOrExpression)

		if tokenOrExpression.isExpression() {  // if it is a previously detected expression, there is nothing to do
			tokensOrExpressionsNew_01_numsDetected= append(tokensOrExpressionsNew_01_numsDetected, tokenOrExpression)
			continue
		}


		isNum := false
		//// isNum? /////////////////////////////////////
		// https://www.erlang.org/doc/reference_manual/data_types.html

		/*
			Number variations - the + or - unary operators are NOT detected here
			12
			12_34


			$A      1 char after $
			$\n     2 char after $

			2#101   base#value, integer result
			16#1f   base#value, characters can be interpreted as num elems (f)
			16#1F	base#value, CAPITAL chars are interpreted, too

			12.34
			12_34.56
			12_34.56_78

			2.3e3
			2.3e+3
			2.3e-3
			2_3.4e+3
			2_3.4e+3_0   result: 2.34e31

		*/














		if tokenOrExpression.token.TokenType == "tokenTextBlockQuotedSingle" {
			isNum = true
		}

		if tokenOrExpression.token.TokenType == "tokenAbcFullWith_At_numbers" {
			if tokenOrExpression.token.charFirstRuneValIsSmallCapsAtomStarter() {
				isNum = true
			}
		}





		//// isNum? /////////////////////////////////////

		if isNum {
			tokenOrExpression.elemType = tokenOrExpression_thisIsAnExpression
			tokenOrExpression.expression = ErlExpression{
				ExpressionType:			expression_num,
				SimpleTokenValue:      tokenOrExpression.token,
			}
			// put back tokenOrExpression with modified elemType and expression
			tokensOrExpressionsNew_01_numsDetected = append(tokensOrExpressionsNew_01_numsDetected, tokenOrExpression)

		} else {  // not an atom - put back the tokenOrExpression without any extra change/modification
			tokensOrExpressionsNew_01_numsDetected = append(tokensOrExpressionsNew_01_numsDetected, tokenOrExpression)
		}
	} // FOR

	return tokensOrExpressionsNew_01_numsDetected
}

