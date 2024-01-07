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

import "fmt"

type Token struct {
	tokenType string

	// "quoted" string's first " is the tokens' first position!
	// 'atoms'  fist ' char is the tokens first pos!
	// so ALL character is included from the src into the token position range
	positionCharFirst int
	positionCharLast  int

	sourceFile    string
	charsInErlSrc []rune
}


const tokenType_TextBlockQuotedSingle = "tokenTextBlockQuotedSingle"
const tokenType_Comment = "tokenComment"
const tokenType_TextBlockQuotedDouble = "tokenTextBlockQuotedDouble"




/*Tokens represent the Erlang source code - so the int-key is the first char's position in the source code*/
type Tokens map[int]Token

/*
	Token detection rules - in every step, a group of tokens are removed only, so all tokens
							are removed in MORE steps.

	At the beginning, the input is the original Erlang source code, and tokensTable is empty
	in every step, Erlang source code has more and more empty sections, as charsInErlSrc are detected,
    and the token table has more and more elems with detected tokens

	input:   Actual, maybe cleaned Erlang source code, previously detected tokens.
	process: from the source code, the characters of detected tokens are relocated into Tokens
	output:  Token-less source code and updated token structure

	With this solution, different token groups can be removed in different layers - all of them
    has well-defined input state, and output.
*/

/*
Receives Erlang source code - return with non-detected source code and detected Tokens.
*/
func Tokens_detect_text_blocks(erlSrc string, tokensTable Tokens) (string, Tokens){
	endOfLine := '\n'

	tokenNow := Token{}
	saveToken := false

	for charPos, charRune := range erlSrc {


		fmt.Println("pos:", charPos, charRune)

		if charRune == '"' { // strings ////////////////////////////////////////
			if tokenNow.tokenType == "" {  // if there is NO active detection
				tokenNow = Token{positionCharFirst: charPos, tokenType: tokenType_TextBlockQuotedDouble}
				continue
			}
			if tokenNow.tokenType == tokenType_TextBlockQuotedDouble {
				saveToken = true
			}
		}

		if charRune == '\'' { // quoted atoms /////////////////////////////////
			if tokenNow.tokenType == "" {  // if there is NO active detection
				tokenNow = Token{positionCharFirst: charPos, tokenType: tokenType_TextBlockQuotedSingle}
				continue
			}
			if tokenNow.tokenType == tokenType_TextBlockQuotedSingle{
				saveToken = true
			}
		}

		if charRune == '%' { // comments  /////////////////////////////////////
			if tokenNow.tokenType == "" { // if there is NO active detection
				tokenNow = Token{positionCharFirst: charPos, tokenType: tokenType_Comment}
				continue
			}
		}
		if charRune == endOfLine && tokenNow.tokenType == tokenType_Comment {
			saveToken = true
		}

		/////////////////////////////////////////////////////////////////////
		if ! saveToken {
			if tokenNow.tokenType != "" { // so if there is an active detection
				tokenNow.charsInErlSrc = append(tokenNow.charsInErlSrc, charRune)
			}
		}

		if saveToken { // for me, it is better han an 'else' because it is readable
			tokenNow.positionCharLast = charPos
			tokensTable[tokenNow.positionCharFirst] = tokenNow

			// restore default values
			tokenNow = Token{}
			saveToken = false
		}

	}
	return erlSrc, tokensTable
}
