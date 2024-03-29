/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package old2

import (
	"fmt"
	"testing"
)

// this fun cannot be in a test file, because in that case it is unavailable from other test file

func Expression_detection_for_tests(erlSrc string, wantedExpressionDetectionTypesCommaSeparated string) ErlExpressions{

	// sourcesTokensExecutables_all can be empty (like here), or it can have existing elements - in a running system new src can be loaded, next to the existing ones
	sourcesTokensExecutables_all := SourcesTokensExecutables_map{}
	sourcesTokensExecutables_all = step_01_tokens_from_passed_source_codes_without_files(erlSrc, sourcesTokensExecutables_all)

	fileNamesOfErlangSources := []string{erlSrc} // if a source code doesn't have source file, the identifier is himself
	sourcesTokensExecutables_all = step_02_expressions_from_tokens_from_lot_of_sources(sourcesTokensExecutables_all, fileNamesOfErlangSources, wantedExpressionDetectionTypesCommaSeparated )

	fmt.Println("num of expressions:", len(sourcesTokensExecutables_all[erlSrc].Expressions))
	return sourcesTokensExecutables_all[erlSrc].Expressions
}

func testCheck_isAtomExpression(testName string, erlExpression ErlExpression, t *testing.T) {
	testCheck_compareExpressionWithWantedType(testName, erlExpression, expression_atom, t)
}

func testCheck_isNumberExpression(testName string, erlExpression ErlExpression, t *testing.T) {
	testCheck_compareExpressionWithWantedType(testName, erlExpression, expression_num, t)
}

func testCheck_isVariableNameExpression(testName string, erlExpression ErlExpression, t *testing.T) {
	testCheck_compareExpressionWithWantedType(testName, erlExpression, expression_variableName, t)
}

// receives list of types and expressions, and compares all expressions with all types
func testCheck_expressions_all_Wanteds(callerTestFunName string, erlExpressions ErlExpressions, typesWanted []int, t *testing.T) {

	for i, erlExpression := range(erlExpressions) {
		if i >= len(typesWanted){
			continue // if the passed typesWanted smaller, check only the available wanted types
		}
		expressionTypeWanted := typesWanted[i]
		testCheck_compareExpressionWithWantedType(callerTestFunName, erlExpression, expressionTypeWanted, t)
	}
}

func testCheck_compareExpressionWithWantedType(testName string, erlExpression ErlExpression, typeWanted int, t *testing.T) {
	if erlExpression.ExpressionType !=  typeWanted {
		t.Fatalf("\nError (%s): incorrect expression type: %s, wanted type: %v %s", testName, erlExpression.expressionTypeForHuman(), typeWanted, ExpressionName_from_num[typeWanted])
	}
}
