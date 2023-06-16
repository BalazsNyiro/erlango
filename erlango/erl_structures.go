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
A language elem is something that has a meaning.

if we have a list, it can contain other lists, or numbers, for example,
and a list is a recursive structure. A lot of tokens are used to build up
one logical object.

A list from the programmer's perspective is one thing: a container.
from the interpreter's perspective: it is ONE LANGUAGE ELEM, one logical unit.

So this is the point when the tokens will be transformed to LanguageElems,
one Elem represents one language objects ( a list, a number, an atom, a keyword)
*/
type LanguageElem struct {
	Type                  string
	IncludedLanguageElems []LanguageElem
	// if a language elem has tokens only, these are the leaves in the graph
	ErlTokensIfNoLangElems []ErlSrcToken
}

// an interesting Python article: https://www.scaler.com/topics/expression-in-python/
// at this point we have a processed chars data structure with detected tokens
func expressionsDetect(chars []ErlSrcChar) {

	/* txt := `A = (1 + 2) * 3.`
	   chars := ErlSrcChars_from_str(txt)
	   chars, _ = ParseErlangSourceCode(chars, "__all__")
	   ^^^^^ this chars is the state, where we are now. we detected the basic tokens.

	  0 posInFile:  0 val:   M  ...  type->Token_type_digits_baseDefined<-
	  1 posInFile:  1 val:      ...  type->separator_whitespace<-
	  2 posInFile:  2 val:   =  ...  type->Token_type_binding_matching<-
	  3 posInFile:  3 val:      ...  type->separator_whitespace<-
	  4 posInFile:  4 val:   #  ...  type->Token_type_bracket_map_open<-
	  5 posInFile:  5 val:   {  ...  type->Token_type_bracket_map_open<-
	  6 posInFile:  6 val:   9  ...  type->Token_type_digits_base10_form<-
	  7 posInFile:  7 val:      ...  type->separator_whitespace<-
	  8 posInFile:  8 val:   =  ...  type->Token_type_arrow_singleToRight<-
	  9 posInFile:  9 val:   >  ...  type->Token_type_arrow_singleToRight<-
	 10 posInFile: 10 val:      ...  type->separator_whitespace<-
	 11 posInFile: 11 val:   "  ...  type->txt_quoted_double<-
	 12 posInFile: 12 val:   n  ...  type->txt_quoted_double<-
	 13 posInFile: 13 val:   i  ...  type->txt_quoted_double<-
	 14 posInFile: 14 val:   n  ...  type->txt_quoted_double<-
	 15 posInFile: 15 val:   e  ...  type->txt_quoted_double<-
	 16 posInFile: 16 val:   "  ...  type->txt_quoted_double<-
	 17 posInFile: 17 val:   }  ...  type->bracket_curly_close<-
	 18 posInFile: 18 val:   ,  ...  type->separator_comma<-
	 19 posInFile: 19 val:      ...  type->separator_whitespace<-
	 20 posInFile: 20 val:   I  ...  type->Token_type_digits_baseDefined<-
	 21 posInFile: 21 val:   D  ...  type->Token_type_digits_baseDefined<-
	 22 posInFile: 22 val:      ...  type->separator_whitespace<-
	 23 posInFile: 23 val:   =  ...  type->Token_type_binding_matching<-
	 24 posInFile: 24 val:      ...  type->separator_whitespace<-
	 25 posInFile: 25 val:   (  ...  type->bracket_round_open<-
	 26 posInFile: 26 val:   1  ...  type->Token_type_digits_base10_form<-
	 27 posInFile: 27 val:   +  ...  type->Token_type_math_binary_add<-
	 28 posInFile: 28 val:   2  ...  type->Token_type_digits_base10_form<-
	 29 posInFile: 29 val:   )  ...  type->bracket_round_close<-
	 30 posInFile: 30 val:   *  ...  type->Token_type_math_binary_mul<-
	 31 posInFile: 31 val:   3  ...  type->Token_type_digits_base10_form<-
	 32 posInFile: 32 val:   ,  ...  type->separator_comma<-
	 33 posInFile: 33 val:      ...  type->separator_whitespace<-
	 34 posInFile: 34 val:   m  ...  type->Token_type_atom_quoteless<-
	 35 posInFile: 35 val:   a  ...  type->Token_type_atom_quoteless<-
	 36 posInFile: 36 val:   p  ...  type->Token_type_atom_quoteless<-
	 37 posInFile: 37 val:   s  ...  type->Token_type_atom_quoteless<-
	 38 posInFile: 38 val:   :  ...  type->separator_colon<-
	 39 posInFile: 39 val:   f  ...  type->Token_type_atom_quoteless<-
	 40 posInFile: 40 val:   i  ...  type->Token_type_atom_quoteless<-
	 41 posInFile: 41 val:   n  ...  type->Token_type_atom_quoteless<-
	 42 posInFile: 42 val:   d  ...  type->Token_type_atom_quoteless<-
	 43 posInFile: 43 val:   (  ...  type->bracket_round_open<-
	 44 posInFile: 44 val:   I  ...  type->Token_type_digits_baseDefined<-
	 45 posInFile: 45 val:   D  ...  type->Token_type_digits_baseDefined<-
	 46 posInFile: 46 val:   ,  ...  type->separator_comma<-
	 47 posInFile: 47 val:      ...  type->separator_whitespace<-
	 48 posInFile: 48 val:   M  ...  type->Token_type_digits_baseDefined<-
	 49 posInFile: 49 val:   )  ...  type->bracket_round_close<-
	 50 posInFile: 50 val:   .  ...  type->separator_dot<-



		from ^^^^ this char list, here we create expressions.

	An expression has a result. it can be a number or an atom - they mean themselves. (expression_single)
	If you use more numbers and operators, then that is an expression_embedded, because after the operator
	execution and func execution you can calculate the result (which can be replaced by an expression_single
	at the moment of code execution.

	in Erlang, everything has a return value. So the code is ready to be executed if all operators and operands
	are evaluated one by one, until the whole code becomes to one function execution.


	The task: detect embedded expressions and operators, so from the interpreter's perspective a long list is only one language elem.

	if all elem is detected, the operators can be executed.




	*/




	// TODO: collect Tokens from char list.
	// TODO2: if a char hasn't got a token, create a default 'undetecteds' tokens
}
