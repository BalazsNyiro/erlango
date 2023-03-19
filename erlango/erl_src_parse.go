/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/

package erlango

import (
	"fmt"
)

func ParseErlangSourceFile() ([]ErlSrcChar, error) {
	chars, err := ErlSrcChars_from_file("test/parse/hello.erl")
	if err != nil { return []ErlSrcChar{}, err}

	// detect "strings" or 'atoms' - quoted texts
	ErlSrcTokens_Quoted__connect_to_chars(chars, true)

	// detect comments

	// detect whitespaces

	// detect numbers
	// detect

	return chars, nil
}

///////////////////// Globals////////////////////////////////////////////
// this is a perfect theoretical example for an atom, because
// the value here is not important, useless.
// maybe in debugging it's easier to see something instead of a flag
const Token_type_txt_quoted_double string = "txt_quoted_double"  // "abc"
const Token_type_txt_quoted_single string = "txt_quoted_single"  // 'abc'
////////////////////////////////////////////////////////////////////////

// ErlSrcToken : independent language unit, formed by one or more char
// they are character holders, they group the characters,
// if the characters form one meaning.
// for example '123' text has 3 symbols, and they are
// represented by 3 ErlSrcChar elems,
// and they are stored in one Token because they represent one number
//
// Same token can have a totally different meaning at the end,
// depends on the environment.
// for example "name" can be a key in a map, a string, or be a binary elem, too.
// so these token's don't have any meaning at this point
type ErlSrcToken struct {
	PrevToken *ErlSrcToken
	NextToken *ErlSrcToken
	Chars     []*ErlSrcChar
	Type      string
}
func (token ErlSrcToken) CharAppend(charPtr *ErlSrcChar) {
	token.Chars = append(token.Chars, charPtr)
}

//////////////////////////////////////////////////////////////////////
type ErlSrcTokens []ErlSrcToken
func (tokens ErlSrcTokens) IdLast() int {
	return len(tokens) - 1 // it always has minimum 1 value because of Pre-init:
	                       // in tokensForChars__preInitialized()
}

func (tokens ErlSrcTokens) LastPtr() *ErlSrcToken {
	return &(tokens[tokens.IdLast()])
}
//////////////////////////////////////////////////////////////////////

// ErlSrcChar represents one char in the Erlang source codes
type ErlSrcChar struct {
	NextChar   *ErlSrcChar
	PrevChar   *ErlSrcChar
	PosInFile  int
	Value      rune
	Token      *ErlSrcToken
	SourcePath string
}


// Type a char's type is the parent Token's type
func (char ErlSrcChar) Type () string {
	if char.Token == nil {
		return ""
	}
	return char.Token.Type
}

func ErlSrcChars_from_file(filePath string) ([]ErlSrcChar, error) {
	runes, err := file_read_runes(filePath, "ErlSrcChars_from_file")
	if err != nil { return []ErlSrcChar{}, err}
	erlChars := ErlSrcChars_from_runes(runes, filePath)
	// Test_what_happens_with_struct_pointers
	// fmt.Printf("ErlSrcChars_from_file, chars pointer before return: %p\n", erlChars)
	return erlChars, nil
}

func ErlSrcChars_from_str(txt string) []ErlSrcChar {
	runes := runes_from_str(txt)
	return ErlSrcChars_from_runes(runes, "direct_txt_input")
}

func ErlSrcChars_from_runes(runes []rune, sourcePath string) []ErlSrcChar {
	var erlChars []ErlSrcChar
	for posInFile, runeInFile := range runes {
		erlChars = append(erlChars, ErlSrcChar{
			Value:      runeInFile,
			PosInFile:  posInFile,
			SourcePath: sourcePath,
		})
	}
	// after the first for loop exec, the slice size is finalised.
	// when I used one for loop first time, the slice was changed
	// when it reached the capacity limit, and the pointers were incorrect.

	// the slice pointers won't be changed after this point,
	// there is no capacity change later.
	// if we do this from the 'previous linking position'
	// then because of the capacity limit reach, the pointers
	// will be incorrect in the early elements
	for id, _ := range erlChars {
		if id > 0 {
			erlChars[id].PrevChar = &erlChars[id-1]
			erlChars[id-1].NextChar = &erlChars[id]
		}
	}
	return erlChars
}

/* ErlSrcTokens_Quoted__connect_to_chars fun processes the chars one by one:
    - if this is in a Quote: char->Token pointing happens.
    - more than one char can be connected to the same token.

    char_1  ↘
    char_2 → Token - collector, a lot of chars are linked into one Token
    char_3  ↗

	The function connects new Tokens to the existing characters,
    this is the reason why there is no return value here.

    arrows: https://en.wikipedia.org/wiki/Arrows_(Unicode_block)


    ### newline handling in quoted texts ###
    This implementation eats everything between '...' or "..." pairs.
    so here it works:

				 A := " line 1, not closed with quota
						line 2, finished with quota sign "
    The programmer can insert newline into strings with "line1..." ++ "\nline2"
    So now this behaviour is not a problem.
 */
func ErlSrcTokens_Quoted__connect_to_chars(chars []ErlSrcChar, verbose bool) {
	ErlSrcTokens_rangeDetect__connectToChars(
		chars,
		quoteConditionOpener,
		quoteConditionCloser,
		quoteConditionEscape,
		quoteTokenTypeSet,
		verbose)
}
func ErlSrcTokens_rangeDetect__connectToChars(
		chars []ErlSrcChar,
	 	conditionOpener func([]ErlSrcChar, int, *conditionMemory) bool,
		conditionCloser func([]ErlSrcChar, int, *conditionMemory) bool,
		conditionEscape func([]ErlSrcChar, int, *conditionMemory) bool,
	    tokenTypeSetter func(*ErlSrcTokens, *conditionMemory),
		verbose bool) {

	tokenInfo := func (position int, chars []ErlSrcChar, tokens ErlSrcTokens, inCharRange bool, memory conditionMemory ) {
		fmt.Println("ErlSrcTokens_Quoted__connect_to_chars", position, string(chars[position].Value),
			fmt.Sprintf("tokenPtr: %p", chars[position].Token),
			"type->",chars[position].Type(), "<>", (*tokens.LastPtr()).Type, "<- ",
			bool_to_str(inCharRange, "in Quote:"+string(memory.runes["actualQuoteChar"]), ""))
	}

	tokens := tokensForChars__preInitialized()
	conditionMemory := conditionMemoryEmpty()
	inCharRange, escapeOn := false, false

	for position, _ := range chars {
		nowOpened, nowEscaped := false, false

		if !inCharRange && conditionOpener(chars, position, &conditionMemory) {
			tokenTypeSetter(&tokens, &conditionMemory)
			inCharRange, nowOpened = true, true
		}

		if !escapeOn && inCharRange && conditionEscape(chars, position, &conditionMemory) {
			escapeOn, nowEscaped = true, true // escaping is important for the closing condition
		}

		if inCharRange {
			chars[position].Token = tokens.LastPtr()
			chars[position].Token.CharAppend(&(chars[position]))
		}
		if verbose { tokenInfo(position, chars, tokens, inCharRange, conditionMemory) }

		if nowOpened || nowEscaped { continue }
		// ##### stop here ^^^^ the char processing in these 2 cases ###########
		// if nowOpened == true, the sign is '\' and I don't want to turn it off if it was turned on just now
		// if it's nowEscaped, I don't want to turn it off too because it has effect on the next char

		if !escapeOn && inCharRange && conditionCloser(chars, position, &conditionMemory) {
			inCharRange = false  // active escape blocks the conditionCloser()
			tokens = append(tokens, tokenEmpty())
		}
		escapeOn = false // if not now escaped, the escape disappearing at the next char.
	}
}

///////////////// token opener/closer //////////////////
func conditionMemoryEmpty() conditionMemory {
	return conditionMemory{runes: map[string]rune{}}
}
type conditionMemory struct {
	nums map[string]int
	strings map[string]string
	runes map[string]rune
}

// this is the first char in the range if returns with true
func quoteConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	result := isSingleOrDoubleQuoteRune(chars[position].Value)
	if result {
		memory.runes["actualQuoteChar"] = chars[position].Value
	}
	return result
}

// this is the last char in the range
func quoteConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return chars[position].Value == memory.runes["actualQuoteChar"]
}

// skip the next char if it returns with true
func quoteConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return chars[position].Value == '\\'
}

func quoteTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	tokenIdLast := len(*tokens) - 1
	if isSingleQuoteRune(memory.runes["actualQuoteChar"]) {
		(*tokens)[tokenIdLast].Type = Token_type_txt_quoted_single
	} else {
		(*tokens)[tokenIdLast].Type = Token_type_txt_quoted_double
	}

}
///////////////// token opener/closer //////////////////

////////////////////////////////// token funs ////////////////////////////////////
func isSingleQuoteRune(r rune) bool { return r == '\''}
func isDoubleQuoteRune(r rune) bool { return r == '"'}
func isSingleOrDoubleQuoteRune (r rune) bool {return isSingleQuoteRune(r) || isDoubleQuoteRune(r)}
func tokenEmpty() ErlSrcToken                       { return ErlSrcToken{} }
func tokensForChars__preInitialized() ErlSrcTokens { return ErlSrcTokens{tokenEmpty()} }
//  ^^^^ // in Go, a variable's memory address stay the same when you assign a new value.
// so, I can use a token only once - it's necessary to generate always new tokens,
// and a simple 'tokenActual = tokenEmpty()' can't work, if the variable is always the same,
// because if I pass its pointer, later I can overwrite the value behind the variable.
// the current solution generates new tokens into a list, and the last elem is always
// updated, so it will have a new address after each update
////////////////////////////////// token funs ////////////////////////////////////