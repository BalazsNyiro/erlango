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

// ############# PARSER ELEMS #############################
// use minimal set of types. Don't overcomplicate the parser.
type ErlToken struct {

	// don't create a pointer to the prev/next Token.
	// it is a dangerous error source. Use only a positionId in a collector of ErlTokens

	/* Possible token types:
	 - notDetectedFromCharacters
	 - atom
	 - string
	 - binaryDoubleOpener      <<
	 - binaryDoubleCloser      >>
	 - operatorMathPlusPlus    ++
	 - operatorMathMinusMinus  --
	 - operatorBoolEqualDouble ==
	 - operatorBoolBigger      >
	 - operatorBoolSmaller     <
	 - codeBlockStart         ->

	a few sign can have different meanings in different places:
	a comma can be expression separator, or parameter separator - so it's later decided how is it handled
	 - langComma ,
	 - langDot   .
	 - langColon :
	 - langBraceRoundOpen   (
	 - langBraceRoundClose  )
	 - langBraceSquaredOpen [
	 - langBraceSqaredClose ]
	 - langVariableBound    =

	*/
	TokenType  string
	TokenId int

	// one char can have a meaning alone, for example: (
	// but sometime more than one few char can form a token, for example: '->' which has his own meaning, but represented by 2 chars
	SourceCodeChars Chars
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

type Chars []Char

type Char struct {
	PositionInFile  int
	Value      rune
	FilePath string
	TokenDetected bool
	TokenId int
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

	PathErlFile          string
	CharsFromErlFile     Chars
	Tokens               ErlTokens

	ExecutableFromTokens string // FIXME: This is a pile of executable objects, not a string
}

type SourcesTokensExecutables_map map[string]SourceTokensExecutables

func (sourceTokensExecutables SourceTokensExecutables) tokens_print()  {
	print("=== Tokens in a file", sourceTokensExecutables.PathErlFile, "===\n")

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
		stringRepresentation := token.stringRepresentation()
		if stringRepresentation[0] == '"' {
			stringRepresentation = "\\" + stringRepresentation
		}
		if stringRepresentation[len(stringRepresentation)-1] == '"' {
			stringRepresentation = stringRepresentation[:len(stringRepresentation)-1] + "\\\""
		}
		// TODO: stringrepresentation needs to be escaped " signs?
		fmt.Printf("{\"%s\", %v, %v, \"%s\"}\n", token.TokenType, tokenPosFirst, tokenPosLast, stringRepresentation)
	}
}
