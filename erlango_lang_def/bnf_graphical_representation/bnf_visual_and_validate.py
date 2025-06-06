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

    def expandPossibilities(self):
        # all possible options are given back, as possible expansions

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

    errors = list()
    symbols, errors = symbol_detect(filePathBnf, errors)

    for tooBigUseSmallerSetForGrammarCheck in ["<letterSmall>", "<letterCapital>", "<digit>"]:
        symbols[tooBigUseSmallerSetForGrammarCheck].limitExpandPossibilitiesInTooBigSets = True


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
    if not missingSymbols:
        for symbolName, symbol in symbols.items():
            if symbolName != "<float>": continue

            print(f"\n=================== {symbolName} Expand ================================")
            display_possible_accepted_language_elems(symbolName, symbols)

    ################################################
    if not errors:
        print(f"No problem detected in the BNF")

    for err in errors:
        print(f"ERROR: {err}")



def  count_non_terminatings_is_under_repetition_limit(nonTerminatingSymbols: [str], allowedSymbolReuseInSamePossibility) -> dict[str, int]:
    """count non-terminating symbols. To avoid neverending recursion, stop if the same elem has more than allowed repetitions"""
    stats = dict()

    for nonTermSymbol in nonTerminatingSymbols:
        stats.setdefault(nonTermSymbol, 0)
        stats[nonTermSymbol] += 1

    for symbolName, counted in stats.items():
        if counted > allowedSymbolReuseInSamePossibility:   # in case of <float>, there are 2 <digits> immediatelly in the grammar, so use 3 here.
            return False  # not under limit

    return True


def display_possible_accepted_language_elems(symbolName: str, symbols: dict[str, Symbol], allowedSymbolReuseInSamePossibility=2):
    """Expand all possible matching elems. To block neverending code generation, max 2 recursive call is allowed."""


    symbol = symbols[symbolName]
    print(f"display possible accepted language elems in this symbol: {symbolName} -> {symbol.expandPossibilities()}")

    expandTheseSymbolsUntilTerminationIsNotReached = symbol.expandPossibilities()

    report = []
    while expandTheseSymbolsUntilTerminationIsNotReached:

        onePossibilitySymbolChangingList = expandTheseSymbolsUntilTerminationIsNotReached.pop(0)
        expandedOnlyTerminatings = []

        # expand one symbol in the possibility. if it has more than one options, insert all of them back into the list
        while onePossibilitySymbolChangingList:

            # be careful, if you print this, helps to understand what is happening,
            # but you will get more thousand extra lines in the output
            # print(f"one possibility symbols, in expansion process: {onePossibilitySymbolChangingList}")

            symbolInPossibility = onePossibilitySymbolChangingList.pop(0)

            isTerminating = symbolInPossibility.startswith('"') and symbolInPossibility.endswith('"')

            if isTerminating:
                expandedOnlyTerminatings.append(symbolInPossibility)
            else:
                for nonTerminatingExpansion in symbols[symbolInPossibility].expandPossibilities():
                    oneStepExpansionHappened = expandedOnlyTerminatings + nonTerminatingExpansion + onePossibilitySymbolChangingList
                    if count_non_terminatings_is_under_repetition_limit(oneStepExpansionHappened, allowedSymbolReuseInSamePossibility=allowedSymbolReuseInSamePossibility):
                        expandTheseSymbolsUntilTerminationIsNotReached.append(oneStepExpansionHappened)
                break

        # there is no more symbol that can be converted in the possibility, add it to the report
        if expandedOnlyTerminatings and len(onePossibilitySymbolChangingList) == 0:
            quotesRemovedFromTerminatingSimbols = []
            for terminatingSymbol in expandedOnlyTerminatings:
                quotesRemovedFromTerminatingSimbols.append(terminatingSymbol[1:-1])
            report.append("".join(quotesRemovedFromTerminatingSimbols))

    for expanded in report:
        print(f"expanded, language accepted terminating symbol: {expanded}")





def get_symbolname_and_definition_in_line(line, errors):
    """<newSymbol> ::= .....definition....
    in a line, there is only definition, or if it is a new symbol, a symbolName and definition.

    detect them.
    """
    acceptedSymbolChars = "_-abcdefghijklmnopqrstuvwxyZABCDEFGHIJKLMNOPQRSTUVWXYZ"

    newSymbolNameInLine = ""
    definitionInLine = line

    if "::=" in line:
        # wanted: "<symbol>::="
        lineClean = line.strip().replace(" ", "").replace("\t", "")
        maybeSymbol = lineClean.split("::=")[0]
        if maybeSymbol.startswith("<") and maybeSymbol.endswith(">"):
            # it can have only a-zA-Z_- chars

            allLettersAreAcceptedInSymbolName = True
            for letter in maybeSymbol[1:-1]:
                if letter not in acceptedSymbolChars:
                    allLettersAreAcceptedInSymbolName = False
                    errors.append(f"maybe human error: strange character '{letter}' in '<symbol> ::=' definition:\n---> {maybeSymbol} ")
                    break

            if allLettersAreAcceptedInSymbolName:
                newSymbolNameInLine = maybeSymbol

                # keep the lenght of indentation WITHOUT the '<symbol> ::=' part
                # split only at the first ::=
                elems = line.split("::=", 1)  # the split is executed on the original line, so every char is kept
                definitionInLine = " " * (len(elems[0])+3) + elems[1]  # the '<..> ::=' part, filled with space, and the definition

    return newSymbolNameInLine, definitionInLine

def symbol_detect(filePathBnf: str, errors: [str]):
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



def file_src_lines(path: str) -> [str]:
    lines = []
    with open(path, 'r') as file:
        lines = file.readlines()  # Pycharm highlight if I return directly in this line
    return lines

def file_validation(path : str):
    if not os.path.exists(path):
        print(f"ERROR: invalid file path: {path}")
        sys.exit(1)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog='BNF visualiser')
    parser.add_argument("--file_bnf_path", type=str, default="../erlango_lang.bnf", required=False)
    args = parser.parse_args()

    file_validation(args.file_bnf_path)

    main(args.file_bnf_path)