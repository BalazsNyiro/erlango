/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite

*/

package erlango

import "fmt"

func is_empty_token_block_name__textBlockDetection(blockName string) bool {
	return blockName == ""
}

func token_detect_comments_textblocks(chars []Char, tokens []ErlToken) ([]Char, []ErlToken){
	// the "wrapper" quotes around the string values or 'atoms' are the part of the tokens,
	// they are necessary to define a text block (single or double qoted texts)
	// but not part of the value of the token

	/* Erlang accepts newlines in atoms:

	Erlang/OTP 25 [erts-13.1.5] [source] [64-bit] [smp:4:4] [ds:4:4:10] [async-threads:1] [jit:ns]

	Eshell V13.1.5  (abort with ^G)
	1> A = 'atom\n2'.
	'atom\n2'
	2> A.
	'atom\n2'
	3> A2 = 'atom\\n'.
	'atom\\n'
	4>

	discussion: https://erlang.org/pipermail/erlang-questions/2014-February/077922.html

	*/

	fmt.Println("token detext comments, quoted textblocks")
	blockName := ""

	tokenActual := ErlToken{}
	commentLineCloser := "\n"

	for charPos := 0; charPos < len(chars); charPos += 1 {
		tokenActualId := len(tokens) // len(..) is always represent the next free, unused elem Id in the slice

		//charTxtPrev1 := char_txt_value_get(charPos-1, chars)
		charTxtNow := char_txt_value_get(charPos, chars)
		//charTxtNext1 := char_txt_value_get(charPos+1, chars)

		// block Start detection is important, when the opener and closer patterns are the same: " or ' chars
		blockStarted := false
		blockLastElemDetected__saveCompleteDetectedToken := false

		///////// this section can be refactored to a separated fun.
		// BUT: that step means high complexity - it is longer, a little repetitive,
		// plus readable and easier to follow

		////////////// double quoted text detect //////////
		if charTxtNow == "\""{

			if is_empty_token_block_name__textBlockDetection(blockName) {
				tokenActual = token_empty_obj("tokenTextBlockQuotedDouble", tokenActualId)
				blockName = "inTextBlockQuotedDouble"
				blockStarted = true
			}

			if blockName == "inTextBlockQuotedDouble" && ! blockStarted {
				if ! is_char_escaped_in_text_block(charPos, chars) {
					blockLastElemDetected__saveCompleteDetectedToken = true
				} // char is not escaped
			}
		}

		////////////// single quoted text detect //////////
		if charTxtNow == "'"{

			if is_empty_token_block_name__textBlockDetection(blockName) {
				tokenActual = token_empty_obj("tokenTextBlockQuotedSingle", tokenActualId)
				blockName = "inTextBlockQuotedSingle"
				blockStarted = true
			}

			if blockName == "inTextBlockQuotedSingle" && ! blockStarted {
				if ! is_char_escaped_in_text_block(charPos, chars) { // an atom can have a ' char in it's content, too
					blockLastElemDetected__saveCompleteDetectedToken = true
				} // char is not escaped
			}
		}


		////////////// for comment detect, blockStarted var is not important,
		// because the opener '%' and the closer '\n' patterns are different,
		// the opening or closing situations can be detected easily.
		// in previous cases, for 'atom', or "string", the opener and closer patterns are same,
		// so the blockStart var is necessary to know: have we started or closed a block?
		if is_empty_token_block_name__textBlockDetection(blockName) {
			if charTxtNow == "%"{
				tokenActual = token_empty_obj("tokenComment", tokenActualId)
				blockName = "inComment"
			}
		}
		if blockName == "inComment" {
			if charTxtNow == commentLineCloser {
				blockLastElemDetected__saveCompleteDetectedToken = true
			}
		} // comment detect... /////////////////////////////////////////////


		if ! is_empty_token_block_name__textBlockDetection(blockName) { // if we are in a token block, save the current char into the token
			chars[charPos].TokenId = tokenActual.TokenId
			chars[charPos].TokenDetected = true
			tokenActual.SourceCodeChars = append(tokenActual.SourceCodeChars, chars[charPos])
		}

		if blockLastElemDetected__saveCompleteDetectedToken {
			blockName = ""
			tokens = append(tokens, tokenActual)
		}
	}

	return chars, tokens
}
func is_char_escaped_in_text_block(posChar int, chars []Char) bool {
	// char is escaped if there are 'odd' num of escape char before that.
	escaped := false
	escapeCharCounter := 0

	posTestedChar := posChar - 1
	for posTestedChar >= 0 {
		if chars[posTestedChar].Value == '\\' {
			escapeCharCounter += 1
			posTestedChar -= 1
		} else { // if the char is not a backslash, leave the loop
			break
		}
	}
	escaped = (escapeCharCounter % 2) == 1 // odd escape chars are before the current char
	return escaped
}

