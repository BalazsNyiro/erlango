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

	/**/
	fmt.Println("Erlango")

	erlSrc := ""
	tokensTable := tokens.Tokens{}
	tokens.Tokens_detect_text_blocks(erlSrc, tokensTable)
}
