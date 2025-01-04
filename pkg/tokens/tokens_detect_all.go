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
	charactersInErlSrc, tokensInErlSrc = tokens_detect_comments_strings_quotedatoms(charactersInErlSrc, tokensInErlSrc)
	return charactersInErlSrc, tokensInErlSrc
}

func Tokens_detection_print_verbose(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) {

	lineNumInErlSrc := 0

	type reportLine []rune
	type reportLineMore []reportLine

	// one erland line -> more lines are printed with token type infos
	printedReportLinesForOneErlangSrcLine := 3 //  for every erlang line, 2 report lines are printed
	reportLine_2_separator := []rune("============================")
	reportLine_1_token_type := reportLine{}
	reportLine_0_erl_src_chars := reportLine{}

	reportLines := reportLineMore{}

	for _, charInErlSrc := range charactersInErlSrc {

		if charInErlSrc.runeInErlSrc == '\n' { // newline chars
			lineNumInErlSrc += 1
			reportLines = append(reportLines, reportLine_2_separator)
			reportLines = append(reportLines, reportLine_1_token_type)
			reportLines = append(reportLines, reportLine_0_erl_src_chars)
			reportLine_1_token_type = reportLine{}
			reportLine_0_erl_src_chars = reportLine{}
		} else { // non-newline char
			oneCharWideTokenTypeRepresentation := TokenTypeReprShort(charInErlSrc.tokenDetectedType)
			reportLine_1_token_type = append(reportLine_1_token_type, oneCharWideTokenTypeRepresentation)
			reportLine_0_erl_src_chars = append(reportLine_0_erl_src_chars, charInErlSrc.runeInErlSrc)
		}
	}

	// add the possible last elems, too, without newline chars
	reportLines = append(reportLines, reportLine_2_separator)
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
