/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.3, third total rewrite
*/

package tokens

import (
	"fmt"
	"testing"
)


func Test_parse_comments_textDoubleQuoted_textSingleQuoted(t *testing.T) {

	erlSrc :=`	
		VarAtomQuoted = 'atomValue1', 
		VarStr  = "string"	
		[A, B, StringValue] = [1, 2, "This is a line with a string"]

	`
	tokensTable := Tokens{}

	erlSrc_commentsStringsAtomsquoted_removed, tokensTableXXX := Tokens_detect_comments_strings_quotedatoms(erlSrc, tokensTable)

	fmt.Println("erlSrc, without strings, quoted atoms", erlSrc_commentsStringsAtomsquoted_removed)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTableXXX)

}

/*
func compare_tokenDetected_tokenWanted(callerInfo string, tokensDetected ErlTokens, tokenWanted tokenWanted, t *testing.T) {
	tokenDetected, tokenWantedIsInDetected:= tokensDetected[tokenWanted.positionFirst]

	if tokenWantedIsInDetected {
		// theoretically the charPosFirst is always ok here, because the key in map was the same position
		tokenDetected_charPosFirst, tokenDetected_charPosLast := tokenDetected.charPositionFirstLast()
		if tokenDetected_charPosFirst != tokenWanted.positionFirst {
			t.Fatalf("\nErr First: %s : detected posFirst: %v  is different from wanted posFirst:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosFirst, tokenWanted.positionFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected_charPosLast != tokenWanted.positionLast {
			t.Fatalf("\nErr Last: %s : detected posLast: %v  is different from wanted posLast:  %v, error:\n'%s'\n'%s'\n\n",
				callerInfo, tokenDetected_charPosLast, tokenWanted.positionLast, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
		if tokenDetected.stringRepresentation() != tokenWanted.textRepresentation {
			t.Fatalf("\nErr repr %s : startPos:%v  detected string representation: %v  is different from wanted representation:  %v, error",
				callerInfo, tokenDetected_charPosFirst, tokenDetected.stringRepresentation(), tokenWanted.textRepresentation)
		}
	} else {
		t.Fatalf("\nErr %s : wanted tokenPosFirst %v is not in detecteds - error", callerInfo, tokenWanted.positionFirst)
	}

}

 */
