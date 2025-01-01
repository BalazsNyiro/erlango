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

// This fun can be called from other functions, where you can get dynamically the actual function name
func GetCurrentGoFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	funName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	if strings.Contains(funName, ".") {
		funName = strings.Split(funName, ".")[1]
	}
	return funName
}
