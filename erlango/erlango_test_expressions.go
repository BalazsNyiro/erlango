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



// go test -run Test_expression_detection
// expressions are in focus here:
func Test_expression_detection_simple_1(t *testing.T) {

	/*
	Eshell V13.1.5  (abort with ^G)
	1> A = [1, 2, 3].
	[1,2,3]
	*/
	erlSrc := "A = [1, 2, 3].\n"

	sourcesTokensExecutables_all := Experssion_detection_for_tests(erlSrc)
	fmt.Println("TODO: test expressions from string", sourcesTokensExecutables_all)
}

func Experssion_detection_for_tests(erlSrc string) SourcesTokensExecutables_map {

	// sourcesTokensExecutables_all can be empty (like here), or it can have existing elements - in a running system new src can be loaded, next to the existing ones
	sourcesTokensExecutables_all := SourcesTokensExecutables_map{}

	sourcesTokensExecutables_all = step_01_tokens_from_passed_source_codes_without_files(erlSrc, sourcesTokensExecutables_all)

	fileNamesOfErlangSources := []string{erlSrc} // if a source code doesn't have source file, the identifier is himself
	step_02_expressions_from_tokens(sourcesTokensExecutables_all, fileNamesOfErlangSources)

	return sourcesTokensExecutables_all
}

