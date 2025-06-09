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

    for fileName in [f for f in os.listdir('.') if os.path.isfile(f) and f.startswith('grammar_') and f.endswith('.bnf')]:
        linesLocalsRemoved = []
        
        for line in bnf_lib.file_src_lines(fileName):
            if line.startswith(patternLocalSection):
                break # the local section won't be added into the merged grammar

            # comments are not added
            if not line.strip().startswith("#"):
                linesLocalsRemoved.append(line)
            
        src[fileName] = linesLocalsRemoved
    return src

if __name__ == '__main__':  # pragma: no cover
    main()
