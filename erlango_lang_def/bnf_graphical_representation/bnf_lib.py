#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os

class Symbol:
    symbolNameMax = 0

    def __init__(self, symbolName):
        self.name = symbolName
        self.definitionInBnf = ""
        self.definitionCounterInBnf = 0  # in lucky case, the symbol is defined only once in the file

        Symbol.symbolNameMax = max(Symbol.symbolNameMax, len(symbolName))

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

        return expanded

    def grammar_elems_nonterminating_collect(self):
        return symbols_nonterminating_collect(self.definitionInBnf)


def symbols_nonterminating_count(grammarDefInBnf: str):
    """collect all <non-terminating> elems only in the given grammar"""
    tokensNonTerminating = re.findall(r'<[^<>]+>', grammarDefInBnf)
    return tokensNonTerminating


def symbolnames_concate_simple_str(names: [str]) -> str:
    """concatenate list of symbolnames into one str"""
    return ", ".join(names)


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


def is_terminating_symbolname(symbolName: str) -> bool:
    """is it a terminating symbol name?"""
    return symbolName.startswith('"') and symbolName.endswith('"')


def print_symbols_elem_stats(symbols: dict[str, Symbol]):
    """display statistics about symbols to see where do we have too big set"""
    for symbolName, symbol in symbols.items():
        possibilities = symbols[symbolName].expandPossibilities()
        maxElemInOnePossibility = 0
        for onePossible in possibilities:
            maxElemInOnePossibility = max(maxElemInOnePossibility, len(onePossible))
        print(f"{symbolName:>50} possib: {len(possibilities)}",  possibilities)


def file_src_lines(path: str) -> [str]:
    """read all lines of a file"""
    lines = []
    with open(path, 'r') as file:
        lines = file.readlines()  # Pycharm highlight if I return directly in this line
    return lines


def file_is_exists(path: str):
    """check if a file exists or not, generate an exception if there is a problem"""
    if not os.path.exists(path):
        msg = f"ERROR: invalid file path: {path}"
        print(msg)
        raise vauleError(msg)

def file_write(path: str, content: str):
    with open(path, mode="w") as f:
        f.write(content)

def files_collect_in_dir(directory, prefix="grammar_"):
    files_and_path = []
    files = [f for f in os.listdir(directory) if f.startswith(prefix) and os.path.isfile(os.path.join(directory, f))]

    for f in files:
        files_and_path.append(os.path.join(directory, f))

    return files_and_path


def symbols_detect_in_file(filePathBnf: str, errors: [str]) -> dict[str, Symbol]:
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
                msg = f"ERROR: the symbol is defined more than once in the bnf grammar: {symbolName}, defCount: {symbols[symbolName].definitionCounterInBnf} "
                print(msg)
                errors.append(msg)

        # one symbol definition is max a few lines long, not a long string,
        # so this naive string concatenate is not a problem.
        if symbolName:  # not the empty non-defined:
            symbols[symbolName].definitionInBnf += definitionInLine

    return symbols, errors

