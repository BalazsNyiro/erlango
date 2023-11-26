/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango

// general storage for settings - strings, ints, bools are typical config elems
type SettingsTable struct {
	SettingStr map[string][]string
	SettingInt  map[string][]int
	SettingBool map[string]bool
}


/* 	very rarely modified, and often read program-wide data structure.
Stored infos:
	- erlango command line call parameters (good to know everywhere)

in program start arguments, there can be strings, ints, bools
*/
type ProgramWideStateVariable struct {
	ArgumentsErlangoCliStart SettingsTable
}

// const debbuggerVerboseMode = true
func new_program_state(verboseForErlangoInvestigations__useFalseInProdEnv bool) ProgramWideStateVariable {
	return ProgramWideStateVariable {
			SettingsTable{
				SettingStr:  map[string][]string{},
				SettingInt:  map[string][]int{},
				SettingBool: map[string]bool{"verboseForErlangoInvestigations__useFalseInProdEnv": verboseForErlangoInvestigations__useFalseInProdEnv},
			},
	}
}
