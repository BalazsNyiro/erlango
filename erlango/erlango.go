/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

package erlango

import "fmt"

/*

*/


func Erlang_program_exec() {

	prg := new_program_state()
	prg = cli_argument_detect(prg)

	sourcesTokensExecutables_list := SourcesTokensExecutables_list{}
	// step_01_tokens_from_source_code()
	// step_02_executables_from_tokens()
	// step_03_exec_main()
	fmt.Print(sourcesTokensExecutables_list)
	fmt.Print("prg", prg)

}


func main() {
	Erlang_program_exec()
}