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


func tokens_in_file(fileName string, parentChannel chan SourceTokensExecutables) {
	fmt.Println("Tokens from file:", fileName)
	sourceTokensExecutables := SourceTokensExecutables{}

	parentChannel <- sourceTokensExecutables
}

func step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_list SourcesTokensExecutables_list, fileNamePaths []string) SourcesTokensExecutables_list {
	fmt.Println("filenames to detect tokens", fileNamePaths)

	SourceTokensExecutables__list := []SourceTokensExecutables{}

	returnFromTokenDetection := make(chan SourceTokensExecutables)

	for _, fileName := range(fileNamePaths) {
		go tokens_in_file(fileName, returnFromTokenDetection)
	}

	for len(SourceTokensExecutables__list) <  len(fileNamePaths) {
		sourceTokensExecutables := <- returnFromTokenDetection
		fmt.Println("Token detection returned structure:", sourceTokensExecutables)
		SourceTokensExecutables__list = append(SourceTokensExecutables__list, sourceTokensExecutables)
	}

	return sourcesTokensExecutables_list
}