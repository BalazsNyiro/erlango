#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import unittest, bnf_validator, bnf_lib

filePath_prefix = "testfile_bnf_validator_output__"


# python3 bnf_validator_test.py Test_possible_accepted_language_elems_save
class Test_possible_accepted_language_elems_save(unittest.TestCase):

    # python3 bnf_validator_test.py Test_possible_accepted_language_elems_save.test_possible_accepted_language_elems_save___happy_path
    def test_possible_accepted_language_elems_save___happy_path(self):
        """use the happy path to check manually special cases,
        this one is not a real test.

        """
        errors: list[str] = []
        filePathBnf = "../grammar_40_simple_types.bnf"

        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        symbolName = "<atom>"
        fname_bnf_accepted, fname_log, errors = bnf_validator.possible_accepted_language_elems_save(
            symbolName, symbolsTable, filePath_prefix, limitOfSymbolsForTestCases=4,
            errors=errors, filePathBnf=filePathBnf
        )

        symbolName = "<string>"
        fname_bnf_accepted, fname_log, errors = bnf_validator.possible_accepted_language_elems_save(
            symbolName, symbolsTable, filePath_prefix, limitOfSymbolsForTestCases=4,
            errors=errors, filePathBnf=filePathBnf
        )

        # symbolName = "<stringQuoteTriple>"
        # fname_bnf_accepted, fname_log, errors = bnf_validator.possible_accepted_language_elems_save(
        #     symbolName, symbolsTable, filePath_prefix, limitOfSymbolsForTestCases=4,
        #     errors=errors, filePathBnf=filePathBnf
        # )



    # python3 bnf_validator_test.py Test_possible_accepted_language_elems_save.test_possible_accepted_language_elems_save___missingDefInSymbol
    def test_possible_accepted_language_elems_save___missingDefInSymbol(self):
        """the wanted symbol doesn't have a definition"""
        testName = "test_possible_accepted_language_elems_save___missingDefInSymbol"
        print(f"======= {testName} ============")

        errors: list[str] = []
        filePathBnf = "bnf_lib_test_grammar_missing_definition.bnf"

        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        print(f"==== symbol table, {testName} ==========")
        bnf_lib.symbols_table_print(symbolsTable)


        symbolName = "<missing_definition_in_right_side>"

        fname_bnf_accepted, fname_log, errors = bnf_validator.possible_accepted_language_elems_save(
            symbolName, symbolsTable, filePath_prefix, limitOfSymbolsForTestCases=10,
            errors=errors, filePathBnf=filePathBnf
        )

        print("missing definition detection, ERRORS:", errors)
        self.assertTrue(errors[0].startswith("ERROR: <missing_definition_in_right_side> symbol: undefined expansion rules after ::="))


    # python3 bnf_validator_test.py Test_possible_accepted_language_elems_save.test_possible_accepted_language_elems_save___missingDefInChildrenSymbol
    def test_possible_accepted_language_elems_save___missingDefInChildrenSymbol(self):
        """the wanted symbol's children doesn't have a definition"""
        errors: list[str] = []
        filePathBnf = "bnf_lib_test_grammar_missing_definition.bnf"

        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        symbolName = "<missing_definition_in_child_elem>"

        fname_bnf_accepted, fname_log, errors = bnf_validator.possible_accepted_language_elems_save(
            symbolName, symbolsTable, filePath_prefix, limitOfSymbolsForTestCases=10,
            errors=errors, filePathBnf=filePathBnf
        )

        self.assertTrue(errors[0].startswith("ERROR: <missing_definition_in_right_side> symbol (reached in child expansion): missing expansion rules, nothing is defined after ::= "))




if __name__ == '__main__':  # pragma: no cover
    unittest.main()         # pragma: no cover
