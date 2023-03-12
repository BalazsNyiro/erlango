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
	Chars     []ErlSrcChar
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
			lastId := len(erlChars) - 1
			if lastId > 0 {
				erlChars[lastId].PrevChar = &erlChars[lastId-1]
				erlChars[lastId-1].NextChar = &erlChars[lastId]
			}
		}
	}
	return erlChars, nil
}
