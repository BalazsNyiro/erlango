/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package tokens

import (
	"erlango.org/erlango/pkg/base_toolset"
	"fmt"
)

// 'convert characters -> tokens'
func Tokens_detect_in_erl_src(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) (CharacterInErlSrcCollector, TokenCollector, errorMessages) {

	charactersInErlSrc = character_block_detect__01_erlang_strings__quoted_atoms__comments(charactersInErlSrc)
	charactersInErlSrc = character_block_detect__02_erlang_whitespaces(charactersInErlSrc)
	charactersInErlSrc = character_block_detect__03_erlang_alphanumerics(charactersInErlSrc)
	charactersInErlSrc = character_block_detect__04_erlang_braces__dotsCommas__operatorBuilders(charactersInErlSrc)

	errors := character_blocks_validations___unknownSections__nonClosedSections(charactersInErlSrc)
	return charactersInErlSrc, tokensInErlSrc, errors
}

// convert 'fileErl -> characters -> tokens'
func Tokens_detect_in_erl_file(interpreterHostMachineCoord string, fileErl string, callerFun string) (CharacterInErlSrcCollector, TokenCollector, errorMessages) {
	erlSrcRunes, _ := base_toolset.File_read_runes(fileErl, callerFun)
	charactersInErlSrc := Runes_to_character_structs(erlSrcRunes, "file:"+interpreterHostMachineCoord+":"+fileErl)
	tokensInErlSrc := TokenCollector{}
	charactersInErlSrc2, tokensInErlSrc2, errors := Tokens_detect_in_erl_src(charactersInErlSrc, tokensInErlSrc)
	return charactersInErlSrc2, tokensInErlSrc2, errors
}

// detect unknown character sections in erlang source
// detect non-closed sections
func character_blocks_validations___unknownSections__nonClosedSections(charactersInErlSrc CharacterInErlSrcCollector) errorMessages {
	errors := errorMessages{}

	counterOpener := 0
	counterCloser := 0

	for charPositionNowInSrc := 0; charPositionNowInSrc < len(charactersInErlSrc); charPositionNowInSrc++ {

		charStructNow := charactersInErlSrc[charPositionNowInSrc]

		if charStructNow.charBlockIsNotDetected() {
			errors = append(errors, "char block type is not detected: "+charStructNow.stringReprDetailed())
		}

		if charStructNow.charBlockOpenerCharacter {
			counterOpener++

			if counterOpener-counterCloser != 1 {
				errors = append(errors, "char block closer was not detected before this opener: "+charStructNow.stringReprDetailed())
			}
		}
		if charStructNow.charBlockCloserCharacter {
			counterCloser++
			if counterOpener-counterCloser != 0 {
				errors = append(errors, "char block closer: too many detected: "+charStructNow.stringReprDetailed())
			}
		}

		// this is the last character in the source
		if charPositionNowInSrc == len(charactersInErlSrc)-1 {
			if counterOpener-counterCloser != 0 {
				errors = append(errors, "char block closer: The last block closer is missing: "+charStructNow.stringReprDetailed())
			}
		}

	}

	return errors
}

func Tokens_detection_print_one_char_per_line(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector, displayOnlyUnknownChars bool) {

	lineNumInErlSrc := 1

	for _, charInErlSrc := range charactersInErlSrc {
		display := true
		if displayOnlyUnknownChars {
			display = false
			if charInErlSrc.charBlockIsNotDetected() {
				display = true
			}
		}

		if display {
			fmt.Println("line:", lineNumInErlSrc, string(charInErlSrc.runeInErlSrc), charInErlSrc.charBlockDetectedType)
		}

		// fmt.Println("charCounter: ", charCounter)
		if charInErlSrc.runeInErlSrc == '\n' { // newline chars
			lineNumInErlSrc += 1
		}
	}

}

func Tokens_detection_print_verbose(charactersInErlSrc CharacterInErlSrcCollector, tokensInErlSrc TokenCollector) {

	lineNumInErlSrc := 1

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
		if charInErlSrc.charBlockOpenerCharacter {
			openerCloserStatus = 'o'
		}
		if charInErlSrc.charBlockCloserCharacter {
			openerCloserStatus = 'c'
		}
		if charInErlSrc.charBlockOpenerCharacter && charInErlSrc.charBlockCloserCharacter {
			openerCloserStatus = '2' // closer AND opener same time
		}
		reportLine_2_opener_closer = append(reportLine_2_opener_closer, openerCloserStatus)

		oneCharWideTokenTypeRepresentation := CharBlockReprShort(charInErlSrc.charBlockDetectedType)
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
