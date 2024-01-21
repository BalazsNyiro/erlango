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

import "slices"

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
func (token Token) stringRepr() string {
	return string(token.charsInErlSrc)
}


const tokenType_TextBlockQuotedSingle = "tokenTextBlockQuotedSingle"
const tokenType_Comment = "tokenComment"
const tokenType_TextBlockQuotedDouble = "tokenTextBlockQuotedDouble"

const tokenType_Num_digitsZeroNine = "tokenNumDigitsZeroNine"

// when the first char is a digit, then later underscores can be between digits
const tokenType_Num_digitsZeroNine_underscoreMaybeLater = "tokenNumDigitsZeroNine_underscoreMaybeLater"

const tokenType_Num_charLiterals = "tokenNumCharLiteral"

/* Tokens represent the Erlang source code - so the int-key is the first char's position in the source code */
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
	if you pass more character groups, are all of they are matching?
*/
func charsGroupsAreMatching(charPosFirstToTest int, erlSrcRunes []rune, allowedChars_sets []([]rune), direction string) int {
	inSetCounter := 0

	charPosToTest := charPosFirstToTest
	for _, allowedCharsInSet := range allowedChars_sets {
		countedCharNum := charsHowManyAreInTheGroup(charPosToTest, erlSrcRunes, allowedCharsInSet, direction)
		if countedCharNum == 0 {
			inSetCounter = 0  // if in any set, there is no matching elems, the whole group Matching is unsuccessful
			continue
		}

		// here the countedCharNum is > 0
		inSetCounter += countedCharNum
		charPosToTest += countedCharNum // the nextTested char position is after the countedCharNum
	}
	return inSetCounter
}

/*
	charPosFirstToTest: the first tested character position in erlSrcRunes
	allowedCharsInSet: one or more runes: if the actual char is in the set,
	func detectCharacterSetLength(charPosFirstToTest int, erlSrcRunes []rune, allowedCharsInSet []rune, direction string){

	it tests one Char position.
		if it is in ghe set, counter++, and step to the next one.
		if it is not in the set, stop the validation and return with the counter


	if the slices.Contains() is too slow, a map can be used instead of the slice
*/

func charsHowManyAreInTheGroup(charPosFirstToTest int, erlSrcRunes []rune, allowedCharsInSet []rune, direction string) int {
	inSetCounter := 0

	delta := +1 // direction: right, add +1 in all steps
	conditionFun := func(position int) bool { return position <= len(erlSrcRunes) - 1 }

	if direction == "left" {
		delta = -1 // to go left, pos has to be decreased
		conditionFun = func(position int) bool { return position >= 0}
	}

	for pos := charPosFirstToTest; conditionFun(pos); pos+=delta {
		charNow := erlSrcRunes[pos]
		if slices.Contains(allowedCharsInSet, charNow) {
			inSetCounter += 1
		} else { // if the tested position's char is not in the set, leave the validation loop
			break
		}
	}

	return inSetCounter
}


// return with the next Nth char (relative to the actual char.
// charPosRelative == 1 means: the next char
// charPosRelative == 0 means: the actual char
// charPosRelative == -1 means: the prev char
func charRuneNext(charPosActual, charPosRelative int, erlSrcRunes []rune) (rune, bool) {
	charRuneWanted := ' '                   // if the wanted position is not in range, this is the default value
	wantedCharInSrcRunesIndexRange := false // it is possible that the calculated position is outside of the range or runes.

	charPosCalculated := charPosActual + charPosRelative
	if charPosCalculated < len(erlSrcRunes) {
		if charPosCalculated >= 0 {

				wantedCharInSrcRunesIndexRange = true
				charRuneWanted = erlSrcRunes[charPosCalculated]
		}
	}
	return charRuneWanted, wantedCharInSrcRunesIndexRange
}

