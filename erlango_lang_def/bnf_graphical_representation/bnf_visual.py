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

def main(filePathBnf: str):
    print(f"BNF def file: {filePathBnf}")

    for line in file_src_lines(filePathBnf):
        print(line)

        if line.startswith("#"):
            print("--> commented")
            continue  # comment line








def file_src_lines(path: str) -> [str]:
    lines = []
    with open(path, 'r') as file:
        for line in file.readlines():
            lines.append(line.strip())
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