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
	"Token_type_txt_quoted_double"   : Token_type_txt_quoted_double,
	"Token_type_txt_quoted_single"   : Token_type_txt_quoted_single,
	"Token_type_comment"             : Token_type_comment,
	"Token_type_not_detected"        : Token_type_not_detected,
	"Token_type_whitespace"          : Token_type_whitespace,

	// these are important to describe a char that you can't write
	// in a wantedChar table
	"empty_string"                   : "",
	"space"                          : " ",
	"tabulator"                      : "\t",
	"newline_unix"                   : "\n",
} //////////////////////////////////////////////////////////////////////

func Test_ErlSrcRead(t *testing.T) {
	chars, _ := ErlSrcChars_from_file("test/parse/hello.erl")
	compare_rune_pair("val m", chars[1].Value, 'm', t)
	compare_rune_pair("val o", chars[2].Value, 'o', t)
	compare_rune_pair("val d", chars[3].Value, 'd', t)

	compare_rune_pair("val next 1", chars[1].NextChar.Value, 'o', t)
	compare_rune_pair("val prev 3", chars[3].PrevChar.Value, 'o', t)

	compare_int_pair("pos 3", chars[3].PosInFile, 3, t)
	compare_char_pointer_pair_are_same("compare char_0 and nil_prev_char", chars[0].PrevChar, nil, t)

	compare_int_pair("pos_0.next, pos2.prev", chars[0].NextChar.PosInFile, chars[2].PrevChar.PosInFile, t)
	// debug_print_ErlSrcChars(chars)
	// compare_char_pointer_pair_are_same("compare char_0.next and char_2.prev", chars[0].NextChar, chars[2].PrevChar, t)
}

func Test_ErlSrcTokens_Quoted(t *testing.T) {
	// The first column is the char, the second column is the type, others are comments
	// a single char means himself. Keywords has special meanings
	wantedTable1 := `  a       Token_type_not_detected 
                       b       Token_type_not_detected
                       c       Token_type_not_detected
                       '       Token_type_txt_quoted_single
                       d       Token_type_txt_quoted_single
                       "       Token_type_txt_quoted_single     <- this " is in the ''pair  
                       \       Token_type_txt_quoted_single     
                       '       Token_type_txt_quoted_single     <- escaped ' sign, not valid exit
                       f       Token_type_txt_quoted_single
                       '       Token_type_txt_quoted_single     <- valid closing pair
                       g       Token_type_not_detected
                       h       Token_type_not_detected
                       i       Token_type_not_detected                       `

	txt := str_joined_from_wantedCharsTable_char_column(wantedTable1)
	chars := ErlSrcChars_from_str(txt)
	ErlSrcTokensDetect__string_atom__connect_to_chars(chars, true)
	debug_print_ErlSrcChars(chars)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars, wantedTable1,  t)

	// token checking
	debug_print_ErlSrcChars(chars)
	compare_tokenPointers___are_______nil("test1 token check A", &chars, []int{0,1,2,10,11,12}, t)
	compare_tokenPointers___are______same("test1 token check B", &chars, []int{3,4,5,6,7,8,9}, t)

	///////////////////////////////////////////////////////////////////////////



	// here we search the "..." sections only, so the '...' is not detected
	wantedTable2 := `  a       Token_type_not_detected 
                       "       Token_type_txt_quoted_double 
                     space     Token_type_txt_quoted_double
                       '       Token_type_txt_quoted_double     <- embedded ' char in the "" pair
                       \       Token_type_txt_quoted_double
                       "       Token_type_txt_quoted_double     <- escaped " char, not a valid exit
                       1       Token_type_txt_quoted_double
                       "       Token_type_txt_quoted_double     <- valid closing " char
                       f       Token_type_not_detected
                       '       Token_type_txt_quoted_single     <- valid opening ' char
                       i       Token_type_txt_quoted_single
	                   h       Token_type_txt_quoted_single
	                   '       Token_type_txt_quoted_single     <- valid closing ' char
                     `

	txt2:= str_joined_from_wantedCharsTable_char_column(wantedTable2)
	chars2 := ErlSrcChars_from_str(txt2)
	ErlSrcTokensDetect__string_atom__connect_to_chars(chars2, true)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars2, wantedTable2,  t)
	// debug_print_ErlSrcChars(chars2)

	compare_tokenPointers___are_______nil("test2 token check AA1", &chars2, []int{0}, t)
	compare_tokenPointers___are__not__nil("test2 token check BB1", &chars2, []int{1}, t)
	compare_tokenPointers___are______same("test2 token check BB2", &chars2, []int{1,2,3,4,5,6,7}, t)
	compare_tokenPointers___are_______nil("test2 token check CC1", &chars2, []int{8}, t)
	compare_tokenPointers___are______same("test2 token check DD1", &chars2, []int{9,10,11,12}, t)
	compare_tokenPointerPair_is_different("test2 token check DD2", chars2[7].Token, chars2[9].Token, t)


	// mixed test
	wantedTable3 := `  a       Token_type_not_detected 
                       "       Token_type_txt_quoted_double 
                     space     Token_type_txt_quoted_double
                       '       Token_type_txt_quoted_double
                       \       Token_type_txt_quoted_double
                       "       Token_type_txt_quoted_double
                       1       Token_type_txt_quoted_double
                       "       Token_type_txt_quoted_double
                       f       Token_type_not_detected
                       '       Token_type_txt_quoted_single
                       "       Token_type_txt_quoted_single
	                   h       Token_type_txt_quoted_single
                       "       Token_type_txt_quoted_single
	                   '       Token_type_txt_quoted_single `

	txt3 := str_joined_from_wantedCharsTable_char_column(wantedTable3)
	chars3 := ErlSrcChars_from_str(txt3)
	ErlSrcTokensDetect__string_atom__connect_to_chars(chars3, true)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars3, wantedTable3,  t)
	// debug_print_ErlSrcChars(chars3)
}


func Test_ErlSrcTokens_Comments(t *testing.T) {
	wantedCharsTable1 := `  a       Token_type_not_detected 
	                        b       Token_type_not_detected
                            %       Token_type_comment
                            %       Token_type_comment
                            '       Token_type_comment  <- comment detection is AFTER str/atom detect
                            a       Token_type_comment     but: you can have a string in a comment, too!
                            t       Token_type_comment
                            o       Token_type_comment
                            m       Token_type_comment
                            '       Token_type_comment
                            "       Token_type_comment  <- string in the comment
                            s       Token_type_comment
                            t       Token_type_comment
                            r       Token_type_comment
                            "       Token_type_comment
                            n       Token_type_comment
                            o       Token_type_comment
                            t       Token_type_comment
                            e       Token_type_comment
                      newline_unix  Token_type_not_detected  <- newline is the closer of comments
                            t       Token_type_not_detected   
                            x       Token_type_not_detected
                            t       Token_type_not_detected
                            "       Token_type_txt_quoted_double
                            %       Token_type_txt_quoted_double   <- comment sign in a string
                            s       Token_type_txt_quoted_double
                            t       Token_type_txt_quoted_double
                            r       Token_type_txt_quoted_double
                            "       Token_type_txt_quoted_double
                            '       Token_type_txt_quoted_single
                            %       Token_type_txt_quoted_single   <- comment sign in an atom
                            a       Token_type_txt_quoted_single
                            t       Token_type_txt_quoted_single
                            o       Token_type_txt_quoted_single
                            m       Token_type_txt_quoted_single
                            '       Token_type_txt_quoted_single
    `

	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars1 := ErlSrcChars_from_str(srcFromChars1)
	ParseErlangSourceCode(chars1, "strings_atoms,comments")
	debug_print_ErlSrcChars(chars1)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Comments_1", chars1, wantedCharsTable1,  t)


	// basic comment test
	wantedCharTable2 := `  a       Token_type_not_detected 
	                       b       Token_type_not_detected
                           %       Token_type_comment
                           n       Token_type_comment
                           o       Token_type_comment
                           t       Token_type_comment
                           e       Token_type_comment
                     newline_unix  Token_type_not_detected
                           t       Token_type_not_detected   
                           x       Token_type_not_detected
                           t       Token_type_not_detected
    `

	srcFromChars2 := str_joined_from_wantedCharsTable_char_column(wantedCharTable2)
	chars2 := ErlSrcChars_from_str(srcFromChars2)
	ParseErlangSourceCode(chars2, "strings_atoms,comments")
	debug_print_ErlSrcChars(chars2)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Comments_2", chars2, wantedCharTable2,  t)
}

func Test_ErlSrcTokens_whitespaces_separators(t *testing.T) {
	fmt.Println(">>> Test_ErlSrcTokens_whitespaces_separators")

    // naive tests
	wantedCharsTable1 := `  a       Token_type_not_detected 
	                      space     Token_type_whitespace
	                        b       Token_type_not_detected
	                      space     Token_type_whitespace
	                        c       Token_type_not_detected
                        tabulator   Token_type_whitespace
	                        d       Token_type_not_detected
                       newline_unix Token_type_whitespace
	                        d       Token_type_not_detected
    `

	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars1 := ErlSrcChars_from_str(srcFromChars1)
	ParseErlangSourceCode(chars1, "strings_atoms,comments,whitespaces")
	debug_print_ErlSrcChars(chars1)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_whitespace_naive", chars1, wantedCharsTable1,  t)


	/*
	wantedCharTable := wantedCharsTable_from_src_file("test/parse/erlang_whitespaces_separators.erl", 2, 7)
	srcWithoutTestdata := str_joined_from_wantedCharsTable_char_column(wantedCharTable)

	fmt.Println("=============== 1 src ===================")
	fmt.Println(srcWithoutTestdata)
	fmt.Println("=============== 2 - chars ===============")
	chars := ErlSrcChars_from_str(srcWithoutTestdata)
	fmt.Println("=============== 3 - parse ===============")
	ParseErlangSourceCode(chars, "strings_atoms,comments,whitespaces")
	fmt.Println("=============== 4 - compare ===============")
	compare_ErlSrcChar_with_wantedCharsTable("Test_ErlSrcTokens_whitespaces_separators", chars, wantedCharTable,  t)

	 */
}

// //////// test tools /////////////

// There is a normal source code + TEST DATA in the src file.
// This fun builds a wantedCharsTable from the src
// lineRange: in escripts, the first line is for Bash - don't process.
// lineRange 2: if you want to focus on selected lines only
func wantedCharsTable_from_src_file(filePath string, lineRangeStart, lineRangeEnd int) string   {
	// this is the src+test wanted data.
	funcName := "wantedCharsTable_from_src_file"

	isTestLine := func(txt string) bool {
		txt = strings.TrimSpace(txt)
		if strings.Contains(txt, "%") {
			if txt[0] == '%' {
				if strings.Contains(txt, "Token_type_") {
					return true
				}
			}
		}
		return false
	}
	wantedTokenTypeInTestline := func(lineActual string) string {
		for _, word := range strings.Split(lineActual, " ") {
			if strings.Contains(word, "Token_type") {
				return word
			}
		}
		return "wantedTestResult!"
	}

	var charsAndWantedTestResult []string
	testData := map[int]string{}

	lines, _ := file_read_lines(filePath, funcName)
	for lineNumInSrc, line:= range lines { // %INFO: skipt these lines too, not real comments.
		if lineNumInSrc < lineRangeStart-1 || lineNumInSrc > lineRangeEnd-1 || strings.Contains(line, "%INFO"){
			continue // process only the selected line range: -1 because of internal 0 based line numbering
		}
		line = line + string('\n') // restore original line ending
		if isTestLine(line) {
			wantedTestResult := wantedTokenTypeInTestline(line)
			// modify the matching chars' wanted test Result
			// collect the test lines from the source code
			inTokenMatchArea := false
			for positionInLine, runeNow := range line {
				if runeNow == '%' && !inTokenMatchArea { inTokenMatchArea = true }
				if runeNow != '%' && inTokenMatchArea { break }
				if runeNow == '%' && inTokenMatchArea {
					testData[positionInLine] = wantedTestResult
				}
			}
		} else {
			// convert the line to list of chars and append the wanted test results
			for posInLine, char := range line {
				prefix := "                 "
				postfix:= "    "
				wantedTestResult := "Token_type_not_detected"
				if realWantedTestResult, ok := testData[posInLine]; ok {
					wantedTestResult = strings.TrimSpace(realWantedTestResult)
				}

				insertedStr := string(char)
				if char == ' '  { insertedStr = "       space"}
				if char == '\t' { insertedStr = "   tabulator"}
				if char == '\n' { insertedStr = "newline_unix"}
				if insertedStr != string(char) { // so if it's 'space', 'tabulator' or other...
					prefix = "      "
				}
				// all linenum and position is 0 based so they are incremented in the output because in the original sources the editors use 1 based numbering
				charsAndWantedTestResult = append(charsAndWantedTestResult,
					prefix + insertedStr + postfix + wantedTestResult +
					"   line: " + strconv.Itoa(lineNumInSrc+1) + " pos:" + strconv.Itoa(posInLine+1))
			}
			testData = map[int]string{}
		}
	}
	return strings.Join(charsAndWantedTestResult, "\n")
}

/*
   the first column can contain one character, or a keyword, that is translated to a char.
   the second column is the type of a bounded token.

   everything from a possible third column can be a comment.
*/
func wantedCharsTable_line_cleaning(line string) string {
	return str_double_space_remove(strings.TrimSpace(line))
}
func str_joined_from_wantedCharsTable_char_column(wantedTable string) string {
	var chars []string
	for _, line := range strings.Split(wantedTable, "\n") {
		// at this point there is a CHAR-TYPE pair with one space only:
		line = wantedCharsTable_line_cleaning(line)
		charOrKey := strings.Split(line, " ")[0]
		if val, ok := TestGlobals[charOrKey]; ok {
			chars = append(chars, val)  // use the value from the Global table, or:
		} else {
			if len(charOrKey) > 1 {
				panic("ERROR in test - use one char only: " + charOrKey)
			}
			chars = append(chars, charOrKey) // use the original chars
		}
	}
	return strings.Join(chars, "")
}


func compare_ErlSrcChar_with_wantedCharsTable(caller string, chars []ErlSrcChar, wantedTable string,  t *testing.T) {
	wantedTableLines := strings.Split(wantedTable, "\n")
	for charId, charObj := range chars {
		line := wantedCharsTable_line_cleaning(wantedTableLines[charId])
		// fmt.Println("DEBUG line compare:", line)
		typeKey := strings.Split(line, " ")[1]
		// fmt.Println("DEBUG      typeKey:", typeKey)
		wantedType, keyExists:= TestGlobals[typeKey]
		if ! keyExists {
			print("ERROR: ", typeKey, " not in TestGlobals")
		}
		compare_str_pair(
			caller+":compare_ErlSrcChar:"+strconv.Itoa(charId)+"->" + string(charObj.Value)+"<-",
			charObj.Type(), wantedType, t)
	}
}

// the problem: if the received is NOT similar with the wanted
func compare_char_pointer_pair_are_same(callerInfo string, received, wanted *ErlSrcChar, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nErr: %s received: %v\n  wanted: %v, error", callerInfo, received, wanted)
	}
}

// the problem is if received is similar with wanted
func compare_char_pointer_pair_are_different(callerInfo string, received, wanted *ErlSrcChar, t *testing.T) {
	if received == wanted {
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
		t.Fatalf("\nErr: %s received string ->%s<-, wanted ->%s<-, error", callerInfo, received, wanted)
	}
}

func compare_tokenPointer_is_nil(callerInfo string, receivedPtr *ErlSrcToken, t *testing.T) {
	fmt.Println("receivedPtr:", receivedPtr)
	if receivedPtr != nil {
		t.Fatalf("\nErr, PTR is not nil: %s receivedPtr: %v  wanted: nil, error", callerInfo, receivedPtr)
	}
}

func compare_tokenPointer_is_not_nil(callerInfo string, receivedPtr *ErlSrcToken, t *testing.T) {
	fmt.Println("receivedPtr:", receivedPtr)
	if receivedPtr == nil {
		t.Fatalf("\nErr, PTR is nil: %s receivedPtr: %v  wanted: nil, error", callerInfo, receivedPtr)
	}
}

func compare_tokenPointerPair_is__same(callerInfo string, receivedPtr *ErlSrcToken, wantedPtr *ErlSrcToken, t *testing.T) {
	if receivedPtr != wantedPtr {
		t.Fatalf("\nErr, different PTRs: %s receivedPtr: %v  wanted: %v, error", callerInfo, receivedPtr, wantedPtr)
	}
}

func compare_tokenPointerPair_is_different(callerInfo string, receivedPtr *ErlSrcToken, wantedPtr *ErlSrcToken, t *testing.T) {
	if receivedPtr == wantedPtr {
		t.Fatalf("\nErr, same PTRs: %s receivedPtr: %v  wanted: %v, error", callerInfo, receivedPtr, wantedPtr)
	}
}

func compare_tokenPointers___are_______nil(callerInfo string, charsPtr *[]ErlSrcChar, positions []int, t *testing.T) {
	fmt.Println("compare_tokenPointers___are_______nil positions:", positions)
	for _, charPos := range positions {
		fmt.Println("char position:", charPos)
		compare_tokenPointer_is_nil(callerInfo + fmt.Sprintf(" (charId:%d)", charPos), (*charsPtr)[charPos].Token, t)
	}
}

func compare_tokenPointers___are__not__nil(callerInfo string, charsPtr *[]ErlSrcChar, positions []int, t *testing.T) {
	fmt.Println("compare_tokenPointers___are__not__nil positions:", positions)
	for _, charPos := range positions {
		fmt.Println("char position:", charPos, " charValStr:", string((*charsPtr)[charPos].Value),  "  char's token: ", (*charsPtr)[charPos].Token)
		compare_tokenPointer_is_not_nil(callerInfo + fmt.Sprintf(" (charId:%d)", charPos), (*charsPtr)[charPos].Token, t)
	}
}

// a lot of pointers has the same value - it uses the pair comparison
func compare_tokenPointers___are______same(callerInfo string, charsPtr *[]ErlSrcChar, positions []int, t *testing.T) {
	fmt.Println(" compare_tokenPointers___are______same, position 0: ", positions[0])
	wantedTokenPtr := (*charsPtr)[positions[0]].Token  // read the first elem's token - and check that the others have the same?
	fmt.Println(" compare_tokenPointers___are______same, position 0 token value: ", wantedTokenPtr)

	for _, charPos := range positions {
		fmt.Println(" compare_tokenPointers___are______same, checked position: ", charPos)
		compare_tokenPointerPair_is__same(callerInfo + fmt.Sprintf(" (charId:%d)", charPos), (*charsPtr)[charPos].Token, wantedTokenPtr, t)
	}
}


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
// /////// go experimental tests //////////////////////////////////////////////////////////////

