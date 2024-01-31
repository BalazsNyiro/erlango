/*
Erlang - Go implementation.

Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package old2

import "fmt"

func Erlang_program_exec(prg ProgramWideStateVariable) {  // in program plan

	sourcesTokensExecutables_list := SourcesTokensExecutables_map{}
	// step_01_tokens_from_source_code_of_files(false, prg["verboseForErlangoInvestigations__useFalseInProdEnv"])
	print("FIXME: check and display possible ERRORS after token detection")
	//sourcesTokensExecutables_all =  step_02_expressions_from_tokens_from_lot_of_sources(sourcesTokensExecutables_all, "detectAllExpressions")
	// step_03_exec_main()
	fmt.Print(sourcesTokensExecutables_list)
	fmt.Print("prg", prg)

}


func main() {  // in program plan
	prg := new_program_state(false)
	prg = cli_argument_detect(prg)  // all arguments are parsed, placed in prg
	Erlang_program_exec(prg)
}