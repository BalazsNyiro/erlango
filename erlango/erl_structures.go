/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

const LanguageElemBlockList string = "LanguageElemBlockList "
const LanguageElemBlockParenthesesRound string = "LanguageElemBlockParenthesesRound "
const LanguageElemBlockParenthesesSquare string = "LanguageElemBlockParenthesesSquare "

/*
A language elem is one LOGICAL UNIT, a thing with one meaning: a list, a map, a case structure.
it's typically built by a lot of tokens: opening/closing (...) pairs, and it can have a lot of internal elems
*/

// The graph of Language Elems is an Erlang program, internally.
// So this is the core of the interpreter :-)
// a language elem can contain more other language elems,

/*
A language elem is something that has a value

if we have a list, it can contain other lists, or numbers, for example,
and a list is a recursive structure. A lot of tokens are used to build up
one logical object.

A list from the programmer's perspective is one thing: a container.
from the interpreter's perspective: it is ONE LANGUAGE ELEM, one logical unit.

So this is the point when the tokens will be transformed to LanguageElems,
one Elem represents one language objects ( a list, a number, an atom, a keyword)

The whole program is built from LanguageElems, too, because with operators
values are bind to variables, and the whole program becomes an executable graph.

A LanguageElem from the interpreter's perspective: something that you can evaluate,
and it has a return value, or it can run in a neverending loop

*/

type LanguageElem struct {
	Type                  string
	IncludedLanguageElems []LanguageElem
	// if a language elem has tokens only, these are the leaves in the graph
	ErlTokensIfNoLangElems []ErlSrcToken
}

// an interesting Python article: https://www.scaler.com/topics/expression-in-python/
// at this point we have a processed chars data structure with detected tokens
func languageElemsDetect(chars []ErlSrcChar) {

	/* From chars, the ParseErlangSourceDoce detected Tokens.
	   Here from tokens, generate executable LanguageElems. */

	prg := &Prg{callStackDisplay: true}
	chars, _ = ParseErlangSourceCode(prg, chars, "__all__")
	debug_print_ErlSrcChars(chars)

}
