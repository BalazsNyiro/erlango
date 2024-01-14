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
func Tokens_detect_text_blocks(erlSrc string, tokensTable Tokens) (string, Tokens){

	tokenCloserDetected__saveTheToken := "tokenCloserDetected__saveTheToken "
	tokenOpenerDetected__tokenNew := "tokenOpenerDetected__tokenNew "

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

		charRuneNext1, _ := charRuneNext(charPos, +1, erlSrcRunes)

		// closers.......... (before openers, to avoid tokenType set side effect)....
		if charRune == '"' && tokenNow.tokenType == tokenType_TextBlockQuotedDouble {
			event = tokenCloserDetected__saveTheToken
		}

		if charRune == '\'' && tokenNow.tokenType == tokenType_TextBlockQuotedSingle {
			event = tokenCloserDetected__saveTheToken
		}

		if charRuneNext1 == '\n' && tokenNow.tokenType == tokenType_Comment {
			// the endOfLine cannot be removed from original src,
			// comment is finished BEFORE the end of line
			event = tokenCloserDetected__saveTheToken
		}

		// openers...................................................................
		if noActiveTokenDetection__tokenTypeIsEmpty(tokenNow) {

			if charRune == '"' { // string
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_TextBlockQuotedDouble}
				event = tokenOpenerDetected__tokenNew
			}

			if charRune == '\'' { // quoted atom
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_TextBlockQuotedSingle}
				event = tokenOpenerDetected__tokenNew
			}

			if charRune == '%' { // comments
				tokenNow = Token{ positionCharFirst: charPos,
					              tokenType: tokenType_Comment}
				event = tokenOpenerDetected__tokenNew
			}
		} // not in active token detection



		/////////////////////////////////////////////////////////////////////
		if event == tokenOpenerDetected__tokenNew {
			// the opening/ending chars are removed from the original src, too
			// empty char is added instead of the original one
			erlSrcTokenDetectionsRemoved = append(erlSrcTokenDetectionsRemoved, ' ')
			event = ""
			continue
		}

		/////////////////////////////////////////////////////////////////////
		if event == tokenCloserDetected__saveTheToken {
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
