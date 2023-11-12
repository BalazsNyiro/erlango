/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango

import (
	"fmt"
)

func TokensDetectFromFile() {
// read a file
// return with TokensDetectFromString() result
}

func TokensDetectFromString() {
// convert String to ErlTokens
// return TokensDetect()
}
func TokensDetect(erlSrc []ErlToken) {
//
}
// ############# PARSER ELEMS #############################

func tokens_in_file(fileName string, parentChannel chan int) {
	fmt.Println("Tokens from file:", fileName)
	retCode := 0 // there was no error during token detecton

	// TODO: give back the Tokens, too
	parentChannel <- retCode
}

func step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_list SourcesTokensExecutables_list, fileNamePaths []string) SourcesTokensExecutables_list {
	fmt.Println("filenames to detect tokens", fileNamePaths)
	goRoutineStarted := 0
	goRoutineFinished := 0

	// 0: detection was fine - bigger than 0: detection had an error
	returnCodesTokeDetectionFinished := make(chan int)

	for _, fileName := range(fileNamePaths) {
		go tokens_in_file(fileName, returnCodesTokeDetectionFinished)
		goRoutineStarted += 1
	}

	for goRoutineFinished < goRoutineStarted {
		retCodeMsg := <- returnCodesTokeDetectionFinished
		fmt.Println("Token detection ret code:", retCodeMsg)
		goRoutineFinished += 1
	}

	return sourcesTokensExecutables_list
}