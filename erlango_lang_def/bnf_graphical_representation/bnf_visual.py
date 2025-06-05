#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os, argparse

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


def main(filePathBnf: str):

    symbols, errors = symbol_detect(filePathBnf)

    for symbolName, symbol in symbols.items():
        print()
        print(f"detected symbol: {symbolName:>{Symbol.symbolNameMax}}")
        for defLine in symbol.definitionInBnf.split("\n"):
            print(f"    {defLine.strip()}")  # to use standard indentation


    ################################################
    if not errors:
        print(f"No problem detected in the BNF")

    for err in errors:
        print(f"ERROR: {err}")


def symbol_detect(filePathBnf: str):
    """collect symbols and definitions from the bnf file"""

    print(f"BNF def file: {filePathBnf}")

    errors = list()
    symbols = dict()
    ################################################

    symbolName = ""

    for line in file_src_lines(filePathBnf):
        if line.startswith("#"):
            continue  # comment line

        symbolDefInLine = line
        if "::=" in line:
            symbolName, symbolDefInLine = line.split("::=")
            symbolName = symbolName.strip()

            if symbolName not in symbols:
                symbols[symbolName] = Symbol(symbolName)
            else:
                symbols[symbolName].definitionCounterInBnf += 1
                errors.append(f"problem: the symbol is defined more than once in the bnf grammar: {symbolName}, defCount: {symbols[symbolName].definitionCounterInBnf} ")


        # one symbol definition is max a few lines long, not a long string,
        # so this naive string concatenate is not a problem.
        if symbolName:  # not the empty non-defined:
            symbols[symbolName].definitionInBnf += symbolDefInLine

    return symbols, errors



def file_src_lines(path: str) -> [str]:
    lines = []
    with open(path, 'r') as file:
        lines = file.readlines()
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