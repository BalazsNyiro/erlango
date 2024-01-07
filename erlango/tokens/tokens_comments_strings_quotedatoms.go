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

type Token struct {
	tokenName          string
	positionCharFirst  int
	positionCharLast   int
	sourceFile         string
	textRepresentation string
}

/*Tokens represent the Erlang source code - so the int-key is the first char's position in the source code*/
type Tokens map[int]Token

/*
	Token detection rules - in every step, a group of tokens are removed only, so all tokens
							are removed in MORE steps.

	input:   Actual, maybe cleaned Erlang source code, previously detected tokens.
	process: from the source code, the characters of detected tokens are relocated into Tokens
	output:  Token-less source code and updated token structure

	With this solution, different token groups can be removed in different layers - all of them
    has well-defined input state, and output.
*/

/*
Receives Erlang source code - return with non-detected source code and detected Tokens.
*/
func Tokens_detect_comments_strings_quotedatoms(erlSrc string, tokensTable Tokens) (string, Tokens){
	return erlSrc, tokensTable
}
