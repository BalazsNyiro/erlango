/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango

import (
	"os"
	"slices"
)

/*
Main rule:
	all argument name starts with '--' prefix
	--argument_name_without_whitespace param1 param2  --secondArg 1 2

	original Erlang erl params can be passed with --erl "original erlang params"

	an argument can have zero, one ore more values, but everything is stored in a slice list, as a general solution.

	a program flag: where the argument name is defined only, without any values. :-)

	argument parsing is at erlango start, the program can die if a param is incorrect (invalid int value, for example)
*/

func cli_argument_detect(prg ProgramWideStateVariable) ProgramWideStateVariable {
	prg.ArgumentsErlangoCliStart.SettingStr["interpreterName"] = []string{os.Args[0]}

	argumentNamesWithIntegerValues := []string{"--argWithIntValues"}

	prg = prg_cli_argument_append_from_list(prg, os.Args[1:], argumentNamesWithIntegerValues )
	return prg
}

func prg_cli_argument_append_from_list(prg ProgramWideStateVariable, argumentElems []string, argumentNamesWithIntegerValues []string) ProgramWideStateVariable {
	funName := "prg_cli_argument_append_from_list"

	argumentNameLastDetected := ""

	for i := 0; i < len(argumentElems); i++ {

		arg := argumentElems[i]

		if arg[:2] == "--" {
			// set emtpy value set for the argument if the --argName pattern is detected
			argumentNameLastDetected = arg
			if slices.Contains(argumentNamesWithIntegerValues, argumentNameLastDetected) {
				prg.ArgumentsErlangoCliStart.SettingInt[argumentNameLastDetected] = []int{}
			} else {
				prg.ArgumentsErlangoCliStart.SettingStr[argumentNameLastDetected] = []string{}
			}
			continue
		}

		if argumentNameLastDetected != "" {
			if slices.Contains(argumentNamesWithIntegerValues, argumentNameLastDetected) {
				intVal, _ := int_from_str(arg, true, funName)
				prg.ArgumentsErlangoCliStart.SettingInt[argumentNameLastDetected] =
					append(prg.ArgumentsErlangoCliStart.SettingInt[argumentNameLastDetected], intVal)
			} else {

				prg.ArgumentsErlangoCliStart.SettingStr[argumentNameLastDetected] =
					append(prg.ArgumentsErlangoCliStart.SettingStr[argumentNameLastDetected], arg)
			}
			continue
		}
	}
	return prg
}

/*  in program params, filenames can be directly defined to load/start them
	--filenames src1.erl other/src2.erl

	and because in the future, maybe directory names will be handled, too,
	and other options to specify the wanted files, this function will collect
	the file list from arguments:
*/
func filenames_erlang_sources_collect_from_cli_params(prg ProgramWideStateVariable) []string {
	filePathList := []string{}

	filePathListInArgs, ok := prg.ArgumentsErlangoCliStart.SettingStr["--files"]
	if ok {
		for _, filePathInArg := range(filePathListInArgs) {
			filePathList = append(filePathList, filePathInArg)
		}
	}
	return filePathList
}