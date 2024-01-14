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
	tokenType string

	// "quoted" string's first " is the tokens' first position!
	// 'atoms'  fist ' char is the tokens first pos!
	// so ALL character is included from the src into the token position range
	positionCharFirst int
	positionCharLast  int

	sourceFile    string
	charsInErlSrc []rune
}

func (token Token) emptyType() bool {
	return token.tokenType == ""
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

func (tokensInMap Tokens) deepCopy() Tokens {
	tokensTableUpdated := Tokens{}
	for _, token := range tokensInMap {
		tokensTableUpdated[token.positionCharFirst] = token
	}
	return tokensTableUpdated
}



/*
Receives Erlang source code - return with non-detected source code and detected Tokens.
*/
func Tokens_detect_text_blocks(erlSrc string, tokensTable Tokens) (string, Tokens){

	tokenClosingDetected__saveTheToken := "tokenClosingDetected__saveTheToken "
	tokenOpeningDetected__tokenNew := "tokenOpeningDetected__tokenNew "

	tokensTableUpdated := tokensTable.deepCopy()
	var erlSrcTokenDetectionsRemoved []rune

	noActiveTokenDetection__tokenTypeIsEmpty := func (token Token) bool {
		// if the token is emtpy, then there is no active detection
		return token.emptyType()
	}

	activeTokenDetection__typeNotEmpty := func (token Token) bool {
		// if the token is emtpy, then there is no active detection
		return ! noActiveTokenDetection__tokenTypeIsEmpty(token)
	}

	tokenNow := Token{}
	event := ""

	erlSrcRunes := []rune(erlSrc)

	for charPos, charRune := range erlSrcRunes {

		charRuneNext := ' '
		if charPos < len(erlSrcRunes) -1 {
			charRuneNext = erlSrcRunes[charPos+1]
		}

		// closers...................................................................
		if charRune == '"' && tokenNow.tokenType == tokenType_TextBlockQuotedDouble {
			event = tokenClosingDetected__saveTheToken
		}

		if charRune == '\'' && tokenNow.tokenType == tokenType_TextBlockQuotedSingle {
			event = tokenClosingDetected__saveTheToken
		}

		if charRuneNext == '\n' && tokenNow.tokenType == tokenType_Comment {
			// the endOfLine cannot be removed from original src,
			// comment is finished BEFORE the end of line
			event = tokenClosingDetected__saveTheToken
		}

		// openers...................................................................
		if noActiveTokenDetection__tokenTypeIsEmpty(tokenNow) {

			if charRune == '"' { // string
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_TextBlockQuotedDouble}
				event = tokenOpeningDetected__tokenNew
			}

			if charRune == '\'' { // quoted atom
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_TextBlockQuotedSingle}
				event = tokenOpeningDetected__tokenNew
			}

			if charRune == '%' { // comments
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_Comment}
				event = tokenOpeningDetected__tokenNew
			}
		} // not in active token detection



		/////////////////////////////////////////////////////////////////////
		if event == tokenOpeningDetected__tokenNew {
			// the opening/ending chars are removed from the original src, too
			// empty char is added instead of the original one
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, ' ')
			event = ""
			continue
		}

		/////////////////////////////////////////////////////////////////////
		if event == tokenClosingDetected__saveTheToken {
			tokenNow.positionCharLast = charPos
			tokensTableUpdated[tokenNow.positionCharFirst] = tokenNow
			tokenNow = Token{} // restore default values

			// the token closer last char is removed, too, from the original source code
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, ' ')
			event = ""
			continue
		}

		/////////////////////////////////////////////////////////////////////
		// not opening/closing event:
		if activeTokenDetection__typeNotEmpty(tokenNow) {
			// then this is an active detection, between Opening/Closing elems
			tokenNow.charsInErlSrc = append(tokenNow.charsInErlSrc, charRune)
			// if the current char is part of a token, remove if fromm src:
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, ' ')
		}

		if noActiveTokenDetection__tokenTypeIsEmpty(tokenNow) {
			// active or !active: it is clearer than an else {} block,
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, charRune)
		}

	}

	return string(erlSrcTokenDetectionsRemoved), tokensTableUpdated
}
