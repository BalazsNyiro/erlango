/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.3, third total rewrite
*/

package main

import "fmt"
import "erlango.org/erlango/tokens"

func main() {  // in program plan

	fmt.Println("Erlango")

	erlSrc := ""
	tokensTable := tokens.Tokens{}


	erlSrc_cleaned__stringsAtomsCommens_removed, tokensTable_stringsAtomsCommentsAdded :=
		tokens.Tokens_0_detect_text_blocks(erlSrc, tokensTable)

	erlSrc_cleaned__numbersRemoved, tokensTable_numbersAdded :=
		tokens.Tokens_1_detect_numbers(erlSrc_cleaned__stringsAtomsCommens_removed, tokensTable_stringsAtomsCommentsAdded)

	erlSrc_cleaned__atomsVariablesRemoved, tokensTable_atomsVariablesAdded :=
		tokens.Tokens_2_detect_atoms_variableNames(erlSrc_cleaned__numbersRemoved, tokensTable_numbersAdded)

	fmt.Println("END", erlSrc_cleaned__atomsVariablesRemoved, tokensTable_atomsVariablesAdded)

}
