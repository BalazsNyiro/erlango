#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os, argparse

"""
Analyse the passed BNF file to support the development.
"""

def main(filePathBnf: str):
    print(f"BNF def file: {filePathBnf}")



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