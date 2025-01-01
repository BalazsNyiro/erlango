/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package base_toolset

import (
	"bufio"
	"io"
	"os"
)

func File_read_runes(filePath, caller string) ([]rune, error) {
	f, err := os.Open(filePath)
	if err != nil {
		LogError(err, caller+" "+filePath)
		return []rune{}, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	runes := []rune{}
	for {
		if runeInFile, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				LogError(err, caller+" read, Rune problem: "+filePath)
			}
		} else {
			runes = append(runes, runeInFile)
		}
	}
	return runes, nil
}
