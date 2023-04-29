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
	"strings"
)

func ParseErlangSourceFile() ([]ErlSrcChar, error) {
	chars, err := ErlSrcChars_from_file("test/parse/hello.erl")
	if err != nil { return []ErlSrcChar{}, err}
	return ParseErlangSourceCode(chars, "__all__")
}

/*  What is a token? (my terminology).
    Expression: the smallest usable language unit with meaning:
                example a: 1
                example b: 23.456

    but in the source code: the float 2.3 has 3 tokens:
                token 1: 23
                token 2: .
                token 3: 456

     Why do we have 3 tokens in that float?
       because the dot (.) can have different meanings in different environments.
       so the interpreter detects the <digits+dot+digits>
       and later the interpreter realises that if we have a digit-dot-digit combo, than that is a float num

       in other situations the dot has different semantic meaning, so the interpreter has to understand the circumstances.

	So: a token is a unit, that maybe has, maybe doesn't have standalone language meanings,
        and more token can build an expression.

    with other words: a token is a building block of expressions.

*/


/* New Token type definition steps:

	- in ParseErlangSourceCode():
       if execStep(<FIND-A-GOOD-EXECUTION-STEP-NAME) { <CREATE-STANDARD-TOKEN-DETECTOR-FUN-NAME> }

	- define a new constant, example:
        const Token_type_variable  string = "Token_type_digits_baseDefined"

    - Detector creation:
        ErlSrcTokensDetect_______variables_______connect_to_chars

    - CREATE funs: opener, closer, escape, typesetter

    - in tests create a string -> constant binding to detect token type names from tests, for example:
      "Token_type_always_accepted" : "Token_type_always_accepted",

    - create tests
*/

// https://www.erlang.org/doc/reference_manual/expressions.html#expression-evaluation
// https://www.tutorialspoint.com/erlang/erlang_operators.htm
// the order of stepsWanted doesn't important, because the here defined execStep() order is the competent.

// important: this whole function is based on Shallow-Copy. chars is passed and copied in the called functions,
// but inside the structs are the same. So, all func changes only structs, and chars are not changed.
func ParseErlangSourceCode(chars []ErlSrcChar, stepsWanted string) ([]ErlSrcChar, error) {
	// detect "strings" or 'atoms' - quoted texts

	execStep := func (stepName string) bool {
		if strings.Contains(stepsWanted, "__all__") || strings.Contains(stepsWanted, stepName) {
			return true
		} else {return false}
	}

	// when you call ParseErlangSourceCode(), you can pass which steps do you want to execute
	// so different steps can be executed from different tests
	verbose := false
	if execStep("strings_atoms_quotes") { ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(chars, verbose) }
	if execStep("comments")             { ErlSrcTokensDetect________comments_______connect_to_chars(chars, verbose) }
	if execStep("whitespaces")          { ErlSrcTokensDetect______whitespaces______connect_to_chars(chars, verbose) }
	if execStep("commas")               { ErlSrcTokensDetect________commas_________connect_to_chars(chars, verbose) }
	if execStep("dots")                 { ErlSrcTokensDetect__________dot__________connect_to_chars(chars, verbose) }
	if execStep("semicolons")           { ErlSrcTokensDetect_______semicolon_______connect_to_chars(chars, verbose) }

	if execStep("bracket_round_opener") { ErlSrcTokensDetect____bracketRoundOp_____connect_to_chars(chars, verbose) }
	if execStep("bracket_round_closer") { ErlSrcTokensDetect____bracketRoundCl_____connect_to_chars(chars, verbose) }

	if execStep("variables")            { ErlSrcTokensDetect_______variables_______connect_to_chars(chars, verbose) }
	if execStep("atoms_quoteless")      { ErlSrcTokensDetect____atoms_quoteless____connect_to_chars(chars, verbose) }

	// TODO: detect nums
	// baseDefined is a more wider range than base10
	// // digits can be in atoms/variableNames too, so this section has to be after variables/atoms
	if execStep("digits_baseDefined") { }
	// here we can't detect dots/float nums, because . is not a digit and the float detection exists before the first dot
	if execStep("digits_base10_form")   { ErlSrcTokensDetect_____digits_base10_____connect_to_chars(chars, verbose) }
	if execStep("numbers_floats")       { chars = multi_token_detect(chars, verbose, detector_tokens_floats)        }

	// arrows:  ->    <-    =>
	if execStep("arrow_singleToRight")  { ErlSrcTokensDetect__arrow_singleToRight__connect_to_chars(chars, verbose) } // ->
	if execStep("arrow_singleToLeft")   { ErlSrcTokensDetect__arrow_singleToLeft___connect_to_chars(chars, verbose) } // <-
	if execStep("arrow_doubleToRight")  { ErlSrcTokensDetect__arrow_doubleToRight__connect_to_chars(chars, verbose) } // =>


	if execStep("binding_matching")     { ErlSrcTokensDetect____binding_matching___connect_to_chars(chars, verbose) }

	if execStep("math_binary_add")      { ErlSrcTokensDetect____math_binary_add____connect_to_chars(chars, verbose) }
	if execStep("math_binary_sub")      { ErlSrcTokensDetect____math_binary_sub____connect_to_chars(chars, verbose) }
	if execStep("math_binary_mul")      { ErlSrcTokensDetect____math_binary_mul____connect_to_chars(chars, verbose) }
	if execStep("math_binary_div")      { ErlSrcTokensDetect____math_binary_div____connect_to_chars(chars, verbose) }

	return chars, nil
}

///////////////////// Globals////////////////////////////////////////////
// this is a perfect theoretical example for an atom, because
// the value here is not important, useless.
// maybe in debugging it's easier to see something instead of a flag
const Token_type_txt_quoted_double     string = "txt_quoted_double"     // "abc"
const Token_type_txt_quoted_single     string = "txt_quoted_single"     // 'abc'
const Token_type_comment               string = "comment"               // % abc
const Token_type_not_detected          string = "noTokenConnected"
const Token_type_whitespace            string = "separator_whitespace"  // \t ' ' \n
const Token_type_comma                 string = "separator_comma"       // ,
const Token_type_dot                   string = "separator_dot"         // .
const Token_type_semicolon             string = "separator_semicolon"   // ;
const Token_type_bracket_round_open    string = "bracket_round_open"    // (
const Token_type_bracket_round_close   string = "bracket_round_close"   // )
const Token_type_bracket_square_open   string = "bracket_square_open"   // [
const Token_type_bracket_square_close  string = "bracket_square_close"  // ]
const Token_type_bracket_curly_open    string = "bracket_curly_open"    // {
const Token_type_bracket_curly_close   string = "bracket_curly_close"   // }

const Token_type_digits_base10_form    string = "Token_type_digits_base10_form"  // 1234567890
const Token_type_digits_baseDefined    string = "Token_type_digits_baseDefined"  // 16#af6bfa23, only whole nums, there is no 16 based float
const Token_type_float_dotInDigits     string = "Token_type_float_dotInDigits"   // 12.34

const Token_type_variable              string = "Token_type_digits_baseDefined"  // ErlangVariableName :-)
const Token_type_atom_quoteless        string = "Token_type_atom_quoteless"      // erlang_atom_defined_without_quotes

const Token_type_arrow_singleToRight   string = "Token_type_arrow_singleToRight" // ->
const Token_type_arrow_singleToLeft    string = "Token_type_arrow_singleToLeft"  // <-
const Token_type_arrow_doubleToRight   string = "Token_type_arrow_singleToRight" // ->

const Token_type_binding_matching      string = "Token_type_binding_matching"    // =

const Token_type_math_binary_add      string = "Token_type_math_binary_add"    // +
const Token_type_math_binary_sub      string = "Token_type_math_binary_sub"    // +
const Token_type_math_binary_mul      string = "Token_type_math_binary_mul"    // +
const Token_type_math_binary_div      string = "Token_type_math_binary_div"    // +


const Token_type_deleted_dont_use      string = "Token_type_deleted_dont_use"    // in some situations from more Tokens the intrepreter creates one.
                                                                                 // in float detection for example: digit.digit

const ABC_Eng_Upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ABC_Eng_Lower string = "abcdefghijklmnopqrstuvwxyz"
const ABC_Eng_digits string = "0123456789"
const ABC_Eng_alphanum string = ABC_Eng_Upper + ABC_Eng_Lower + ABC_Eng_digits

const ErlangVariableOpener = ABC_Eng_Upper + "_"
const ErlangVariableBody = ABC_Eng_alphanum + "_@"
const ErlangAtomNoQuotesOpener = ABC_Eng_Lower + "_"
const ErlangAtomNoQuotesBody = ABC_Eng_alphanum + "_@"
// //////////////////////////////////////////////////////////////////////

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

func (token *ErlSrcToken) CharAppend(charPtr *ErlSrcChar) {
	token.Chars = append(token.Chars, charPtr)
}

func (token ErlSrcToken) StrValueFromChars() string {
	ret := ""
	for _, chrPtr := range token.Chars {
		ret += string(chrPtr.Value)
	}
	return ret
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
	if ! char.TokenConnected() {
		return Token_type_not_detected
	}
	return char.Token.Type
}

// true: if the char is connected to a token
func (char ErlSrcChar) TokenConnected () bool {
	return char.Token != nil
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

/* ErlSrcTokensDetect___string_atom_quotes__connect_to_chars fun processes the chars one by one:
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
func ErlSrcTokensDetect___string_atom_quotes__connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		quoteConditionOpener,
		quoteConditionCloser,
		quoteConditionEscape,
		quoteTokenTypeSet,
		false,
		verbose,
		"parse_strings_atoms",
		false,  // the opener char cannot be the closer: we find pairs: "..."
	)
}

func ErlSrcTokensDetect________comments_______connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		commentConditionOpener,
		commentConditionCloser,
		commentConditionEscape,
		commentTokenTypeSet,
		false,
		verbose,
		"parse comments",
		false,  // the opener char cannot be the closer: comment...newLine
		)
}

func ErlSrcTokensDetect______whitespaces______connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		whitespacesConditionOpener,
		whitespacesConditionCloser,
		whitespacesConditionEscape,
		whitespacesTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse whitespaces",
		true,  // the opener char is the closer char in same time
	)
}

func ErlSrcTokensDetect________commas_________connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		commaConditionOpener,
		commaConditionCloser,
		commaConditionEscape,
		commaTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse commas",
		true,  // the opener char is the closer char in same time
	)
}

func ErlSrcTokensDetect__________dot__________connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		dotConditionOpener,
		dotConditionCloser,
		dotConditionEscape,
		dotTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse dots",
		true,  // the opener char is the closer char in same time
	)
}

func ErlSrcTokensDetect_______semicolon_______connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		semicolonConditionOpener,
		semicolonConditionCloser,
		semicolonConditionEscape,
		semicolonTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse semicolons",
		true,  // the opener char is the closer char in same time
	)
}

func ErlSrcTokensDetect____bracketRoundOp_____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		bracketRoundOpConditionOpener,
		bracketRoundOpConditionCloser,
		bracketRoundOpConditionEscape,
		bracketRoundOpTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse bracket Round opener",
		true,  // the opener char is the closer char in same time
	)
}

func ErlSrcTokensDetect____bracketRoundCl_____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		bracketRoundClConditionOpener,
		bracketRoundClConditionCloser,
		bracketRoundClConditionEscape,
		bracketRoundClTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse bracket Round opener",
		true,  // the opener char is the closer char in same time
	)
}



func ErlSrcTokensDetect_____digits_base10_____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		digitsBase10ConditionOpener,
		digitsBase10ConditionCloser,
		digitsBase10ConditionEscape,
		digitsBase10TokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse digits base10",
		true,  // the opener char can be closer, too
	)
}

func ErlSrcTokensDetect_______variables_______connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		variablesConditionOpener,
		variablesConditionCloser,
		variablesConditionEscape,
		variablesTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse variables",
		true,  // variable name can be 1 char long, too
	)
}

func ErlSrcTokensDetect____atoms_quoteless____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		atomsConditionOpener,
		atomsConditionCloser,
		atomsConditionEscape,
		atomsTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse digits base10",
		true,  // atom name can be 1 char long
	)
}

//////////// ARROWS  ->   <-  =>
func  ErlSrcTokensDetect__arrow_singleToRight__connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		arrowSingleToRightConditionOpener,
		arrowSingleToRightConditionCloser,
		arrowSingleToRightConditionEscape,
		arrowSingleToRightTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse arrow single to right",
		false,  // because it's longer than 1 char
	)
}
func ErlSrcTokensDetect__arrow_singleToLeft___connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		arrowSingleToLeftConditionOpener,
		arrowSingleToLeftConditionCloser,
		arrowSingleToLeftConditionEscape,
		arrowSingleToLeftTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse arrow single to left",
		false,  // because it's longer than 1 char
	)
}
func ErlSrcTokensDetect__arrow_doubleToRight__connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		arrowDoubleToRightConditionOpener,
		arrowDoubleToRightConditionCloser,
		arrowDoubleToRightConditionEscape,
		arrowDoubleToRightTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse arrow double to right",
		false,  // because it's longer than 1 char
	)
}


/// binding_matching

func  ErlSrcTokensDetect____binding_matching___connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		bindingMatchingConditionOpener,
		bindingMatchingConditionCloser,
		bindingMatchingConditionEscape,
		bindingMatchingTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse =",
		true,  // =  is 1 char long
	)
}

/////////////////  math binary operators //////////////////////

func ErlSrcTokensDetect____math_binary_add____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		mathBinaryAddConditionOpener,
		mathBinaryAddConditionCloser,
		mathBinaryAddConditionEscape,
		mathBinaryAddTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse math binary add",
		true,  // +  is 1 char long
	)
}
func ErlSrcTokensDetect____math_binary_sub____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		mathBinarySubConditionOpener,
		mathBinarySubConditionCloser,
		mathBinarySubConditionEscape,
		mathBinarySubTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse math binary sub",
		true,  // - is one char long
	)
}
func ErlSrcTokensDetect____math_binary_mul____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		mathBinaryMulConditionOpener,
		mathBinaryMulConditionCloser,
		mathBinaryMulConditionEscape,
		mathBinaryMulTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse math binary mul",
		true,  // *  is 1 char long
	)
}


func ErlSrcTokensDetect____math_binary_div____connect_to_chars(chars []ErlSrcChar, verbose bool) {
	erlSrcTokens_rangeDetect__connectToChars(
		chars,
		mathBinaryDivConditionOpener,
		mathBinaryDivConditionCloser,
		mathBinaryDivConditionEscape,
		mathBinaryDivTokenTypeSet,
		true, // skip chars with tokens
		verbose,
		"parse math binary div",
		true,  //  /(division)  is 1 char long
	)
}


///////////////////////////////////////////////////////////////////////
func erlSrcTokens_rangeDetect__connectToChars(
		chars []ErlSrcChar,
	 	conditionOpener func([]ErlSrcChar, int, *conditionMemory) bool,
		conditionCloser func([]ErlSrcChar, int, *conditionMemory) bool,
		conditionEscape func([]ErlSrcChar, int, *conditionMemory) bool,
	    tokenTypeSetter func(*ErlSrcTokens, *conditionMemory),
		skip_chars_with_tokens bool,
		verbose bool, caller string,
		canBeOneCharWideTokenDetection bool,
		) {

	tokenInfo := func (position int, chars []ErlSrcChar, tokens ErlSrcTokens, inCharRange bool, memory conditionMemory ) {
		fmt.Println("ErlSrcTokensDetect", caller, position, string(chars[position].Value),
			fmt.Sprintf("tokenPtr: %p", chars[position].Token),
			"type->",chars[position].Type(), "<>", (*tokens.LastPtr()).Type, "<- ",
			bool_to_str(inCharRange, "in Quote:"+string(memory.runes["actualQuoteChar"]), ""))
	}

	tokens := tokensForChars__preInitialized()
	conditionMemoryTemporaryWorkspace := conditionMemoryEmpty()
	inTokenDetection_activeCharsFound, escapeOn := false, false

	for position, charNow := range chars {
		if verbose { fmt.Println("Token to char, charNow:", charNow) }

		if skip_chars_with_tokens && chars[position].TokenConnected() { continue } // modify only the unprocessed chars, without Tokens
		nowOpened, nowEscaped := false, false

		if !inTokenDetection_activeCharsFound && conditionOpener(chars, position, &conditionMemoryTemporaryWorkspace) {
			tokenTypeSetter(&tokens, &conditionMemoryTemporaryWorkspace)
			inTokenDetection_activeCharsFound, nowOpened = true, true
		}

		if !escapeOn && inTokenDetection_activeCharsFound && conditionEscape(chars, position, &conditionMemoryTemporaryWorkspace) {
			escapeOn, nowEscaped = true, true // escaping is important for the closing condition
		}

		if inTokenDetection_activeCharsFound {
			chars[position].Token = tokens.LastPtr()
			chars[position].Token.CharAppend(&(chars[position]))
		}
		if verbose { tokenInfo(position, chars, tokens, inTokenDetection_activeCharsFound, conditionMemoryTemporaryWorkspace) }

		if !canBeOneCharWideTokenDetection { // here we know that the wanted token leng is bigger than 1
			if nowOpened || nowEscaped { continue } // the opener cannot be the closer: ".." pairs for example
		} // else: if oneCharWideToken == true, the char can be a closer char, too


		// ##### stop here ^^^^ the char processing in these 2 cases ###########
		// if nowOpened == true, the sign is '\' and I don't want to turn it off if it was turned on just now
		// if it's nowEscaped, I don't want to turn it off too because it has effect on the next char

		if !escapeOn && inTokenDetection_activeCharsFound && conditionCloser(chars, position, &conditionMemoryTemporaryWorkspace) {
			inTokenDetection_activeCharsFound = false // active escape blocks the conditionCloser()
			tokens = append(tokens, tokenEmpty())
		}
		escapeOn = false // if not now escaped, the escape disappearing at the next char.
	} // for
}

///////////////// token opener/closer //////////////////
// conditionMemory is a place where the opener/closer/other funs can save their infos during the detection.
func conditionMemoryEmpty() conditionMemory {
	return conditionMemory{bools: map[string]bool{}, runes: map[string]rune{}}
}
type conditionMemory struct {
	nums map[string]int
	bools map[string]bool
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
	if isSingleQuoteRune(memory.runes["actualQuoteChar"]) {
		generalTokenTypeSetThis(tokens, memory, Token_type_txt_quoted_single)
	} else {
		generalTokenTypeSetThis(tokens, memory, Token_type_txt_quoted_double)
	}
}

///////////////////////////////////////
func commentConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	if chars[position].TokenConnected() { return false } // "in text, % is not a comment"
	return chars[position].Value == '%'
}


// this is a special situation, the general closer is a little different.
func commentConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	lenChars := len(chars)
	nextPos := position+1
	lastPos := lenChars-1
	if position == lastPos { return true}	// this is the last char, we won't find a newline.

	if nextPos <= lastPos {

		// if the next char is not in token and the next char is a newline
		if (! chars[nextPos].TokenConnected()) &&  chars[nextPos].Value == '\n' {
			return true	 // the newline is not part of the comment Token
		} else {
			return false
		}
	}
	return true
}

func commentConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in comments
}

func commentTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_comment)
}
///////////////////////////////////////

func whitespacesConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, " \r\n\t")
}

func whitespacesConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func whitespacesConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in whitespaces
}

func whitespacesTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_whitespace)
}

/////////////////// comma ////////////////////

func commaConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ",")
}

func commaConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func commaConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in comma
}

func commaTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_comma)
}

/////////////////// dot ////////////////////

func dotConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ".")
}

func dotConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func dotConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in dot
}

func dotTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_dot)
}

/////////////////// semicolon ////////////////////

func semicolonConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ";")
}

func semicolonConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func semicolonConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in semicolon
}

func semicolonTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_semicolon)
}

/////////////////// bracket round opener ////////////////////

func bracketRoundOpConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "(")
}

func bracketRoundOpConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func bracketRoundOpConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in bracketRoundOp
}

func bracketRoundOpTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_bracket_round_open)
}

/////////////////// bracket round closer ////////////////////

func bracketRoundClConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ")")
}

func bracketRoundClConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true // the opener is a closer in same time
}

func bracketRoundClConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in bracketRoundCl
}

func bracketRoundClTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_bracket_round_close)
}
/////////////////// bracket round closer ////////////////////


func digitsBase10ConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ABC_Eng_digits)
}

func digitsBase10ConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionCloser(chars, position, ABC_Eng_digits)
}

func digitsBase10ConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in digitsBase10
}

func digitsBase10TokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_digits_base10_form)
}




/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////  variables, quoteless-atoms //////////////////////////

// Erlang variable can start with a capital letter or underscore
func variablesConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ErlangVariableOpener)
}

func variablesConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionCloser(chars, position, ErlangVariableBody)
}

func variablesConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in variables
}

func variablesTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_variable)
}

func atomsConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, ErlangAtomNoQuotesOpener)
}

func atomsConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionCloser(chars, position, ErlangAtomNoQuotesBody)
}

func atomsConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in atoms
}

func atomsTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_atom_quoteless)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////




////////////// ARROWS ->  <-  <= ////////////
func arrowSingleToRightConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerMultiCharsInPattern(chars, position, memory, "->")
}

func arrowSingleToRightConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return chars[position].Value == '>'
}

func arrowSingleToRightConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false
}

func arrowSingleToRightTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_arrow_singleToRight)
}

func arrowDoubleToRightConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerMultiCharsInPattern(chars, position, memory, "=>")
}

func arrowDoubleToRightConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return chars[position].Value == '>'
}

func arrowDoubleToRightConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false
}

func arrowDoubleToRightTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_arrow_doubleToRight)
}

func arrowSingleToLeftConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerMultiCharsInPattern(chars, position, memory, "<-")
}

func arrowSingleToLeftConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return chars[position].Value == '-'
}

func arrowSingleToLeftConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false
}

func arrowSingleToLeftTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_arrow_singleToLeft)
}



////////////// binding-matching ////////////
func bindingMatchingConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "=")
}

func bindingMatchingConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true
}

func bindingMatchingConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in bindingMatching
}

func bindingMatchingTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_binding_matching)
}




/// math binary operators ///

func mathBinaryAddConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "+")
}

func mathBinaryAddConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true
}

func mathBinaryAddConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in math operators
}

func mathBinaryAddTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_math_binary_add)
} ///


func mathBinarySubConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "-")
}

func mathBinarySubConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true
}

func mathBinarySubConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in math operators
}

func mathBinarySubTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_math_binary_sub)
} ///

func mathBinaryMulConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "*")
}

func mathBinaryMulConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true
}

func mathBinaryMulConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in math operators
}

func mathBinaryMulTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_math_binary_mul)
} ///

func mathBinaryDivConditionOpener(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return generalConditionOpenerCharInPattern(chars, position, memory, "/")
}

func mathBinaryDivConditionCloser(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return true
}

func mathBinaryDivConditionEscape(chars []ErlSrcChar, position int, memory *conditionMemory) bool {
	return false // there is no meaning of an escape in math operators
}

func mathBinaryDivTokenTypeSet(tokens *ErlSrcTokens, memory *conditionMemory) {
	generalTokenTypeSetThis(tokens, memory, Token_type_math_binary_div)
}



//////////////  general opener, type setter /////////////////////////

// test one char: can it be an opener?
func generalConditionOpenerCharInPattern(chars []ErlSrcChar, position int, memory *conditionMemory, pattern string) bool {
	lenChars := len(chars)
	lastPos := lenChars-1

	if position > lastPos { return false } // can't look over the end of the chars
	if chars[position].TokenConnected() { return false }

	charNow := string(chars[position].Value)
	return strings.Contains(pattern, charNow)
}

// test series of chars, so in one step we can say if the first char can be an opener for a multi-char wide opener
func generalConditionOpenerMultiCharsInPattern(chars []ErlSrcChar, position int, memory *conditionMemory, pattern string) bool {
	for patternRelativePos, charInPattern := range pattern {
		// if any of the relative Opener test is false, the whole test is failed
		if ! generalConditionOpenerCharInPattern(chars, position+patternRelativePos, memory, string(charInPattern)) {
			return false
		}
	}
	return true
}

func generalConditionCloser(chars []ErlSrcChar, position int, validPossibleBodyPattern string) bool {
	lenChars := len(chars)
	nextPos := position+1
	lastPos := lenChars-1

	if nextPos <= lastPos {
		if chars[nextPos].TokenConnected() { return true }
		if ! strings.Contains(validPossibleBodyPattern, string(chars[nextPos].Value)) {
			return true	 // this is a closer, because nextPos hasn't got a valid body pattern
		} else {
			return false
		}
	}
	return true  // because nextpos > lastPos


}

func generalTokenTypeSetThis(tokens *ErlSrcTokens, memory *conditionMemory, typeNew string) {
	tokenIdLast := len(*tokens) - 1
	(*tokens)[tokenIdLast].Type = typeNew
}

///////////////// token opener/closer //////////////////


func tokens_from_chars(chars []ErlSrcChar) []*ErlSrcToken {
	TokensPtrs := []*ErlSrcToken{}
	for _, charNow := range chars {
		if charNow.Token != nil {
			TokensPtrs = append(TokensPtrs, charNow.Token)
		}
	}
	return TokensPtrs
}

// digit-dot-digit token combo -> this is a float num.
// this fun works with Tokens, not with chars.
// so take the first Token
func detector_tokens_floats(tokenPtrs []*ErlSrcToken, tokenIdActual int, verbose bool) {

		tokenId_prev_1 := tokenIdActual - 1
		tokenId_prev_2 := tokenId_prev_1 - 1

		if tokenId_prev_2 < 0 { return }

		if tokenPtrs[tokenId_prev_2].Type == Token_type_digits_base10_form &&
	 	   tokenPtrs[tokenId_prev_1].Type == Token_type_dot &&
			tokenPtrs[tokenIdActual].Type == Token_type_digits_base10_form {

			tokenPtrs[tokenId_prev_2].Type = Token_type_deleted_dont_use
			tokenPtrs[tokenId_prev_1].Type = Token_type_deleted_dont_use
			tokenPtrs[tokenIdActual].Type = Token_type_float_dotInDigits

			charsAllThree := []*ErlSrcChar{}
			charsAllThree = append(charsAllThree, tokenPtrs[tokenId_prev_2].Chars...)
			charsAllThree = append(charsAllThree, tokenPtrs[tokenId_prev_1].Chars...)
			charsAllThree = append(charsAllThree, tokenPtrs[tokenIdActual].Chars...)
			tokenPtrs[tokenIdActual].Chars = charsAllThree
			for _, Char := range tokenPtrs[tokenIdActual].Chars {
				Char.Token = tokenPtrs[tokenIdActual]
			}
		}
}

func multi_token_detect(chars []ErlSrcChar, verbose bool, detector func([]*ErlSrcToken, int, bool) )[]ErlSrcChar {
	// fmt.Println("---------- MULTI-TOKEN DETECT 1")
	// debug_print_ErlSrcChars(chars)
	tokenPtrs := tokens_from_chars(chars)
	for tokenId, _ := range tokenPtrs { detector(tokenPtrs, tokenId, verbose) }
	return chars // this is not necessary, the orig chars are modified from the detector
}



////////////////////////////////// token funs ////////////////////////////////////
func isSingleQuoteRune(r rune) bool { return r == '\''}
func isDoubleQuoteRune(r rune) bool { return r == '"'}
func isSingleOrDoubleQuoteRune (r rune) bool {return isSingleQuoteRune(r) || isDoubleQuoteRune(r)}
func tokenEmpty() ErlSrcToken { return ErlSrcToken{Type: "???"} }
func tokensForChars__preInitialized() ErlSrcTokens { return ErlSrcTokens{tokenEmpty()} }
//  ^^^^ // in Go, a variable's memory address stay the same when you assign a new value.
// so, I can use a token only once - it's necessary to generate always new tokens,
// and a simple 'tokenActual = tokenEmpty()' can't work, if the variable is always the same,
// because if I pass its pointer, later I can overwrite the value behind the variable.
// the current solution generates new tokens into a list, and the last elem is always
// updated, so it will have a new address after each update
////////////////////////////////// token funs ////////////////////////////////////