#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys, os
sys.path.append("bnf_validator")
import bnf_lib

"""convert BNF grammar files to a merged Lark Parser file"""

def main():

    grammarFileNames_langDef = filenames_and_grammar_src_collect__without_locals()
    print(f"Erlang BNF files:")
    print(grammarFileNames_langDef)
    

def filenames_and_grammar_src_collect__without_locals() -> tuple[str, list[str]]:
    """filename and src collector.

    The local symbols are limited terminating symbol sets,
    to validate the separated grammar files.

    They are fully defined in other .bnf files, so the locals aren't inserted
    into the merged grammar.

    """
    patternLocalSection = "# LOCAL SYMBOLS"

    src = dict()
    errors = []

    symbolsTableAll = dict()

    # Collect BNF based symbol definitions
    for filePathBnf in [f for f in os.listdir('.') if os.path.isfile(f) and f.startswith('grammar_') and f.endswith('.bnf')]:
        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        symbolsTableAll.update(symbolsTable)  # insert new keys into global collector

    for symbolName, symbol in symbolsTableAll.items():
        print(f"global symbol table, collected name: {symbolName}", symbol.expandPossibilitiesInBnf())
        print(f"convert the bnf to LARK: ")

    return src

if __name__ == '__main__':  # pragma: no cover
    main()
