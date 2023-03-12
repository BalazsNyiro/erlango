/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package erlango

import (
	"fmt"
	"testing"
)

func Test_ParseErlangSourceFile(t *testing.T) {
	received := ParseErlangSourceFile()
	wanted := 0
	compare_int_pair("fake parse test", received, wanted, t)
}

func Test_ErlSrcRead(t *testing.T) {
	chars, _ := ErlSrcRead("test/parse/hello.erl")
	compare_rune_pair("val m", chars[1].Value, 'm', t)
	compare_rune_pair("val o", chars[2].Value, 'o', t)
	compare_rune_pair("val d", chars[3].Value, 'd', t)

	compare_rune_pair("val next 1", chars[1].NextChar.Value, 'o', t)
	compare_rune_pair("val prev 3", chars[3].PrevChar.Value, 'o', t)

	compare_int_pair("pos 3", chars[3].PosInFile, 3, t)
	/*
		if chars[0].PrevChar != nil {
			t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
		}

	*/
	compare_char_pointer_pair("compare char_0 and nil_prev_char", chars[0].PrevChar, nil, t)

	compare_int_pair("pos_0.next, pos2.prev", chars[0].NextChar.PosInFile, chars[2].PrevChar.PosInFile, t)
	debug_print_ErlSrcChars(chars)
	// compare_char_pointer_pair("compare char_0.next and char_2.prev", chars[0].NextChar, chars[2].PrevChar, t)
}

// //////// test tools /////////////
func debug_print_ErlSrcChars(chars []ErlSrcChar) {
	fmt.Println("")
	for i, _ := range chars {
		fmt.Printf("%3d posInFile:%3d val:%4v ", i, chars[i].PosInFile, chars[i].Value)

		prevPos := 0
		if chars[i].PrevChar != nil {
			prevPos = chars[i].PrevChar.PosInFile
		}
		fmt.Printf(" PrevPosInFile:%3d ", prevPos)

		fmt.Printf(" %p <- %p -> %p", chars[i].PrevChar, &chars[i], chars[i].NextChar)
		fmt.Println("")
	}
}

func compare_char_pointer_pair(callerInfo string, received, wanted *ErlSrcChar, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

func compare_int_pair(callerInfo string, received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

func compare_rune_pair(callerInfo string, received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received rune = %v, wanted %v, error", callerInfo, received, wanted)
	}
}

///////// go experimental tests

func Test_struct_modifications(_ *testing.T) {

	// this is a trial: what happens with structs after modifications?
	var chars []ErlSrcChar
	A := ErlSrcChar{Value: 'a'}

	chars = append(chars, A)

	// question 1: is chars[0] object same with A?
	fmt.Printf("Address of struct       A = %+v: %p\n", A, &A)
	fmt.Printf("Address of struct chars[0]= %+v: %p\n", chars[0], &chars[0])
	/*  at this point the address of the two objects are different:
	    Address of struct       A = {NextChar:<nil> PrevChar:<nil> PosInFile:0 Value:97 Token:<nil>}: 0xc000014510
	    Address of struct chars[0]= {NextChar:<nil> PrevChar:<nil> PosInFile:0 Value:97 Token:<nil>}: 0xc000014540
	*/

	A.PosInFile = 1
	fmt.Printf("after position in file change in A:\n")
	fmt.Printf("Address of struct       A = %+v: %p\n", A, &A)
	fmt.Printf("Address of struct chars[0]= %+v: %p\n", chars[0], &chars[0])

	// result: when I append an elem into a slice, a copy is inserted.
}
