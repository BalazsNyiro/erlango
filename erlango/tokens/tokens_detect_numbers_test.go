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
	"strings"
	"testing"
)

func Test_parse_numbers_int_simple(t *testing.T) {
	testName := "parse_numbers 01, simple int detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C = 1, 22, 345.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C =  ,   ,    .`

	erlSrc_received_after_tokenDetect, tokensTable_02_intDetected := Tokens_1_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table " + testName, tokensTable_02_intDetected)


	fmt.Println("after simple int detection, erlSrc: ", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_intDetected)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_string_string(testName, "1", tokensTable_02_intDetected[10].stringRepr(), t)
	compare_string_string(testName, "22", tokensTable_02_intDetected[13].stringRepr(), t)
	compare_string_string(testName, "345", tokensTable_02_intDetected[17].stringRepr(), t)

}

func Test_parse_numbers_int_with_underscore(t *testing.T) {
	testName := "parse_numbers 02, int and underscore detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C = 1_1, 2_2_2, 34_5.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C =    ,      ,     .`

	erlSrc_received_after_tokenDetect, tokensTable_02_intDetected := Tokens_1_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table: " + testName, tokensTable_02_intDetected)


	fmt.Println("after int detection, erlSrc", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_intDetected)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_string_string(testName, "1_1", tokensTable_02_intDetected[10].stringRepr(), t)
	compare_string_string(testName, "2_2_2", tokensTable_02_intDetected[15].stringRepr(), t)
	compare_string_string(testName, "34_5", tokensTable_02_intDetected[22].stringRepr(), t)

	fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)
}


func Test_parse_numbers_floats(t *testing.T) {
	testName := "parse_numbers 03, float detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C, D = 11.2_2, 33.44, 5_5.66, 7_7.8_8.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C, D =       ,      ,       ,        .`

	erlSrc_received_after_tokenDetect, tokensTable_02_detected := Tokens_1_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table: " + testName, tokensTable_02_detected)


	fmt.Println("after float detection, erlSrc:", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_detected)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_string_string(testName, "11.2_2", tokensTable_02_detected[13].stringRepr(), t)
	compare_string_string(testName, "33.44", tokensTable_02_detected[21].stringRepr(), t)
	compare_string_string(testName, "5_5.66", tokensTable_02_detected[28].stringRepr(), t)
	compare_string_string(testName, "7_7.8_8", tokensTable_02_detected[36].stringRepr(), t)

	fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)
}



func Test_parse_numbers_hexa_nondecimal(t *testing.T) {
	testName := "parse_numbers 04, hexa/nondecimal detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C, D = 16#4f, 1_6#4f, 1_6#4_f, 16#4_f.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C, D =      ,       ,        ,       .`

	erlSrc_received_after_tokenDetect, tokensTable_02_detected := Tokens_1_detect_numbers(erlSrcOrig, tokensTable)
	// print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	// print_tokens("tokens table: " + testName, tokensTable_02_detected)


	// fmt.Println("after hexa/nondecimal detection, erlSrc:", erlSrc_received_after_tokenDetect)

	// fmt.Println("tokensTableOriginal:", tokensTable)
	// fmt.Println("tokensTableUpdated:", tokensTable_02_detected)

	compare_string_string(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_string_string(testName, "16#4f", tokensTable_02_detected[13].stringRepr(), t)
	compare_string_string(testName, "1_6#4f", tokensTable_02_detected[20].stringRepr(), t)
	compare_string_string(testName, "1_6#4_f", tokensTable_02_detected[28].stringRepr(), t)
	compare_string_string(testName, "16#4_f", tokensTable_02_detected[37].stringRepr(), t)

	// fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)
}


// there are preparation steps for testing - this is a wrapper func to simplify tests

//  go test -v -run tokens_detectNumbers_simpleTest
func tokens_detectNumbers_simpleTest(erlExpression, tokenTypeWanted string, t *testing.T) {
	funName := "tokens_detectNumbers_simpleTest"
	tokensTable := Tokens{}
	erlSrcTokenRemoved , tokensTable_detected := Tokens_1_detect_numbers(erlExpression, tokensTable)
	// in the tests, the token is checked only now,
	// so the cleaned src is not used now.
	_ = erlSrcTokenRemoved


	// TYPE checking: typeWanted VS typeDetected
	fmt.Println("\n\nnumTokenDetect:", erlExpression, " <- in Erlang source code")
	compare_string_string(funName + ": " + erlExpression, tokenTypeWanted, tokensTable_detected[0].tokenType, t)



	////////////////////// what is Erlang output if expression is executed? ////////////
	erlangOutput, erlErr := erlBinExpressionParse(erlExpression)
	erlangOutput = strings.TrimSpace(erlangOutput) // whitespace is removed
	fmt.Println("ERLANG output:", erlangOutput, erlErr)


	////////////////////// what is bigNum from the token? //////////////////////////////
	bigNumFromToken, errBn := bigNum_from_token(tokensTable_detected[0])
	if errBn != nil {
		errMsg := fmt.Sprintf("bignum is not detected in token: %s", errBn.Error())
		fmt.Println(errMsg)
		// maybe this is not a problem, if we wanted to test an incorrect token
		// t.Fatalf(errMsg)
	} else {
		fmt.Println("bigNumFromToken detected from token: ", bigNumFromToken)
	}

	//////// compare erlang and bigNumFromToken values ///////////////
	if erlErr == nil {

		// we have a string based Erlang output here. Accept as an ETALON
		bigNum_as_string := bigNumFromToken.stringRepresentation()

		if bigNum_as_string != erlangOutput {
			t.Fatalf("ERROR: different num representations: erlang %s <>  %s bigNum", erlangOutput, bigNum_as_string)
		}
	} // erl error == nil, compare the num with bigNumFromToken


	if erlErr != nil { // error happend in erlang binary
		if tokenTypeWanted == tokenType_SyntaxError {
			fmt.Println("OK! error detected in erlang binary, and in erlango parser, too")
			// fmt.Println("erlang error: ", erlErr)
		} else {
			t.Fatalf("NUM DETECTION PROBLEM: (%s)   error detected in erlang binary, but not in erlango parser", erlExpression)
		}
	} // erlang error happened

}

//  go test -v -run  Test_anyNumSystem_charsSelectScientificPart
func Test_anyNumSystem_charsSelectScientificPart(t *testing.T) {
	testName := "Test_anyNumSystem_charsSelectScientificPart_"

	txt := "123"
	scientificEsignDetected, beforeScientificPart, scientificPart, splitter := anyNumSystem_charsSelectScientificPart([]rune(txt))
	fmt.Println("DEBUG:>>>", scientificEsignDetected, beforeScientificPart, scientificPart, splitter)
	compare_bool_bool(testName+txt, false, scientificEsignDetected, t)
	compare_runes_runes(testName+txt, []rune{}, scientificPart, t)
	compare_runes_runes(testName+txt, []rune(""), beforeScientificPart, t)
	compare_string_string(testName+txt, "", splitter, t)

	txt = "123e+45"
	scientificEsignDetected, beforeScientificPart, scientificPart, splitter = anyNumSystem_charsSelectScientificPart([]rune(txt))
	compare_bool_bool(testName+txt, true, scientificEsignDetected, t)
	compare_runes_runes(testName+txt, []rune("45"), scientificPart, t)
	compare_runes_runes(testName+txt, []rune("123"), beforeScientificPart, t)
	compare_string_string(testName+txt, "e+", splitter, t)

	txt = "123e-45"
	scientificEsignDetected, beforeScientificPart, scientificPart, splitter = anyNumSystem_charsSelectScientificPart([]rune(txt))
	compare_bool_bool(testName+txt, true, scientificEsignDetected, t)
	compare_runes_runes(testName+txt, []rune("45"), scientificPart, t)
	compare_runes_runes(testName+txt, []rune("123"), beforeScientificPart, t)
	compare_string_string(testName+txt, "e-", splitter, t)

	txt = "1_6_7#4eE+89"
	scientificEsignDetected, beforeScientificPart, scientificPart, splitter = anyNumSystem_charsSelectScientificPart([]rune(txt))
	compare_bool_bool(testName+txt, true, scientificEsignDetected, t)
	compare_runes_runes(testName+txt, []rune("89"), scientificPart, t)
	compare_string_string(testName+txt, "e+", splitter, t) // E- or e- return with e- as type, that is not case sensitive

	// here I don't care about numberSystemPart, so everything before the scientific has to be returned
	compare_runes_runes(testName+txt, []rune("1_6_7#4e"), beforeScientificPart, t)

	txt = "1_6_7#4eE-89"
	scientificEsignDetected, beforeScientificPart, scientificPart, splitter = anyNumSystem_charsSelectScientificPart([]rune(txt))
	compare_bool_bool(testName+txt, true, scientificEsignDetected, t)
	compare_runes_runes(testName+txt, []rune("89"), scientificPart, t)
	compare_runes_runes(testName+txt, []rune("1_6_7#4e"), beforeScientificPart, t)
	compare_string_string(testName+txt, "e-", splitter, t)  // E- or e- return with e- as type, that is not case sensitive

}



//  go test -v -run Test_anyNumSystem_detectNumSystem
func Test_anyNumSystem_detectNumSystem(t *testing.T) {
	testName := "Test_anyNumSystem_charsSelectScientificPart_"

	txt := "1_6_7#4ee+89"
	numberSystemType, isHashMarkDetected, charsAfterHashMark, numberSystemDetectionError := anyNumSystem_detectNumSystem([]rune(txt))
	fmt.Println("numberSystemTypeDetected in the test:", numberSystemType)
	compare_bool_bool(testName+txt, true, numberSystemDetectionError==nil, t)
	compare_bigNum_bigNum(testName+txt, bigNum_create_from_int(167), numberSystemType, t)
	compare_bool_bool(testName+txt, true, isHashMarkDetected, t)
	compare_runes_runes(testName+txt, []rune("4ee+89"), charsAfterHashMark, t)

	txt = "2a#4ee+789"
	numberSystemType, isHashMarkDetected, charsAfterHashMark, numberSystemDetectionError = anyNumSystem_detectNumSystem([]rune(txt))
	// error has to be detected
	compare_bool_bool(testName+txt, false, numberSystemDetectionError==nil, t)
	compare_bool_bool(testName+txt, true, isHashMarkDetected, t)
	compare_runes_runes(testName+txt, []rune("4ee+789"), charsAfterHashMark, t)

}



///////////////////////////////////////////////
//  go test -v -run Test_charsCopySplitWithChars
func Test_charsCopySplitAtFirstWithChars(t *testing.T) {
	testName := "Test_charsCopySplitAtFirstWithChars "

	txt := "abc-012"
	charsOrig := []rune(txt)
	_, leftDetected, rightDetected := charsCopySplitAtFirstWithChars(charsOrig, []rune("-"))

	leftWanted := []rune("abc")
	rightWanted := []rune("012")
	compare_runes_runes(testName+txt, leftWanted, leftDetected, t)
	compare_runes_runes(testName+txt, rightWanted, rightDetected, t)


	txt = "abc-012-345"
	charsOrig = []rune(txt)
	_, leftDetected, rightDetected = charsCopySplitAtFirstWithChars(charsOrig, []rune("-"))

	leftWanted = []rune("abc")
	rightWanted = []rune("012-345")
	compare_runes_runes(testName+txt, leftWanted, leftDetected, t)
	compare_runes_runes(testName+txt, rightWanted, rightDetected, t)


	txt = "abc-012-345"
	charsOrig = []rune(txt)
	_, leftDetected, rightDetected = charsCopySplitAtFirstWithChars(charsOrig, []rune("01"))

	leftWanted = []rune("abc-")
	rightWanted = []rune("2-345")
	compare_runes_runes(testName+txt, leftWanted, leftDetected, t)
	compare_runes_runes(testName+txt, rightWanted, rightDetected, t)


	txt = "abc-012-012-345"
	charsOrig = []rune(txt)
	_, leftDetected, rightDetected = charsCopySplitAtFirstWithChars(charsOrig, []rune("01"))

	leftWanted = []rune("abc-")
	rightWanted = []rune("2-012-345")
	compare_runes_runes(testName+txt, leftWanted, leftDetected, t)
	compare_runes_runes(testName+txt, rightWanted, rightDetected, t)
}

//  go test -v -run Test_mass_number_detection
func Test_mass_number_detection(t *testing.T) {
	// this test calls ERLANG BINARY to check the values, too!
	tokens_detectNumbers_simpleTest(`16         `, tokenType_Num_int, t)
	tokens_detectNumbers_simpleTest(`1_6        `, tokenType_Num_int, t)
	tokens_detectNumbers_simpleTest(`1_6_7      `, tokenType_Num_int, t)

	tokens_detectNumbers_simpleTest(`1__6      `, tokenType_SyntaxError, t)
	tokens_detectNumbers_simpleTest(`1_6_      `, tokenType_SyntaxError, t)
	tokens_detectNumbers_simpleTest(`1_6__     `, tokenType_SyntaxError, t)

	tokens_detectNumbers_simpleTest(`16#4f   `, tokenType_Num_maybeNonDecimal, t)
	tokens_detectNumbers_simpleTest(`1_6#4f  `, tokenType_Num_maybeNonDecimal, t)
	tokens_detectNumbers_simpleTest(`1_6#4_f `, tokenType_Num_maybeNonDecimal, t)

	tokens_detectNumbers_simpleTest(`1_6_#4_f`, tokenType_SyntaxError, t)
	tokens_detectNumbers_simpleTest(`1__6#4_f`, tokenType_SyntaxError, t)

	tokens_detectNumbers_simpleTest(`1.6`, tokenType_Num_float, t)
}


//  go test -v -run Test_index_reverse
func Test_index_reverse(t *testing.T) {
	testName := "Test_index_reverse"

	text := []rune("abcdefghijkl")

	indexLast := len(text)-1

	index := 0  // the reverse index is (LEN-1)
	indexReverse, _ := indexReverse_get__worksWithNonEmptySlicesOnly(len(text), index)
	compare_int_int(testName, indexLast, indexReverse, t)

	index = 1
	indexReverse, _ = indexReverse_get__worksWithNonEmptySlicesOnly(len(text), index)
	compare_int_int(testName, indexLast-index, indexReverse, t)

	index = 4
	indexReverse, _ = indexReverse_get__worksWithNonEmptySlicesOnly(len(text), index)
	compare_int_int(testName, indexLast-index, indexReverse, t)

	index = indexLast
	indexReverse, _ = indexReverse_get__worksWithNonEmptySlicesOnly(len(text), index)
	compare_int_int(testName, indexLast-index, indexReverse, t)
}

//  go test -v -run Test_bigNum_string_representation
func Test_bigNum_string_representation(t *testing.T) {
	testName := "Test_bigNum_string_representation"

	bn := bigNum_ten()
	representation := bn.stringRepresentation()

	compare_string_string(testName, "10", representation, t)

	bn = bigNum_operator_add(bn, bigNum_one())
	bn.exponent = -1
	representation = bn.stringRepresentation()
	compare_string_string(testName, "1.1", representation, t)

	bn.exponent = -2
	representation = bn.stringRepresentation()
	compare_string_string(testName, "0.11", representation, t)

	bn.exponent = -3
	representation = bn.stringRepresentation()
	compare_string_string(testName, "0.011", representation, t)

	bn.exponent = -4
	representation = bn.stringRepresentation()
	compare_string_string(testName, "0.0011", representation, t)
}


