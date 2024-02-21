/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.3, third total rewrite
*/

package tokens


/*
the ret value is always this, in a Token-detect step:
 - the cleaned src, where the detected elems are removed,
 - the tokens table where the detected elems are inserted as Tokens

*/
func Tokens_2_detect_atoms_variableNames(erlSrc string, tokensTable Tokens) (string, Tokens) {
	tokensTableUpdated := tokensTable.deepCopy()
	var erlSrcTokenDetectionsRemoved []rune
	///////////////////////////////////////////////////////////////



	///////////////////////////////////////////////////////////////
	return string(erlSrcTokenDetectionsRemoved), tokensTableUpdated
}