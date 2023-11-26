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
	"sort"
)


type ErlTokens map[int] ErlToken 	// list AND map, same time. the token's first char position is the key,
									// so later it is easy to add new token, and by nature it is a list, too
									// and in the tests it is easy to check (if it is not a map in tests it is difficult
									// to check the result. if it is a map, the token positions can be checked directly



func (tokens ErlTokens) keysListOfPositions() []int {
	// thank you and respect for the faster implementation than append:
	// https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map
	positionKeys := make([]int, len(tokens))
	i := 0
	for k := range tokens {
		positionKeys[i] = k
		i++
	}
	sort.Ints(positionKeys)
	return positionKeys
}


// ############# PARSER ELEMS #############################
// use minimal set of types. Don't overcomplicate the parser.
type ErlToken struct {

	// don't create a pointer to the prev/next Token.
	// it is a dangerous error source. Use only a positionId in a collector of ErlTokens

	TokenType  string
	TokenId int

	DebugStringRepresentation string // the string representation can be asked with a function, but in the debugger it is easier if it is stored in an attribute

	// one char can have a meaning alone, for example: (
	// but sometime more than one few char can form a token, for example: '->' which has his own meaning, but represented by 2 chars
	SourceCodeChars Chars
	TokenIsDetectedAsPartOfExpression bool
}

func (token ErlToken) typeIsEmpty() bool {
	return token.TokenType == ""
}
func (token ErlToken) charPositionFirstLast() (int, int) {
	charPosFirst := -1 // in a source code, 0 is the smallest position, so -1 means: no real position
	charPosLast := -1
	if len(token.SourceCodeChars) > 0 {
		charPosFirst = token.SourceCodeChars[0].PositionInFile
		charPosLast = token.SourceCodeChars[len(token.SourceCodeChars)-1].PositionInFile
	}
	return charPosFirst, charPosLast
}

func (token ErlToken) charPosFirst() int {
	charPosFirst, _ := token.charPositionFirstLast()
	return charPosFirst
}
func (token ErlToken) charPosLast() int {
	_, charPosLast := token.charPositionFirstLast()
	return charPosLast
}

func (token ErlToken) stringRepresentation() string {
	runes := []rune{}
	for _, charNow := range(token.SourceCodeChars) {
		runes = append(runes, charNow.Value)
	}
	return string(runes)
}

func (token ErlToken) stringRepresentation_escapedSourceForTests() string {
	runes := []rune{}
	for _, charNow := range(token.SourceCodeChars) {

		if charNow.Value =='"' {
			runes = append(runes, '\\')
			runes = append(runes, '"')
			continue
		}

		if charNow.Value =='\\' {
			runes = append(runes, '\\')
			runes = append(runes, '\\')
			continue
		}

		if charNow.Value =='\n' {
			runes = append(runes, '\\')
			runes = append(runes, 'n')
			continue
		}

		if charNow.Value =='\t' {
			runes = append(runes, '\\')
			runes = append(runes, 't')
			continue
		}


		runes = append(runes, charNow.Value)
	}
	return string(runes)
}

type Chars []Char

func (chars Chars) print_with_tokens(tokens ErlTokens) {
	/*

	to check the tokens, one character wide token signs are used.
	so the long %%%% means that where you see %, that is a comment

	    token type: %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	source    orig: % "double \"quoted\" comment, with a single quoted 'atom'"
	positionInFile: 01234567890123456789012345678901234567890123456789012345678

	a string example - here the strings and comments are detected only:
	??????????????"""""""""""""""""""""""""""""""""""""""""""""""??????????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
	example(1) -> "case 1 \\\" complex string \" with \n newline";                     % comment in 'example' function
	0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234


	% comments
	" strings
	a alphabet
	v variables
	8 numbers

	*/
	tokenElems := []rune{}
	charElems := []rune{}
	idElems := []rune{}

	tokenFlags := map[string]rune{
		"tokenTextBlockQuotedDouble": '"',
		"tokenComment": '%',
		"tokenTextBlockQuotedSingle": '\'',
		"tokenAbcFullWith_At_numbers": 'C',  // abc lower, upper plus atom builder extra chars
		                             // to avoid mixing it with atoms, C abC is the flag.
									 // and it stands for character, too ;-)
									 // it can be an atom, if the first char is uppercase,
									 // but in first token detection step, only alphabet is checked,
									 // which has capital letters, too
		"tokenOtherPunctuation": ':',
		"tokenWhiteSpace": 'w',
		"tokenCharLiteral": 'L',
	}

	for _, char := range chars {

		idRune := rune(str_from_int(char.PositionInLine % 10)[0])
		idElems = append(idElems, idRune)


		if char.Value == '\n' {
			charElems = append(charElems, ' ')
		} else {
			charElems = append(charElems, char.Value)
		}

		tokenFlag := '?'

		// tokens are stored by their position, not by their id.
		// if a token is detected for the char position, save the last position too, so we can print the whole token range
		// the token ID is not the best id to find it for chars, because tokens are stored by a position in file, not by id.
		if char.TokenDetected {
			tokenDetected, tokenInTheMap := tokens[char.TokenFirstCharPositionInFile]
			if tokenInTheMap {
				value, flagDetected := tokenFlags[tokenDetected.TokenType]
				if flagDetected {
					tokenFlag = value
				}
			}
		}


		tokenElems = append(tokenElems, tokenFlag)

		if char.Value == '\n' {
			fmt.Println()
			fmt.Println()
			runes_print(tokenElems)
			runes_print(charElems)
			runes_print(idElems)
			tokenElems = []rune{}
			charElems = []rune{}
			idElems = []rune{}
		}
	}
	fmt.Println()
	fmt.Println()
	runes_print(tokenElems)
	runes_print(charElems)
	runes_print(idElems)
}


type Char struct {
	PositionInFile  int // Nth char in the whole file
	Value                rune

	WhereTheCharIsStored string
	ErlangSourceWithoutFilePath bool

	TokenDetected bool
	TokenId int
	TokenFirstCharPositionInFile int  // the Token's first char's position in file

	LineNum int
	PositionInLine int  // Nth char in the line
}

/*
A lot of different versions of a module can be loaded and prepared to be executed same time.

	example situation: we have a 'database', 'financial' and 'logging' modules.

	# the module version attrib is a string, totally free, how do you define them. The teams can find out their own versioning
	# in the example, different teams are working with different modules,
	# and they are unable to use the same version naming convention, to represent an unperfect situation
	in the filesystem, the next versioned source codes are availabe (version name is a comment in the file)
	'database'  module available source code versions: "1.0.0 init", "1.1.0 connection fixes", "1.2.0 new db version"
	'financial' module available source code versions: "1.0a draft", "1.0b first working"
	'logging'   module available source code versions: "2023_10_22_a original", "2023_10_25_b", "2023_12_01_a", "2023_12_01_c"

	The architect defines the next executable version groups, where the 3 modules are able to work together:

	in the example you can see, that the 'financial' module is same in both release,
    but the database '1.0.0 init' can work only with 'logging' module's "2023_10_22_a original",

	So you can load/define more, than 2 different executable version sets, for your modules.



	FIXME: DEFINE EXACTLY WHERE the groups can be defined, and how is it handled
	Basically Erlango works like Erlang: you can load your source codes, and exec them.
	But if you define more version groups, more versions can be compiled/loaded/ready in the memmory,
	and in your incoming requests you can define which version group can execute them.


	With this solution, you can keep a lot of versions in the memory, and if you need,
	you can revert your changes.

	This versioning is different from hot code loading, because in that situation
    you modify an 'Executable_version_group' in a running system, and you change one element,
    which means that at the next func usage a new version will be executed.


	Executable_version_groups = #{
			"1_freely_defined_group_name" => [	['database',  "1.0.0 init"],
												['financial', "1.0b first working"],
												['logging',   "2023_10_22_a original"]
                                             ],
			"2_release_with_db_fixes"     => [  ['database',  "1.1.0 connection fixes"],
												['financial', "1.0b first working"],
												['logging', "2023_12_01_a"]
											]
	}

*/
type SourceTokensExecutables struct {
	ModuleName string

	// TODO: in source code loading, detect this, and set this attribute
	ModuleVersion string    // same module's different versions can be compiled and loaded same time
							// in your erl files, you can set it:
							// % ERLANGO_MODULE_VERSION your_version_definition

	WhereTheCodeIsStored string
	CharsFromErlFile     Chars
	Tokens               ErlTokens  // this section is filled by token detection part (step_01_tokens_from_source_code_of_files)

	Expressions ErlExpressions // this section is filled by expression detection (step_02_expressions_from_tokens_from_lot_of_sources)
	Errors errorsDetected
	ErlangSourceWithoutFilePath bool // if this is true, the source was not coming from a file
}

type SourcesTokensExecutables_map map[string]SourceTokensExecutables

func (sourceTokensExecutables SourceTokensExecutables) tokens_print()  {
	// mode: testCasePrint
	// mode: humanReading

	print("=== Tokens in a file: ", sourceTokensExecutables.WhereTheCodeIsStored, "===\n")

	// TODO: use a generic here
	keys := make([]int, 0, len(sourceTokensExecutables.Tokens))
	for k := range sourceTokensExecutables.Tokens {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		token := sourceTokensExecutables.Tokens[key]
		tokenPosFirst, tokenPosLast := token.charPositionFirstLast()
		// fmt.Println(key, "token:", token.TokenId, token.TokenType, tokenPosFirst, tokenPosLast, token.stringRepresentation() )
		// This format can be used in tests, immediatelly
		stringRepresentation := token.stringRepresentation_escapedSourceForTests()


		// TODO: stringrepresentation needs to be escaped " signs?
		fmt.Printf("{\"%s\", %v, %v, \"%s\"},\n", token.TokenType, tokenPosFirst, tokenPosLast, stringRepresentation)
	}
}


type errorsDetected []errorDetected

type errorDetected struct {
	filePath string
	lineNum  int
	charPosInLine int
	charPosInFile int
	errMsg string
}