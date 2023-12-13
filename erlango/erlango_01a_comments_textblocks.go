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

func token_detect_comments_textblocks_alphanums_whitespaces_literals(chars Chars, tokens ErlTokens, verboseForErlangoInvestigations__useFalseInProdEnv bool) ([]Char, ErlTokens, errorsDetected){ // in program plan
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
	tokenActual := ErlToken_empty_obj("", tokenActualId)

	commentLineCloser := "\n"

	for charPos := 0; charPos < len(chars); charPos += 1 {

		charTxtPrev2 := charTxtValueGet_for_tokenDetection(charPos-2, chars)
		charTxtPrev1 := charTxtValueGet_for_tokenDetection(charPos-1, chars)
		charTxtNow := charTxtValueGet_for_tokenDetection(charPos, chars)
		charTxtNext1 := charTxtValueGet_for_tokenDetection(charPos+1, chars)

		saveCompleteDetectedToken := false

		///////// this section can be refactored to a separated fun.
		// BUT: that step means high complexity - it is longer, a little repetitive,
		// plus readable and easier to follow

		////////////// double quoted text detect //////////
		if charTxtNow == "\""{

			if tokenActual.typeIsEmpty() {
				tokenActual.TokenType = tokenType_TextBlockQuotedDouble

			} else { // TokenType is set before this " detection:
				if tokenActual.TokenType == tokenType_TextBlockQuotedDouble {
					if ! isCharEscapedInTextBlock__tokenDetectionQuoteds(charPos, chars) {
						saveCompleteDetectedToken = true
					} // char is not escaped
				}
			} // TokenType was not empty
		}

		////////////// single quoted text detect //////////
		if charTxtNow == "'"{

			if tokenActual.typeIsEmpty() {
				tokenActual.TokenType = tokenType_TextBlockQuotedSingle

			} else { // TokenType is set before this ' detection:
				if tokenActual.TokenType == tokenType_TextBlockQuotedSingle {
					if ! isCharEscapedInTextBlock__tokenDetectionQuoteds(charPos, chars) { // an atom can have a ' char in its content, too
						saveCompleteDetectedToken = true
					} // char is not escaped
				}
			} // TokenType was not empty
		}


		if tokenActual.typeIsEmpty() {
			if charTxtNow == "%"{
				tokenActual.TokenType = tokenType_Comment
			}
		}

		if tokenActual.TokenType == tokenType_Comment {
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
				tokenActual.TokenType = tokenType_AbcFullWith_Underscore_At_numbers
			}
		}
		if tokenActual.TokenType == tokenType_AbcFullWith_Underscore_At_numbers {
			// if the next char is not in abc, then the current one is the closer.
			if ! strings.Contains(abcEngLowerUpper_underscore_at_digits_Underscore_At_digits__atomFormerChars, charTxtNext1) {
				saveCompleteDetectedToken = true
			}
		} // ABC detect



		///// OTHER PUNCTUATION DETECT - they are 1 char wide elems in the source code //////////////////////////////////////////////////////////
		if tokenActual.typeIsEmpty() {
			if strings.Contains(otherPunctuation, charTxtNow) {
				tokenActual.TokenType = tokenType_OtherPunctuation
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}
		} // OTHER PUNCTUATION DETECT


		///// white space  DETECT - they are 1 char wide elems in the source code //////////////////////////////////////////////////////////
		if tokenActual.typeIsEmpty() {
			if is_whitespace_only(charTxtNow) {
				tokenActual.TokenType = tokenType_WhiteSpace
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}
		} // whitespace DETECT


		// Character literals. Example: $∑
		if tokenActual.typeIsEmpty() {
			// $A: A is literal, prev is $
			if charTxtPrev1 == "$" && charTxtNow != "\\" {  // so this is not an escaped char literal, \n for example
				tokenActual.TokenType = tokenType_CharLiteral
				// because they are 1 char wide elems, the block is closed at the first char
				saveCompleteDetectedToken = true
			}

			if charTxtPrev2 == "$" && charTxtNow == "\\" {  // escaped literal
			// $\n \n is literal, prev2 is $, prev1 is escape
				tokenActual.TokenType = tokenType_CharLiteral
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
			if tokenActual.TokenType != tokenType_WhiteSpace && tokenActual.TokenType != tokenType_Comment {

				if verboseForErlangoInvestigations__useFalseInProdEnv {
					// the string representation can be asked with a function,
					// so this value supports the Erlango debugging only:
					// in a debugger it is easier to follow a simple value
					tokenActual.DebugStringRepresentation = tokenActual.stringRepresentation()
				}

				// save tokenActual into tokens - skip comments and whitespace tokens
				tokens[tokenActual.charPosFirst()] = tokenActual
			}
			tokenActualId := len(tokens) // len(..) is always represent the next free, unused elem Id
			tokenActual = ErlToken_empty_obj("", tokenActualId)
		}
	}

	return chars, tokens, errors
}
func isCharEscapedInTextBlock__tokenDetectionQuoteds(posChar int, chars Chars) bool {
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

func charTxtValueGet_for_tokenDetection(pos int, chars Chars) string {  // in program plan
	ret := "" 	// I would like to handle empty values, too, so runes cannot be given back.
	// empty value means: there is no real character in the wanted position
	// the position has a real value only if it is in the valid range
	if pos >= 0 && pos < len(chars) {
		ret = string(chars[pos].Value)
	}
	return ret
}

