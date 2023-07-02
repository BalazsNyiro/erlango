/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

import "fmt"

func debug_print_ErlSrcChars(prg *Prg, chars []ErlSrcChar) {
	fmt.Printf("inside >>>  debug_print_ErlSrcChars %p \n", chars)
	for id, _ := range chars {
		debug_print_ErlSrcChar(prg, id, chars)
	}
}

// func debug_print_ErlSrcChar(prg *Prg, id int, charPtr *ErlSrcChar, ) {
func debug_print_ErlSrcChar(prg *Prg, id int, chars []ErlSrcChar) {
	charPtr := &(chars[id])
	// fmt.Printf("charPtr:%p  %3d posInFile:%3d val:%4s ",charPtr, id, charPtr.PosInFile, string(charPtr.Value))
	fmt.Printf("%3d %4s ", charPtr.PosInFile, string(charPtr.Value))

	prevPos := -1
	if charPtr.PrevChar != nil {
		prevPos = charPtr.PrevChar.PosInFile
	}
	_ = prevPos
	/*
	fmt.Printf(" PrevPosInFile:%3d ", prevPos)

	fmt.Printf(" %p <- %p -> %p tokenPtr: %p type->%s<-", charPtr.PrevChar, charPtr, charPtr.NextChar, charPtr.Token, (*charPtr).Type())
	fmt.Println("")

	 */
	// charPtr.Token, (*charPtr).Type())
	tokenDisplayed := fmt.Sprintf("token: %p", charPtr.Token)
	if id > 0 {
		if charPtr.Token == (&(chars[id-1])).Token {
			tmp := " "
			for len(tmp) < len(tokenDisplayed) { tmp = tmp + " "}
			tokenDisplayed = tmp
		}
	}
	fmt.Printf(" %s ", tokenDisplayed)
	fmt.Printf(" (%s)", (*charPtr).Type())
	fmt.Println("")
}

func debug_print_mem_address(msg string, object any ){
	fmt.Printf("%s %p \n", msg, object)
}

