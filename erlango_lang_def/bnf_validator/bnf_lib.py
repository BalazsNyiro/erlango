#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os, re

class Symbol:
    symbolNameLenmax = 0

    def __init__(self, symbolName):
        self.name = symbolName
        self.definitionInBnf = ""

        self.definitionCounterInBnf = 1
        # basically the symbol is defined only once in the file

        Symbol.symbolNameLenmax = max(Symbol.symbolNameLenmax, len(symbolName))

    def expandPossibilities(self) -> list[list[str]]:
        """collect all possible expansions
        one possibility can have one-or-more symbols (list of lists)
        """
        possibilities = symbol_names_collect_from_grammar_def(self.definitionInBnf)
        return possibilities


    # to detect missing symbol definitions,
    # collect all non-terminating symbols from all possibilites
    def grammar_elems_nonterminating_collect_in_all_possibilities(self):
        possibilities = self.expandPossibilities()
        symbolsNonTerms = set()
        for onePossibility in possibilities:
            nonTerms = symbols_nonterminating_collect(onePossibility)
            for nonTerm in nonTerms:
                symbolsNonTerms.add(nonTerm)

        return sorted(symbolsNonTerms)


def file_exists___alert_if_not(path: str, raiseException=True) -> bool:
    """check if a file exists or not, generate an exception if there is a problem"""
    if not os.path.exists(path):
        msg = f"ERROR: invalid file path: {path}"
        print(msg)

        if raiseException:
            raise ValueError(msg)
        return False

    return True


def file_src_lines(path: str) -> list[str]:
    """read all lines of a file"""
    print(f"read file: {path}")
    lines = []
    with open(path, 'r') as file:
        lines = file.readlines()  # Pycharm highlight if I return directly in this line
    return lines


def file_write(path: str, content: str) -> None:
    with open(path, mode="w") as f:
        f.write(content)


def files_collect_in_dir(directory: str, prefix: str="grammar_") -> list[str]:
    files_and_path = []
    files = [f for f in os.listdir(directory) if f.startswith(prefix) and os.path.isfile(os.path.join(directory, f))]

    for f in files:
        files_and_path.append(os.path.join(directory, f))

    return files_and_path




def symbol_names_collect_from_grammar_def(grammarDefInBnf: str, verbose: bool=False) -> list[list[str]]:
    """collect all symbols (non-terminating AND terminating)
    from grammar def"""

    possibilities: list[list[str]] = []
    allSymNames: list[str] = []

    # the full def, maybe multiple lines and | separators
    grammarDefInBnf = grammarDefInBnf.strip()

    inTerminalDetection = False
    inNonTerminalDetection = False
    symbolName = ""

    def debug(msg):
        if verbose:
            print(msg)

    def inDetection():
        return inTerminalDetection or inNonTerminalDetection

    backSlashCounter = 0   # continuous backslash counter
    def isEscapedLetter():  # if the num of backSlashes is ODD
        return backSlashCounter % 2 != 0

    for letter in grammarDefInBnf:

        if letter == "\\":
            backSlashCounter += 1


        if not inDetection():
            if letter == "<": inNonTerminalDetection = True
            if letter == '"': inTerminalDetection = True

            # < or " is detected now, one Detection -> True
            if inDetection():  # detection is started now
                symbolName += letter
                debug(f"-> inDetection")
                debug(f"  -> add: '{letter}'")
                continue

            if letter == "|":
                possibilities.append(allSymNames)
                allSymNames = []
            continue


        if inDetection():         # naive string concatenate is enough
            symbolName += letter  # symbol names are typically short,
            debug(f"  -> add: '{letter}'")

            if not isEscapedLetter():
                if letter == ">": inNonTerminalDetection = False
                if letter == '"': inTerminalDetection = False

            if not inDetection():   # inDetection->endOfDetection
                debug(f"  -> END")
                allSymNames.append(symbolName)
                symbolName = ""


        if letter != "\\":
            backSlashCounter = 0


    if allSymNames:
        possibilities.append(allSymNames)

    return possibilities



def symbols_nonterminating_collect(symbolNamesAll: list[str]) -> list[str]:
    """collect all <non-terminating> elems only in the given grammar,
    one string, or list of symbols
    """
    nonTerminatings = []

    for symbolName in symbolNamesAll:
        if symbolName.startswith('<') and symbolName.endswith('>'):
            nonTerminatings.append(symbolName)

    return nonTerminatings


def symbolnames_concate_simple_str(names: list[str]) -> str:
    """concatenate list of symbolnames into one str"""
    return ", ".join(names)


def symbolname_and_grammar_definition_in_line__get(line: str, errors: list[str]) -> tuple[str, str, list[str]]:
    """<newSymbol> ::= .....definition....
    in a line, there is only definition, or if it is a new symbol, a symbolName and definition.

    detect them.

    """
    acceptedSymbolNameChars = "_-abcdefghijklmnopqrstuvwxyZABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
                if letter not in acceptedSymbolNameChars:
                    allLettersAreAcceptedInMaybeSymbolName = False
                    errors.append(f"maybe human error: strange character '{letter}' in '<symbol> ::=' definition:\n---> {maybeSymbol} ")
                    break

            if allLettersAreAcceptedInMaybeSymbolName:
                newSymbolNameInLine = maybeSymbol

                # the definition's indentation is kept, because if the grammar has multiple lines,
                # the multi-line display can keep the formatting then
                # keep the lenght of indentation WITHOUT the '<symbol> ::=' part
                # split only at the first ::=
                elems = line.split("::=", 1)  # the split is executed on the original line, so every char is kept
                definitionInLine = " " * (len(elems[0])+3) + elems[1]  # the '<..> ::=' part, filled with space, and the definition

        else:
            # <....> missing brackets
            errors.append(
                f"missing symbol brackets ---> {maybeSymbol} <- {line}")

    return newSymbolNameInLine, definitionInLine, errors


def symbolname_terminating(symbolName: str) -> bool:
    """is it a terminating symbol name?"""
    return symbolName.startswith('"') and symbolName.endswith('"')




def symbols_detect_in_file(filePathBnf: str, errors: list[str]
    ) -> tuple[dict[str, Symbol], list[str], list[str]]:

    """collect symbols and definitions from the bnf file
    errors is returned to represent on caller level that it is modified here
    """

    print(f"BNF def file: {filePathBnf}")

    symbols = dict()
    symbolNamesInLocalDefinition = set()
    localSymbolDefinitionSection = False
    
    ################################################

    symbolName = ""

    for line in file_src_lines(filePathBnf):
        # print(f"debug sym detect: {line}")

        if line.startswith("#"):
            if "LOCAL SYMBOLS" in line:
                localSymbolDefinitionSection = True

            continue  # comment line

        symbolNameNewDetected, definitionInLine, errors = symbolname_and_grammar_definition_in_line__get(line, errors)

        if symbolNameNewDetected:
            symbolName = symbolNameNewDetected

            if localSymbolDefinitionSection:
                symbolNamesInLocalDefinition.add(symbolName)


            if symbolName not in symbols:
                symbols[symbolName] = Symbol(symbolName)
            else:
                symbols[symbolName].definitionCounterInBnf += 1
                msg = f"ERROR: the symbol is defined more than once in the bnf grammar: {symbolName}, defCount: {symbols[symbolName].definitionCounterInBnf} "
                print(msg)
                errors.append(msg)

        # one symbol definition is max a few lines long, not a long string,
        # so this naive string concatenate is not a problem.
        if symbolName:  # not the empty non-defined:
            symbols[symbolName].definitionInBnf += definitionInLine

    return symbols, sorted(symbolNamesInLocalDefinition), errors


def symbols_table_print(symbols: dict[str, Symbol]):
    """display statistics about symbols to see where do we have too big set"""
    for symbolName, symbol in symbols.items():
        possibilities = symbols[symbolName].expandPossibilities()
        maxElemInOnePossibility = 0
        for onePosssibility in possibilities:
            maxElemInOnePossibility = max(maxElemInOnePossibility, len(onePosssibility))
        print(f"{symbolName:>50} possib: {len(possibilities)}",  possibilities)


def symbolnames_possibilities_print(possibilities: list[list[str]], prefix="", caller=""):
    """print list of possibilities"""
    for onePosssibility in possibilities:
        print(f"{prefix} one possibility ({caller}): {onePosssibility}")
