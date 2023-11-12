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

func step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_list SourcesTokensExecutables_list, fileNamePaths []string) SourcesTokensExecutables_list {
	// parallel token detection from erl sources
	fmt.Println("filenames to detect tokens", fileNamePaths)

	SourceTokensExecutables__list := []SourceTokensExecutables{}

	returnFromTokenDetection := make(chan SourceTokensExecutables)

	for _, fileName := range(fileNamePaths) {
		go step_01a_tokens_detect_in_file(fileName, returnFromTokenDetection)
	}

	for len(SourceTokensExecutables__list) <  len(fileNamePaths) {
		sourceTokensExecutables := <- returnFromTokenDetection
		// fmt.Println("Token detection returned structure:", sourceTokensExecutables)
		SourceTokensExecutables__list = append(SourceTokensExecutables__list, sourceTokensExecutables)
	}

	return sourcesTokensExecutables_list
}


func step_01a_tokens_detect_in_file(filePath string, parentChannel chan SourceTokensExecutables) {
	fmt.Println("Tokens from file:", filePath)
	funName := "step_01a_tokens_detect_in_file"

	runes, errFileReadingRunes := file_read_runes(filePath, funName)

	TokensDetected := []ErlToken{}
	charsFromErlFile := []Char{}
	if errFileReadingRunes == nil {

		// ##### step A: read all chars from Erlang source #########
		for posInFile, runeInFile := range(runes) {
			fmt.Println("rune in file:", string(runeInFile), runeInFile)
			charNow := Char{PositionInFile: posInFile, Value: runeInFile, FilePath: filePath}
			charsFromErlFile = append(charsFromErlFile, charNow)
		}

		// ##### step B: Tokens detect ########################


	} else {
		// FIXME: what to do if file_read_runes has a problem?
	}

	sourceTokensExecutables := SourceTokensExecutables{
		PathErlFile: filePath,
		ModuleVersion: "not-detected-version",
		CharsFromErlFile: charsFromErlFile,
		Tokens: TokensDetected,
	}

	parentChannel <- sourceTokensExecutables
}