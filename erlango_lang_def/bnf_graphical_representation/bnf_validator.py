#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import argparse
import bnf_lib

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



def main(filePathBnf: str):
    """
     - collect all symbols from the grammar
     - reduce the too big grammar sets to a smaller set, to get a manageable set of symbols for generating possibilities
     - detect missing symbol definitions
     - display detected errors in the grammar
    """

    errors = list()
    symbols, errors = bnf_lib.symbols_detect_in_file(filePathBnf, errors)

    ################################################
    missingSymbols = []
    print(f"=============== DETECT MISSING SYMBOL DEFINITIONS (not in left side of ::= operator)  =========")
    for symbolName, symbol in symbols.items():
        print()
        print(f"detected symbol: {symbolName:>{bnf_lib.Symbol.symbolNameMax}}")
        print(symbol.definitionInBnf)

        for nonTerminatingSymbolInDefinition in symbol.grammar_elems_nonterminating_collect():
            if nonTerminatingSymbolInDefinition not in symbols:
                missingSymbols.append(nonTerminatingSymbolInDefinition)
                errors.append(f"non-defined symbol:  {symbolName} ::= .... {nonTerminatingSymbolInDefinition} <===== not defined in the grammar ")

    ################################################
    if not missingSymbols:
        for symbolName, symbol in symbols.items():
            print(f"\n=================== {symbolName}  ================================")
            possible_accepted_language_elems_save(symbolName, symbols)

    ################################################
    if not errors:
        print(f"No problem detected in the BNF")

    for err in errors:
        print(f"ERROR: {err}")




def possible_accepted_language_elems_save(symbolName: str, symbols: dict[str, bnf_lib.Symbol], limitOfNonTerminatingSymbols=40):
    """Expand all possible matching elems.
    To avoid neverending recursive loops, there is a limitation against the maximum number of Symbols that will be expanded.
    """

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
        log(f"{loopCounter:>5}. loop === first possibility === :", bnf_lib.symbolnames_concate_simple_str(onePossibilitySymbolChangingList), extraLineBefore=True)

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

            if bnf_lib.is_terminating_symbolname(symbolInPossibility):
                expandedOnlyTerminatingsPossibilities.append(symbolInPossibility)
            else:
                # in the first record of Possibility, there is a non-terminating symbol.
                # Expand it and pack it back to the first position, to continue the expanding totally.

                insertBack__oneExpansionHappened = []
                for nonTerminatingExpansion in symbols[symbolInPossibility].expandPossibilities():
                    oneStepExpansionHappened = expandedOnlyTerminatingsPossibilities + nonTerminatingExpansion + onePossibilitySymbolChangingList
                    # log("oneStepExpanded before SymbolReuseCheck", oneStepExpansionHappened)

                    insertBack = True
                    # to avoid neverending loops
                    if  len(bnf_lib.symbols_nonterminating_collect(oneStepExpansionHappened)) < limitOfNonTerminatingSymbols:
                        #log("oneStepExpanded, the num of non-terminating symbols are too high, don't expand it ")
                        insertBack = False

                    if insertBack:
                        insertBack__oneExpansionHappened.append(oneStepExpansionHappened)

                expandTheseSymbolsUntilTerminationIsNotReached = insertBack__oneExpansionHappened + expandTheseSymbolsUntilTerminationIsNotReached
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
    bnf_lib.file_write(f"{fname}.bnf_accepted", "\n".join(reportAcceptedLangExamples))
    bnf_lib.file_write(f"{fname}.log", "\n".join(logs))



if __name__ == '__main__':
    defaultFiles = bnf_lib.files_collect_in_dir("..", prefix="grammar_")
    print(f"default files: {defaultFiles}")

    parser = argparse.ArgumentParser(prog='BNF validator')
    parser.add_argument("--file_bnf_path", type=str, default="../grammar_50_basic.bnf", help="one file, or more comma separated filenames to check/validate", required=False)
    args = parser.parse_args()

    for file in args.file_bnf_path.split(","):
        bnf_lib.file_is_exists(file)
        main(args.file_bnf_path)