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
