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
	"strconv"
	"strings"
	"testing"
)
///////////////////// TEST Globals //////////////////////////////////////
/* sometimes the direct GLOBAL names are used, sometime I can use it as a key only.
   so these key-value pairs are defined in two ways, the meaning is same.
   the calling mode is different.

   They are test supporter variables, used in strings - with these values
   the testing is much easier
*/
var TestGlobals = map[string]string{  // used from tests
	"Token_type_txt_quoted_double": Token_type_txt_quoted_double,
	"Token_type_txt_quoted_single": Token_type_txt_quoted_single,
	"no_type"                     : "", // in ErlSrcChar, "" means: no-type
	"empty_string"                : "",
	"space"                       : " ",
	"tabulator"                   : "\t",
	"newline_unix"                : "\n",
} //////////////////////////////////////////////////////////////////////

// TODO: at the end do a normal test for a complete parse
func Test_ParseErlangSourceFile(t *testing.T) {
	received := ParseErlangSourceFile()
	wanted := 0
	compare_int_pair("fake parse test", received, wanted, t)
}

func Test_ErlSrcRead(t *testing.T) {
	chars, _ := ErlSrcChars_from_file("test/parse/hello.erl")
	compare_rune_pair("val m", chars[1].Value, 'm', t)
	compare_rune_pair("val o", chars[2].Value, 'o', t)
	compare_rune_pair("val d", chars[3].Value, 'd', t)

	compare_rune_pair("val next 1", chars[1].NextChar.Value, 'o', t)
	compare_rune_pair("val prev 3", chars[3].PrevChar.Value, 'o', t)

	compare_int_pair("pos 3", chars[3].PosInFile, 3, t)
	compare_char_pointer_pair("compare char_0 and nil_prev_char", chars[0].PrevChar, nil, t)

	compare_int_pair("pos_0.next, pos2.prev", chars[0].NextChar.PosInFile, chars[2].PrevChar.PosInFile, t)
	debug_print_ErlSrcChars(chars)
	// compare_char_pointer_pair("compare char_0.next and char_2.prev", chars[0].NextChar, chars[2].PrevChar, t)
}

func Test_ErlSrcTokens_Quoted(t *testing.T) {
	// rules: the left column is the char, the right column is the type
	// a single char means himself. Keywords has special meanings
	wantedTable := `   a       no_type 
                       b       no_type
                       c       no_type
                       '       Token_type_txt_quoted_single
                       d       Token_type_txt_quoted_single
                       e       Token_type_txt_quoted_single
                       \       Token_type_txt_quoted_single
                       '       Token_type_txt_quoted_single
                       f       Token_type_txt_quoted_single
                       '       Token_type_txt_quoted_single
                       g       no_type
                       h       no_type
                       i       no_type                       `

	txt:= str_joined_from_wanted_table_char_column(wantedTable)
	chars := ErlSrcChars_from_str(txt)
	ErlSrcTokens_Quoted__connect_to_chars('\'', chars, true)
	compare_ErlSrcChar_with_wantedTable("ErlSrcTokens_Quoted", chars, wantedTable,  t)
	debug_print_ErlSrcChars(chars)
}

// //////// test tools /////////////

func wanted_table_line_cleaning(line string) string {
	return str_double_space_remove(strings.TrimSpace(line))
}
func str_joined_from_wanted_table_char_column(wantedTable string) string {
	var chars []string
	for _, line := range strings.Split(wantedTable, "\n") {
		// at this point there is a CHAR-TYPE pair with one space only:
		line = wanted_table_line_cleaning(line)
		charOrKey:= strings.Split(line, " ")[0]
		if val, ok := TestGlobals[charOrKey]; ok {
			chars = append(chars, val)
		}
		chars = append(chars, charOrKey)
	}
	return strings.Join(chars, "")
}


func debug_print_ErlSrcChars(chars []ErlSrcChar) {
	fmt.Println("")
	for i, _ := range chars {
		fmt.Printf("%3d posInFile:%3d val:%4v ", i, chars[i].PosInFile, chars[i].Value)

		prevPos := -1
		if chars[i].PrevChar != nil {
			prevPos = chars[i].PrevChar.PosInFile
		}
		fmt.Printf(" PrevPosInFile:%3d ", prevPos)

		tokenType := ""
		if chars[i].Token != nil {
			tokenType = chars[i].Token.Type
		}
		fmt.Printf(" %p <- %p -> %p token: %p %s", chars[i].PrevChar, &chars[i], chars[i].NextChar, chars[i].Token, tokenType)
		fmt.Println("")
	}
}

func compare_ErlSrcChar_with_wantedTable(caller string, chars []ErlSrcChar, wantedTable string,  t *testing.T) {
	wantedTableLines := strings.Split(wantedTable, "\n")
	for charId, charObj := range chars {
		line := wanted_table_line_cleaning(wantedTableLines[charId])
		typeKey := strings.Split(line, " ")[1]
		wantedType, _ := TestGlobals[typeKey]
		compare_str_pair(
			caller+":compare_ErlSrcChar:"+strconv.Itoa(charId),
			charObj.Type(), wantedType, t)
	}
}

func compare_char_pointer_pair(callerInfo string, received, wanted *ErlSrcChar, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

func compare_int_pair(callerInfo string, received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received int: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

func compare_rune_pair(callerInfo string, received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received rune = %v, wanted %v, error", callerInfo, received, wanted)
	}
}

func compare_str_pair(callerInfo, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received string = %s, wanted %s, error", callerInfo, received, wanted)
	}
}

// /////// go experimental tests - guys, I am learning Go too,
// so sometime I do a few language tests :-)

// /// pointer address checks
func Test_what_happens_with_struct_pointers(_ *testing.T) {
	/*  in ErlSrcChars_from_file, chars variable's pointer is similar with
	    the current one here, so at this point there is no copyyin:
		ErlSrcChars_from_file, chars pointer before return: 0xc000158a80
		ErlSrcChars_from_file->Test chars pointer: 0xc000158a80
	*/
	chars, _ := ErlSrcChars_from_file("test/parse/hello.erl")
	fmt.Printf("ErlSrcChars_from_file->Test chars pointer: %p\n", chars)
	_what_happens_with_the_address_simple_obj_pass(chars)
	_what_happens_with_the_address_pointer_pass(&chars)
	/*
			when I pass a simple list, the called fun receives the same object,
		    so the data is not copied if we use a list
	*/
}
func _what_happens_with_the_address_simple_obj_pass(obj []ErlSrcChar) {
	fmt.Printf("ErlSrcChars_from_file, Test chars ojb passed, pointer: %p\n", obj)
}
func _what_happens_with_the_address_pointer_pass(obj *[]ErlSrcChar) {
	fmt.Printf("ErlSrcChars_from_file, Test chars PTR passed, pointer: %p\n", *obj)
}

///// pointer address checks

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
