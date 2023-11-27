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

== Definitions for Erlang language elem detection ==

https://www.erlang.org/doc/reference_manual/expressions.html

Erlang term (simple data types):
	- an integer, float, atom, string, list, map, or tuple.

Erlang variables
	Variables start with an uppercase letter or underscore (_). Variables can contain alphanumeric characters, underscore and @.
	Variables starting with underscore (_), for example, _Height, are normal variables, not anonymous.

	Special chars are NOT allowed:
	Eshell V13.1.5  (abort with ^G)
	1> Aáéői = 3.
	* 1:4: illegal character

Operators:
	whitespaces are not important in operator detection (+ is addition, - is unary operator):
	1> A = 2+-1.
	1


== TOKENS ==
 - meaningful characters that can be interpreted in different environments.
   Tokens have to be interpreted with their environments/positions
   For example: a . have a different meanings in these situations
		- 1.2
		- #Name.Field
		- . at the end of a function.


== Language elem detection steps: ==

 - detect simple terms:
    - integers: https://www.erlang.org/doc/reference_manual/data_types.html
		- 1234.
		- 1_234_567_890.
		- $A.
		- 16#1f.
		- 16#4865_316F_774F_6C64.
		- 2e-3.

	- floats:
		- 2.3.
		- 2.3e-3.
		- 1_234.333_333.

	- atom-quoted
	- atom

	- string
	- bit-string https://www.erlang.org/doc/man/binary.html

 - detect complex terms: list, map, tuple

 - detect complex structures (language elems):
	- functions
	- conditions
	- exceptions




== LANGUAGE ELEMS


== LINKS, REFERENCES ==
           data types:  https://www.erlang.org/doc/reference_manual/data_types
dot, colon, semicolon:  https://stackoverflow.com/questions/1110601/in-erlang-when-do-i-use-or-or
   terms, expressions:  https://www.erlang.org/doc/reference_manual/expressions#terms


== NUMBERS ==
1> 1_000 * 3.
3000


*/

/* simple types: 1-10, complex types: 10+ */
const expression_atom = 1
const expression_num = 2
const expression_stringDoubleQuoted = 3

const expression_list = 10
const expression_tuple = 11
const expression_map = 12

const expression_parentheseRoundedGroup = 20

const expression_operator = 30
// list operator: <- <=      [X*2 || X <- [1,2,3]].
// map operators: =>
// blockStart operator ->  (after fun, case, if)


// file name passing is important, because maybe the expression detection
// is running in an old elem, where once the expressions were detected,
// and now it is re-detected, based on new tokens
func step_02_expressions_from_tokens_from_lot_of_sources(sourcesTokensExecutables_all SourcesTokensExecutables_map, fileNamePathsWhereExpressionsWillBeDetected []string)  SourcesTokensExecutables_map {
	fmt.Println("filenames or sourcePassedInString_notFromFile to detect expressions", fileNamePathsWhereExpressionsWillBeDetected)

	returnFromExpressionDetection := make(chan SourceTokensExecutables)
	for _, filePath := range(fileNamePathsWhereExpressionsWillBeDetected) {
		go step_02a_expressions_detect_in_one_erlang_source(filePath, returnFromExpressionDetection, sourcesTokensExecutables_all[filePath])
	}

	// because of the parallel expression detection, it is simpler
	// if the whole sourceTokensExecutables structure is updated, and I don't use any pointer,
	// anywhere - that is a risk.
	// in prod env there are 30-40 or more cores, the parsing cannot be a problem.
	numOfReceivedReply := 0
	for numOfReceivedReply < len(fileNamePathsWhereExpressionsWillBeDetected) {
		sourceTokensExecutables := <- returnFromExpressionDetection
		sourcesTokensExecutables_all[sourceTokensExecutables.WhereTheCodeIsStored] = sourceTokensExecutables
		numOfReceivedReply += 1
	}

	return sourcesTokensExecutables_all // it can have errors, too!
}


/* 	Hi anybody who reads this: expression detection was my most fearful part of the interpreter :-)

	an expression can be formed by one token (an atom or a string for example),
	or by more tokens. The first token position will be used as the expression start position,
	and as a general ID in the file for the expression

	Why is it tricky?

	because same things can be represented more way.
	Numbers for example - different num representations:

			- integers: https://www.erlang.org/doc/reference_manual/data_types.html
				- 1234
				- 1_234_567_890
				- $A
				- 16#1f
				- 16#4865_316F_774F_6C64
				- 2e-3

			- floats:
				- 2.3
				- 2.3e-3
				- 1_234.333_333


	the integers/strings/atoms are maybe the friendly part of the story,
	but lists, tuples can be recursive structures, so a list can have tuples which has lists....
	but the recursive structures has to be finished once.

	First I will focus on recursive structures: lists, tuples, maps
	(functions can be recursive structures, too :-) because funs can have embedded funs, too

	if a block of tokens is detected as the part of a recursive structure,
	the tokens are taken and the expression detection are executed in there again.

	with this solution the expressions can be embedded into each other,
	and the calling structure will be represented by the newly created
	embedded expression structure.

	Second big problem: because we are in Golang, and expressions are self-recursive structures,
	I can use one data type to represent every expression.

	// And at the beginning, God created the expressions... :-)

	recursive expressions
		- tuple {...}
		- list  [...]
		- map  #{...}
        - parentheses (...)  # a parentheses content needed to be evaluated as a term at the end

		- block elems:
			- function
			- condition (if, case)
			- receive

	Because this is the critical part of the whole parsing, I will use general rules to describe and find sections.
*/

func step_02a_expressions_detect_in_one_erlang_source(filePath string, parentChannel chan SourceTokensExecutables, sourceTokensExecutables SourceTokensExecutables){
	fmt.Println("Expression detect in file:", filePath)

	tokensOrExpressions := tokens_copy_to_tokensOrExpressions(sourceTokensExecutables.Tokens)
	fmt.Println("tokens OR expressions ", tokensOrExpressions)

	tokensOrExpressions = expression_detect_from_tokens(tokensOrExpressions, 0)


	parentChannel <- sourceTokensExecutables
}

func expression_detect_from_tokens(tokensOrExpressionsOld TokensOrExpressions, expressionCallstackLevel int) TokensOrExpressions {
	// we have list of tokens.
	// select a group of tokens, replace them with an expression,
	// and the selected tokens are inserted INTO the expression.
	// from a flat structure an embedded expression structure will be created

	tokensOrExpressionsNew := TokensOrExpressions{}

	for _, tokenOrExpression := range(tokensOrExpressionsOld) {
		fmt.Println("token expression", tokenOrExpression)
		//
		// tokensOrExpressions = append(tokensOrExpressions, TokenOrExpression{token: tokenNow, elemType: "token"})
	}


	return tokensOrExpressionsNew
}

type TokensOrExpressions []TokenOrExpression
type TokenOrExpression struct {
	// this is a Token OR an Expression storage.
	// originally everything is a token - then slowly all of them will be replaced with Expressions
	elemType string  // "token" or "expression"
	token ErlToken
	expression ErlExpression
}
func (tokenOrExpression TokenOrExpression) isExpression() bool {
	return tokenOrExpression.elemType == "expression"
}

type ErlExpressions map[int] ErlExpression
type ErlExpression struct {
	/*  This is the heart of the interpreter */

	Level int
	/* level 0: the root level of a file, the base namespace.
		every function is the part of this, as sub-expressions.
	*/

	Tokens ErlTokens

	ExpressionType int // expression_atom, expression_num...
	Children ErlExpressions  // lists, tuples, maps, functions have children
}



func tokens_copy_to_tokensOrExpressions(tokens ErlTokens) TokensOrExpressions {
	/* 	One tokenOrExpression can be both: a token Or an expression.
	At the beginning, everything is a token - but as the expression detector works,
	more and more elems will be removed, and replaced by expressions
	*/
	tokensOrExpressions	:= TokensOrExpressions{}
	for _, tokenPosition := range(tokens.keysListOfPositions()) {
		tokenNow := tokens[tokenPosition]
		tokensOrExpressions = append(tokensOrExpressions, TokenOrExpression{token: tokenNow})
	}
	return tokensOrExpressions
}
