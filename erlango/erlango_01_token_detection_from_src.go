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


// sourcesTokensExecutables_all can be empty, or it can have existing elements - and maybe only newer ones are added.
func step_01_tokens_from_source_code_of_files(  // in program plan
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

func step_01_tokens_from_passed_source_codes_without_files(erlSrc string, sourcesTokensExecutables_all SourcesTokensExecutables_map) SourcesTokensExecutables_map { // in program plan
	// detect tokens from a passed string src, without file path
	returnFromTokenDetection := make(chan SourceTokensExecutables)
	go step_01a_tokens_detect(erlSrc, "", returnFromTokenDetection, true)
	sourceTokensExecutables := <- returnFromTokenDetection
	sourcesTokensExecutables_all[sourceTokensExecutables.WhereTheCodeIsStored] = sourceTokensExecutables
	return sourcesTokensExecutables_all
}

func step_01a_tokens_detect(erlangSource string, filePath string, parentChannel chan SourceTokensExecutables, verboseForErlangoInvestigations__useFalseInProdEnv bool) { // in program plan
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
		// charsFromErlFile and tokensDetected are updated during token detection
		charsFromErlFile, tokensDetected, errors = token_detect_comments_textblocks_alphanums_whitespaces_literals(charsFromErlFile, tokensDetected, verboseForErlangoInvestigations__useFalseInProdEnv)

	} else {
		// FIXME: what to do if file_read_runes has a problem?
		// FIXME: create a log directory, and a default log file where errors can be collected
		// FIXME: logs_errors_file_reading_problems.txt ?
	}

	// The most important elems are: Chars|Tokens|Expressions.
	// From Chars, Tokens (building elements) are detected.
	// Expressions are empty here, later from Tokens->Expressions are detected
	sourceTokensExecutables := SourceTokensExecutables{
		WhereTheCodeIsStored:        whereCodeIsStored,
		ModuleVersion:               "not-detected-version",
		CharsFromErlFile:            charsFromErlFile,
		Tokens:                      tokensDetected,
		Errors:                      errors,
		ErlangSourceWithoutFilePath: erlangSourceWithoutFilePath,
		Expressions: ErlExpressions{}, // at this point the expressions are not detected, so it is empty
	}

	parentChannel <- sourceTokensExecutables
}

/////////////////////////////////////////////////////////////////////////
