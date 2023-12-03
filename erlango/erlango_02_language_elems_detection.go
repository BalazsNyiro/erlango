/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite

*/

package erlango

import (
	"fmt"
	"strings"
)

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

// BOOKMARK expressionTypes
const expression_nonDetectedFromToken = 0

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

// reverse conversion: from the code, know, what is the type
var ExpressionName_from_num map[int]string = map[int]string {
	expression_nonDetectedFromToken: "expression_nonDetectedFromToken",
	expression_atom: "expression_atom",
	expression_num: "expression_num",
	expression_stringDoubleQuoted: "expression_stringDoubleQuoted",
	expression_list: "expression_list",
	expression_tuple: "expression_tuple",
	expression_map: "expression_map",
	expression_parentheseRoundedGroup: "expression_parentheseRoundedGroup",
	expression_operator:"expression_operator",
}

// file name passing is important, because maybe the expression detection
// is running in an old elem, where once the expressions were detected,
// and now it is re-detected, based on new tokens
func step_02_expressions_from_tokens_from_lot_of_sources(
	sourcesTokensExecutables_all SourcesTokensExecutables_map,
	fileNamePathsWhereExpressionsWillBeDetected []string,
	wantedExpressionDetectionTypesCommaSeparated string) SourcesTokensExecutables_map {
	fmt.Println("filenames or sourcePassedInString_notFromFile to detect expressions", fileNamePathsWhereExpressionsWillBeDetected)

	returnFromExpressionDetection := make(chan SourceTokensExecutables)
	for _, filePath := range(fileNamePathsWhereExpressionsWillBeDetected) {
		go step_02a_expressions_detect_in_one_erlang_source(
			filePath,
			returnFromExpressionDetection,
			sourcesTokensExecutables_all[filePath],
			wantedExpressionDetectionTypesCommaSeparated,
			)
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

func step_02a_expressions_detect_in_one_erlang_source(
	filePath string,
	parentChannel chan SourceTokensExecutables,
	sourceTokensExecutables SourceTokensExecutables,
	wantedExpressionDetectionTypesCommaSeparated string){

	/* Task: prepare tokensOrExperssions structure, then start the detection */
	fmt.Println("Expression detect:", filePath)

	tokensOrExpressions := tokens_copyTo_tokensOrExpressions(sourceTokensExecutables.Tokens)
	tokensOrExpressions = expressionDetectAllType_from_tokens(tokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated)

	// at the end, move back the tokensOrExpressions into simple expression list?
	fmt.Println("and maybe throw error, if a tokenOrExpression is not converted to be an expression")
	for _, tokenOrExpression := range(tokensOrExpressions) {
		if tokenOrExpression.isExpression() {
			sourceTokensExecutables.Expressions = append(sourceTokensExecutables.Expressions, tokenOrExpression.expression)
		} else {
			fmt.Println("ERROR: missing EXPRESSION CONVERSION: a tokenOrExpression is not converted to be an expression - nonDetected Expresssion inserted")
			errorExpression := ErlExpression{
				ExpressionType: expression_nonDetectedFromToken,
				SimpleTokenValue: tokenOrExpression.token,
			}
			sourceTokensExecutables.Expressions = append(sourceTokensExecutables.Expressions, errorExpression)
		}
	}

	parentChannel <- sourceTokensExecutables
}

const tokenOrExpression_elemTypeExpression = "expression"

func expressionDetectAllType_from_tokens(
	tokensOrExpressionsOld TokensOrExpressions,
	wantedExpressionDetectionTypesCommaSeparated string,
	) TokensOrExpressions {
	// we have list of tokens.
	// select a group of tokens, replace them with an expression,
	// and the selected tokens are inserted INTO the expression.
	// from a flat structure an embedded expression structure will be created

	// Named function definitions =======================================================
	/*  https://www.erlang.org/doc/reference_manual/functions.html */

	tokensOrExpressionsNew_01_atomsDetected := expression_detect_atoms(tokensOrExpressionsOld, wantedExpressionDetectionTypesCommaSeparated)

	return tokensOrExpressionsNew_01_atomsDetected
}

////////////////////////////////////////////////////////////////////////////////

func expression_detect_atoms(tokensOrExpressionsOld TokensOrExpressions, wantedExpressionDetectionTypesCommaSeparated string) TokensOrExpressions {
	if ! (strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "atomsQuoted") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "atomsSimple") ||
		strings.Contains(wantedExpressionDetectionTypesCommaSeparated, "detectAllExpressions")) {
		// if atom detection is not a wanted operation, then don't do that
		return tokensOrExpressionsOld
	}

	tokensOrExpressionsNew_01_atomsDetected := TokensOrExpressions{}

	// first detect basic types (atom, string, numbers)
	for _, tokenOrExpression := range(tokensOrExpressionsOld) {
		fmt.Println("detect atoms - token expression", tokenOrExpression)

		if tokenOrExpression.isExpression() {  // if it is a previously detected expression, there is nothing to do
			tokensOrExpressionsNew_01_atomsDetected = append(tokensOrExpressionsNew_01_atomsDetected, tokenOrExpression)
			continue
		}

		isAtom := false

		// 'quoted atom' - honestly the string based type checking is maybe slower here, than the int based in expressions.
		// The token->expression conversation is not a runtime operation, so in this level now it's fine.
		if tokenOrExpression.token.TokenType == "tokenTextBlockQuotedSingle" {
			isAtom = true
		}

		// atom
		if tokenOrExpression.token.TokenType == "tokenAbcFullWith_At_numbers" {
			if tokenOrExpression.token.charFirstRuneValIsSmallCapsAtomStarter() {
				isAtom = true
			}
		}

		if isAtom {
			tokenOrExpression.elemType = tokenOrExpression_elemTypeExpression
			tokenOrExpression.expression = ErlExpression{
				ExpressionType:        expression_atom,
				SimpleTokenValue:      tokenOrExpression.token,
			}
			// put back tokenOrExpression with modified elemType and expression
			tokensOrExpressionsNew_01_atomsDetected = append(tokensOrExpressionsNew_01_atomsDetected, tokenOrExpression)

		} else {  // not an atom - put back the tokenOrExpression without any extra change/modification
			// I know this is same with the isAtom's last append - but it is more readable to see the two sections
			tokensOrExpressionsNew_01_atomsDetected = append(tokensOrExpressionsNew_01_atomsDetected, tokenOrExpression)
		}
	} // FOR

	return tokensOrExpressionsNew_01_atomsDetected
}

////////////////////////////////////////////////////////////////////////////////

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


type ErlExpressions []ErlExpression
type ErlExpression struct {
	/*  This is the heart of the interpreter */
	ExpressionType int // expression_atom, expression_num... (BOOKMARK-labeled in source code)

	TokensOrExpressions TokensOrExpressions

	SimpleTokenValue ErlToken
	// if the expression is a simple term, number, atom, string,
	// the token value is stored here
}

// give back the human representation of type
func (erlExpression ErlExpression) expressionTypeForHuman() string {
	return ExpressionName_from_num[erlExpression.ExpressionType]
}


func tokens_copyTo_tokensOrExpressions(tokens ErlTokens) TokensOrExpressions {
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
