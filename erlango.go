/*
Erlang - Go implementation.

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.4, fourth rewrite
*/

package main

import (
	tokens "erlango.org/erlango/pkg/tokens"
	"fmt"
)

func main() {

	fmt.Println("Erlango")

	// Placeholder
	erlSrcRunes := []rune("Num = 7.")
	charactersInErlSrc := tokens.Runes_to_character_structs(erlSrcRunes, "on-the-fly-defined")

	charactersInErlSrc, errors := tokens.Character_block_detect_in_erl_src(charactersInErlSrc)

	// RESULT 1: every important char is detected, and a simple token is assigned to them
	if len(errors) == 0 {
		// https://stackoverflow.com/questions/5639177/what-exactly-is-a-token-in-relation-to-parsing
		// TODO: tokensInErlSrc = ....append() ? erlang_language_objects_creation_from_character_blocks
		tokens.CharacterCollectorPrint(charactersInErlSrc)
	}

}
