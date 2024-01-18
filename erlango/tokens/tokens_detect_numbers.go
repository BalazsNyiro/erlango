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
Receives Erlang source code - return with non-detected source code and detected Tokens.
*/
func Tokens_detect_numbers(erlSrc string, tokensTable Tokens) (string, Tokens) {

	tokensTableUpdated := tokensTable.deepCopy()
	var erlSrcTokenDetectionsRemoved []rune

	return string(erlSrcTokenDetectionsRemoved), tokensTableUpdated
}
