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
	"testing"
)

//  go test -v -run Test_numbers_integers
func Test_numbers_integers(t *testing.T) {
	funName := "Test_numbers_integers"
	fmt.Println(funName)


	erlSrc :=`Atom = 'atom', String = "str", Int = 1, Float = 1.1.`

	wantedExpressionDetectionTypes := "atomsQuoted,atomsSimple,numbers"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

}
