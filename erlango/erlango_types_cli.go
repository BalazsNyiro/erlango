/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango


/* 	very rarely modified, and often read program-wide data structure.
Stored infos:
	- erlango command line call parameters (good to know everywhere)

*/
type ArgumentsErlangoCliStart struct {
	ArgStr  map[string][]string
	ArgInt  map[string][]int
	ArgBool map[string]bool
}

type ProgramWideStateVariable struct {
	ArgumentsErlangoCliStart ArgumentsErlangoCliStart
}

func new_program_state() ProgramWideStateVariable {
	return ProgramWideStateVariable {
			ArgumentsErlangoCliStart {
				ArgStr: map[string][]string{},
				ArgInt: map[string][]int{},
				ArgBool: map[string]bool{},
			},
	}
}
