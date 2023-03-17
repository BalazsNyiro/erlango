/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// logLevel: Info, Warning, Error
func LogInfo(msg string) {
	log("Info", msg)
}
func LogWarning(msg string) {
	log("Warning", msg)
}
func LogError(err error, msg string) {
	// err.Error() the string representation of the error
	log("Error", err.Error()+" - "+msg)
}
func log(logLevel, msg string) {
	fmt.Println(logLevel, "->", msg)
}

func bool_to_str(val bool, trueTxt, falseTxt string) string {
	if val { return trueTxt } else { return falseTxt}
}

func file_read_runes(filePath, caller string) ([]rune, error) {
	f, err := os.Open(filePath)
	if err != nil {
		LogError(err, caller + " "+filePath)
		return []rune{}, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	runes := []rune{}
	for {
		if runeInFile, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				LogError(err, caller + " read, Rune problem: "+filePath)
			}
		} else {
			runes = append(runes, runeInFile)
		}
	}
	return runes, nil
}

func runes_from_str(txt string) []rune {
	var runes []rune
	for _, runeNow := range txt {
		runes = append(runes, runeNow)
	}
	return runes
}

func str_double_space_remove(txt string) string {
	for strings.Contains(txt, "  ") {
		txt = strings.Replace(txt, "  ", " ", -1)
	}
	return txt
}





/////////////////////////// DEBUG //////////////////////////////////////////////
// I need it from tests and normal code too, so this is the best place
func debug_print_ErlSrcChars_original(chars []ErlSrcChar) {
	fmt.Println("")
	for i, _ := range chars {
		fmt.Printf("%3d posInFile:%3d val:%4s ", i, chars[i].PosInFile, string(chars[i].Value))

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

func debug_print_ErlSrcChars(charsPtr *([]ErlSrcChar)) {
	fmt.Printf("inside >>>  debug_print_ErlSrcChars %p \n", charsPtr)
	for id, _ := range *charsPtr {
		debug_print_ErlSrcChar(id, &( (*charsPtr)[id])  )
	}
}

func debug_print_ErlSrcChar(id int, charPtr *ErlSrcChar) {
		fmt.Printf("charPtr:%p  %3d posInFile:%3d val:%4s ",charPtr, id, charPtr.PosInFile, string(charPtr.Value))

		prevPos := -1
		if charPtr.PrevChar != nil {
			prevPos = charPtr.PrevChar.PosInFile
		}
		fmt.Printf(" PrevPosInFile:%3d ", prevPos)

		tokenType := "<?>"
		if charPtr.Token != nil {
			tokenType = charPtr.Token.Type
		}
		fmt.Printf(" %p <- %p -> %p tokenPtr: %p type->%s<-", charPtr.PrevChar, charPtr, charPtr.NextChar, charPtr.Token, tokenType)
		fmt.Println("")
}

func debug_print_mem_address(msg string, object any ){
	fmt.Printf("%s %p \n", msg, object)
}
/////////////////////////// DEBUG //////////////////////////////////////////////
