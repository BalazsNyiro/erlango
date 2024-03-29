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
	"Token_type_comma"               : Token_type_comma,
	"Token_type_dot"                 : Token_type_dot,
	"Token_type_semicolon"           : Token_type_semicolon,
	"Token_type_colon"               : Token_type_colon,

	"Token_type_bracket_round_open"  : Token_type_bracket_round_open,
	"Token_type_bracket_round_close" : Token_type_bracket_round_close,
	"Token_type_bracket_square_open" : Token_type_bracket_square_open,
	"Token_type_bracket_square_close": Token_type_bracket_square_close,
	"Token_type_bracket_curly_open"  : Token_type_bracket_curly_open,
	"Token_type_bracket_curly_close" : Token_type_bracket_curly_close,
	"Token_type_bracket_map_open"    : Token_type_bracket_map_open,

	"Token_type_digits_base10_form"  : Token_type_digits_base10_form,
	"Token_type_digits_baseDefined"  : Token_type_digits_baseDefined,
	"Token_type_float_dotInDigits"   : Token_type_float_dotInDigits,

	"Token_type_variable"            : Token_type_variable,
	"Token_type_atom_quoteless"      : Token_type_atom_quoteless,


	"Token_type_arrow_singleToRight" : Token_type_arrow_singleToRight,
	"Token_type_arrow_singleToLeft"  : Token_type_arrow_singleToLeft,
	"Token_type_arrow_doubleToRight" : Token_type_arrow_singleToRight,

	"Token_type_binding_matching"    : Token_type_binding_matching,

	"Token_type_math_binary_add"     : Token_type_math_binary_add,
	"Token_type_math_binary_sub"     : Token_type_math_binary_sub,
	"Token_type_math_binary_mul"     : Token_type_math_binary_mul,
	"Token_type_math_binary_div"     : Token_type_math_binary_div,

	"Token_type_always_accepted"     : "Token_type_always_accepted",
	"Token_type_deleted_dont_use"    : Token_type_deleted_dont_use,

	// these are important to describe a char that you can't write
	// in a wantedChar table
	"empty_string"                   : "",
	"space"                          : " ",
	"tabulator"                      : "\t",
	"newline_unix"                   : "\n",
	"carriage_return"                : "\r",
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
	prg := &Prg{callStackDisplay: true}
	ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(prg, chars, false)
	// debug_print_ErlSrcChars(chars)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars, wantedTable1, t)

	// token checking
	// debug_print_ErlSrcChars(chars)
	compare_tokenPointers___are_______nil("test1 token check A", &chars, []int{0, 1, 2, 10, 11, 12}, t)
	compare_tokenPointers___are______same("test1 token check B", &chars, []int{3, 4, 5, 6, 7, 8, 9}, t)

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

	txt2 := str_joined_from_wantedCharsTable_char_column(wantedTable2)
	chars2 := ErlSrcChars_from_str(txt2)
	ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(prg, chars2, false)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars2, wantedTable2, t)
	// debug_print_ErlSrcChars(chars2)

	compare_tokenPointers___are_______nil("test2 token check AA1", &chars2, []int{0}, t)
	compare_tokenPointers___are__not__nil("test2 token check BB1", &chars2, []int{1}, t)
	compare_tokenPointers___are______same("test2 token check BB2", &chars2, []int{1, 2, 3, 4, 5, 6, 7}, t)
	compare_tokenPointers___are_______nil("test2 token check CC1", &chars2, []int{8}, t)
	compare_tokenPointers___are______same("test2 token check DD1", &chars2, []int{9, 10, 11, 12}, t)
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
	ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(prg, chars3, false)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Quoted", chars3, wantedTable3, t)
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

	prg := &Prg{callStackDisplay: true}
	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars1 := ErlSrcChars_from_str(srcFromChars1)
	ParseErlangSourceCode(prg, chars1, "strings_atoms_quotes,comments")
	// debug_print_ErlSrcChars(chars1)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Comments_1", chars1, wantedCharsTable1, t)

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
	ParseErlangSourceCode(prg, chars2, "strings_atoms_quotes,comments")
	// debug_print_ErlSrcChars(chars2)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_Comments_2", chars2, wantedCharTable2, t)
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
                            "       Token_type_txt_quoted_double
                            %       Token_type_txt_quoted_double   <- comment sign in a string
	                      space     Token_type_txt_quoted_double   
                            s       Token_type_txt_quoted_double
                            t       Token_type_txt_quoted_double
                            r       Token_type_txt_quoted_double
                       newline_unix Token_type_txt_quoted_double
                            "       Token_type_txt_quoted_double
                            '       Token_type_txt_quoted_single
                            %       Token_type_txt_quoted_single   <- comment sign in an atom
                            a       Token_type_txt_quoted_single
                            t       Token_type_txt_quoted_single
                            o       Token_type_txt_quoted_single
                            m       Token_type_txt_quoted_single
                            '       Token_type_txt_quoted_single
                            %       Token_type_comment
	                      space     Token_type_comment
                            n       Token_type_comment
                            o       Token_type_comment
                            t       Token_type_comment
                            e       Token_type_comment  <- no newline at the end of the comment
    `

	prg := &Prg{callStackDisplay: true}
	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars1 := ErlSrcChars_from_str(srcFromChars1)
	ParseErlangSourceCode(prg, chars1, "strings_atoms_quotes,comments,whitespaces")
	// debug_print_ErlSrcChars(chars1)
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_whitespace_naive", chars1, wantedCharsTable1, t)

	wantedCharTable := wantedCharsTable_from_src_file("test/parse/erlang_whitespaces_separators.erl", 2, 999)
	srcWithoutTestdata := str_joined_from_wantedCharsTable_char_column(wantedCharTable)

	fmt.Println("=============== 1 src ===================")
	fmt.Println(srcWithoutTestdata)
	fmt.Println("=============== 2 - chars ===============")
	chars := ErlSrcChars_from_str(srcWithoutTestdata)
	fmt.Println("=============== 3 - parse ===============")
	ParseErlangSourceCode(prg, chars, "strings_atoms_quotes,comments,whitespaces,commas,dots,semicolons,bracket_round_opener,bracket_round_closer,digits_base10_form")
	debug_print_ErlSrcChars(prg, chars)
	fmt.Println(" wantedCharTable:\n", wantedCharTable)
	fmt.Println("=============== 4 - compare ===============")
	compare_ErlSrcChar_with_wantedCharsTable("Test_ErlSrcTokens_whitespaces_separators", chars, wantedCharTable, t)
}

func Test_ErlSrcTokens_numbers_variables(t *testing.T) {
	fmt.Println(">>> Test_ErlSrcTokens_numbers_variables")

	wantedCharsTable1 := `  "       Token_type_txt_quoted_double
                            s       Token_type_txt_quoted_double
                            t       Token_type_txt_quoted_double
                            r       Token_type_txt_quoted_double
                            "       Token_type_txt_quoted_double
                            N       Token_type_variable
                            u       Token_type_variable
                            m       Token_type_variable
	                      space     Token_type_not_detected
	                        =       Token_type_binding_matching
	                      space     Token_type_not_detected
	                        1       Token_type_digits_base10_form
	                        2       Token_type_digits_base10_form
	                        3       Token_type_digits_base10_form
	                        4       Token_type_digits_base10_form
	                        ,       Token_type_comma
                            A       Token_type_variable                <- variable, 1 char long
	                        =       Token_type_binding_matching
	                        5       Token_type_digits_base10_form
	                        ,       Token_type_comma
                            A       Token_type_variable                <- variable, 1 char long
                            T       Token_type_variable                <- variable, 1 char long
                            O       Token_type_variable                <- variable, 1 char long
                            M       Token_type_variable                <- variable, 1 char long
	                        =       Token_type_binding_matching
	                        a       Token_type_atom_quoteless
	                        ,       Token_type_comma
                            B       Token_type_variable                <- variable, 1 char long
	                        =       Token_type_binding_matching
	                        b       Token_type_atom_quoteless
	                        ,       Token_type_comma
                            X       Token_type_variable                <- variable, 1 char long
	                        =       Token_type_binding_matching
	                        5       Token_type_digits_base10_form
	                        +       Token_type_math_binary_add
	                        6       Token_type_digits_base10_form
	                        -       Token_type_math_binary_sub
	                        7       Token_type_digits_base10_form
	                        *       Token_type_math_binary_mul
	                        8       Token_type_digits_base10_form
	                        /       Token_type_math_binary_div
	                        9       Token_type_digits_base10_form
    `

	prg := &Prg{callStackDisplay: true}
	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars := ErlSrcChars_from_str(srcFromChars1)
	chars, _ = ParseErlangSourceCode(prg, chars, "strings_atoms_quotes,digits_base10_form,variables,atoms_quoteless,commas,binding_matching,math_binary_add,math_binary_sub,math_binary_mul,math_binary_div")
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_numbers_naive", chars, wantedCharsTable1, t)
	// debug_print_ErlSrcChars(chars1)
	compare_str_pair("ErlSrcTokens_numbers_naive", chars[12].Token.StrValueFromChars(), "1234", t)
}

func Test_ErlSrcTokens_arrows_floats_brackets(t *testing.T) {
	fmt.Println(">>> Test_ErlSrcTokens_arrows_floats_brackets")

	wantedCharsTable1 := `          =       Token_type_arrow_doubleToRight      <- from this point it is not valid Erlang code,
			                        >       Token_type_arrow_doubleToRight      <- but from token detection point it's fine
			                        ,       Token_type_comma
			                        -       Token_type_arrow_singleToRight
			                        >       Token_type_arrow_singleToRight
			                        ,       Token_type_comma
			                        <       Token_type_arrow_singleToLeft
			                        -       Token_type_arrow_singleToLeft
			                        ,       Token_type_comma
		                            X       Token_type_variable                <- variable, 1 char long
			                        =       Token_type_binding_matching
			                        5       Token_type_float_dotInDigits
			                        .       Token_type_float_dotInDigits
			                        6       Token_type_float_dotInDigits
			                        ,       Token_type_comma
		                            Y       Token_type_variable                <- variable, 1 char long
			                        =       Token_type_binding_matching
			                        5       Token_type_digits_base10_form
			                        ,       Token_type_comma
			                        #       Token_type_bracket_map_open      <- be careful: MAP OPEN, not curly_brace
			                        {       Token_type_bracket_map_open
			                        ,       Token_type_comma
			                        {       Token_type_bracket_curly_open       <- bracket detections
			                        }       Token_type_bracket_curly_close
			                        [       Token_type_bracket_square_open
			                        ]       Token_type_bracket_square_close
			                        :       Token_type_colon

		    `
	prg := &Prg{callStackDisplay: true}

	srcFromChars1 := str_joined_from_wantedCharsTable_char_column(wantedCharsTable1)
	chars := ErlSrcChars_from_str(srcFromChars1)
	chars, _ = ParseErlangSourceCode(prg, chars, "__all__")
	compare_ErlSrcChar_with_wantedCharsTable("ErlSrcTokens_arrows_floats_brackets", chars, wantedCharsTable1, t)
}

// //////// test tools /////////////

// There is a normal source code + TEST DATA in the src file.
// This fun builds a wantedCharsTable from the src
// lineRange: in escripts, the first line is for Bash - don't process.
// lineRange 2: if you want to focus on selected lines only
func wantedCharsTable_from_src_file(filePath string, lineRangeStart1based, lineRangeEnd1based int) string {
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
		return "no detected token type in TEST LINE!"
	}

	var charsAndWantedTestResult []string
	testData := map[int]string{}

	lines, _ := file_read_lines(filePath, funcName)
	for lineNumInSrc0based, line := range lines { // %INFO: skip these lines too, not real comments.
		lineNum1basedThatYouSeeInAnEditor := lineNumInSrc0based + 1
		if lineNum1basedThatYouSeeInAnEditor < lineRangeStart1based || lineNum1basedThatYouSeeInAnEditor > lineRangeEnd1based || strings.Contains(line, "%INFO") {
			continue
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
			fmt.Println("<< testData, positionInLine >>", line)
			map_print_keysorted__int_str(testData)


			// convert the line to list of chars and append the wanted test results
			postfix := "yyyy"
			wantedTestResult := "xxx"
			prefix := "                 "
			bigSpace := "                                          "

			for posInLine, char := range line {
				if realWantedTestResult, ok := testData[posInLine]; ok {
					wantedTestResult = strings.TrimSpace(realWantedTestResult)
				} else {
					wantedTestResult = "Token_type_always_accepted"
				}
				postfix = bigSpace[:len(bigSpace) - len(wantedTestResult)]

				insertedStr := string(char)

				if char == ' '  { insertedStr = "          space"}
				if char == '\t' { insertedStr = "      tabulator"}
				if char == '\n' { insertedStr = "   newline_unix"}
				if char == '\r' { insertedStr = "carriage_return"}
				// if char == '"' { insertedStr = "              \""; wantedTestResult = "Token_type_txt_quoted_double" }
				prefix= bigSpace[:len(bigSpace) - len(insertedStr)]

				// all linenum and position is 0 based
				charsAndWantedTestResult = append(charsAndWantedTestResult,
					prefix + insertedStr + postfix + wantedTestResult +
					"   line: " + strconv.Itoa(lineNum1basedThatYouSeeInAnEditor) + " pos:" + strconv.Itoa(posInLine))
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
		// fmt.Println("DEBUG line compare 1:",wantedTableLines[charId])
		line := wantedCharsTable_line_cleaning(wantedTableLines[charId])
		// fmt.Println("DEBUG line compare 2:", line)
		typeKey := strings.Split(line, " ")[1]
		// fmt.Println("DEBUG      typeKey:", typeKey)
		wantedType, keyExists := TestGlobals[typeKey]

		if ! keyExists {
			print("ERROR: ", typeKey, " not in TestGlobals")
		}

		if ! strings.Contains(wantedType, "Token_type_always_accepted") {
			compare_str_pair(
				caller+":compare_ErlSrcChar:"+strconv.Itoa(charId)+"->" + string(charObj.Value)+"<-",
				charObj.Type(), wantedType, t)
		}
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
