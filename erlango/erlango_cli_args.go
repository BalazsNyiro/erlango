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
)

/*
Main rule:
	all argument name starts with '--' prefix
	--argument_name_without_whitespace param1 param2  --secondArg 1 2

	original Erlang erl params can be passed with --erl "original erlang params"

	an argument can have zero, one ore more values, but everything is stored in a slice list, as a general solution.

	a program flag: where the argument name is defined only, without any values. :-)

*/
func cli_argument_detect(prg ProgramWideStateVariable) ProgramWideStateVariable {

	prg.ArgumentsErlangoCliStart.ArgStr["interpreterName"] = []string{os.Args[0]}

	argumentNameLastDetected := ""

	for i := 0; i < len(os.Args); i++ {

		arg := os.Args[i]

		if arg[:2] == "--" {
			argumentNameLastDetected = arg
			prg.ArgumentsErlangoCliStart.ArgStr[argumentNameLastDetected] = []string{}
			continue
		}

		if argumentNameLastDetected != "" {
			prg.ArgumentsErlangoCliStart.ArgStr[argumentNameLastDetected] =
				append(prg.ArgumentsErlangoCliStart.ArgStr[argumentNameLastDetected], arg)
			continue
		}

	}
	return prg
}