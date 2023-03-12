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
	compare_int_pair(received, wanted, t)
}

func Test_ErlSrcRead(t *testing.T) {
	chars, _ := ErlSrcRead("test/parse/hello.erl")
	compare_rune_pair(chars[1].Value, 'm', t)
	compare_rune_pair(chars[2].Value, 'o', t)
	compare_rune_pair(chars[3].Value, 'd', t)

	compare_rune_pair(chars[1].NextChar.Value, 'o', t)
	compare_rune_pair(chars[3].PrevChar.Value, 'o', t)

	compare_int_pair(chars[3].PosInFile, 3, t)
	/*
		if chars[0].PrevChar != nil {
			t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
		}

	*/
	compare_char_pointer_pair(chars[0].PrevChar, nil, t)

	compare_int_pair(chars[0].NextChar.PosInFile, chars[2].PrevChar.PosInFile, t)
	// compare_char_pointer_pair(chars[0].NextChar, chars[2].PrevChar, t)
}

// //////// test tools /////////////
func compare_char_pointer_pair(received, wanted *ErlSrcChar, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
	}
}

func compare_int_pair(received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nreceived: %v\n  wanted: %v, error", received, wanted)
	}
}

func compare_rune_pair(received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf(`received rune = %v, wanted %v, error`, received, wanted)
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
