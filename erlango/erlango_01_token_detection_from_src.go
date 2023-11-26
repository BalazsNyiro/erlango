/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite


Token: the smallest block/region of characthers that has his own meaning.

a simple case: 123  is an integer num, and it is one token, because it has one meaning. (if we accept that this is a decimal num)

but  2.3e3 or 16#1f are numbers too - and they have more tokens.
in case of '2.3e3' the '2', the '.', the '3' and 'e'+'3' has their own meaning,
and if we interpret this expression, the calculated result give back the real number
(so this is a num representation, with more tokens: the 'e'+'3' has a math meaning,
the '2','.',3' has it's own meaning, too)


*/

package erlango

import (
	"fmt"
)


func step_01_tokens_from_source_code_of_files(sourcesTokensExecutables_all SourcesTokensExecutables_map, fileNamePaths []string, verboseForErlangoInvestigations__useFalseInProdEnv bool) SourcesTokensExecutables_map {
	// parallel token detection from erl sources
	fmt.Println("filenames to detect tokens", fileNamePaths)

	returnFromTokenDetection := make(chan SourceTokensExecutables)

	for _, fileName := range(fileNamePaths) {
		go step_01a_tokens_detect_in_file(fileName, returnFromTokenDetection, verboseForErlangoInvestigations__useFalseInProdEnv)
	}

	for len(sourcesTokensExecutables_all) < len(fileNamePaths) {
		sourceTokensExecutables := <- returnFromTokenDetection
		// fmt.Println("Token detection returned structure:", sourceTokensExecutables)
		sourcesTokensExecutables_all[sourceTokensExecutables.PathErlFile] = sourceTokensExecutables
	}

	return sourcesTokensExecutables_all // it can have errors, too!
}

func step_01a_tokens_detect_in_file(filePath string, parentChannel chan SourceTokensExecutables, verboseForErlangoInvestigations__useFalseInProdEnv bool) {
	fmt.Println("Tokens from file:", filePath)
	funName := "step_01a_tokens_detect_in_file"

	newLineSeparator := '\n'

	runes, errFileReadingRunes := file_read_runes(filePath, funName)

	tokensDetected := ErlTokens{}
	charsFromErlFile := Chars{}
	errors := errorsDetected{}

	if errFileReadingRunes == nil {

		// ##### step A: read all chars from Erlang source #########
		lineNum := 0
		positionInLine := 0

		for posInFile, runeInFile := range(runes) {
			fmt.Println(posInFile, "rune in file:", string(runeInFile), runeInFile)
			charNow := Char{PositionInFile: posInFile, Value: runeInFile, FilePath: filePath, LineNum: lineNum, PositionInLine: positionInLine}
			charsFromErlFile = append(charsFromErlFile, charNow)

			positionInLine += 1

			if runeInFile == newLineSeparator {
				lineNum += 1
				positionInLine = 0
			}
		}

		// ##### step B: Tokens detect ########################
		charsFromErlFile, tokensDetected, errors = token_detect_comments_textblocks_alphanums(charsFromErlFile, tokensDetected, verboseForErlangoInvestigations__useFalseInProdEnv)

	} else {
		// FIXME: what to do if file_read_runes has a problem?
	}

	sourceTokensExecutables := SourceTokensExecutables{
		PathErlFile: filePath,
		ModuleVersion: "not-detected-version",
		CharsFromErlFile: charsFromErlFile,
		Tokens: tokensDetected,
		Errors: errors,
	}

	parentChannel <- sourceTokensExecutables
}

func char_txt_value_get(pos int, chars Chars) string {
	ret := "" 	// I would like to handle empty values, too, so runes cannot be given back.
	// empty value means: there is no real character in the wanted position
	// the position has a real value only if it is in the valid range
	if pos >= 0 && pos < len(chars) {
		ret = string(chars[pos].Value)
	}
	return ret
}

func token_empty_obj(tokenType string, tokenId int) ErlToken {
	return ErlToken{ TokenType: tokenType, TokenId: tokenId, SourceCodeChars: Chars{}, }
}

/////////////////////////////////////////////////////////////////////////
