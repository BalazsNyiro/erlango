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
func str_from_runes(runes []rune) string {
	return string(runes)
}

func str_double_space_remove(txt string) string {
	for strings.Contains(txt, "  ") {
		txt = strings.Replace(txt, "  ", " ", -1)
	}
	return txt
}

func int_from_str(txt string, exitIfError bool, callerFunName string) (int, error) {
	i, err := strconv.Atoi(txt)
	if err != nil {
		LogError(err, callerFunName + " convert string->int problem: " + txt)
		if exitIfError {
			panic(err)
		}
	}
	return i, err
}



/////////////////////////// DEBUG //////////////////////////////////////////////
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

/*
func log_fun(prg *Prg, msg, funName string) {
	fname := "devlog/log.txt"
	f, _ := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	text := fmt.Sprintf("%s %s", msg, funName) + "\n"

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
		formattedOutput := indentation + text
		fmt.Printf(formattedOutput)
		f.WriteString(formattedOutput+"\n")

		if msg == "<-" {
			// don't check the length, because it can hide open/close mismatching errors
			// if len(prg.callStackFunNames) > 0 { }
			prg.callStackFunNames = prg.callStackFunNames[:len(prg.callStackFunNames)-1]
		}

	}
}
*/

func getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	// complex fun name: command_x2dline_x2darguments.ParseErlangSourceCode
	funName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	if strings.Contains(funName, ".") {
		funName = strings.Split(funName, ".")[1]
	}
	return funName
}
