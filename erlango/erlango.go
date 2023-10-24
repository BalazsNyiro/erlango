/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango



// ############# PARSER ELEMS #############################
// use minimal set of types. Don't overcomplicate the parser.
type ErlToken struct {

	// don't create a pointer to the prev/next Token.
	// it is a dangerous error source. Use only a positionId in a collector of ErlTokens

	/* Possible token types:
	 - notDetected
	 - atom
	 - string
	 - binaryDoubleOpener      <<
	 - binaryDoubleCloser      >>
	 - operatorMathPlusPlus    ++
	 - operatorMathMinusMinus  --
	 - operatorBoolEqualDouble ==
	 - operatorBoolBigger      >
	 - operatorBoolSmaller     <
	 - codeBlockStart         ->

	a few sign can have different meanings in different places:
	a comma can be expression separator, or parameter separator - so it's later decided how is it handled
	 - langComma ,
	 - langDot   .
	 - langColon :
	 - langBraceRoundOpen   (
	 - langBraceRoundClose  )
	 - langBraceSquaredOpen [
	 - langBraceSqaredClose ]
	 - langVariableBound    =

	*/
	tokenType  string

	// one char can have a meaning alone, for example: (
	// but sometime more than one few char can form a token, for example: '->' which has his own meaning, but represented by 2 chars
	PosFirstCharInFile int
	PosLastCharInFile int

	stringRepresentation      string   // in a source code we typically use short words. So instead of an optimization (list of chars) a simple string is used here
	SourceFilePath string
}

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
