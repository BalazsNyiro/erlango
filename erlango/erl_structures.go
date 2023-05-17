/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

const LanguageElemTypeBlockList string = "LanguageElemTypeBlockList "
const LanguageElemTypeBlockParenthesesRound string = "LanguageElemTypeBlockParenthesesRound "
const LanguageElemTypeBlockParenthesesSquare string = "LanguageElemTypeBlockParenthesesSquare "

/*
A language elem is one thing with one meaning: a list, a map, a case structure.

A language elem typically built from a lot of tokens: opening/closing (...) pairs,

	and it can have a lot of internal elems
*/
type LanguageElem struct {
	Type                     string
	IncludedLanguageElems    []LanguageElem
	OneErlTokenIfNoLangElems ErlSrcToken
}

func embeddedStructuresDetectFromFlatChars(chars []ErlSrcChar) {
}
