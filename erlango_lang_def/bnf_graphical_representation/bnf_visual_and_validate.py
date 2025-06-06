#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os, argparse, re

"""
Analyse the passed BNF file to support the development.

BNF uses a specific set of symbols to define the syntax of a language. 
These symbols, sometimes called metasymbols, act as operators within the grammar definition. 

Here's a breakdown of the common operators:

Basic Operators:

::=  This is the definition operator. It separates a non-terminal symbol 
     (on the left) from its possible expansions (on the right). 
     For example, <expression> ::= <term>.

|    The vertical bar represents an "or" operator. 
     It indicates that a non-terminal can be replaced by one of several alternatives. 
     For example, <digit> ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9".
     
< >  Angle brackets enclose non-terminal symbols. 
     Non-terminals are placeholders that represent a syntactic category. 
     For example, <expression>, <term>, <digit>.
     
" "  Double quotes enclose terminal symbols. 
     Terminals are the actual characters or keywords that appear in the language. 
     For example, "if", "+", "1".


specials here: if # sign is the first non-whitespace char in a line, that line is a comment.
"""

class Symbol:
    symbolNameMax = 0

    def __init__(self, symbolName):
        self.name = symbolName
        self.definitionInBnf = ""
        self.definitionCounterInBnf = 0  # in lucky case, the symbol is defined only once in the file

        Symbol.symbolNameMax = max(Symbol.symbolNameMax, len(symbolName))

        # the ABC/numbers is too wide. To see the grammar working, a few chars are more than enough.
        self.limitExpandPossibilitiesInTooBigSets = False
        self.limitExpandPossibilities = 2


    def symbols_nonterminating_all_used_in_bnf_grammar(self):
        """collect all <non-terminating> elems only"""
        tokensNonTerminating = re.findall(r'<[^<>]+>', self.definitionInBnf)
        return tokensNonTerminating

    def expandPossibilities(self) -> [[str]]:
        """collect all possible expansions"""
        # one expansion: list of one or more symbol series
        # so the return value is 'list of series-of-symbols'

        expanded = []

        for onePossibility in self.definitionInBnf.replace("\n", " ").strip().split("|"):
            # print(f"one possibility in one string: {onePossibility}")
            tokens_in_one_possib = re.findall(r'<[^<>]+>|"[^"]+"', onePossibility)
            # print(f"one possibility, separated tokens in one elem: {tokens_in_one_elem}")
            expanded.append(tokens_in_one_possib)

        if self.limitExpandPossibilitiesInTooBigSets:
            return expanded[:self.limitExpandPossibilities]
        return expanded



def main(filePathBnf: str):
    """
     - collect all symbols from the grammar
     - reduce the too big grammar sets to a smaller set, to get a manageable set of symbols for generating possibilities
     - detect missing symbol definitions
     - display detected errors in the grammar
    """

    errors = list()
    symbols, errors = level0_symbol_detect(filePathBnf, errors)



    ####### USE SMALL SETS FOR TOO WIDE ELEMS #################
    print(f"=== DISPLAY ELEM NUMS TO SEE TOO WIDE SETS before reduction.... ===")

    tooManyElems = ["<letterSmall>", "<letterCapital>", "<digit>", "<letterSmallCapital>"]
    for tooBigUseSmallerSetForGrammarCheck in tooManyElems:
        symbols[tooBigUseSmallerSetForGrammarCheck].limitExpandPossibilitiesInTooBigSets = True

    symbols["<digit>"].definitionInBnf = '"1"'

    symbols["<variableLetterOrDigitOrUnderscore>"].definitionInBnf = '"i"'
    symbols["<variableLetterCapitalOrUnderscore>"].definitionInBnf = '"V"'
    symbols["<atomPossibleCharAfterFirstPosition>"].definitionInBnf = '"t"'
    symbols["<escapeCharInSeq>"].definitionInBnf = '"n"'


    # num of operators are too high, use one only
    symbols["<expression>"].definitionInBnf = '''
                     <expression> "+" <term>
                   | <term>
    '''


    symbols["<term>"].definitionInBnf = '''
           <term> "*" <factor>
         | <factor>
    '''

    # simplified factor is necessary
    symbols["<factor>"].definitionInBnf = '''
             <number>
           | <variable>
           | "(" <expression> ")"
    '''

    print(f"=== ELEM NUMS after reduction.... ===")
    print_symbols_elem_stats(symbols)
    # input(f"press ENTER to continue")
    ####### USE SMALL SETS FOR TOO WIDE ELEMS #################



    ################################################
    missingSymbols = []
    print(f"=============== DETECT MISSING SYMBOL DEFINITIONS (not in left side of ::= operator)  =========")
    for symbolName, symbol in symbols.items():
        print()
        print(f"detected symbol: {symbolName:>{Symbol.symbolNameMax}}")
        print(symbol.definitionInBnf)

        for nonTerminatingSymbolInDefinition in symbol.symbols_nonterminating_all_used_in_bnf_grammar():
            if nonTerminatingSymbolInDefinition not in symbols:
                missingSymbols.append(nonTerminatingSymbolInDefinition)
                errors.append(f"non-defined symbol:  {symbolName} ::= .... {nonTerminatingSymbolInDefinition} <===== not defined in the grammar ")


    ################################################
    expand_these_symbols = {
        "<anyUnicodeChars>",
        "<anyUnicodeCharsExceptNonEscapedSingleQuote>",
        "<atomPossibleCharAfterFirstPosition>",
        "<digit>",
        "<empty>",
        "<escapeCharInSeq>",
        "<escapeSequence>",
        "<float>",
        "<integer>",
        "<letterCapital>",
        "<letterSmall>",
        "<letterSmallCapital>",
        "<number>",
        "<variable>",
        "<variableLetterCapitalOrUnderscore>",
        "<variableLetterOrDigitOrUnderscore>",
        "<variableTail>",
        "<numberTail>",
        "<atomCharList_inQuotes>",
        "<atomInQuotes>",            
        
        

        # "<argumentList>",
        # "<argumentListTail>",
        "<atom>",
        "<atomSmallFirstChar>",
        "<atomSmallFirstChar_tail>",
        "<atom_or_variable_name_in_paramlist>",
        # "<caseClause>",
        # "<caseClauseList>",
        # "<caseClauseTail>",
        # "<caseExpression>",
        # "<dashArg>",
        # "<dashPrgAttrib>",
        # "<exportEntry>",
        # "<exportEntryList>",
        # "<exportEntryListTail>",
        # "<exportList>",
        "<expression>",
        # "<expressionList>",
        # "<expressionListTail>",
        # "<factor>",
        # "<functionCall>",
        # "<functionDefinition>",
        # "<functionName>",
        # "<list>",
        # "<map>",
        # "<mapEntry>",
        # "<mapEntryList>",
        # "<mapEntryListTail>",
        # "<parameterList>",
        # "<parameterListTail>",
        # "<pattern>",
        "<pid>",
        # "<program>",
        # "<programTail>",
        "<string>",
        "<stringInSigil>",
        "<stringQuoteOne>",
        "<stringQuoteOneTail>",
        "<stringSigilTail>",
        # "<term>",
        # "<tuple>",


        
    }
    if not missingSymbols:
        for symbolName, symbol in symbols.items():

            expand = True
            # one symbol process only, for testing
            if expand_these_symbols:
                if symbolName not in expand_these_symbols:
                    expand = False

            print(f"\n=================== {symbolName} Expand: {expand} ================================")
            if expand:
                level0_possible_accepted_language_elems_save(symbolName, symbols)


    ################################################
    if not errors:
        print(f"No problem detected in the BNF")

    for err in errors:
        print(f"ERROR: {err}")


def level0_symbol_detect(filePathBnf: str, errors: [str]):
    """collect symbols and definitions from the bnf file
    errors is returned to represent on caller level that it is modified here
    """

    print(f"BNF def file: {filePathBnf}")

    symbols = dict()
    ################################################

    symbolName = ""

    for line in file_src_lines(filePathBnf):
        if line.startswith("#"):
            continue  # comment line

        symbolNameNewDetected, definitionInLine = get_symbolname_and_definition_in_line(line, errors)

        if symbolNameNewDetected:
            symbolName = symbolNameNewDetected

            if symbolName not in symbols:
                symbols[symbolName] = Symbol(symbolName)
            else:
                symbols[symbolName].definitionCounterInBnf += 1
                errors.append(f"problem: the symbol is defined more than once in the bnf grammar: {symbolName}, defCount: {symbols[symbolName].definitionCounterInBnf} ")


        # one symbol definition is max a few lines long, not a long string,
        # so this naive string concatenate is not a problem.
        if symbolName:  # not the empty non-defined:
            symbols[symbolName].definitionInBnf += definitionInLine

    return symbols, errors


def symbolnames_simple_str(names: [str]) -> str:
    """concatenate list of symbolnames"""
    return ", ".join(names)

def level0_possible_accepted_language_elems_save(symbolName: str, symbols: dict[str, Symbol], allowedSymbolReuseInSamePossibility=2):
    """Expand all possible matching elems. To block neverending code generation, max N recursive call is allowed."""

    symbol = symbols[symbolName]
    print(f"display possible accepted language elems in this symbol: {symbolName} -> {symbol.expandPossibilities()}")

    expandTheseSymbolsUntilTerminationIsNotReached = symbol.expandPossibilities()

    reportAcceptedLangExamples = []
    logs = []
    def log(msg, val, extraLineBefore=False, extraLineAfter=False):
        if extraLineBefore:
            logs.append("")
        out = f"{msg:>50} -> {val}"
        print(out)
        logs.append(out)
        if extraLineAfter:
            logs.append("")

    loopCounter = 0
    while expandTheseSymbolsUntilTerminationIsNotReached:
        loopCounter += 1

        onePossibilitySymbolChangingList = expandTheseSymbolsUntilTerminationIsNotReached.pop(0)
        log(f"{loopCounter:>5}. loop === first possibility === :", symbolnames_simple_str(onePossibilitySymbolChangingList), extraLineBefore=True)

        expandedOnlyTerminatingsPossibilities = []

        # expand one symbol/one-word only in the possibility.
        # if it has more than one options, insert all of them back into the list
        while onePossibilitySymbolChangingList:

            # be careful, if you print this, helps to understand what is happening,
            # but you will get more thousands extra lines in the output
            # print(f"one possibility symbols, in expansion process: {onePossibilitySymbolChangingList}")

            ###############################################################################
            # get first word of the possibility
            symbolInPossibility = onePossibilitySymbolChangingList.pop(0)
            # log(f"{loopCounter:>5}. loop - one symbol in possibility:", symbolInPossibility)

            if is_terminating_symbolname(symbolInPossibility):
                expandedOnlyTerminatingsPossibilities.append(symbolInPossibility)
            else:
                # in the first record of Possibility, there is a non-terminating symbol.
                # Expand it and pack it back to the first position, to continue the expanding totally.
                insertTheseAfterOneExpand = []
                for nonTerminatingExpansion in symbols[symbolInPossibility].expandPossibilities():
                    oneStepExpansionHappened = expandedOnlyTerminatingsPossibilities + nonTerminatingExpansion + onePossibilitySymbolChangingList
                    # log("oneStepExpanded before SymbolReuseCheck", oneStepExpansionHappened)

                    underRepetitionLimit, symbolNamesOverLimit = count_non_terminatings_are_under_repetition_limit(oneStepExpansionHappened, allowedSymbolReuseInSamePossibility=allowedSymbolReuseInSamePossibility)

                    if not underRepetitionLimit:
                        pass
                        # in an expression, this is too long, don't display
                        # log(f"overRepetition > {allowedSymbolReuseInSamePossibility} here:", symbolnames_simple_str(oneStepExpansionHappened))

                        # log("overRepetitionLimit:", symbolNamesOverLimit)

                    if underRepetitionLimit:
                        #log("oneStepExpanded after SymbolReuseCheck", oneStepExpansionHappened)
                        insertTheseAfterOneExpand.append(oneStepExpansionHappened)
                expandTheseSymbolsUntilTerminationIsNotReached = insertTheseAfterOneExpand + expandTheseSymbolsUntilTerminationIsNotReached
                break
            ###############################################################################


        # there is no more symbol that can be converted in the possibility, add it to the reportAcceptedLangExamples
        if expandedOnlyTerminatingsPossibilities and len(onePossibilitySymbolChangingList) == 0:
            quotesRemovedFromTerminatingSimbols = []
            for terminatingSymbol in expandedOnlyTerminatingsPossibilities:
                quotesRemovedFromTerminatingSimbols.append(terminatingSymbol[1:-1])
            reportAcceptedLangExamples.append("".join(quotesRemovedFromTerminatingSimbols))
            log(" only terminating symbolname", "".join(quotesRemovedFromTerminatingSimbols), extraLineAfter=True)

    fname = f"grammar_{symbolName[1:-1]}"
    file_write(f"{fname}.bnf_accepted", "\n".join(reportAcceptedLangExamples))
    file_write(f"{fname}.log", "\n".join(logs))


def is_terminating_symbolname(symbolName: str) -> bool:
    """is it a terminating symbol name?"""
    return symbolName.startswith('"') and symbolName.endswith('"')


def  count_non_terminatings_are_under_repetition_limit(symbolsInPossibility: [str], allowedSymbolReuseInSamePossibility) -> dict[str, int]:
    """count non-terminating symbols. To avoid neverending recursion, stop if the same elem is repeated more than allowed"""
    stats = dict()

    # build full statistics about nonterminals... (relative small operation, full statistics can be calculated)
    for symbolName in symbolsInPossibility:
        stats.setdefault(symbolName, 0)
        stats[symbolName] += 1

    overLimit = []
    for symbolName, counted in stats.items():
        if counted > allowedSymbolReuseInSamePossibility:   # in case of <float>, there are 2 <digits> immediatelly in the grammar, so use 3 here.
            overLimit.append(symbolName)

    if overLimit:
        return False, overLimit

    return True, []


def get_symbolname_and_definition_in_line(line, errors):
    """<newSymbol> ::= .....definition....
    in a line, there is only definition, or if it is a new symbol, a symbolName and definition.

    detect them.
    """
    acceptedSymbolChars = "_-abcdefghijklmnopqrstuvwxyZABCDEFGHIJKLMNOPQRSTUVWXYZ"

    # If there is no definition in a line, return with empty string.
    newSymbolNameInLine = ""
    definitionInLine = line

    if "::=" in line:
        # wanted: "<symbol>::="
        lineClean = line.strip().replace(" ", "").replace("\t", "")
        maybeSymbol = lineClean.split("::=", 1)[0]
        if maybeSymbol.startswith("<") and maybeSymbol.endswith(">"):
            # it can have only a-zA-Z_- chars

            allLettersAreAcceptedInMaybeSymbolName = True
            for letter in maybeSymbol[1:-1]:
                if letter not in acceptedSymbolChars:
                    allLettersAreAcceptedInMaybeSymbolName = False
                    errors.append(f"maybe human error: strange character '{letter}' in '<symbol> ::=' definition:\n---> {maybeSymbol} ")
                    break

            if allLettersAreAcceptedInMaybeSymbolName:
                newSymbolNameInLine = maybeSymbol

                # keep the lenght of indentation WITHOUT the '<symbol> ::=' part
                # split only at the first ::=
                elems = line.split("::=", 1)  # the split is executed on the original line, so every char is kept
                definitionInLine = " " * (len(elems[0])+3) + elems[1]  # the '<..> ::=' part, filled with space, and the definition

    return newSymbolNameInLine, definitionInLine


def print_symbols_elem_stats(symbols):
    """display statistics about symbols to see where do we have too big set"""
    for symbolName, symbol in symbols.items():
        possibilities = symbols[symbolName].expandPossibilities()
        maxElemInOnePossibility = 0
        for onePossible in possibilities:
            maxElemInOnePossibility = max(maxElemInOnePossibility, len(onePossible))
        print(f"{symbolName:>50} possib: {len(possibilities)}",  possibilities)


def file_src_lines(path: str) -> [str]:
    lines = []
    with open(path, 'r') as file:
        lines = file.readlines()  # Pycharm highlight if I return directly in this line
    return lines


def file_is_exists(path: str):
    if not os.path.exists(path):
        print(f"ERROR: invalid file path: {path}")
        sys.exit(1)

def file_write(path: str, content: str):
    with open(path, mode="w") as f:
        f.write(content)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog='BNF visualiser')
    parser.add_argument("--file_bnf_path", type=str, default="../erlango_lang.bnf", required=False)
    args = parser.parse_args()

    file_is_exists(args.file_bnf_path)

    main(args.file_bnf_path)