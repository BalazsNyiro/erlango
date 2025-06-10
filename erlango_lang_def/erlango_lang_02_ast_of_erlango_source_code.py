#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys, os
sys.path.append("bnf_validator")
import bnf_lib

from lark import Lark

def main(erlangoSrc: str):

    grammarLines = """
        value: dict
         | list
         | ESCAPED_STRING
         | SIGNED_NUMBER
         | "true" | "false" | "null"

    list : "[" [value ("," value)*] "]"

    dict : "{" [pair ("," pair)*] "}"
    pair : ESCAPED_STRING ":" value
    
    DIGIT: "0" | "1"

    %import common.ESCAPED_STRING
    %import common.SIGNED_NUMBER
    %import common.WS
    %ignore WS

    """
    grammarLines = "".join(bnf_lib.file_src_lines("erlango_lang_def.lark"))
    print(grammarLines)
    erlango_parser = Lark(grammarLines, start="integer")

    parsed = erlango_parser.parse(erlangoSrc)

if __name__ == '__main__':  # pragma: no cover

    erlangoSrc = """1"""

    main(erlangoSrc)
