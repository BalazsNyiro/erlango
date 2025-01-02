/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import (
	"fmt"
)

func tokens_detect_in_erl_src(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	charactersInErlSrc2, tokensInErlSrc2 := tokens_detect_comments_strings_quotedatoms(charactersInErlSrc, tokensInErlSrc)
	return charactersInErlSrc2, tokensInErlSrc2
}

func tokens_detect_comments_strings_quotedatoms(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {
	funTokenOpener := token_opener_detect_quote_double
	funTokenCloser := token_closer_quote_double
	charactersInErlSrc, tokensInErlSrc = character_loop(charactersInErlSrc, tokensInErlSrc, funTokenOpener, funTokenCloser)
	return charactersInErlSrc, tokensInErlSrc
}

func character_loop(
	charactersInErlSrc CharacterInErlSrcCollector,
	tokensInErlSrc TokenCollector,
	tokenOpenerConditionFun func() bool,
	tokenCloserConditionFun func() bool) (CharacterInErlSrcCollector, TokenCollector) {

	for charPositionInSrc, charStruct := range charactersInErlSrc {
		fmt.Printf("charPosition: %d, characterLoop: %s\n", charPositionInSrc, charStruct.stringRepr())
	}

	return charactersInErlSrc, tokensInErlSrc
}
func token_opener_detect_quote_double() bool {
	return true
}

func token_closer_quote_double() bool {
	return true
}
