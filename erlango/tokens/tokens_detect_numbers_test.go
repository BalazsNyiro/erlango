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
