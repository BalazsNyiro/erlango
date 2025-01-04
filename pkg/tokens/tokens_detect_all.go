/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import "fmt"

func Tokens_detect_in_erl_src(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	charactersInErlSrc2, tokensInErlSrc2 := tokens_detect_comments_strings_quotedatoms(charactersInErlSrc, tokensInErlSrc)
	return charactersInErlSrc2, tokensInErlSrc2
}

func Tokens_detection_print_verbose(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) {

	for charPos, charInErlSrc := range charactersInErlSrc {
		fmt.Printf("print_verbose >>> charPosition: %d, characterLoop: %s\n", charPos, string(charInErlSrc.runeInErlSrc))
	}

}
