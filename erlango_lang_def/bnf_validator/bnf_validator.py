#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import argparse, os, time
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


def limit_of_symbol_length_to_reduce_options(fname: str, symbolName: str, default=10):
    """if there are too many options, possible language options cannot be generated,
    the option reduction is necessary to limit the number of cases

    the limitations for language elem generator are here.

    lower limitation is used where the variation of possible/accepted words are too wide,
    typically higher/more complex symbols have wider sets of expanded symbols.

    The local symbols in grammar files:
    to block/limit the too high sets of expanded symbols,
    locally the children symbols sets are limited.
    (a few alphabet elements are enough to test, there is no need
    to use the full lowercase+uppercase combo, for example)

    Later all grammar files will be merged, and they can use
    the whole set of children
    """
    limits = {   # bnf source file: symbol->allowed_length_of_expanded_expression
        "grammar_40_simple_types.bnf": {
            "<atom>": 5,
            "<atomCharList_inQuotes>": 4,
            "<atomInQuotes>": 4,
            "<atomPossibleCharAfterFirstPosition>": 4,
            "<atomSmallFirstChar>": 4,
            "<atomSmallFirstChar_tail>": 4,
            "<float>": 6,
            "<integer>": 6,
            "<number>": 6,
            "<numberTail>": 6,
            "<pid>": 10,
            "<triple_anyCharExceptQuoteOrEmtpy>": 6,
            "<stringQuoteTriple_safeChar>": 6,
            "<string>": 4,
            "<stringQuoteOne>": 6,
            "<stringQuoteOneTail>": 6,
            "<stringQuoteTriple>": 6,
            "<stringQuoteTripleTail>": 6,

        },
        "grammar_50_basic.bnf": {
            "<anyUnicodeCharExceptDoubleQuote>": 6,
            "<anyUnicodeCharsExceptSingleQuote>": 6,
            "<digit>": 4,
            "<empty>": 4,
            "<letter>": 4,
            "<letterCapital>": 6,
            "<letterSmall>": 6,
        }
    }

    return limits.get(fname, dict()).get(symbolName, default)


def main(filePathBnf: str, symbolNamesAnalyseOnly: list[str]):
    """
     - collect all symbols from the grammar
     - reduce the too big grammar sets to a smaller set, to get a manageable set of symbols for generating possibilities
     - detect missing symbol definitions
     - display detected errors in the grammar
    """

    errors: list[str] = list()
    symbolsTable, symbolNamesInLocalDefinition, errors = bnf_lib.symbols_detect_in_file(filePathBnf, errors)

    for sName, obj in symbolsTable.items():
        print(f"symbol: {sName:>50} {obj.expandPossibilitiesInBnf() }")
    # input("press ENTER to continue")

    print(f"local symbols: {symbolNamesInLocalDefinition}")


    for symbolWanted in symbolNamesAnalyseOnly:
        if symbolWanted not in symbolsTable:
            raise ValueError(f"a wanted symbol name {symbolWanted} is unknown in 'symbols'")

    ################################################
    missingSymbols = []
    print(f"=============== DETECT MISSING SYMBOL DEFINITIONS (not in left side of ::= operator)  =========")
    for symbolName, symbol in symbolsTable.items():
        print()
        print(f"detected symbol: {symbolName:>{bnf_lib.Symbol.symbolNameLenmax}}")
        print(symbol.definitionInBnf)

        for nonTerminatingSymbolInDefinition in symbol.grammar_elems_nonterminating_collect_in_all_possibilities():
            if nonTerminatingSymbolInDefinition not in symbolsTable:
                missingSymbols.append(nonTerminatingSymbolInDefinition)
                errors.append(f"ERROR in: '{filePathBnf}' non-defined symbol:  {symbolName} ::= .... {nonTerminatingSymbolInDefinition} <===== not defined in the grammar ")

    ################################################
    timeReport = []

    if not missingSymbols:
        for symbolName, symbol in symbolsTable.items():
            print(f"\n=================== {symbolName}  ================================")
            filePath_prefix = os.path.basename(filePathBnf)

            if symbolName not in symbolNamesInLocalDefinition:

                if symbolNamesAnalyseOnly and symbolName not in symbolNamesAnalyseOnly:
                    continue


                limitOfSymbolLengthInValidationToAvoidNeverendingLoop = (
                    limit_of_symbol_length_to_reduce_options(os.path.basename(filePathBnf), symbolName))

                timeStart = time.time()

                _, _, errors = possible_accepted_language_elems_save(
                    symbolName, symbolsTable, filePath_prefix,
                    limitOfSymbolsForTestCases=limitOfSymbolLengthInValidationToAvoidNeverendingLoop,
                    errors=errors,
                    filePathBnf=filePathBnf)

                timeKey =f"{filePathBnf}.{symbolName:<{symbol.symbolNameLenmax}}"
                timeReport.append(f"{timeKey} {time.time() - timeStart:.2f}   sec")

    bnf_lib.file_write(f"report_time_{os.path.basename(filePathBnf)}", "\n".join(timeReport))
    ################################################
    if not errors:
        print(f"No problem detected in the BNF")

    for err in errors:
        print(f"ERROR: {err}")




def possible_accepted_language_elems_save(symbolName: str, symbolsTable: dict[str, bnf_lib.Symbol],
                                          fileNamePrefixOfGrammar="x_prefix__",
                                          limitOfSymbolsForTestCases=10,
                                          errors: list[str]=[],
                                          filePathBnf: str="bnf_source_is_unknown") -> tuple[str, str, list[str]]:
    """Expand all possible matching elems.
    To avoid neverending recursive loops, there is a limitation against the maximum number of Symbols that will be expanded.

    return with the saved file names (they used in tests)
    """

    symbol = symbolsTable[symbolName]
    print(f"display possible accepted language elems in this symbol: {symbolName} -> {symbol.expandPossibilitiesInBnf()}")

    expandTheseSymbolsUntilTerminationIsNotReached = symbol.expandPossibilitiesInBnf()
    if not expandTheseSymbolsUntilTerminationIsNotReached :
        errMsg = f"ERROR: {symbolName} symbol: undefined expansion rules after ::= in file: {filePathBnf}"
        errors.append(errMsg)


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

    log(f"params: limitOfSymbols:", limitOfSymbolsForTestCases)

    while expandTheseSymbolsUntilTerminationIsNotReached:
        loopCounter += 1


        # DEBUG
        bnf_lib.symbolnames_possibilities_print(
            expandTheseSymbolsUntilTerminationIsNotReached, "possib-debug:", caller="while loop in main")


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
            log(f"{loopCounter:>5}. loop - one symbol in possibility:", symbolInPossibility)

            if bnf_lib.symbolname_terminating(symbolInPossibility):
                log(f"{loopCounter:>5}. loop - terminating symbol name:", symbolInPossibility)
                expandedOnlyTerminatingsPossibilities.append(symbolInPossibility)
            else:
                # in the first record of Possibility, there is a non-terminating symbol.
                # Expand it and pack it back to the first position, to continue the expanding totally.

                insertBack__oneExpansionHappened = []
                possibsOfCurrentSymbol = symbolsTable[symbolInPossibility].expandPossibilitiesInBnf()

                # the non-terminating symbol cannot be expanded to other symbols,
                # for example:  <empty> ::=       <------missing symbols here on the right side

                if not possibsOfCurrentSymbol:
                    errMsg = f"ERROR: {symbolInPossibility} symbol (reached in child expansion): missing expansion rules, nothing is defined after ::= in file: {filePathBnf}"
                    errors.append(errMsg)


                for nonTerminatingExpansion in possibsOfCurrentSymbol:
                    oneStepExpansionHappened = expandedOnlyTerminatingsPossibilities + nonTerminatingExpansion + onePossibilitySymbolChangingList
                    log("oneStepExpanded before SymbolReuseCheck", oneStepExpansionHappened)

                    insertBack = True
                    # to avoid neverending loops
                    if  len(oneStepExpansionHappened) > limitOfSymbolsForTestCases:
                        log("oneStepExpanded, number of non-terminating symbols are too high, don't expand it ", oneStepExpansionHappened)
                        insertBack = False

                    if insertBack:
                        insertBack__oneExpansionHappened.append(oneStepExpansionHappened)

                expandTheseSymbolsUntilTerminationIsNotReached = insertBack__oneExpansionHappened + expandTheseSymbolsUntilTerminationIsNotReached
                break
            ###############################################################################


        # there is no more symbol that can be converted in the possibility, add it to the reportAcceptedLangExamples
        if expandedOnlyTerminatingsPossibilities and len(onePossibilitySymbolChangingList) == 0:
            quotesRemovedFromTerminatingSimbols = [terminatingSymbol[1:-1] for terminatingSymbol in expandedOnlyTerminatingsPossibilities]
            reportAcceptedLangExamples.append("".join(quotesRemovedFromTerminatingSimbols))
            log(" only terminating symbolname", "".join(quotesRemovedFromTerminatingSimbols), extraLineAfter=True)

    fname_prefix = f"{fileNamePrefixOfGrammar}___{symbolName[1:-1]}"
    fname_bnf_accepted = f"{fname_prefix}.bnf_accepted"
    fname_log = f"{fname_prefix}.log"
    bnf_lib.file_write(f"{fname_bnf_accepted}", "\n".join(reportAcceptedLangExamples))
    bnf_lib.file_write(f"{fname_log}", "\n".join(logs))

    return fname_bnf_accepted, fname_log, errors




if __name__ == '__main__':  # pragma: no cover
    defaultFiles = ",".join(bnf_lib.files_collect_in_dir("..", prefix="grammar_"))

    parser = argparse.ArgumentParser(prog='BNF validator')
    parser.add_argument("--file_bnf_path", type=str, default=defaultFiles, help="one file, or more comma separated filenames to check/validate", required=False)
    parser.add_argument("--symbol_names_analyse_only", type=str, default="", help="analyse only these symbol names", required=False)
    args = parser.parse_args()

    print(f"validate these files: {args.file_bnf_path}")

    symbolNamesAnalyseOnly = []
    if args.symbol_names_analyse_only:
        symbolNamesAnalyseOnly = args.symbol_names_analyse_only.split(",")

    for file in args.file_bnf_path.split(","):
        bnf_lib.file_exists___alert_if_not(file)
        main(file, symbolNamesAnalyseOnly)