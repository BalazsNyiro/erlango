/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import "fmt"

func Tokens_detect_in_erl_src(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector) {

	charactersInErlSrc, tokensInErlSrc = tokens_detect_01_erlang_strings__quoted_atoms__comments(charactersInErlSrc, tokensInErlSrc)
	charactersInErlSrc, tokensInErlSrc = tokens_detect_02_erlang_whitespaces(charactersInErlSrc, tokensInErlSrc)

	// TODO: whitespace detection, operator detection.
	// atom detection
	// variable detection
	// operators detection?
	// comma, dot, :,
	return charactersInErlSrc, tokensInErlSrc
}

func Tokens_detection_print_verbose(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) {

	lineNumInErlSrc := 0

	type reportLine []rune
	type reportLineMore []reportLine

	// one erland line -> more lines are printed with token type infos
	printedReportLinesForOneErlangSrcLine := 5 //  for every erlang line, X report lines are printed
	reportLine_4_separator := []rune("============================")
	reportLine_3_char_counter := reportLine{}
	reportLine_2_opener_closer := reportLine{}
	reportLine_1_token_type := reportLine{}
	reportLine_0_erl_src_chars := reportLine{}

	reportLines := reportLineMore{}

	for charCounter, charInErlSrc := range charactersInErlSrc {

		openerCloserStatus := ' '
		if charInErlSrc.tokenOpenerCharacter {
			openerCloserStatus = 'o'
		}
		if charInErlSrc.tokenCloserCharacter {
			openerCloserStatus = 'c'
		}
		if charInErlSrc.tokenOpenerCharacter && charInErlSrc.tokenCloserCharacter {
			openerCloserStatus = '2' // closer AND opener same time
		}
		reportLine_2_opener_closer = append(reportLine_2_opener_closer, openerCloserStatus)

		oneCharWideTokenTypeRepresentation := TokenTypeReprShort(charInErlSrc.tokenDetectedType)
		reportLine_1_token_type = append(reportLine_1_token_type, oneCharWideTokenTypeRepresentation)

		runePrinted := charInErlSrc.runeInErlSrc // this will be printed/displayed,
		if charInErlSrc.runeInErlSrc == '\t' {   //and long tabulators needs to be replaced
			runePrinted = ' '
		}
		if charInErlSrc.runeInErlSrc == '\n' { // newline representation
			runePrinted = unicode_stop_table
		}

		reportLine_0_erl_src_chars = append(reportLine_0_erl_src_chars, runePrinted)
		reportLine_3_char_counter = append(reportLine_3_char_counter, []rune(fmt.Sprintf("%d", charCounter%10))[0])

		// fmt.Println("charCounter: ", charCounter)
		if charInErlSrc.runeInErlSrc == '\n' { // newline chars

			reportLines = append(reportLines, reportLine_4_separator)
			reportLines = append(reportLines, reportLine_3_char_counter)
			reportLines = append(reportLines, reportLine_2_opener_closer)
			reportLines = append(reportLines, reportLine_1_token_type)
			reportLines = append(reportLines, reportLine_0_erl_src_chars)

			lineNumInErlSrc += 1
			reportLine_3_char_counter = reportLine{}
			reportLine_2_opener_closer = reportLine{}
			reportLine_1_token_type = reportLine{}
			reportLine_0_erl_src_chars = reportLine{}
		}
	}

	// add the possible last elems, too, without newline chars
	reportLines = append(reportLines, reportLine_4_separator)
	reportLines = append(reportLines, reportLine_3_char_counter)
	reportLines = append(reportLines, reportLine_2_opener_closer)
	reportLines = append(reportLines, reportLine_1_token_type)
	reportLines = append(reportLines, reportLine_0_erl_src_chars)

	lineNum := 0
	for _, oneReportLine := range reportLines {

		fmt.Printf("line %3d >>> ", lineNum/printedReportLinesForOneErlangSrcLine)
		for _, oneReportChar := range oneReportLine {
			fmt.Printf("%c", oneReportChar)
		}
		fmt.Printf("\n")
		lineNum += 1
	}

}
