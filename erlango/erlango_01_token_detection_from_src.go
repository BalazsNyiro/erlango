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


// sourcesTokensExecutables_all can be empty, or it can have existing elements - and maybe only newer ones are added.
func step_01_tokens_from_source_code_of_files(
		sourcesTokensExecutables_all SourcesTokensExecutables_map,
		fileNamePaths []string,
		verboseForErlangoInvestigations__useFalseInProdEnv bool) SourcesTokensExecutables_map {

	// parallel token detection from erl sources
	fmt.Println("filenames to detect tokens", fileNamePaths)

	returnFromTokenDetection := make(chan SourceTokensExecutables)

	for _, fileName := range(fileNamePaths) {
		go step_01a_tokens_detect("", fileName, returnFromTokenDetection, verboseForErlangoInvestigations__useFalseInProdEnv)
	}

	numOfReceivedReply := 0
	for numOfReceivedReply < len(fileNamePaths) {
		sourceTokensExecutables := <- returnFromTokenDetection
		// fmt.Println("Token detection returned structure:", sourceTokensExecutables)
		sourcesTokensExecutables_all[sourceTokensExecutables.WhereTheCodeIsStored] = sourceTokensExecutables
		numOfReceivedReply += 1
	}

	return sourcesTokensExecutables_all // it can have errors, too!
}

func step_01_tokens_from_passed_source_codes_without_files(erlSrc string, sourcesTokensExecutables_all SourcesTokensExecutables_map) SourcesTokensExecutables_map {
	// detect tokens from a passed string src, without file path
	returnFromTokenDetection := make(chan SourceTokensExecutables)
	go step_01a_tokens_detect(erlSrc, "", returnFromTokenDetection, true)
	sourceTokensExecutables := <- returnFromTokenDetection
	sourcesTokensExecutables_all[sourceTokensExecutables.WhereTheCodeIsStored] = sourceTokensExecutables
	return sourcesTokensExecutables_all
}

func step_01a_tokens_detect(erlangSource string, filePath string, parentChannel chan SourceTokensExecutables, verboseForErlangoInvestigations__useFalseInProdEnv bool) {
	/*
		if erlangSource is empty, source is read from filePath
	*/
	fmt.Println("Tokens detect filePath:", filePath, "passedErlangSourceLength:", len(erlangSource))
	funName := "step_01a_tokens_detect"

	newLineSeparator := '\n'

	whereCodeIsStored := filePath
	runes := []rune(erlangSource)
	var errFileReadingRunes error = nil
	erlangSourceWithoutFilePath :=  erlangSource != ""

	if erlangSourceWithoutFilePath {
		// if erlangSource is passed without a filename, the passed source code can be used as ID of himself.
		// the tokens and expressions are bind to a file where they can be found, but in this situation
		// the tokens can be connected to the passed nameless source code
		// filePath has to be unique, because source code blocks are stored based on their source

		// why is it a correct theory? Tokens -> and later Expressions are skeletons only.
		// when the code is executed with different inputs, the skeleton will be filled with variables and inputs,
		// so the executed code has the same logic - so same source code will be represented with same
		// expression structure.
		whereCodeIsStored = erlangSource
	} else {
		runes, errFileReadingRunes = file_read_runes(filePath, funName)
	}

	tokensDetected := ErlTokens{}
	charsFromErlFile := Chars{}
	errors := errorsDetected{}

	if errFileReadingRunes == nil {

		// ##### step A: read all chars from Erlang source #########
		lineNum := 0
		positionInLine := 0

		for posInFile, runeInFile := range(runes) {
			fmt.Println(posInFile, "rune in file:", string(runeInFile), runeInFile)
			charNow := Char{	PositionInFile:       posInFile,
								Value:                runeInFile,
								WhereTheCharIsStored: whereCodeIsStored,
								LineNum:              lineNum,
								PositionInLine:       positionInLine,
								ErlangSourceWithoutFilePath: erlangSourceWithoutFilePath,
			}
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
		WhereTheCodeIsStored:        whereCodeIsStored,
		ModuleVersion:               "not-detected-version",
		CharsFromErlFile:            charsFromErlFile,
		Tokens:                      tokensDetected,
		Errors:                      errors,
		ErlangSourceWithoutFilePath: erlangSourceWithoutFilePath,
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
	return ErlToken{ TokenType: tokenType, TokenId: tokenId, SourceCodeChars: Chars{}, TokenIsDetectedAsPartOfExpression: false}
}

/////////////////////////////////////////////////////////////////////////
