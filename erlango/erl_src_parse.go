/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	var erlChars []ErlSrcChar
	f, err := os.Open(filePath)
	if err != nil {
		LogError(err, "Erl Src read: "+filePath)
		return erlChars, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	posInFile := -1
	for {
		if runeInFile, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				LogError(err, "Erl Src read, Rune problem: "+filePath)
			}
		} else {
			posInFile += 1
			// fmt.Printf("%q [%d]\n", string(runeInFile), runeSize)
			erlChars = append(erlChars, ErlSrcChar{
				Value:     runeInFile,
				PosInFile: posInFile,
			})
			// Place: previous linking position
		}
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

	// Test_what_happens_with_struct_pointers
	// fmt.Printf("ErlSrcRead, chars pointer before return: %p\n", erlChars)

	return erlChars, nil
}

func ErlSrcTokens_Quoted(wanted rune, chars []ErlSrcChar) {
	tokenActual := ErlSrcToken{}
	inQuote := false
	escapeOn := false

	for id, char := range chars {
		nowOpened, nowEscaped := false, false

		if !inQuote && (char.Value == wanted) {
			inQuote = true
			nowOpened = true
		}

		if !escapeOn && inQuote && (char.Value == '\\') {
			escapeOn = true
			nowEscaped = true
		}

		info := ""
		if inQuote {
			tokenActual.Chars = append(tokenActual.Chars, &chars[id])
			chars[id].Token = &tokenActual
			info = "inQuote"
		}
		fmt.Println("ErlSrcTokens_Quoted", id, string(char.Value), info)

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
