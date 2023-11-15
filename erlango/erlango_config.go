/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite
*/

// follow XDG Base Directory specification
// https://neovim.io/doc/user/starting.html#xdg
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

// TODO: create config file, and parse that

package erlango

const abcEngLower = "abcdefghijklmnopqrstuvwxyz"
const abcEngUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const abcFullWith_At = "_@" + abcEngLower + abcEngUpper

const abcDigits = "0123456789"

const whiteSpaces = " \t\n\r"

/*These one char wide elems are part of Erlang language,
they can have more meanings, depends on their position,
so in token detection's first step they don't have deeper meaning

the other punctuation is coming from : unicode char's Category:
https://www.compart.com/en/unicode/U+003A


? - for conditional match operator:
The conditional match operator in {ok, A} ?= a() fails to match,
https://www.erlang.org/doc/reference_manual/expressions.html
*/
const otherPunctuation = "=.:,;(){}[]+-*/%<>#!?"
