/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

import (
	"fmt"
)

/* iterative parsing:
   - from a flat list of chars the parser builds up
     an embedded token structure
*/

// Token: indepentend language unit, formed by one or more char
type ErlSrcToken struct {
	PrevToken *ErlSrcToken
	NextToken *ErlSrcToken
	Chars     []*ErlSrcChar
}

// represents one char in the Erlang source codes
type ErlSrcChar struct {
	NextChar  *ErlSrcChar
	PrevChar  *ErlSrcChar
	PosInFile int
	Value     rune
	Token     *ErlSrcToken
}

func ParseErlangSourceFile() int {
	return 0
}

func ErlSrcRead(filePath string) ([]ErlSrcChar, error) {
	runes, err := file_read_runes(filePath, "ErlSrcRead")
	if err != nil { return []ErlSrcChar{}, err}
	erlChars := ErlSrcChars_from_runes(runes)
	// Test_what_happens_with_struct_pointers
	// fmt.Printf("ErlSrcRead, chars pointer before return: %p\n", erlChars)
	return erlChars, nil
}

func ErlSrcChars_from_runes(runes []rune) []ErlSrcChar {
	var erlChars []ErlSrcChar
	for posInFile, runeInFile := range runes {
		erlChars = append(erlChars, ErlSrcChar{
			Value:     runeInFile,
			PosInFile: posInFile,
		})
	}
	// the slice pointers won't be change after this point,
	// there is no capacity change later.
	// if we do this from the 'previous linking position'
	// then because of the capacity limit reach, the pointers
	// will be incorrect in the early elements
	for id, _ := range erlChars {
		if id > 0 {
			erlChars[id].PrevChar = &erlChars[id-1]
			erlChars[id-1].NextChar = &erlChars[id]
		}
	}
	return erlChars
}

func ErlSrcTokens_Quoted(wanted rune, chars []ErlSrcChar, verbose bool) {
	tokenActual := ErlSrcToken{}
	inQuote, escapeOn := false, false

	for id, char := range chars {
		nowOpened, nowEscaped := false, false

		if !inQuote && (char.Value == wanted) {
			inQuote, nowOpened = true, true
		}

		if !escapeOn && inQuote && (char.Value == '\\') {
			escapeOn, nowEscaped = true, true
		}

		if inQuote {
			tokenActual.Chars = append(tokenActual.Chars, &chars[id])
			chars[id].Token = &tokenActual
		}
		if verbose { fmt.Println("ErlSrcTokens_Quoted", id, string(char.Value), bool_to_str(inQuote, "in Quote", "")) }

		if nowOpened || nowEscaped { continue }
		// ##### stop here ^^^^ the char processing in these 2 cases ###########
		// if nowOpened == true, the sign is '\' and I don't want to turn it off if it was turned on just now
		// if it's nowEscaped, I don't want to turn it off too because it has effect on the next char

		if !escapeOn && inQuote && (char.Value == wanted) { // active escape blocks the next char detection: \", \'
			inQuote = false
			tokenActual = ErlSrcToken{}
		}

		escapeOn = false // if not now escaped, the escape disappearing at the next char.
	}
}
