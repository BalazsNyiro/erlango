#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys, os
sys.path.append("bnf_validator")
import bnf_lib

"""convert BNF grammar files to a merged Lark Parser file"""

def main():
    errors = []

    symbolsTableAll = filenames_and_grammar_src_collect__without_locals(errors)

    for symbolName, symbol in symbolsTableAll.items():
        larkGrammar = bnf_to_lark_converter(symbolName, symbolsTableAll)
        print(f"Lark conversion: ", larkGrammar)

    if errors:
        for err in errors:
            print(err)


def filenames_and_grammar_src_collect__without_locals(errors) -> tuple[str, list[str]]:
    """

    The local symbols are limited terminating symbol sets,
    to validate the separated grammar files.

    They are fully defined in other .bnf files, so the locals aren't inserted
    into the merged grammar.

    """

    symbolsTableAll = dict()

    # Collect BNF based symbol definitions
    for filePathBnf in [f for f in os.listdir('.') if os.path.isfile(f) and f.startswith('grammar_') and f.endswith('.bnf')]:
        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        symbolsTableAll.update(symbolsTable)  # insert new keys into global collector

    return symbolsTableAll

def bnf_to_lark_converter(symbolName: str, symbols: dict[str, bnf_lib.Symbol]):
    """convert bnf grammar to lark """
    larkGrammar = []

    possibilitiesBnf = symbols[symbolName].expandPossibilitiesInBnf()

    def symbolNameConvertToLark(symbolNameBnf: str) -> str:
        """convert non-terminating symbol names to lark. Terminatings are similar"""
        if symbolNameBnf.startswith("<") and symbolNameBnf.endswith(">"):
            return symbolNameBnf[1:-1]
        return symbolNameBnf  # there is no change for terminatings...

    for possibBnf in possibilitiesBnf:
        if len(larkGrammar) > 0:
            larkGrammar.append("|")

        for symbolNameInPossibility in possibBnf:
            larkGrammar.append(symbolNameConvertToLark(symbolNameInPossibility))

    return f"{symbolNameConvertToLark(symbolName)}: {" ".join(larkGrammar)}"


if __name__ == '__main__':  # pragma: no cover
    main()
