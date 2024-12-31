/*
Erlang - Go implementation.

Copyright (c) 2024, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package base_toolset

import (
	"fmt"
	"runtime"
	"strings"
)

func GetCurrentGoFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	// complex fun name: command_x2dline_x2darguments.ParseErlangSourceCode
	funName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	if strings.Contains(funName, ".") {
		funName = strings.Split(funName, ".")[1]
	}
	return funName
}
