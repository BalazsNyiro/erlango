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

func isCharEscaped(posChar int, chars []Char) bool {
	// char is escaped if there are 'odd' num of escape char before that.
	escaped := false
	escapeCharCounter := 0

	posTestedChar := posChar - 1
	for posTestedChar >= 0 {
		if chars[posTestedChar].Value == '\\' {
			escapeCharCounter += 1
			posTestedChar -= 1
		} else { // if the char is not a backslash, leave the loop
			break
		}
	}
	escaped = (escapeCharCounter % 2) == 1 // odd escape chars are before the current char
	return escaped
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

func token_empty(tokenType string, tokenId int) ErlToken {
	return ErlToken{ TokenType: tokenType, TokenId: tokenId, SourceCodeChars: []Char{}, }
}

func token_detect_comments_textblocks(chars []Char, tokens []ErlToken) ([]Char, []ErlToken){
	// the "wrapper" quotes around the string values or 'atoms' are the part of the tokens,
	// they are necessary to define a text block (single or double qoted texts)
	// but not part of the value of the token

	/* Erlang accepts newlines in atoms:

		Erlang/OTP 25 [erts-13.1.5] [source] [64-bit] [smp:4:4] [ds:4:4:10] [async-threads:1] [jit:ns]

		Eshell V13.1.5  (abort with ^G)
		1> A = 'atom\n2'.
		'atom\n2'
		2> A.
		'atom\n2'
		3> A2 = 'atom\\n'.
		'atom\\n'
		4>

		discussion: https://erlang.org/pipermail/erlang-questions/2014-February/077922.html

	*/

	fmt.Println("token detext comments, quoted textblocks")
	inBlock := ""

	tokenActual := ErlToken{}

	for charPos := 0; charPos < len(chars); charPos += 1 {
		tokenActualId := len(tokens) // len(..) is always represent the next free, unused elem Id in the slice

		//charTxtPrev1 := charTxtGet(charPos-1, chars)
		charTxtNow := charTxtGet(charPos, chars)
		//charTxtNext1 := charTxtGet(charPos+1, chars)

		blockStarted := false
		blockFinished := false

		////////////// quoted text detect //////////
		if charTxtNow == "\""{

			if inBlock == "" {
				tokenActual = token_empty("tokenTextBlockQuotedDouble", tokenActualId)
				inBlock = "inTextBlockQuotedDouble"
				blockStarted = true
			}

			if inBlock == "inTextBlockQuotedDouble" && ! blockStarted {
				if ! isCharEscaped(charPos, chars) {
					blockFinished = true
				} // char is not escaped
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
