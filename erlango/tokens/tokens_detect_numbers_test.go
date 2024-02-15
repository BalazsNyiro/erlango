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
	"strconv"
	"testing"
)

func Test_parse_numbers_int_simple(t *testing.T) {
	testName := "parse_numbers 01, simple int detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C = 1, 22, 345.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C =  ,   ,    .`

	erlSrc_received_after_tokenDetect, tokensTable_02_intDetected := Tokens_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table " + testName, tokensTable_02_intDetected)


	fmt.Println("after simple int detection, erlSrc: ", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_intDetected)

	compare_strings(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_strings(testName, "1", tokensTable_02_intDetected[10].stringRepr(), t)
	compare_strings(testName, "22", tokensTable_02_intDetected[13].stringRepr(), t)
	compare_strings(testName, "345", tokensTable_02_intDetected[17].stringRepr(), t)

}

func Test_parse_numbers_int_with_underscore(t *testing.T) {
	testName := "parse_numbers 02, int and underscore detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C = 1_1, 2_2_2, 34_5.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C =    ,      ,     .`

	erlSrc_received_after_tokenDetect, tokensTable_02_intDetected := Tokens_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table: " + testName, tokensTable_02_intDetected)


	fmt.Println("after int detection, erlSrc", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_intDetected)

	compare_strings(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_strings(testName, "1_1", tokensTable_02_intDetected[10].stringRepr(), t)
	compare_strings(testName, "2_2_2", tokensTable_02_intDetected[15].stringRepr(), t)
	compare_strings(testName, "34_5", tokensTable_02_intDetected[22].stringRepr(), t)

	fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)
}


func Test_parse_numbers_floats(t *testing.T) {
	testName := "parse_numbers 03, float detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C, D = 11.2_2, 33.44, 5_5.66, 7_7.8_8.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C, D =       ,      ,       ,        .`

	erlSrc_received_after_tokenDetect, tokensTable_02_detected := Tokens_detect_numbers(erlSrcOrig, tokensTable)
	print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	print_tokens("tokens table: " + testName, tokensTable_02_detected)


	fmt.Println("after float detection, erlSrc:", erlSrc_received_after_tokenDetect)

	fmt.Println("tokensTableOriginal:", tokensTable)
	fmt.Println("tokensTableUpdated:", tokensTable_02_detected)

	compare_strings(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_strings(testName, "11.2_2", tokensTable_02_detected[13].stringRepr(), t)
	compare_strings(testName, "33.44", tokensTable_02_detected[21].stringRepr(), t)
	compare_strings(testName, "5_5.66", tokensTable_02_detected[28].stringRepr(), t)
	compare_strings(testName, "7_7.8_8", tokensTable_02_detected[36].stringRepr(), t)

	fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)
}



func Test_parse_numbers_hexa_nondecimal(t *testing.T) {
	testName := "parse_numbers 04, hexa/nondecimal detection"

	tokensTable := Tokens{}
	erlSrcOrig :=                    `A, B, C, D = 16#4f, 1_6#4f, 1_6#4_f, 16#4_f.`
	erlSrcWantedAfterTokenDetect :=  `A, B, C, D =      ,       ,        ,       .`

	erlSrc_received_after_tokenDetect, tokensTable_02_detected := Tokens_detect_numbers(erlSrcOrig, tokensTable)
	// print_string_runes_diff(erlSrcOrig, erlSrc_received_after_tokenDetect)
	// print_tokens("tokens table: " + testName, tokensTable_02_detected)


	// fmt.Println("after hexa/nondecimal detection, erlSrc:", erlSrc_received_after_tokenDetect)

	// fmt.Println("tokensTableOriginal:", tokensTable)
	// fmt.Println("tokensTableUpdated:", tokensTable_02_detected)

	compare_strings(testName, erlSrcWantedAfterTokenDetect, erlSrc_received_after_tokenDetect, t)
	compare_strings(testName, "16#4f", tokensTable_02_detected[13].stringRepr(), t)
	compare_strings(testName, "1_6#4f", tokensTable_02_detected[20].stringRepr(), t)
	compare_strings(testName, "1_6#4_f", tokensTable_02_detected[28].stringRepr(), t)
	compare_strings(testName, "16#4_f", tokensTable_02_detected[37].stringRepr(), t)

	// fmt.Println(testName, erlSrcWantedAfterTokenDetect, erlSrcOrig, tokensTable)

	fmt.Println("FIXME: erlang native hivasok, szamokkal, eredmenyt olvasd vissza")
}


// there are preparation steps for testing - this is a wrapper func to simplify tests
func tokens_detectNumbers_simpleTest(erlExpression, tokenTypeWanted string, t *testing.T) {
	funName := "tokens_detectNumbers_simpleTest"
	tokensTable := Tokens{}
	erlSrcTokenRemoved , tokensTable_detected := Tokens_detect_numbers(erlExpression, tokensTable)
	// in the tests, the token is checked only now,
	// so the cleaned src is not used now.
	_ = erlSrcTokenRemoved


	fmt.Println("\n\nnumTokenDetect:", erlExpression)
	compare_strings(funName + ": " + erlExpression, tokenTypeWanted, tokensTable_detected[0].tokenType, t)




	//////// compare erlang and bigNum values ///////////////
	erlOutDetectedNumString, erlErr := erlBinExpressionParse(erlExpression)

	if erlErr == nil {
		erlOutDetectedNum, errIntConv := strconv.Atoi(erlOutDetectedNumString)
		if errIntConv != nil {
			fmt.Println("int conversion problem from erlang shell:", erlOutDetectedNumString, errIntConv)

		} else {
			fmt.Println("ERLANG DETECTED NUM:", erlOutDetectedNum)

			bigNum, err := bigNum_from_token(tokensTable_detected[0])
			bigNum_INT := 0
			if err != nil {
				fmt.Println("problem with bignum conversion: ", err)
			} else {
				fmt.Println("bigNum: ", bigNum)
				bigNum_INT = bigNum_convert_to_INT_for_testcases(bigNum)
			}
			fmt.Println("bigNum INT:", bigNum_INT, "erl bin out:", erlOutDetectedNumString, "erl error:", erlErr)

			if bigNum_INT != erlOutDetectedNum {
				t.Fatalf("ERROR: erlang integer %d <> %d bigNum int", erlOutDetectedNum, bigNum_INT)
			}

		} // no problem with erlang output's conversion to int
	} // erl error == nil, compare the num with bigNum



	if erlErr != nil { // error happend in erlang binary
		if tokenTypeWanted == tokenType_SyntaxError {
			fmt.Println("OK! error detected in erlang binary, and in erlango parser, too")
			// fmt.Println("erlang error: ", erlErr)
		} else {
			t.Fatalf("NUM DETECTION PROBLEM: (%s)   error detected in erlang binary, but not in erlango parser", erlExpression)
		}
	} // erlang error happened

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
	// tokens_detectNumbers_simpleTest(`1_6#4f  `, tokenType_Num_maybeNonDecimal, t)
	// tokens_detectNumbers_simpleTest(`1_6#4_f `, tokenType_Num_maybeNonDecimal, t)

	// tokens_detectNumbers_simpleTest(`1_6_#4_f`, tokenType_SyntaxError, t)
	// tokens_detectNumbers_simpleTest(`1__6#4_f`, tokenType_SyntaxError, t)
}
