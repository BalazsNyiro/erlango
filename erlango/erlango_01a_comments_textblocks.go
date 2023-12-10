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
	"strings"
)

func token_detect_comments_textblocks_alphanums(chars Chars, tokens ErlTokens, verboseForErlangoInvestigations__useFalseInProdEnv bool) ([]Char, ErlTokens, errorsDetected){
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
	errors := errorsDetected{}

	fmt.Println("token detest comments, quoted text blocks")

	tokenActualId := 0
	tokenActual := token_empty_obj("", tokenActualId)

	commentLineCloser := "\n"

	for charPos := 0; charPos < len(chars); charPos += 1 {

		charTxtPrev2 := char_txt_value_get(charPos-2, chars)
		charTxtPrev1 := char_txt_value_get(charPos-1, chars)
		charTxtNow := char_txt_value_get(charPos, chars)
		charTxtNext1 := char_txt_value_get(charPos+1, chars)

		saveCompleteDetectedToken := false

		///////// this section can be refactored to a separated fun.
		// BUT: that step means high complexity - it is longer, a little repetitive,
		// plus readable and easier to follow

		////////////// double quoted text detect //////////
		if charTxtNow == "\""{

			if tokenActual.typeIsEmpty() {
				tokenActual.TokenType = "tokenTextBlockQuotedDouble"

			} else { // TokenType is set before this " detection:
				if tokenActual.TokenType == "tokenTextBlockQuotedDouble" {
					if ! is_char_escaped_in_text_block(charPos, chars) {
						saveCompleteDetectedToken = true
					} // char is not escaped
				}
			} // TokenType was not empty
		}

		////////////// single quoted text detect //////////
		if charTxtNow == "'"{

			if tokenActual.typeIsEmpty() {
				tokenActual.TokenType = "tokenTextBlockQuotedSingle"

			} else { // TokenType is set before this ' detection:
				if tokenActual.TokenType == "tokenTextBlockQuotedSingle" {
					if ! is_char_escaped_in_text_block(charPos, chars) { // an atom can have a ' char in its content, too
						saveCompleteDetectedToken = true
					} // char is not escaped
				}
			} // TokenType was not empty
		}


		if tokenActual.typeIsEmpty() {
			if charTxtNow == "%"{
				tokenActual.TokenType = "tokenComment"
			}
		}

		if tokenActual.TokenType == "tokenComment" {
			if charTxtNow == commentLineCloser {
				saveCompleteDetectedToken = true
			}
		} // comment detect...



		/* you can ask this: why is it good to detect abc letters and numbers together?
		because numbers can be mixed in the Erlang code often with letters,
		and later it is easier to analyse one block and decide that is it a number only, or not.

		With other words, if abc+nums are detected together, from this direction
		it is easier to find numbers only, than to detect the numeric and abc chars separately,
		and explain the sitation when characters and numbers are mixed in one condition.
		*/
		//// ABC + numbers block detect  ///////////////////////////////////////////////////////////
		if tokenActual.typeIsEmpty() {
			if strings.Contains(abcEngLowerUpper_underscore_at_digits_Underscore_At_digits__atomFormerChars, charTxtNow) {
				tokenActual.TokenType = "tokenAbcFullWith_Underscore_At_numbers"
			}
		}
		if tokenActual.TokenType == "tokenAbcFullWith_Underscore_At_numbers" {
			// if the next char is not in abc, then the current one is the closer.
			if ! strings.Contains(abcEngLowerUpper_underscore_at_digits_Underscore_At_digits__atomFormerChars, charTxtNext1) {
				saveCompleteDetectedToken = true
			}
		} // ABC detect



		///// OTHER PUNCTUATION DETECT - they are 1 char wide elems in the source code //////////////////////////////////////////////////////////
		if tokenActual.typeIsEmpty() {
			if strings.Contains(otherPunctuation, charTxtNow) {
				tokenActual.TokenType = "tokenOtherPunctuation"
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}
		} // OTHER PUNCTUATION DETECT


		///// white space  DETECT - they are 1 char wide elems in the source code //////////////////////////////////////////////////////////
		if tokenActual.typeIsEmpty() {
			if is_whitespace_only(charTxtNow) {
				tokenActual.TokenType = "tokenWhiteSpace"
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}
		} // whitespace DETECT


		// Character literals. Example: $∑
		if tokenActual.typeIsEmpty() {
			// $A: A is literal, prev is $
			// $\n \n is literal, prev2 is $, prev1 is escape
			if charTxtPrev1 == "$" || (charTxtPrev1 == "\\" && charTxtPrev2 == "$") {
				tokenActual.TokenType = "tokenCharLiteral"
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}
		}


		/////////////////////// TOKEN SAVE, CLOSE ////////////////////////////////////////
		weAreInTokenDetection := ! tokenActual.typeIsEmpty()

		if weAreInTokenDetection { // if we are in a token block, save the current char into the token
			chars[charPos].TokenId = tokenActual.TokenId
			chars[charPos].TokenDetected = true
			tokenActual.SourceCodeChars = append(tokenActual.SourceCodeChars, chars[charPos])
			chars[charPos].TokenFirstCharPositionInFile = tokenActual.SourceCodeChars[0].PositionInFile
			// we can be sure that minimum one char exists, because 2 lines before there is an append()
		} else {
			// if we are NOT in a block, that is a problem, a non-recognised character
			errMsg := errorDetected{
				filePath: chars[charPos].WhereTheCharIsStored,
				lineNum: chars[charPos].LineNum,
				charPosInLine: chars[charPos].PositionInLine,
				charPosInFile: chars[charPos].PositionInFile,
				errMsg: "Unrecognised char ("+string(chars[charPos].Value)+") in token detection",
			}

			errors = append(errors, errMsg)
		}

		if saveCompleteDetectedToken {
			if tokenActual.TokenType != "tokenWhiteSpace" && tokenActual.TokenType != "tokenComment" {

				if verboseForErlangoInvestigations__useFalseInProdEnv {
					// the string representation can be asked with a function,
					// so this saving is support the Erlango debugging only,
					// when in a debugger it is easier to follow the saved value.
					// in production this is not necessary
					tokenActual.DebugStringRepresentation = tokenActual.stringRepresentation()
				}

				// save tokenActual into tokens - skip comments and whitespace tokens
				tokens[tokenActual.charPosFirst()] = tokenActual
			}
			tokenActualId := len(tokens) // len(..) is always represent the next free, unused elem Id
			tokenActual = token_empty_obj("", tokenActualId)
		}
	}

	return chars, tokens, errors
}
func is_char_escaped_in_text_block(posChar int, chars Chars) bool {
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

