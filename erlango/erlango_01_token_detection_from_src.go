/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango

import "fmt"

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

func step_01_tokens_from_source_code(sourcesTokensExecutables_list SourcesTokensExecutables_list, fileNamePaths []string) SourcesTokensExecutables_list {
	fmt.Println("filenames to detect tokens", fileNamePaths)
	return sourcesTokensExecutables_list
}