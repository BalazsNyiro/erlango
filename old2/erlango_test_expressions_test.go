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



// go test -run Test_expression_detection_function
// expressions are in focus here:
func Test_expression_detection_function(t *testing.T) {
	funName := "Test_expression_detection_function"
	fmt.Println(funName)

	/* https://www.erlang.org/doc/reference_manual/functions.html */

	erlSrc :=` 	fact(N) when N>0 ->
					N * fact(N-1);
				fact(0) ->
					1.
			`

	wantedExpressionDetectionTypes := "atomsAndStrings"  // the test focuses on atoms
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

	testCheck_isAtomExpression(funName + "_fact", erlExpressions[0], t)
	testCheck_isAtomExpression(funName + "_when", erlExpressions[4], t)
	testCheck_isAtomExpression(funName + "_fact", erlExpressions[12], t)
}


/*
func Test_expression_detection_list(t *testing.T) {
	funName := "Test_expression_detection_list"
	//
		Eshell V13.1.5  (abort with ^G)
		1> A = [1, 2, 3].
		[1,2,3]
	//

	erlSrc := "A = [1, 2, 3].\n"

	erlExpressions := Expression_detection_for_tests(erlSrc)

	for _, erlExpression := range erlExpressions {
		fmt.Println("TODO: test expression from string", erlExpression)
	}

	t.Fatalf("\nErr repr %s : startPos:%v  detected string representation: %v  is different from wanted representation:  %v, error",
		funName, "aaa", "bbb", "ccc")
}
*/




















func compare_expressionDetected_ExpressionWanted(callerInfo string, expressionDetected ErlExpression, expressionWanted ErlExpression, t *testing.T) {
	/*
	tokenDetected, tokenWantedIsInDetected:= tokensDetected[tokenWanted.positionFirst]

	if tokenWantedIsInDetected {
		// theoretically the charPosFirst is always ok here, because the key in map was the same position
		tokenDetected_charPosFirst, tokenDetected_charPosLast := tokenDetected.charPositionFirstLast()
		if tokenDetected_charPosFirst != tokenWanted.positionFirst {
			t.Fatalf("\nErr First: %s : detected posFirst: %v  is different from wanted posFirst:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected_charPosLast != tokenWanted.positionLast {
			t.Fatalf("\nErr Last: %s : detected posLast: %v  is different from wanted posLast:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosLast, tokenWanted.positionLast, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected.stringRepresentation() != tokenWanted.textRepresentation {
			t.Fatalf("\nErr repr %s : startPos:%v  detected string representation: %v  is different from wanted representation:  %v, error",
				callerInfo, tokenDetected_charPosFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
	} else {
		t.Fatalf("\nErr %s : wanted tokenPosFirst %v is not in detecteds - error", callerInfo, tokenWanted.positionFirst)
	}
	*/
}



