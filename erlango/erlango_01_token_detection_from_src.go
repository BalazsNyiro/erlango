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

	returnFromTokenDetection := make(chan SourceTokensExecutables)

	for _, fileName := range(fileNamePaths) {
		go step_01a_tokens_detect_in_file(fileName, returnFromTokenDetection)
	}

	for len(sourcesTokensExecutables_list) < len(fileNamePaths) {
		sourceTokensExecutables := <- returnFromTokenDetection
		// fmt.Println("Token detection returned structure:", sourceTokensExecutables)
		sourcesTokensExecutables_list = append(sourcesTokensExecutables_list, sourceTokensExecutables)
	}

	return sourcesTokensExecutables_list
}

func step_01a_tokens_detect_in_file(filePath string, parentChannel chan SourceTokensExecutables) {
	fmt.Println("Tokens from file:", filePath)
	funName := "step_01a_tokens_detect_in_file"

	runes, errFileReadingRunes := file_read_runes(filePath, funName)

	tokensDetected := []ErlToken{}
	charsFromErlFile := []Char{}
	if errFileReadingRunes == nil {

		// ##### step A: read all chars from Erlang source #########
		for posInFile, runeInFile := range(runes) {
			fmt.Println(posInFile, "rune in file:", string(runeInFile), runeInFile)
			charNow := Char{PositionInFile: posInFile, Value: runeInFile, FilePath: filePath}
			charsFromErlFile = append(charsFromErlFile, charNow)
		}

		// ##### step B: Tokens detect ########################
		charsFromErlFile, tokensDetected = token_detect_comments_textblocks(charsFromErlFile, tokensDetected)

	} else {
		// FIXME: what to do if file_read_runes has a problem?
	}

	sourceTokensExecutables := SourceTokensExecutables{
		PathErlFile: filePath,
		ModuleVersion: "not-detected-version",
		CharsFromErlFile: charsFromErlFile,
		Tokens: tokensDetected,
	}

	parentChannel <- sourceTokensExecutables
}


/////////////////////////////////////////////////////////////////////////
func charTxtGet(pos int, chars []Char) string {
	ret := "" 	// I would like to handle empty values, too, so runes cannot be given back.
	// empty value means: there is no real character in the wanted position
	// the position has a real value only if it is in the valid range
	if pos >= 0 && pos < len(chars) {
		ret = string(chars[pos].Value)
	}
	return ret
}

func token_detect_comments_textblocks(chars []Char, tokens []ErlToken) ([]Char, []ErlToken){
	fmt.Println("token detext comments, quoted textblocks")
	inBlock := ""

	tokenActual := ErlToken{}

	for charPos := 0; charPos < len(chars); charPos += 1 {
		tokenActualId := len(tokens) // len(..) is always represent the next free, unused elem Id in the slice

		//charTxtPrev2 := charTxtGet(charPos-2, chars) // the often used relative chars are collected here
		//charTxtPrev1 := charTxtGet(charPos-1, chars)
		charTxtNow := charTxtGet(charPos, chars)
		//charTxtNext1 := charTxtGet(charPos+1, chars)
		//charTxtNext2 := charTxtGet(charPos+2, chars)

		blockStarted := false
		blockFinished := false

		if charTxtNow == "\""{

			if inBlock == "" {

				tokenActual = ErlToken{
					TokenType: "tokenTextBlockQuotedDouble",
					TokenId: tokenActualId,
					SourceCodeChars: []Char{},
				}
				inBlock = "inTextBlockQuotedDouble"
				blockStarted = true
			}

			if inBlock == "inTextBlockQuotedDouble" && ! blockStarted {
				blockFinished = true
			}
		}

		if inBlock != "" {
			tokenActual.SourceCodeChars = append(tokenActual.SourceCodeChars, chars[charPos])
		}

		if blockFinished {
			inBlock = ""
			tokens = append(tokens, tokenActual)
		}
	}

	return chars, tokens
}
/////////////////////////////////////////////////////////////////////////
