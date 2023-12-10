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


func expression_detect_variable_names(tokensOrExpressionsOld TokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated string) TokensOrExpressions {
	if ! (strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "variableNames") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "detectAllExpressions")) {
		return tokensOrExpressionsOld
	}

	tokensOrExpressionsNew_variableNamesDetected := TokensOrExpressions{}

	for _, tokenOrExpression := range(tokensOrExpressionsOld) {
		fmt.Println("detect variableNames - token expression", tokenOrExpression)

		if tokenOrExpression.isExpression() {  // if it is a previously detected expression, there is nothing to do
			tokensOrExpressionsNew_variableNamesDetected = append(tokensOrExpressionsNew_variableNamesDetected, tokenOrExpression)
			continue
		}

		isVariableName := false

		if tokenOrExpression.token.TokenType == tokenTypeAbcFullWith_Underscore_At_numbers {
			if tokenOrExpression.token.charFirstRuneValIsVariableStarter() {
				isVariableName = true
			}
		}

		if isVariableName{
			tokenOrExpression.elemType = tokenOrExpression_thisIsAnExpression
			tokenOrExpression.expression = ErlExpression{
				ExpressionType:        expression_variableName,
				SimpleTokenValue:      tokenOrExpression.token,
			}
			// put back tokenOrExpression with modified elemType and expression
			tokensOrExpressionsNew_variableNamesDetected = append(tokensOrExpressionsNew_variableNamesDetected, tokenOrExpression)

		} else {  // expression is not detected now - put back the tokenOrExpression without any extra change/modification
			tokensOrExpressionsNew_variableNamesDetected = append(tokensOrExpressionsNew_variableNamesDetected, tokenOrExpression)
		}
	} // FOR

	return tokensOrExpressionsNew_variableNamesDetected

}

