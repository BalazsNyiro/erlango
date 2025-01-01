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
	"erlango.org/erlango/pkg/base_toolset"
	"fmt"
)

// TODO: func tokens_detect_all(characters, TOKENS)

func tokens_detect_comments_strings_quotedatoms(charactersInErlSrc []CharacterInErlSrc) string {
	for _, charStruct := range charactersInErlSrc {
		fmt.Printf("tokens_detect_comments_strings_quotedatoms: %s\n", charStruct.stringRepr())
	}
	return base_toolset.GetCurrentGoFuncName()
}
