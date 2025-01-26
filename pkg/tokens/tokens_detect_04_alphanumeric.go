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

var tokens_ABC_Eng_Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var tokens_ABC_Eng_Lower = "abcdefghijklmnopqrstuvwxyz"

var tokens_AlphaNums = tokens_ABC_Eng_Upper + tokens_ABC_Eng_Lower + "_0123456789"

func tokens_detect_prepare__03_erlang_alphanumerics(charactersInErlSrc CharacterInErlSrcCollector) CharacterInErlSrcCollector {

	for _, wantedCharInErl := range []rune(tokens_AlphaNums) {
		fmt.Println("wanted char in alphanum:", wantedCharInErl)
		charactersInErlSrc = character_loop__set_one_char_tokentype(wantedCharInErl, charactersInErlSrc, CharBlock_AlphaNumeric)
	}
	character_loop__opener_closer_sections_set__if_more_separated_isolated_elems_are_next_to_each_other(charactersInErlSrc, CharBlock_AlphaNumeric)

	return charactersInErlSrc
}
