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
	"runtime"
	"sort"
	"strconv"
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

func file_read_lines(filePath, caller string) ([]string, error){
	f, err := os.Open(filePath)
	if err != nil {
		LogError(err, caller + " open: "+filePath)
		return []string{}, err
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	err2 := f.Close()
	if err2 != nil {
		LogError(err, caller + " close: "+filePath)
		return []string{}, err2
	}
	return lines, nil
}

func runes_from_str(txt string) []rune {
	var runes []rune
	for _, runeNow := range txt {
		runes = append(runes, runeNow)
	}
	return runes
}
func str_from_runes(runes []rune) string {
	return string(runes)
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
		if chars[i].TokenConnected() {
			tokenType = chars[i].Token.Type
		}
		fmt.Printf(" %p <- %p -> %p token: %p %s", chars[i].PrevChar, &chars[i], chars[i].NextChar, chars[i].Token, tokenType)
		fmt.Println("")
	}
}

func debug_print_ErlSrcChars(chars []ErlSrcChar) {
	fmt.Printf("inside >>>  debug_print_ErlSrcChars %p \n", chars)
	for id, _ := range chars {
		debug_print_ErlSrcChar(id, &(chars[id])  )
	}
}

func debug_print_ErlSrcChar(id int, charPtr *ErlSrcChar) {
		fmt.Printf("charPtr:%p  %3d posInFile:%3d val:%4s ",charPtr, id, charPtr.PosInFile, string(charPtr.Value))

		prevPos := -1
		if charPtr.PrevChar != nil {
			prevPos = charPtr.PrevChar.PosInFile
		}
		fmt.Printf(" PrevPosInFile:%3d ", prevPos)

		fmt.Printf(" %p <- %p -> %p tokenPtr: %p type->%s<-", charPtr.PrevChar, charPtr, charPtr.NextChar, charPtr.Token, (*charPtr).Type())
		fmt.Println("")
}

func debug_print_mem_address(msg string, object any ){
	fmt.Printf("%s %p \n", msg, object)
}

func map_print_keysorted__int_str(m map[int]string) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Println(k, m[k])
	}
}
/////////////////////////// DEBUG //////////////////////////////////////////////

func log_fun(prg *Prg, msg, funName string) {
	fname := "devlog/log.txt"
	f, _ := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	text := fmt.Sprintf("%s %s", msg, funName) + "\n"
	f.WriteString(text);
	f.Close()

	/////////////////////////////////////////////////////

	if prg.callStackDisplay {
		if msg == "->" {
			prg.callStackFunNames = append(prg.callStackFunNames, funName)
		}

		callLevel := len(prg.callStackFunNames)
		indentation := ""
		if callLevel > 0 {
			indentation = fmt.Sprintf("%"+strconv.Itoa(callLevel)+"s", " ")
		}
		fmt.Printf(indentation + text)

		if msg == "<-" {
			// don't check the length, because it can hide open/close mismatching errors
			// if len(prg.callStackFunNames) > 0 { }
			prg.callStackFunNames = prg.callStackFunNames[:len(prg.callStackFunNames)-1]
		}

	}
}

func getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	// complex fun name: command_x2dline_x2darguments.ParseErlangSourceCode
	funName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	if strings.Contains(funName, ".") {
		funName = strings.Split(funName, ".")[1]
	}
	return funName
}